package main

import (
    "context"
    "log"

    "github.com/jackc/pgx/v4/pgxpool"
)

type WalletSelectionModule struct {
    DB     *Database
    Config Config
}

func InitializeWalletSelection(db *Database, config Config) *WalletSelectionModule {
    return &WalletSelectionModule{
        DB:     db,
        Config: config,
    }
}

func (wsm *WalletSelectionModule) SelectTopWallets(limit int) ([]WalletMetrics, error) {
    query := `
        SELECT wallet_address, trade_count, win_rate, average_profit, 
               average_profit_pct, average_loss, average_loss_pct, 
               average_position_size, average_trade_duration
        FROM wallet_metrics
        WHERE trade_count > 50 AND win_rate > $1
        ORDER BY win_rate DESC, average_profit_pct DESC
        LIMIT $2;
    `

    rows, err := wsm.DB.Pool.Query(context.Background(), query, wsm.Config.TargetWinRate, limit)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var wallets []WalletMetrics
    for rows.Next() {
        var wm WalletMetrics
        err := rows.Scan(
            &wm.WalletAddress,
            &wm.TradeCount,
            &wm.WinRate,
            &wm.AverageProfit,
            &wm.AverageProfitPct,
            &wm.AverageLoss,
            &wm.AverageLossPct,
            &wm.AveragePositionSize,
            &wm.AverageTradeDuration,
        )
        if err != nil {
            log.Println("Error scanning row:", err)
            continue
        }
        wallets = append(wallets, wm)
    }

    return wallets, nil
}
