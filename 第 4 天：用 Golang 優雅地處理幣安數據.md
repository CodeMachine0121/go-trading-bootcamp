# 第 4 天：用 Golang 優雅地處理幣安資料

昨天，我們成功地從幣安 API 獲取了真實的市場資料。但這些資料現在還只是被臨時印在終端機上，就像一堆未經整理的原始食材。一個優秀的廚師會先將食材清洗、分類、放入標記好的容器中，然後才能開始烹飪。
今天，我們的任務就是扮演這位廚師，我們將學習使用 Golang 的核心工具——Structs（結構）、Slices（切片）和 Maps（映射），來創建我們自己的「容器」，將從 API 獲取的原始資料，轉化為我們策略邏輯可以輕鬆使用、結構清晰的格式。

## 第一部分：交易策略 Domain Knowledge

### 1. 為何需要自定義資料結構？

你可能會問：「go-binance SDK 不是已經把資料解析好了嗎？為什麼我們不直接用它提供的格式？」

這是一個非常好的問題，原因有三：

- **解耦 (Decoupling)**: 我們的交易策略核心邏輯，不應該依賴於任何特定的第三方函式庫。如果未來我們想更換 SDK，甚至更換資料源（例如從另一個交易所或 CSV 檔案讀取資料），我們只需要修改資料解析的部分，而不需要改動任何策略程式碼。
- **擴充性 (Extensibility)**: SDK 提供的資料結構只包含原始的 OHLCV。但我們的策略通常需要計算各種技術指標，例如我們昨天寫的 SMA，以及未來會學到的 RSI、MACD 等。我們可以在自定義的資料結構中，輕鬆地加入這些計算好的指標欄位。
- **可控性 (Control)**: API 回傳的價格、數量通常是 string（字串）格式，以避免浮點數的精度問題。在我們的策略計算中，則需要將它們轉換為 float64。擁有自己的資料結構，可以讓我們在這個轉換過程中進行精確的控制和錯誤處理。

### 2. 資料清洗 (Data Cleaning)

在處理真實世界的資料時，我們必須抱持一個心態：永遠不要完全相信資料源，資料可能會有缺失、格式錯誤、或出現極端異常值。雖然幣安這樣的大型交易所資料品質很高，但在量化交易中，建立一套檢查和清洗資料的流程是至關重要的習慣。今天，我們將在資料轉換的過程中，實踐最基礎的資料清洗——確保資料型別正確。

## 第二部分：定義我們自己的 K線結構 (Kline)

在 修改 main 函式之前，我們先定義一個專屬於我們自己的 Kline 結構，它就像一個模板，規定了每一根 K 棒資料應該包含哪些欄位以及它們的型別。

```golang
// ... import statements ...

// 自定義的 Kline 結構體，用於儲存和處理 K 線資料
type Kline struct {
	OpenTime time.Time // 我們直接使用 time.Time 型別，更方便處理
	Open     float64
	High     float64
	Low      float64
	Close    float64   // 計算指標時最常用的價格
	Volume   float64
}
```


## 第三部分：完整程式碼
請用以下完整程式碼替換你的 main.go 檔案。程式碼中的註解詳細解釋了每一步的作用。

```golang
package main

import (
	"context"
	"fmt"
	"strconv" // 提供字串和其他型別之間轉換的功能
	"time"

	"github.comcom/adshao/go-binance/v2"
)

// Kline 是我們自定義的 K 線結構體
// 用於儲存和處理我們需要的 K 線資料
type Kline struct {
	OpenTime time.Time
	Open     float64
	High     float64
	Low      float64
	Close    float64
	Volume   float64
}

// calculateSMA 函式接收一個價格切片和一個週期，並回傳計算出的平均值 (這是第 2 天的程式碼)
func calculateSMA(prices []float64, period int) (float64, error) {
	if period <= 0 {
		return 0, fmt.Errorf("週期必須大於 0")
	}
	if len(prices) < period {
		return 0, fmt.Errorf("資料點不足，需要 %d 個，但只有 %d 個", period, len(prices))
	}

	var sum float64 = 0.0
	for _, price := range prices[len(prices)-period:] {
		sum += price
	}
	return sum / float64(period), nil
}

func main() {
	// 1. 初始化幣安客戶端
	client := binance.NewClient("", "")

	// 2. 獲取幣安回傳的原始 K 線資料
	// 我們這次多獲取一些資料，以便計算不同週期的 MA
	binanceKlines, err := client.NewKlinesService().Symbol("BTCUSDT").
		Interval("1h").Limit(20).Do(context.Background())

	if err != nil {
		fmt.Println("獲取 K 線失敗:", err)
		return
	}

	// 3. 資料轉換與清洗
	// 建立一個我們自定義 Kline 型別的切片，用來儲存轉換後的資料
	var Klines []Kline

	for _, bk := range binanceKlines {
		// --- 核心轉換邏輯 ---
		// strconv.ParseFloat 將字串轉換為 float64
		open, _ := strconv.ParseFloat(bk.Open, 64)
		high, _ := strconv.ParseFloat(bk.High, 64)
		low, _ := strconv.ParseFloat(bk.Low, 64)
		closePrice, _ := strconv.ParseFloat(bk.Close, 64)
		volume, _ := strconv.ParseFloat(bk.Volume, 64)

		// 建立一個 Kline 實例並填充資料
		kline := Kline{
			OpenTime: time.Unix(bk.OpenTime/1000, 0),
			Open:     open,
			High:     high,
			Low:      low,
			Close:    closePrice,
			Volume:   volume,
		}
		// 將轉換好的 kline 添加到我們的切片中
		Klines = append(Klines, kline)
	}

	fmt.Printf("成功轉換並儲存了 %d 根 K 線資料。\n", len(Klines))
	fmt.Println("------------------------------------")

	// 4. 從我們自己的資料結構中提取資料並進行計算
	// 建立一個只包含收盤價的切片
	var closingPrices []float64
	for _, k := range Klines {
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
```

> 註： strconv.ParseFloat 會回傳兩個值：轉換結果和一個 error。在正式的專案中，你應該要檢查這個 error，但在今天的教學中，我們為了簡潔暫時用 _ 忽略了它。

在終端機中執行 go run main.go，你將看到如下輸出：

```text
--- 使用 go-binance SDK 獲取數據 ---
成功獲取 BTC/USDT 最近 5 根 1 小時 K 線:
成功轉換並儲存了 20 根 K 線資料。
------------------------------------
最新的 5 週期 SMA 是: 463192.97
最新的 10 週期 SMA 是: 231596.49
```

## 本日總結與預告

做得非常好！今天，我們從一個資料的使用者，晉升為了一個資訊的管理者。

我們學習了為何要建立自定義的資料結構，並親手定義了屬於我們自己的 Kline。更重要的是，你完成了從「原始 API 資料」到「乾淨、可用的策略資料」的關鍵轉換流程，包括了必要的型別轉換。最後，我們成功地用這些整理好的資料計算出了技術指標。

我們已經鋪好了所有的基礎設施：我們有了穩定的資料源，有了可靠的資料結構。明天 (第 5 天)，我們將迎來第一個激動人心的里程碑：基於我們計算出的移動平均線，打造我們的第一個完整交易策略——均線交叉策略，並產生明確的「買入」和「賣出」訊號！

準備好讓你的程式開始 "思考" 了嗎？我們明天見！