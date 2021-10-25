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

//Bot describes a bot from the 3commas api
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
	Name                        string `json:"name"`
	TakeProfit                  string `json:"take_profit"`
	BaseOrderVolume             string `json:"base_order_volume"`
	SafetyOrderVolume           string `json:"safety_order_volume"`
	SafetyOrderStepPercentage   string `json:"safety_order_step_percentage"`
	TakeProfitType              string `json:"take_profit_type"`
	MartingaleVolumeCoefficient string `json:"martingale_volume_coefficient"`
	MartingaleStepCoefficient   string `json:"martingale_step_coefficient"`
}
