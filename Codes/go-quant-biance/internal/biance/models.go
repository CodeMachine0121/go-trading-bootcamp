package biance

import "time"

type Kline struct {
	OpenTime time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}

type FetchKlineDetailDto struct {
	Symbol string
	Long   string
	Limit  int
}
