package main

import (
	"fmt"
	"go-quant-biance/internal/biance"
	"go-quant-biance/internal/services"
)

func main() {
	fmt.Println("--- 使用 go-binance SDK 獲取數據 ---")

	mcadService := services.NewMcadService(biance.NewBianceClient())

	mcadService.Scan(services.MacdDto{
		Symbol:     "BTCUSDT",
		FastPeriod: 12,
		SlowPeriod: 26,
		Long:       "1h",
		Limit:      500,
	})

}
