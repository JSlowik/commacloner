package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jslowik/commacloner/api/threecommas/websockets"
	dobjs2 "github.com/jslowik/commacloner/api/threecommas/websockets/dobjs"
	"io/ioutil"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

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
	//l, err := log.InitWithConfiguration(c.Logging.Level, c.Logging.Format)
	err = log.InitWithConfiguration(c.Logging)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}
	logger := log.NewLogger("serve")
	logger.Info("logging configured")

	//log mappings
	logger.Info("loading bot mappings")
	botMap := make(map[int][]config.BotMapping)
	for _, mapping := range c.Bots {
		botMap[mapping.Source.ID] = append(botMap[mapping.Source.ID], mapping)
	}

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

	messageOut := make(chan *dobjs2.IdentifierMessage)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, err := generateConnection(nil, c.API.WebsocketURL, logger)
	if err != nil {
		logger.Fatalf("could not make connection to websocket: %v", err)
		return err
	}

	//When the program closes close the connection
	defer conn.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		pong := dobjs2.IdentifierMessage{Type: "pong"}

		for {
			msgType, message, readErr := conn.ReadMessage()
			if readErr != nil {
				if !websocket.IsCloseError(readErr, websocket.CloseNormalClosure) {
					logger.Warnf("abonormal close error. trying resubscribe: %v", readErr)
					conn, err = generateConnection(conn, c.API.WebsocketURL, logger)
					if err != nil {
						logger.Fatalf("could not regenerate connection: %v", err)
						return
					}
					logger.Infof("connection restablished")
					continue
				}
			}
			logger.Debugf("recv: type - %d message - %s", msgType, message)

			ctrlMessage := dobjs2.Message{}
			pingMessage := dobjs2.PingMessage{}
			if unmarshalError := json.Unmarshal(message, &ctrlMessage); unmarshalError == nil {
				switch ctrlMessage.Type {
				case "welcome":
					logger.Infof("received welcome, sending subscription: %s", subscriptionMessage)
					messageOut <- subscriptionMessage
				case "confirm_subscription":
					logger.Infof("subscription confirmed : %s", message)
				case "Deal", "Deal::ShortDeal":
					logger.Debugf("received deal %v", ctrlMessage.Message)
					dealMessage := dobjs2.DealsMessage{}
					var dealErr error
					if dealErr = json.Unmarshal(message, &dealMessage); dealErr == nil {
						dealErr = stream.HandleDeal(dealMessage)
					}
					if dealErr != nil {
						logger.Errorf("could not handle message from deals stream: %v", dealErr)
					}
				default:
					logger.Warnf("unsupported message type %s : %v", ctrlMessage.Type, string(message))
				}

			} else if e := json.Unmarshal(message, &pingMessage); e == nil {
				logger.Debugf("received ping, sending pong: %s", message)
				messageOut <- &pong
			}
		}
	}()

	for {
		select {
		case <-done:
			return nil
		case m := <-messageOut:
			logger.Debugf("Send IdentifierMessage %s", m)
			err := conn.WriteJSON(m)
			if err != nil {
				logger.Errorf("write message out failure: %v", err)
				return err
			}
		case <-interrupt:
			logger.Infof("interrupt")
			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				logger.Errorf("write close failure: %v", err)
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

func generateConnection(existingConnection *websocket.Conn, url string, logger *zap.SugaredLogger) (*websocket.Conn, error) {
	if existingConnection != nil {
		logger.Warnf("closing existing connection")
		err := existingConnection.Close()
		if err != nil {
			logger.Errorf("error closing existing connection: %v", err)
		}
	}

	logger.Infof("connecting to %s", url)
	conn, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		logger.Errorf("handshake failed with status %d, error %v", resp.StatusCode, err)
		logger.Fatalf("dial: %v ", err)
	}
	return conn, err
}
