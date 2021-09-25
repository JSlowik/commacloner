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
	logger, err := log.InitWithConfiguration(c.Logging.Level, c.Logging.Format)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}

	if err := c.Validate(); err != nil {
		return err
	}
	logger.Info("logging configured")

	logger.Info("loading bot mappings")
	botMap := make(map[int][]config.BotMapping)
	for _, mapping := range c.Bots {
		botMap[mapping.Source.ID] = append(botMap[mapping.Source.ID], mapping)
	}

	// Set up the websocket
	stream := websockets.DealsStream{
		APIConfig: c.API,
		Bots:      botMap,
	}
	message, err := stream.Build()
	if err != nil {
		logger.Fatal("could not build deal stream", zap.Error(err))
	}

	logger.Debug("subscribing with message", zap.Any("message", message))

	// Connect to the websocket
	done = make(chan interface{})    // Channel to indicate that the receiverHandler is done
	interrupt = make(chan os.Signal) // Channel to listen for interrupt signal to terminate gracefully

	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	conn, _, err := websocket.DefaultDialer.Dial(c.API.WebsocketURL, nil)
	if err != nil {
		logger.Fatal("error connecting to websocket server: ", zap.Error(err))
	}
	defer conn.Close()
	go stream.ReceiveDeal(conn, logger, done)

	e := conn.WriteJSON(message)
	if e != nil {
		logger.Error("cannot subscribe ", zap.Error(e))
	}

	for {
		select {
		case <-interrupt:
			// We received a SIGINT (Ctrl + C). Terminate gracefully...
			logger.Info("Received SIGINT interrupt signal. Closing all pending connections")

			// Close our websocket connection
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Error("Error during closing websocket:", zap.Error(err))
				return nil
			}

			select {
			case <-done:
				logger.Info("Receiver Channel Closed! Exiting....")
			case <-time.After(time.Duration(1) * time.Second):
				logger.Warn("Timeout in closing receiving channel. Exiting....")
			}
		}
	}
}

var (
	done      chan interface{}
	interrupt chan os.Signal
)
