package breeze_models

import "encoding/json"

type CandleData struct {
	Close        float64 `json:"close"`
	Datetime     string  `json:"datetime"`
	ExchangeCode string  `json:"exchange_code"`
	ExpiryDate   string  `json:"expiry_date"`
	High         float64 `json:"high"`
	Low          float64 `json:"low"`
	Open         float64 `json:"open"`
	OpenInterest int     `json:"open_interest"`
	ProductType  string  `json:"product_type"`
	Right        string  `json:"right"`
	StockCode    string  `json:"stock_code"`
	StrikePrice  float64 `json:"strike_price"`
	Volume       int     `json:"volume"`
}

func GetIciciCandleDataFromJson(jsonData []byte) *CandleData {
	var candleData CandleData
	err := json.Unmarshal(jsonData, &candleData)
	if err != nil {
		panic(err.Error())
	}
	return &candleData
}

type HistoricalDataResponse struct {
	Error   string       `json:"Error"`
	Status  int          `json:"Status"`
	Success []CandleData `json:"Success"`
}
