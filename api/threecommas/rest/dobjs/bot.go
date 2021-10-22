package dobjs

import "time"

//Bot is the data structure of a Bot payload from the 3commas api
type Bot struct {
	ID                      int      `json:"id"`
	AccountID               int      `json:"account_id"`
	IsEnabled               bool     `json:"is_enabled"`
	MaxSafetyOrders         int      `json:"max_safety_orders"`
	ActiveSafetyOrdersCount int      `json:"active_safety_orders_count"`
	Pairs                   []string `json:"pairs"`
	StrategyList            []struct {
		Options struct {
		} `json:"options"`
		Strategy string `json:"strategy"`
	} `json:"strategy_list"`
	MaxActiveDeals              int           `json:"max_active_deals"`
	ActiveDealsCount            int           `json:"active_deals_count"`
	Deletable                   bool          `json:"deletable?"`
	CreatedAt                   time.Time     `json:"created_at"`
	UpdatedAt                   time.Time     `json:"updated_at"`
	TrailingEnabled             bool          `json:"trailing_enabled"`
	TslEnabled                  bool          `json:"tsl_enabled"`
	DealStartDelaySeconds       int           `json:"deal_start_delay_seconds"`
	StopLossTimeoutEnabled      bool          `json:"stop_loss_timeout_enabled"`
	StopLossTimeoutInSeconds    int           `json:"stop_loss_timeout_in_seconds"`
	DisableAfterDealsCount      int           `json:"disable_after_deals_count"`
	DealsCounter                int           `json:"deals_counter"`
	AllowedDealsOnSamePair      int           `json:"allowed_deals_on_same_pair"`
	EasyFormSupported           bool          `json:"easy_form_supported"`
	CloseDealsTimeout           int           `json:"close_deals_timeout"`
	URLSecret                   string        `json:"url_secret"`
	Name                        string        `json:"name"`
	TakeProfit                  string        `json:"take_profit"`
	BaseOrderVolume             string        `json:"base_order_volume"`
	SafetyOrderVolume           string        `json:"safety_order_volume"`
	SafetyOrderStepPercentage   string        `json:"safety_order_step_percentage"`
	TakeProfitType              string        `json:"take_profit_type"`
	Type                        string        `json:"type"`
	MartingaleVolumeCoefficient string        `json:"martingale_volume_coefficient"`
	MartingaleStepCoefficient   string        `json:"martingale_step_coefficient"`
	StopLossPercentage          string        `json:"stop_loss_percentage"`
	Cooldown                    string        `json:"cooldown"`
	BtcPriceLimit               string        `json:"btc_price_limit"`
	Strategy                    string        `json:"strategy"`
	MinVolumeBtc24H             string        `json:"min_volume_btc_24h"`
	ProfitCurrency              string        `json:"profit_currency"`
	MinPrice                    string        `json:"min_price"`
	MaxPrice                    string        `json:"max_price"`
	StopLossType                string        `json:"stop_loss_type"`
	SafetyOrderVolumeType       string        `json:"safety_order_volume_type"`
	BaseOrderVolumeType         string        `json:"base_order_volume_type"`
	AccountName                 string        `json:"account_name"`
	TrailingDeviation           string        `json:"trailing_deviation"`
	FinishedDealsProfitUsd      string        `json:"finished_deals_profit_usd"`
	FinishedDealsCount          string        `json:"finished_deals_count"`
	LeverageType                string        `json:"leverage_type"`
	LeverageCustomValue         string        `json:"leverage_custom_value"`
	StartOrderType              string        `json:"start_order_type"`
	ActiveDealsUsdProfit        string        `json:"active_deals_usd_profit"`
	ActiveDeals                 []interface{} `json:"active_deals"`
}
