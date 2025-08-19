package biance

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type IBianceClient interface {
	GetKlineDetail(dto FetchKlineDetailDto) []Kline
}

type BianceClient struct {
	client *binance.Client
}

func (bc *BianceClient) GetKlineDetail(dto FetchKlineDetailDto) []Kline {
	bianceKline, err :=
		bc.client.NewKlinesService().Symbol(dto.Symbol).Interval(dto.Long).Limit(dto.Limit).Do(context.Background())

	klines := make([]Kline, 0)

	if err != nil {
		fmt.Println("透過 SDK 獲取 K 線失敗:", err)
		return klines
	}

	for _, k := range bianceKline {

		openPrice, _ := strconv.ParseFloat(k.Open, 64)
		closingPrice, _ := strconv.ParseFloat(k.Close, 64)
		low, _ := strconv.ParseFloat(k.Low, 64)
		high, _ := strconv.ParseFloat(k.High, 64)
		volume, _ := strconv.ParseFloat(k.Volume, 64)

		klines = append(klines, Kline{
			OpenTime: time.Unix(k.OpenTime/1000, 0),
			Open:     openPrice,
			Close:    closingPrice,
			Low:      low,
			High:     high,
			Volume:   volume,
		})
	}

	// klines 是一個包含了 K 線數據的切片，讓我們遍歷它並印出
	fmt.Println("成功獲取 BTC/USDT 最近 5 根 1 小時 K 線:")

	return klines
}

func NewBianceClient() IBianceClient {
	// 初始化一個幣安客戶端。因為我們只訪問公開數據，所以 API Key 和 Secret Key 可以留空。
	return &BianceClient{
		binance.NewClient("", ""),
	}
}
