package api

// Deal is the object returned by the 3commas Deals websocket
type Deal struct {
	Identifier string `json:"identifier"`
	Details    struct {
		ID    int    `json:"id"`
		Type  string `json:"type"`
		BotID int    `json:"bot_id"`
		// MaxSafetyOrders                  int         `json:"max_safety_orders"`
		// DealHasError                     bool        `json:"deal_has_error"`
		// FromCurrencyID                   int         `json:"from_currency_id"`
		// ToCurrencyID                     int         `json:"to_currency_id"`
		// AccountID                        int         `json:"account_id"`
		// ActiveSafetyOrdersCount          int         `json:"active_safety_orders_count"`
		// CreatedAt                        time.Time   `json:"created_at"`
		// UpdatedAt                        time.Time   `json:"updated_at"`
		// ClosedAt                         interface{} `json:"closed_at"`
		// Finished                         bool        `json:"finished?"`
		// CurrentActiveSafetyOrdersCount   int         `json:"current_active_safety_orders_count"`
		// CurrentActiveSafetyOrders        int         `json:"current_active_safety_orders"`
		CompletedSafetyOrdersCount       int `json:"completed_safety_orders_count"`
		CompletedManualSafetyOrdersCount int `json:"completed_manual_safety_orders_count"`
		// Cancellable                      bool        `json:"cancellable?"`
		// PanicSellable                    bool        `json:"panic_sellable?"`
		// TrailingEnabled                  bool        `json:"trailing_enabled"`
		// TslEnabled                       bool        `json:"tsl_enabled"`
		// StopLossTimeoutEnabled           bool        `json:"stop_loss_timeout_enabled"`
		// StopLossTimeoutInSeconds         int         `json:"stop_loss_timeout_in_seconds"`
		// ActiveManualSafetyOrders         int         `json:"active_manual_safety_orders"`
		Pair   string `json:"pair"`
		Status string `json:"status"`
		//LocalizedStatus                  string      `json:"localized_status"`
		//TakeProfit                       string      `json:"take_profit"`
		//BaseOrderVolume                  string      `json:"base_order_volume"`
		//SafetyOrderVolume                string      `json:"safety_order_volume"`
		//SafetyOrderStepPercentage        string      `json:"safety_order_step_percentage"`
		//LeverageType                     string      `json:"leverage_type"`
		//LeverageCustomValue              interface{} `json:"leverage_custom_value"`
		//BoughtAmount                     string      `json:"bought_amount"`
		//BoughtVolume                     string      `json:"bought_volume"`
		//BoughtAveragePrice               string      `json:"bought_average_price"`
		//BaseOrderAveragePrice            string      `json:"base_order_average_price"`
		//SoldAmount                       string      `json:"sold_amount"`
		//SoldVolume                       string      `json:"sold_volume"`
		//SoldAveragePrice                 string      `json:"sold_average_price"`
		//TakeProfitType                   string      `json:"take_profit_type"`
		//FinalProfit                      string      `json:"final_profit"`
		//MartingaleCoefficient            string      `json:"martingale_coefficient"`
		//MartingaleVolumeCoefficient      string      `json:"martingale_volume_coefficient"`
		//MartingaleStepCoefficient        string      `json:"martingale_step_coefficient"`
		//StopLossPercentage               string      `json:"stop_loss_percentage"`
		//ErrorMessage                     interface{} `json:"error_message"`
		//ProfitCurrency                   string      `json:"profit_currency"`
		//StopLossType                     string      `json:"stop_loss_type"`
		//SafetyOrderVolumeType            string      `json:"safety_order_volume_type"`
		//BaseOrderVolumeType              string      `json:"base_order_volume_type"`
		//FromCurrency                     string      `json:"from_currency"`
		//ToCurrency                       string      `json:"to_currency"`
		//CurrentPrice                     string      `json:"current_price"`
		//TakeProfitPrice                  string      `json:"take_profit_price"`
		//StopLossPrice                    interface{} `json:"stop_loss_price"`
		//FinalProfitPercentage            string      `json:"final_profit_percentage"`
		//ActualProfitPercentage           string      `json:"actual_profit_percentage"`
		//BotName                          string      `json:"bot_name"`
		//AccountName                      string      `json:"account_name"`
		//UsdFinalProfit                   string      `json:"usd_final_profit"`
		//ActualProfit                     string      `json:"actual_profit"`
		//ActualUsdProfit                  string      `json:"actual_usd_profit"`
		//FailedMessage                    interface{} `json:"failed_message"`
		//ReservedBaseCoin                 string      `json:"reserved_base_coin"`
		//ReservedSecondCoin               string      `json:"reserved_second_coin"`
		//TrailingDeviation                string      `json:"trailing_deviation"`
		//TrailingMaxPrice                 interface{} `json:"trailing_max_price"`
		//TslMaxPrice                      interface{} `json:"tsl_max_price"`
		//Strategy                         string      `json:"strategy"`
		//ReservedQuoteFunds               string      `json:"reserved_quote_funds"`
		//ReservedBaseFunds                string      `json:"reserved_base_funds"`
		//BotEvents                        []struct {
		//	Message   string    `json:"message"`
		//	CreatedAt time.Time `json:"created_at"`
		//} `json:"bot_events"`
	} `json:"message"`
}

type Ping struct {
	Type    string `json:"type"`
	Message int64  `json:"message"`
}
