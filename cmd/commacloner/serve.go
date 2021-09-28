package main

import (
	"errors"
	"fmt"
	"github.com/jslowik/commacloner/api/websockets"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gorilla/websocket"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func commandServe() *cobra.Command {
	return &cobra.Command{
		Use:     "serve [ config file ]",
		Short:   "Connect to 3commas and begin managing deals.",
		Long:    ``,
		Example: "commacloner serve config.yaml",
		Run: func(cmd *cobra.Command, args []string) {
			if err := serve(args); err != nil {
				fmt.Fprintln(os.Stderr, err)
				os.Exit(2)
			}
		},
	}
}

func serve(args []string) error {
	switch len(args) {
	default:
		return errors.New("surplus arguments")
	case 0:
		return errors.New("no arguments provided")
	case 1:
	}

	configFile := args[0]
	configData, err := ioutil.ReadFile(configFile)
	if err != nil {
		return fmt.Errorf("failed to read config file %s: %v", configFile, err)
	}

	var c config.Config

	data := []byte(os.ExpandEnv(string(configData)))
	if err := yaml.Unmarshal(data, &c); err != nil {
		return fmt.Errorf("error parse config file %s: %v", configFile, err)
	}
	if err := c.Validate(); err != nil {
		return err
	}


	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	//init logging
	l, err := log.InitWithConfiguration(c.Logging.Level, c.Logging.Format)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}
	logger := l.Sugar()
	logger.Info("logging configured")

	//log mappings
	logger.Info("loading bot mappings")
	botMap := make(map[int][]config.BotMapping)
	for _, mapping := range c.Bots {
		botMap[mapping.Source.ID] = append(botMap[mapping.Source.ID], mapping)
	}

	logger.Infof("connecting to websocket: %s", c.API.WebsocketURL)
	conn, _, err := websocket.DefaultDialer.Dial(c.API.WebsocketURL, nil)
	if err != nil {
		logger.Fatal("error connecting to websocket server: ", zap.Error(err))
	}
	defer conn.Close()
	done := make(chan struct{})    // Channel to indicate that the receiverHandler is done

	//Send the subscription message
	stream := websockets.DealsStream{
		APIConfig: c.API,
		Bots:      botMap,
	}
	subscriptionMessage, err := stream.Build()
	if err != nil {
		logger.Fatalf("could not build deal subscription: %v", err)
		return err
	}

	e := conn.WriteJSON(subscriptionMessage)
	if e != nil {
		logger.Errorf("cannot subscribe to topic %v", e)
		return e
	}

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				logger.Errorf("read error %v", err)
				return
			}
			logger.Infof("recv: %s", message)
			err = stream.HandleDeal(message, l)
			if err != nil {
				logger.Errorf("could not handle message from deals stream: %v", err )
			}
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			logger.Infof("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			closeError := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if closeError != nil {
				logger.Errorf("write close: %v", closeError)
				return closeError
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}

}

