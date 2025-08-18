package main

import "fmt"

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

	// 假設這是我們獲取到的一系列歷史收盤價
	historicalPrices := []float64{
		67100.0, 67300.0, 67500.0, 67400.0, 67800.0, // 第 1-5 天
		68200.0, 68500.0, 68300.0, 69000.0, 69300.0, // 第 6-10 天
	}

	// --- 呼叫函式，計算 5 日均線 ---
	sma5, err := calculateSMA(historicalPrices, 5)
	if err != nil {
		// 如果計算過程中發生錯誤，就印出錯誤訊息
		fmt.Println("計算 5 日均線時出錯:", err)
	} else {
		// 否則，印出計算結果
		fmt.Printf("最新的 5 日移動平均價 (5MA) 是: %.2f\n", sma5)
	}

	// --- 呼叫函式，計算 10 日均線 ---
	sma10, err := calculateSMA(historicalPrices, 10)
	if err != nil {
		fmt.Println("計算 10 日均線時出錯:", err)
	} else {
		fmt.Printf("最新的 10 日移動平均價 (10MA) 是: %.2f\n", sma10)
	}

	// --- 呼叫一個會出錯的例子 ---
	_, err = calculateSMA(historicalPrices, 20)
	if err != nil {
		fmt.Println("預期的錯誤:", err)
	}

}
