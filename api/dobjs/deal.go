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
