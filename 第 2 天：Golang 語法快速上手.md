# 第 2 天：Golang 語法快速上手

歡迎回到我們的量化交易學習旅程！
在第一天，我們設定了開發環境。今天，我們將學習 Golang 的基礎語法。程式語言的語法就像是我們與電腦溝通的文法規則，只有掌握了它，我們才能將腦中的交易策略，精確地轉譯成電腦可以理解並執行的指令。

> 今天的目標不是成為 Go 語言專家，而是快速掌握那些在量化交易中最常用、最核心的語法元素。

## 第一部分：交易策略 Domain Knowledge

在進入程式碼之前，我們先來複習一個簡單卻極其重要的數學概念，它將貫穿我們整個量化交易的學習過程。
交易中的基礎數學概念

- **百分比 (Percentage)**: 用於計算價格變動的幅度，即「漲跌幅」。例如，比特幣從 $68,000 漲到 $69,360，漲幅就是 `(69360 - 68000) / 68000 * 100% = 2%`。這是衡量資產波動性的基本單位。
- **報酬率 (Rate of Return)**: 你投資的收益或虧損佔初始成本的百分比。這是評估你策略績效的最核心指標。
- **移動平均 (Moving Average, MA)**: 這是技術分析中最基礎、最廣泛使用的指標之一。它指的是將過去特定時間週期（例如 5 天、10 天）內的價格加總後再除以該週期數，得出一個平均值。它的作用是平滑價格波動，幫助我們看清市場的趨勢方向。例如，「5日均線」就是最近 5 天收盤價的平均值。當前價格在 5 日均線之上，通常被視為短期強勢。

今天，我們的程式碼實作將圍繞著「計算移動平均」這個具體的目標展開。

## 第二部分：Coding 實作

打開昨天的 go-quant-binance 專案，我們將繼續在 main.go 檔案中進行今天的練習(可以先清空 main 函式內的程式碼)。

### 1. 變數 (Variables) 與資料型別 (Data Types)

變數就像是儲存數據的容器。在 Go 中，每個變數都必須有一個明確的「型別」，告訴電腦這個容器裡裝的是什麼樣的數據。
在量化交易中，我們最常用的型別有：

| 型別 | 說明 | 用途及範例 |
|------|------|-----------|
| string | 文字 | 用於儲存交易對名稱等。例如 "BTCUSDT" |
| int | 整數 | 用於計數或定義週期。例如 10 (代表 10 日均線) |
| float64 | 浮點數 (帶小數點的數字) | 這是最重要的型別，幾乎所有價格、數量、金額都用它來表示，以確保精度 |
| bool | 布林值 | 只有 true (真) 和 false (假) 兩種狀態，用於邏輯判斷。例如 isMarketBullish = true |

#### 範例：

```golang
package main

import "fmt"

func main() {
    // 使用 var 關鍵字宣告變數
    var tradingPair string = "BTC/USDT"
    var period int = 10

    // 使用 := 進行簡短宣告並初始化 (更常用)
    currentPrice := 69360.55
    isTradingAllowed := true

    fmt.Println("交易對:", tradingPair)
    fmt.Println("計算週期:", period)
    fmt.Println("目前價格:", currentPrice)
    fmt.Println("允許交易:", isTradingAllowed)
}
```

### 2. 集合型別：切片 (Slices)

當我們需要處理一系列的數據時，例如一連串的歷史收盤價，單一的變數就不夠用了。這時我們需要使用「切片」，切片可以看作是一個動態長度的陣列。

#### 範例：

我們用一個切片來儲存過去五天的收盤價。

```golang
// ... in main function
// 宣告一個 float64 型別的切片，並存入 5 個價格數據
closingPrices := []float64{68000.0, 68500.5, 68200.0, 69000.0, 69300.5}

fmt.Println("過去五天的收盤價:", closingPrices)
// 存取特定元素 (索引從 0 開始)
fmt.Println("第一天的價格:", closingPrices[0])
fmt.Println("第五天的價格:", closingPrices[4])```
```

### 3. 流程控制：迴圈 (Loops) 與條件判斷 (Conditionals)

| 控制結構 | 功能 | 用途 |
|---------|------|------|
| `for` 迴圈 | 讓我們可以重複執行某段程式碼 | 常用於遍歷切片中的每一個元素 |
| `if/else` 條件判斷 | 讓程式可以根據不同的條件，執行不同的邏輯 | 這是交易決策的核心 |

#### 範例：

讓我們用 `for` 迴圈來計算上面 `closingPrices` 的總和。

```go
// ... in main function
var sum float64 = 0.0 // 初始化一個變數來儲存總和

// for...range 迴圈會遍歷切片中的每一個元素
// _ (底線) 表示我們忽略索引，只需要值 (price)
for _, price := range closingPrices {
    sum = sum + price // 將每個價格累加到 sum 中
}

fmt.Println("價格總和:", sum)

// 使用 if/else 進行判斷
if currentPrice > 69000.0 {
    fmt.Println("決策: 目前價格高於 69000，市場可能看漲。")
} else {
    fmt.Println("決策: 目前價格不高於 69000，保持觀望。")
}
```

### 4. 函式 (Functions)

函式是組織程式碼、使其可重用的基本單位。我們可以把一個特定的功能（例如計算移動平均）封裝在一個函式裡，之後只要呼叫這個函式的名字就可以使用它，而不用重複撰寫相同的程式碼。

## 第三部分：綜合練習：計算簡單移動平均 (SMA)

現在，我們將綜合運用今天學到的所有知識，來撰寫一個非常有用的函式：**計算簡單移動平均 (Simple Moving Average, SMA)**。
請將以下完整程式碼 手打/複製 到你的 main.go 檔案中。

```golang
package main

import "fmt"

// calculateSMA 函式接收一個價格切片和一個週期，並回傳計算出的平均值
// prices: 一系列歷史價格，例如 []float64{10.0, 11.0, 12.0}
// period: 計算週期，例如 5 (代表 5 日均線)
// 回傳值: float64 型別的平均價
func calculateSMA(prices []float64, period int) (float64, error) {
    // --- 條件判斷：處理異常情況 ---
    if period <= 0 {
        return 0, fmt.Errorf("週期必須大於 0")
    }
    if len(prices) < period {
        // 如果數據點的數量小於要計算的週期，則無法計算
        return 0, fmt.Errorf("數據點不足，需要 %d 個，但只有 %d 個", period, len(prices))
    }

    // --- 迴圈：計算所需數據的總和 ---
    // 我們只取最近 `period` 個週期的數據來計算
    var sum float64 = 0.0
    // prices[len(prices)-period:] 是一個切片操作，表示從後面取 period 個元素
    for _, price := range prices[len(prices)-period:] {
        sum += price // `sum += price` 是 `sum = sum + price` 的簡寫
    }

    // --- 計算並回傳結果 ---
    return sum / float64(period), nil // `nil` 表示沒有錯誤發生
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
```

### 運行程式：

在終端機中執行 `go run main.go`，你將會看到以下輸出：

```text
最新的 5 日移動平均價 (5MA) 是: 68660.00
最新的 10 日移動平均價 (10MA) 是: 68040.00
預期的錯誤: 數據點不足，需要 20 個，但只有 10 個
```

> %.2f 是格式化輸出的語法，代表將浮點數格式化為保留兩位小數。

## 本日總結與預告

太棒了！今天我們已經掌握了 Golang 的核心語法，並成功地將它們應用於一個真實的交易指標計算中。
同時學會了如何使用變數儲存數據、如何用切片管理數據序列、如何用迴圈和條件判斷來執行邏輯，以及如何用函式來組織程式碼。

我們今天用的 `historicalPrices` 數據是手動寫死的。明天 (第 3 天)，我們將進入激動人心的一步：透過幣安 API，用程式去獲取真實、即時的市場 K 線數據！
準備好與真實的市場數據打交道了嗎？我們明天見！