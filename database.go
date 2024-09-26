package main

import (
    "context"
    "fmt"
    "log"

    "github.com/jackc/pgx/v4/pgxpool"
)

type Database struct {
    Pool *pgxpool.Pool
}

func InitializeDatabase(config Config) *Database {
    dbURL := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
        config.DBUser,
        config.DBPassword,
        config.DBHost,
        config.DBPort,
        config.DBName,
    )

    pool, err := pgxpool.Connect(context.Background(), dbURL)
    if err != nil {
        log.Fatalf("Unable to connect to database: %v\n", err)
    }

    db := &Database{Pool: pool}
    db.setupSchema()
    return db
}

func (db *Database) setupSchema() {
    // Create tables if they don't exist
    walletMetricsTable := `
    CREATE TABLE IF NOT EXISTS wallet_metrics (
        wallet_address VARCHAR PRIMARY KEY,
        trade_count INTEGER,
        win_rate FLOAT,
        average_profit FLOAT,
        average_profit_pct FLOAT,
        average_loss FLOAT,
        average_loss_pct FLOAT,
        average_position_size FLOAT,
        average_trade_duration INTERVAL
    );`

    dailyPnLTable := `
    CREATE TABLE IF NOT EXISTS daily_pnl_trend (
        id SERIAL PRIMARY KEY,
        wallet_address VARCHAR REFERENCES wallet_metrics(wallet_address),
        date DATE,
        pnl FLOAT
    );`

    _, err := db.Pool.Exec(context.Background(), walletMetricsTable)
    if err != nil {
        log.Fatalf("Failed to create wallet_metrics table: %v", err)
    }

    _, err = db.Pool.Exec(context.Background(), dailyPnLTable)
    if err != nil {
        log.Fatalf("Failed to create daily_pnl_trend table: %v", err)
    }

    // Create indexes
    indexes := []string{
        `CREATE INDEX IF NOT EXISTS idx_win_rate ON wallet_metrics(win_rate);`,
        `CREATE INDEX IF NOT EXISTS idx_average_profit ON wallet_metrics(average_profit);`,
        `CREATE INDEX IF NOT EXISTS idx_trade_count ON wallet_metrics(trade_count);`,
    }

    for _, idx := range indexes {
        _, err := db.Pool.Exec(context.Background(), idx)
        if err != nil {
            log.Fatalf("Failed to create index: %v", err)
        }
    }
}

func UpsertWalletMetrics(db *Database, wm WalletMetrics) error {
    query := `
        INSERT INTO wallet_metrics (
            wallet_address, trade_count, win_rate, average_profit, 
            average_profit_pct, average_loss, average_loss_pct, 
            average_position_size, average_trade_duration
        ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
        ON CONFLICT (wallet_address) 
        DO UPDATE SET
            trade_count = EXCLUDED.trade_count,
            win_rate = EXCLUDED.win_rate,
            average_profit = EXCLUDED.average_profit,
            average_profit_pct = EXCLUDED.average_profit_pct,
            average_loss = EXCLUDED.average_loss,
            average_loss_pct = EXCLUDED.average_loss_pct,
            average_position_size = EXCLUDED.average_position_size,
            average_trade_duration = EXCLUDED.average_trade_duration
    `

    _, err := db.Pool.Exec(context.Background(), query,
        wm.WalletAddress,
        wm.TradeCount,
        wm.WinRate,
        wm.AverageProfit,
        wm.AverageProfitPct,
        wm.AverageLoss,
        wm.AverageLossPct,
        wm.AveragePositionSize,
        wm.AverageTradeDuration,
    )
    return err
}

func InsertDailyPnL(db *Database, walletAddress string, dailyPnLs []DailyPnL) error {
    query := `
        INSERT INTO daily_pnl_trend (wallet_address, date, pnl)
        VALUES ($1, $2, $3)
        ON CONFLICT (wallet_address, date)
        DO UPDATE SET pnl = EXCLUDED.pnl
    `

    for _, pnl := range dailyPnLs {
        _, err := db.Pool.Exec(context.Background(), query, walletAddress, pnl.Date, pnl.PnL)
        if err != nil {
            return err
        }
    }
    return nil
}
