package dobjs

import "time"

//Exchange is the return payload of an API call to get exchange info
type Exchange struct {
	ID                             int         `json:"id"`
	AutoBalancePeriod              int         `json:"auto_balance_period"`
	AutoBalancePortfolioID         interface{} `json:"auto_balance_portfolio_id"`
	AutoBalanceCurrencyChangeLimit interface{} `json:"auto_balance_currency_change_limit"`
	AutobalanceEnabled             bool        `json:"autobalance_enabled"`
	HedgeModeAvailable             bool        `json:"hedge_mode_available"`
	HedgeModeEnabled               bool        `json:"hedge_mode_enabled"`
	IsLocked                       bool        `json:"is_locked"`
	SmartTradingSupported          bool        `json:"smart_trading_supported"`
	SmartSellingSupported          bool        `json:"smart_selling_supported"`
	AvailableForTrading            struct {
	} `json:"available_for_trading"`
	StatsSupported          bool        `json:"stats_supported"`
	TradingSupported        bool        `json:"trading_supported"`
	MarketBuySupported      bool        `json:"market_buy_supported"`
	MarketSellSupported     bool        `json:"market_sell_supported"`
	ConditionalBuySupported bool        `json:"conditional_buy_supported"`
	BotsAllowed             bool        `json:"bots_allowed"`
	BotsTtpAllowed          bool        `json:"bots_ttp_allowed"`
	BotsTslAllowed          bool        `json:"bots_tsl_allowed"`
	GordonBotsAvailable     bool        `json:"gordon_bots_available"`
	MultiBotsAllowed        bool        `json:"multi_bots_allowed"`
	CreatedAt               time.Time   `json:"created_at"`
	UpdatedAt               time.Time   `json:"updated_at"`
	LastAutoBalance         interface{} `json:"last_auto_balance"`
	FastConvertAvailable    bool        `json:"fast_convert_available"`
	GridBotsAllowed         bool        `json:"grid_bots_allowed"`
	APIKeyInvalid           bool        `json:"api_key_invalid"`
	NomicsID                string      `json:"nomics_id"`
	MarketIcon              string      `json:"market_icon"`
	DepositEnabled          bool        `json:"deposit_enabled"`
	SupportedMarketTypes    []string    `json:"supported_market_types"`
	APIKey                  string      `json:"api_key"`
	Name                    string      `json:"name"`
	AutoBalanceMethod       interface{} `json:"auto_balance_method"`
	AutoBalanceError        interface{} `json:"auto_balance_error"`
	CustomerID              interface{} `json:"customer_id"`
	SubaccountName          interface{} `json:"subaccount_name"`
	LockReason              interface{} `json:"lock_reason"`
	BtcAmount               string      `json:"btc_amount"`
	UsdAmount               string      `json:"usd_amount"`
	DayProfitBtc            string      `json:"day_profit_btc"`
	DayProfitUsd            string      `json:"day_profit_usd"`
	DayProfitBtcPercentage  string      `json:"day_profit_btc_percentage"`
	DayProfitUsdPercentage  string      `json:"day_profit_usd_percentage"`
	BtcProfit               string      `json:"btc_profit"`
	UsdProfit               string      `json:"usd_profit"`
	UsdProfitPercentage     string      `json:"usd_profit_percentage"`
	BtcProfitPercentage     string      `json:"btc_profit_percentage"`
	TotalBtcProfit          string      `json:"total_btc_profit"`
	TotalUsdProfit          string      `json:"total_usd_profit"`
	PrettyDisplayType       string      `json:"pretty_display_type"`
	ExchangeName            string      `json:"exchange_name"`
	MarketCode              string      `json:"market_code"`
}
