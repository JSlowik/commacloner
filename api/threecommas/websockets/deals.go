package websockets

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jslowik/commacloner/api/threecommas"
	"github.com/jslowik/commacloner/api/threecommas/rest"
	dobjs2 "github.com/jslowik/commacloner/api/threecommas/websockets/dobjs"
	"github.com/jslowik/commacloner/config"
	"github.com/jslowik/commacloner/log"
)

const dealsEndpoint = "/deals"

// DealsStream a websocket stream to listen for new deals
type DealsStream struct {
	APIConfig config.API
	Bots      map[int][]config.BotMapping
}

// BuildSignature computes the signature for the websocket subscription message
func (d DealsStream) BuildSignature(endpoint string) string {
	return threecommas.ComputeSignature(endpoint, d.APIConfig.Secret)
}

// buildIdentifier builds the stream identifier
func (d DealsStream) buildIdentifier() dobjs2.Identifier {
	signature := d.BuildSignature(dealsEndpoint)
	return dobjs2.Identifier{
		Channel: "DealsChannel",
		Users: []dobjs2.User{
			{
				APIKey:    d.APIConfig.Key,
				Signature: signature,
			},
		},
	}
}

// Build constructs the deals websocket subscription
func (d DealsStream) Build() (*dobjs2.IdentifierMessage, error) {
	identifier := d.buildIdentifier()
	identifierStr, err := json.Marshal(identifier)
	if err != nil {
		return nil, fmt.Errorf("could not marshall identifier: %v", err)
	}

	return &dobjs2.IdentifierMessage{
		Identifier: string(identifierStr),
		Command:    "subscribe",
	}, nil
}

// HandleDeal reads messages from the websocket connection and handles the deal
func (d DealsStream) HandleDeal(deal dobjs2.DealsMessage) error {
	logger := log.NewLogger("deals")

	details := deal.Details
	if details.Status == "bought" && details.CompletedSafetyOrdersCount == 0 && details.CompletedManualSafetyOrdersCount == 0 {
		logger.Infof("got new deal - bot id: %d, pair %s", details.BotID, details.Pair)
		if d.Bots == nil {
			return errors.New("no bots defined")
		}
		// Determine if we have a mapping which uses this source deal
		for _, bot := range d.Bots[details.BotID] {
			logger.Infof("start new deal for bot %d using pair %s", bot.Destination.ID, details.Pair)
			err := rest.StartNewDeal(d.APIConfig, bot, details.Pair)
			if err != nil {
				logger.Warnf("could not start new deal: %v", err)
				if bot.Overrides.CancelUnavailableDeals {
					e2 := rest.CancelDeal(d.APIConfig, details.ID, bot.Overrides.PanicSellUnavailableDeals)
					if e2 != nil {
						return fmt.Errorf("could not cancel deal: %v", e2)
					}
				}
			}
		}
	}
	return nil
}
