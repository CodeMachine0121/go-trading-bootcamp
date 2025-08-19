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

	bianceKline, err := client.NewKlinesService().Symbol("BTCUSDT").Interval("4h").Limit(1000).Do(context.Background())

	// 錯誤處理
	if err != nil {
		fmt.Println("透過 SDK 獲取 K 線失敗:", err)
		return
	}

	var klines []Kline
	var closingPrices []float64

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

		closingPrices = append(closingPrices, close)
	}

	// --- 2. 策略邏輯與訊號產生 ---
	fastPeriod := 10 // 快線週期
	slowPeriod := 20 // 慢線週期

	fmt.Printf("開始檢測週期為 %d 和 %d 的均線交叉訊號...\n\n", fastPeriod, slowPeriod)

	for i := slowPeriod; i < len(closingPrices); i++ {
		// 獲取當前時間點可用的所有歷史價格
		pricesSoFar := closingPrices[:i+1]

		// 計算當前的快線和慢線
		fastMA, _ := calculateSMA(pricesSoFar, fastPeriod)
		slowMA, _ := calculateSMA(pricesSoFar, slowPeriod)

		// 獲取上一時間點的歷史價格
		pricesPrev := closingPrices[:i]

		// 計算上一根 K 線的快線和慢線
		prevFastMA, _ := calculateSMA(pricesPrev, fastPeriod)
		prevSlowMA, _ := calculateSMA(pricesPrev, slowPeriod)

		// --- 核心判斷邏輯 ---
		// 判斷黃金交叉: 上一刻快線在慢線下方，且當前快線在慢線上方
		if prevFastMA < prevSlowMA && fastMA > slowMA {
			fmt.Printf(
				"[買入訊號] 黃金交叉! 時間: %s, 收盤價: %.2f, 快線: %.2f, 慢線: %.2f\n",
				klines[i].OpenTime.Format("2006-01-02 15:04"),
				klines[i].Close,
				fastMA,
				slowMA,
			)
		}

		// 判斷死亡交叉: 上一刻快線在慢線上方，且當前快線在慢線下方
		if prevFastMA > prevSlowMA && fastMA < slowMA {
			fmt.Printf(
				"[賣出訊號] 死亡交叉! 時間: %s, 收盤價: %.2f, 快線: %.2f, 慢線: %.2f\n",
				klines[i].OpenTime.Format("2006-01-02 15:04"),
				klines[i].Close,
				fastMA,
				slowMA,
			)
		}
	}
}
