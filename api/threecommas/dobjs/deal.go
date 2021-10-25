package dobjs

//Deal describes a deal returned by the 3commas api
type Deal struct {
	ID                               int    `json:"id"`
	Type                             string `json:"type"`
	BotID                            int    `json:"bot_id"`
	CompletedSafetyOrdersCount       int    `json:"completed_safety_orders_count"`
	CompletedManualSafetyOrdersCount int    `json:"completed_manual_safety_orders_count"`
	Pair                             string `json:"pair"`
	Status                           string `json:"status"`
}


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
	//MaxActiveDeals              int           `json:"max_active_deals"`
	//ActiveDealsCount            int           `json:"active_deals_count"`
	//Deletable                   bool          `json:"deletable?"`
	//CreatedAt                   time.Time     `json:"created_at"`
	//UpdatedAt                   time.Time     `json:"updated_at"`
	//TrailingEnabled             bool          `json:"trailing_enabled"`
	//TslEnabled                  bool          `json:"tsl_enabled"`
	//DealStartDelaySeconds       interface{}   `json:"deal_start_delay_seconds"`
	//StopLossTimeoutEnabled      bool          `json:"stop_loss_timeout_enabled"`
	//StopLossTimeoutInSeconds    int           `json:"stop_loss_timeout_in_seconds"`
	//DisableAfterDealsCount      interface{}   `json:"disable_after_deals_count"`
	//DealsCounter                interface{}   `json:"deals_counter"`
	//AllowedDealsOnSamePair      int           `json:"allowed_deals_on_same_pair"`
	//EasyFormSupported           bool          `json:"easy_form_supported"`
	//CloseDealsTimeout           interface{}   `json:"close_deals_timeout"`
	//URLSecret                   string        `json:"url_secret"`
	Name                        string        `json:"name"`
	TakeProfit                  string        `json:"take_profit"`
	BaseOrderVolume             string        `json:"base_order_volume"`
	SafetyOrderVolume           string        `json:"safety_order_volume"`
	SafetyOrderStepPercentage   string        `json:"safety_order_step_percentage"`
	TakeProfitType              string        `json:"take_profit_type"`
	//Type                        string        `json:"type"`
	MartingaleVolumeCoefficient string        `json:"martingale_volume_coefficient"`
	MartingaleStepCoefficient   string        `json:"martingale_step_coefficient"`
	//StopLossPercentage          string        `json:"stop_loss_percentage"`
	//Cooldown                    string        `json:"cooldown"`
	//BtcPriceLimit               string        `json:"btc_price_limit"`
	//Strategy                    string        `json:"strategy"`
	//MinVolumeBtc24H             string        `json:"min_volume_btc_24h"`
	//ProfitCurrency              string        `json:"profit_currency"`
	//MinPrice                    interface{}   `json:"min_price"`
	//MaxPrice                    interface{}   `json:"max_price"`
	//StopLossType                string        `json:"stop_loss_type"`
	//SafetyOrderVolumeType       string        `json:"safety_order_volume_type"`
	//BaseOrderVolumeType         string        `json:"base_order_volume_type"`
	//AccountName                 string        `json:"account_name"`
	//TrailingDeviation           string        `json:"trailing_deviation"`
	//FinishedDealsProfitUsd      string        `json:"finished_deals_profit_usd"`
	//FinishedDealsCount          string        `json:"finished_deals_count"`
	//LeverageType                string        `json:"leverage_type"`
	//LeverageCustomValue         interface{}   `json:"leverage_custom_value"`
	//StartOrderType              string        `json:"start_order_type"`
	//ActiveDealsUsdProfit        string        `json:"active_deals_usd_profit"`
	//ActiveDeals                 []interface{} `json:"active_deals"`
}