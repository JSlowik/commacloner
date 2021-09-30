package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jslowik/commacloner/api/rest"
	"go.uber.org/zap"

	"github.com/jslowik/commacloner/api"
	"github.com/jslowik/commacloner/config"
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

// HandleDeal reads messages from the websocket connection and handles the deal
func (d DealsStream) HandleDeal(rawMessage []byte, logger *zap.Logger) error {
	log := logger.Sugar()

	// Ping Messages
	var ping api.Ping
	err := json.Unmarshal(rawMessage, &ping)
	if err == nil {
		return nil
	}

	// Deals Messages
	var deal api.Deal
	err = json.Unmarshal(rawMessage, &deal)
	if err == nil {
		details := deal.Details
		if details.Status == "bought" && details.CompletedSafetyOrdersCount == 0 && details.CompletedManualSafetyOrdersCount == 0 {
			log.Infof("got new deal - bot id: %d, pair %s", details.BotID, details.Pair)
			if d.Bots == nil {
				return errors.New("no bots defined")
			}
			// Determine if we have a mapping which uses this source deal
			for _, bot := range d.Bots[details.BotID] {
				log.Infof("start new deal for bot %d using pair %s)", bot.Destination, details.Pair)
				err = rest.StartNewDeal(d.APIConfig, bot, details.Pair, logger)
				if err != nil {
					log.Warnf("could not start new deal: %v", err)
					if bot.Overrides.CancelUnavailableDeals {
						e2 := rest.CancelDeal(d.APIConfig, details.ID, bot.Overrides.PanicSellUnavailableDeals, logger)
						if e2 != nil {
							return fmt.Errorf("could not cancel deal: %v", e2)
						}
					}
				}
			}
		}
		return nil
	}

	// Unsupported Message Type
	log.Warnf("unsupported message type %v", string(rawMessage))
	return nil
}
