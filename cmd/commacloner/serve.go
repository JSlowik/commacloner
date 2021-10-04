package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/jslowik/commacloner/api/websockets"
	"github.com/jslowik/commacloner/recws"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/ghodss/yaml"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
	"github.com/spf13/cobra"
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

var (
	ws *recws.RecConn
)

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

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt) // Notify the interrupt channel for SIGINT

	//Make the subscription message
	stream := websockets.DealsStream{
		APIConfig: c.API,
		Bots:      botMap,
	}
	subscriptionMessage, err := stream.Build()
	if err != nil {
		logger.Fatalf("could not build deal subscription: %v", err)
		return err
	}

	ctx := context.TODO()
	ws := recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}

	ws.SubscribeHandler = func() error {
		return ws.WriteJSON(subscriptionMessage)
	}

	ws.Dial(c.API.WebsocketURL, nil)
	//go func() {
	//	time.Sleep(2 * time.Second)
	//	cancel()
	//}()

	for {
		select {
		case <-ctx.Done():
			go ws.Close()
			logger.Warnf("Websocket closed %s", ws.GetURL())
			return nil
		default:
			if !ws.IsConnected() {
				logger.Errorf("Websocket disconnected %s", ws.GetURL())
				continue
			}

			_, message, err := ws.ReadMessage()
			if err != nil {
				logger.Errorf("Error: ReadMessage: %v", err)
				continue
			}

			logger.Debugf("recv: %s", message)
			err = stream.HandleDeal(message, l)
			if err != nil {
				logger.Errorf("could not handle message from deals stream: %v", err)
			}


			logger.Infof("Success: %s", message)
		}
	}
}
