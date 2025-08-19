package calculators

import "fmt"

func CalculateSma(prices []float64, period int) (float64, error) {

	if period <= 0 {
		return 0, fmt.Errorf("週期必須大於 0")
	}

	if length := len(prices); length < period {
		return 0, fmt.Errorf("資料量不足，需要 %d 個，但只有 %d 個", period, length)
	}

	var sum float64 = 0

	for _, price := range prices {
		sum += price
	}

	return sum / float64(period), nil
}
