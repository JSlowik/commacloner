package websockets

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
	"github.com/jslowik/commacloner/api"
	"github.com/jslowik/commacloner/api/rest"
	"github.com/jslowik/commacloner/config"
	"go.uber.org/zap"
)

const dealsEndpoint = "/deals"

// DealsStream a websocket stream to listen for new deals
type DealsStream struct {
	APIConfig config.API
	Bots      map[int][]config.BotMapping
}

// BuildSignature computes the signature for the websocket subscription message
func (d DealsStream) BuildSignature(endpoint string) string {
	return api.ComputeSignature(endpoint, d.APIConfig.Secret)
}

// buildIdentifier builds the stream identifier
func (d DealsStream) buildIdentifier() Identifier {
	signature := d.BuildSignature(dealsEndpoint)
	return Identifier{
		Channel: "DealsChannel",
		Users: []User{
			{
				APIKey:    d.APIConfig.Key,
				Signature: signature,
			},
		},
	}
}

// Build constructs the deals websocket subscription
func (d DealsStream) Build() (*Message, error) {
	identifier := d.buildIdentifier()
	identifierStr, err := json.Marshal(identifier)
	if err != nil {
		return nil, fmt.Errorf("could not marshall identifier: %v", err)
	}

	return &Message{
		Identifier: string(identifierStr),
		Command:    "subscribe",
	}, nil
}

// ReceiveDeal reads messages from the websocket connection and handles the deal
func (d DealsStream) ReceiveDeal(connection *websocket.Conn, logger *zap.Logger, done chan interface{}) {
	log := logger.Sugar()
	defer close(done)
	for {
		_, rawMessage, err := connection.ReadMessage()
		if err != nil {
			log.Error("error in receive: ", zap.Error(err))
			break
		}

		// Ping Messages
		var ping api.Ping
		err = json.Unmarshal(rawMessage, &ping)
		if err == nil {
			continue
		}

		// Deals Messages
		var deal api.Deal
		err = json.Unmarshal(rawMessage, &deal)
		if err == nil {
			details := deal.Details
			if details.Status == "bought" && details.CompletedSafetyOrdersCount == 0 && details.CompletedManualSafetyOrdersCount == 0 {
				log.Infof("got new deal - bot id: %d, pair %s", details.BotID, details.Pair)
				if d.Bots == nil {
					log.Error("no bots defined")
				}
				// Determine if we have a mapping which uses this source deal
				for _, bot := range d.Bots[details.BotID] {
					log.Infof("start new deal for bot %d using pair %s)", bot.Destination, details.Pair)
					err := rest.StartNewDeal(d.APIConfig, bot, details.Pair, logger)
					if err != nil {
						log.Errorf("could not start new deal: %v", err)
						if bot.Overrides.CancelUnavailableDeals {
							e2 := rest.CancelDeal(d.APIConfig, details.ID, bot.Overrides.PanicSellUnavailableDeals, logger)
							if e2 != nil {
								log.Errorf("could not cancel deal: %v", e2)
							}
						}
					}
				}
			}
			continue
		} else {
			log.Error("got error", zap.Error(err))
		}

		// Unsupported Message Type
		log.Warnf("unsupported message type %v", rawMessage)
	}
}
