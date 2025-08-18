package main

import (
	"context"
	"fmt"
	"time"

	"github.com/adshao/go-binance/v2"
)

func main() {
	fmt.Println("--- 使用 go-binance SDK 獲取數據 ---")

	// 初始化一個幣安客戶端。因為我們只訪問公開數據，所以 API Key 和 Secret Key 可以留空。
	client := binance.NewClient("", "")

	// 使用 KlinesService 來請求 K 線數據
	// .Symbol("BTCUSDT")  - 指定交易對
	// .Interval("1h")     - 指定 K 線的時間間隔 (1h 代表 1 小時)
	// .Limit(5)          - 指定獲取最近 5 根 K 線
	// .Do(context.Background()) - 執行請求

	klines, err := client.NewKlinesService().Symbol("BTCUSDT").Interval("1h").Limit(5).Do(context.Background())

	// 錯誤處理
	if err != nil {
		fmt.Println("透過 SDK 獲取 K 線失敗:", err)
		return
	}

	// klines 是一個包含了 K 線數據的切片，讓我們遍歷它並印出
	fmt.Println("成功獲取 BTC/USDT 最近 5 根 1 小時 K 線:")
	for _, k := range klines {
		fmt.Printf(
			"開盤時間: %s, 開盤價: %s, 最高價: %s, 最低價: %s, 收盤價: %s\n",
			// k.OpenTime 是毫秒時間戳，我們需要將它轉換成人類可讀的格式
			time.Unix(k.OpenTime/1000, 0).Format("2006-01-02 15:04:05"),
			k.Open,
			k.High,
			k.Low,
			k.Close,
		)
	}
}
