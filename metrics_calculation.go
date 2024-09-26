package main

import (
    "time"

    "github.com/shopspring/decimal"
)

type Trade struct {
    OpenTime        time.Time
    CloseTime       time.Time
    Profit          float64
    ProfitPct       float64
    PositionSize    float64
    Action          string  // "buy" or "sell"
    Token           string
    Quantity        float64
    Price           float64
}

type DailyPnL struct {
    Date time.Time `json:"date,omitempty"`
    PnL  float64   `json:"pnl,omitempty"`
}

type WalletMetrics struct {
    WalletAddress        string        `json:"walletAddress"`
    TradeCount           int           `json:"tradeCount,omitempty"`
    WinRate              float64       `json:"winRate,omitempty"`
    AverageProfit        float64       `json:"averageProfit,omitempty"`
    AverageProfitPct     float64       `json:"averageProfitPct,omitempty"`
    AverageLoss          float64       `json:"averageLoss,omitempty"`
    AverageLossPct       float64       `json:"averageLossPct,omitempty"`
    AveragePositionSize  float64       `json:"averagePositionSize,omitempty"`
    AverageTradeDuration time.Duration `json:"averageTradeDuration,omitempty"`
    DailyPnLTrend        []DailyPnL    `json:"dailyPnLTrend,omitempty"`
}

func CalculateWalletMetrics(walletAddress string, trades []Trade) WalletMetrics {
    var wm WalletMetrics
    wm.WalletAddress = walletAddress
    wm.TradeCount = len(trades)

    if wm.TradeCount == 0 {
        return wm
    }

    var winningTrades, losingTrades int
    var totalProfit, totalLoss, totalPositionSize float64
    var totalDuration time.Duration
    dailyPnLMap := make(map[string]float64)

    for _, trade := range trades {
        totalPositionSize += trade.PositionSize
        totalDuration += trade.CloseTime.Sub(trade.OpenTime)
        
        if trade.Profit > 0 {
            winningTrades++
            totalProfit += trade.Profit
        } else {
            losingTrades++
            totalLoss += trade.Profit
        }

        date := trade.CloseTime.Format("2006-01-02")
        dailyPnLMap[date] += trade.Profit
    }

    wm.WinRate = (float64(winningTrades) / float64(wm.TradeCount)) * 100

    if winningTrades > 0 {
        wm.AverageProfit = totalProfit / float64(winningTrades)
        wm.AverageProfitPct = (wm.AverageProfit / (totalPositionSize / float64(wm.TradeCount))) * 100
    }

    if losingTrades > 0 {
        wm.AverageLoss = totalLoss / float64(losingTrades)
        wm.AverageLossPct = (wm.AverageLoss / (totalPositionSize / float64(wm.TradeCount))) * 100
    }

    wm.AveragePositionSize = totalPositionSize / float64(wm.TradeCount)
    wm.AverageTradeDuration = totalDuration / time.Duration(wm.TradeCount)

    for dateStr, pnl := range dailyPnLMap {
        date, err := time.Parse("2006-01-02", dateStr)
        if err != nil {
            continue
        }
        wm.DailyPnLTrend = append(wm.DailyPnLTrend, DailyPnL{
            Date: date,
            PnL:  pnl,
        })
    }

    return wm
}
