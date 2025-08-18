package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/adshao/go-binance/v2"
)

type Kline struct {
	OpenTime time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}

func calculateSMA(prices []float64, period int) (float64, error) {

	if period <= 0 {
		return 0, fmt.Errorf("週期必須大於 0")
	}

	if length := len(prices); length < period {
		return 0, fmt.Errorf("樹據點不足，需要 %d 個，但只有 %d 個", period, length)
	}

	var sum float64 = 0

	for _, price := range prices {
		sum += price
	}

	return sum / float64(period), nil
}

func main() {
	fmt.Println("--- 使用 go-binance SDK 獲取數據 ---")

	// 初始化一個幣安客戶端。因為我們只訪問公開數據，所以 API Key 和 Secret Key 可以留空。
	client := binance.NewClient("", "")

	// 使用 KlinesService 來請求 K 線數據
	// .Symbol("BTCUSDT")  - 指定交易對
	// .Interval("1h")     - 指定 K 線的時間間隔 (1h 代表 1 小時)
	// .Limit(5)          - 指定獲取最近 5 根 K 線
	// .Do(context.Background()) - 執行請求

	bianceKline, err := client.NewKlinesService().Symbol("BTCUSDT").Interval("1h").Limit(20).Do(context.Background())

	// 錯誤處理
	if err != nil {
		fmt.Println("透過 SDK 獲取 K 線失敗:", err)
		return
	}

	var klines []Kline

	// klines 是一個包含了 K 線數據的切片，讓我們遍歷它並印出
	fmt.Println("成功獲取 BTC/USDT 最近 5 根 1 小時 K 線:")
	for _, k := range bianceKline {

		open, _ := strconv.ParseFloat(k.Open, 64)
		close, _ := strconv.ParseFloat(k.Close, 64)
		low, _ := strconv.ParseFloat(k.Low, 64)
		high, _ := strconv.ParseFloat(k.High, 64)
		volume, _ := strconv.ParseFloat(k.Volume, 64)

		klines = append(klines, Kline{
			OpenTime: time.Unix(k.OpenTime/1000, 0),
			Open:     open,
			Close:    close,
			Low:      low,
			High:     high,
			Volume:   volume,
		})
	}

	fmt.Printf("成功轉換並儲存了 %d 根 K 線資料。\n", len(klines))
	fmt.Println("------------------------------------")

	// 從我們自己的資料結構中提取資料並進行計算
	var closingPrices []float64
	for _, k := range klines {
		closingPrices = append(closingPrices, k.Close)
	}

	// 呼叫第 2 天的函式來計算 SMA
	sma5, err := calculateSMA(closingPrices, 5)
	if err != nil {
		fmt.Println("計算 5MA 失敗:", err)
	} else {
		fmt.Printf("最新的 5 週期 SMA 是: %.2f\n", sma5)
	}

	sma10, err := calculateSMA(closingPrices, 10)
	if err != nil {
		fmt.Println("計算 10MA 失敗:", err)
	} else {
		fmt.Printf("最新的 10 週期 SMA 是: %.2f\n", sma10)
	}
}
