package services

import (
	"fmt"
	"go-quant-biance/internal/biance"
	"go-quant-biance/internal/calculators"

	"github.com/samber/lo"
)

type IMcadService interface {
	Scan(dto MacdDto)
}

type McadService struct {
	client biance.IBianceClient
}

func (ms *McadService) Scan(dto MacdDto) {

	fmt.Printf("開始檢測週期為 %d 和 %d 的均線交叉訊號...\n\n", dto.FastPeriod, dto.SlowPeriod)

	klines := ms.client.GetKlineDetail(biance.FetchKlineDetailDto{
		Symbol: dto.Symbol,
		Long:   dto.Long,
		Limit:  dto.Limit,
	})

	closingPrices := lo.Map(klines, func(kline biance.Kline, i int) float64 {
		return kline.Close
	})

	for i := dto.SlowPeriod; i < len(closingPrices); i++ {
		// 多一個時間點作為當前時間

		// 獲取當前時間點可用的所有歷史價格
		// 計算當前的快線和慢線
		pricesSoFar := closingPrices[:i+1]
		fastMA, _ := calculators.CalculateSma(pricesSoFar, dto.FastPeriod)
		slowMA, _ := calculators.CalculateSma(pricesSoFar, dto.SlowPeriod)

		// 計算上一根 K 線的快線和慢線
		// 獲取一時間點之前的歷史價格
		pricesPrev := closingPrices[:i]
		prevFastMA, _ := calculators.CalculateSma(pricesPrev, dto.FastPeriod)
		prevSlowMA, _ := calculators.CalculateSma(pricesPrev, dto.SlowPeriod)

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

func NewMcadService(client biance.IBianceClient) IMcadService {
	return &McadService{
		client: client,
	}
}
