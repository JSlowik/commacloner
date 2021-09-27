package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/ghodss/yaml"
	"github.com/gorilla/websocket"
	"github.com/jslowik/commacloner/api/websockets"
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

	//connect to the websocket
	logger.Infof("connecting to websocket: %s", c.API.WebsocketURL)
	conn, _, err := websocket.DefaultDialer.Dial(c.API.WebsocketURL, nil)
	if err != nil {
		logger.Fatal("error connecting to websocket server: ", zap.Error(err))
	}
	defer conn.Close()

	done := make(chan interface{})    // Channel to indicate that the receiverHandler is done

	//Send the subscription message
	stream := websockets.DealsStream{
		APIConfig: c.API,
		Bots:      botMap,
	}
	message, err := stream.Build()
	e := conn.WriteJSON(message)
	if e != nil {
		logger.Errorf("cannot subscribe to topic %v", e)
		return e
	}

	//Create Deal Listener
	go stream.ReceiveDeal(conn, l, done)

	for {
		select {
		case <-done:
			return nil
		case <-interrupt:
			logger.Warn("interrupted")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Errorf("error writing close message: %v", err)
				return err
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return nil
		}
	}
}

