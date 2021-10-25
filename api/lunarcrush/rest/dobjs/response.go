package dobjs

//Response contains the LunarCrush API response
type Response struct {
	Config struct {
		Data  string `json:"data"`
		Type  string `json:"type"`
		Limit int    `json:"limit"`
		Page  int    `json:"page"`
	} `json:"config"`
	Usage struct {
		Day   int `json:"day"`
		Month int `json:"month"`
	} `json:"usage"`
	Data []PairData `json:"data"`
}

//PairData contains the information regarding a pair from lunarcrush.
//TODO consider removing commented out fields if never going to use.
type PairData struct {
	ID int    `json:"id"`
	S  string `json:"s"`
	N  string `json:"n"`
	//P          float64 `json:"p"`
	//PBtc       float64 `json:"p_btc"`
	//V          float64 `json:"v"`
	//Vt         float64 `json:"vt"`
	//Pc         float64 `json:"pc"`
	//Pch        float64 `json:"pch"`
	//Mc         float64   `json:"mc"`
	Gs float64 `json:"gs"` //GalaxyScore
	//Ss         int     `json:"ss"`
	//As         float64 `json:"as"`
	//Bl         int     `json:"bl"`
	//Br         int     `json:"br"`
	//Sp         int     `json:"sp"`
	//Na         int     `json:"na"`
	//Md         int     `json:"md"`
	//T          int     `json:"t"`
	//R          int     `json:"r"`
	//Yt         int     `json:"yt"`
	//Sv         int     `json:"sv"`
	//U          int     `json:"u"`
	//C          int     `json:"c"`
	//Sd         float64 `json:"sd"`
	//D          float64 `json:"d"`
	//Cr         float64 `json:"cr,omitempty"`
	Acr int `json:"acr"` //AltRank
	//Tc         int     `json:"tc,omitempty"`
	//Categories string  `json:"categories"`
	//Df         int     `json:"df,omitempty"`
	//Nft        int     `json:"nft,omitempty"`
	//Dot        int     `json:"dot,omitempty"`
	//E2         int     `json:"e2,omitempty"`
	//Meme       int     `json:"meme,omitempty"`
	//Zil        int     `json:"zil,omitempty"`
	//Bsc        int     `json:"bsc,omitempty"`
	//Sc         int     `json:"sc,omitempty"`
}
