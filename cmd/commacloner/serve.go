package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jslowik/commacloner/api/lunarcrush"
	tcApi "github.com/jslowik/commacloner/api/threecommas/rest"
	"github.com/jslowik/commacloner/api/threecommas/websockets"
	wsapi "github.com/jslowik/commacloner/api/threecommas/websockets/dobjs"
	"io/ioutil"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"

	"github.com/ghodss/yaml"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
	"github.com/spf13/cobra"

	"github.com/go-co-op/gocron"
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
	err = log.InitWithConfiguration(c.Logging)
	if err != nil {
		return fmt.Errorf("invalid config: %v", err)
	}
	logger := log.NewLogger("serve")
	logger.Info("logging configured")

	//Init LunarCrush if configured
	blacklist := make(map[string] bool)
	for _ ,v := range c.LunarcrushAPI.Blacklist {
		blacklist[v] = true
	}

	lunarCache := &lunarcrush.LunarCache{
		Logger: log.NewLogger("lunarCache"),
		Config: c.LunarcrushAPI,
		Blacklist: blacklist,
	}

	scheduler := gocron.NewScheduler(time.UTC)
	if c.LunarcrushAPI.Cache.Enabled {
		//Block for first cache load
		logger.Infof("initializing LunarCrush cache")
		lunarCache.UpdateCache()
		logger.Infof("LunarCrush cache initialized")

		//Schedule for Cache updates
		_, _ = scheduler.Every(c.LunarcrushAPI.Cache.Every).WaitForSchedule().Do(lunarCache.UpdateCache)
	}


	//log mappings
	logger.Info("loading bot mappings")
	botMap := make(map[int][]config.BotMapping)
	for _, mapping := range c.Bots {
		botMap[mapping.Source.ID] = append(botMap[mapping.Source.ID], mapping)
		if mapping.Pairs.Mode == "lunarcrush" {
			//Schedule an Updater for the bot
			_, _ = scheduler.Every(mapping.Pairs.Config.Refresh).Tag(fmt.Sprintf("updatePairs_%s", mapping.ID)).Do(UpdateMapping,c.API,mapping,lunarCache)
		}
	}

	//Start Scheduler
	scheduler.StartAsync()


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

	messageOut := make(chan *wsapi.IdentifierMessage)
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	conn, err := generateConnection(nil, c.API.WebsocketURL, logger)
	if err != nil {
		logger.Fatalf("could not make connection to websocket: %v", err)
		return err
	}

	//When the program closes, close the connection
	defer conn.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		pong := wsapi.IdentifierMessage{Type: "pong"}

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

			ctrlMessage := wsapi.Message{}
			pingMessage := wsapi.PingMessage{}
			if unmarshalError := json.Unmarshal(message, &ctrlMessage); unmarshalError == nil {
				switch ctrlMessage.Type {
				case "welcome":
					logger.Infof("received welcome, sending subscription: %s", subscriptionMessage)
					messageOut <- subscriptionMessage
				case "confirm_subscription":
					logger.Infof("subscription confirmed : %s", message)
				case "Deal", "Deal::ShortDeal":
					logger.Debugf("received deal %v", ctrlMessage.Message)
					dealMessage := wsapi.DealsMessage{}
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

func UpdateMapping(apiConfig config.API,mapping config.BotMapping, cache *lunarcrush.LunarCache) error{
	logger := log.NewLogger(fmt.Sprintf("updatePairs_%s",mapping.ID))

	exchangeMap, err := tcApi.GetExchangeAccounts(apiConfig)
	if err != nil {
		logger.Errorf("cannot get exchange accounts: %v",  err)
	}

	//Get the Destination Bot
	destBot, err := tcApi.GetBot(apiConfig, mapping.Destination.ID)
	if err != nil {
		logger.Errorf("cannot get info for bot %d: %v", mapping.Destination, err)
	}

	//Get the Source Bot
	sourceBot, err := tcApi.GetBot(apiConfig, mapping.Source.ID)
	if err != nil {
		logger.Errorf("cannot get info for bot %d: %v", mapping.Destination, err)
	}

	//Get the pairs of the destination exchange
	marketCode := exchangeMap[destBot.AccountID].MarketCode
	destPairs, err := tcApi.GetExchangePairs(apiConfig,marketCode)
	if err != nil {
		return fmt.Errorf("could not get exchange pairs: %v", err)
	}

	pMap := make(map[string][]string)
	for _,p := range destPairs {
		pairSplit := strings.Split(p, "_")
		if pMap[pairSplit[1]] == nil {
			pMap[pairSplit[1]] = make([]string,0)
		}
		pMap[pairSplit[1]] = append(pMap[pairSplit[1]] , pairSplit[0])
	}
	cat := mapping.Pairs.Config.Category

	var validPairs []string
	switch cat {
	case "galaxyscore":
		logger.Infof("update %s by galaxy score", mapping.ID)
		validPairs = cache.GetByGalaxyScore(pMap,mapping.Pairs.Config.MaxPairs)
		logger.Infof("pairs by galaxy score: %v", validPairs)
	case "altrank":
		logger.Infof("update %s by altrank", mapping.ID)
		validPairs = cache.GetByAltRank(pMap,mapping.Pairs.Config.MaxPairs)
		logger.Infof("pairs by altrank: %v", validPairs)
	default:
		logger.Errorf("invalid update pairs criteria: %s", cat)
	}

	if len(validPairs) > 0 {
		//Update Source bot
		err := tcApi.UpdatePairs(apiConfig,sourceBot, mapping.Pairs.Config.QuoteCurrency.Source,validPairs)
		if err != nil {
			logger.Errorf("could not update source bot %d: %v", mapping.Source.ID, err)
		}

		//Update Destination Bot
		err = tcApi.UpdatePairs(apiConfig,destBot, mapping.Pairs.Config.QuoteCurrency.Dest,validPairs)
		if err != nil {
			logger.Errorf("could not update source bot %d: %v", mapping.Source.ID, err)
		}
	}

	return nil
}