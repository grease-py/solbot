package main

import (
    "sync"
    "time"

    "github.com/shopspring/decimal"
)

type Portfolio struct {
    Balance        decimal.Decimal            // Total SOL balance
    Holdings       map[string]decimal.Decimal // Holdings in different shitcoins
    TransactionLog []Transaction
    mutex          sync.Mutex
}

type Transaction struct {
    Timestamp time.Time
    Action    string          // "buy" or "sell"
    Token     string
    Quantity  decimal.Decimal
    Price     decimal.Decimal
    Total     decimal.Decimal
}

func NewPortfolio(initialSOL decimal.Decimal) *Portfolio {
    return &Portfolio{
        Balance:  initialSOL,
        Holdings: make(map[string]decimal.Decimal),
    }
}

func (p *Portfolio) Buy(token string, quantity, price decimal.Decimal) bool {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    totalCost := quantity.Mul(price)
    if p.Balance.LessThan(totalCost) {
        return false // Not enough balance
    }

    p.Balance = p.Balance.Sub(totalCost)
    p.Holdings[token] = p.Holdings[token].Add(quantity)

    p.TransactionLog = append(p.TransactionLog, Transaction{
        Timestamp: time.Now(),
        Action:    "buy",
        Token:     token,
        Quantity:  quantity,
        Price:     price,
        Total:     totalCost,
    })

    return true
}

func (p *Portfolio) Sell(token string, quantity, price decimal.Decimal) bool {
    p.mutex.Lock()
    defer p.mutex.Unlock()

    holding, exists := p.Holdings[token]
    if !exists || holding.LessThan(quantity) {
        return false // Not enough holdings
    }

    totalRevenue := quantity.Mul(price)
    p.Balance = p.Balance.Add(totalRevenue)
    p.Holdings[token] = holding.Sub(quantity)

    p.TransactionLog = append(p.TransactionLog, Transaction{
        Timestamp: time.Now(),
        Action:    "sell",
        Token:     token,
        Quantity:  quantity,
        Price:     price,
        Total:     totalRevenue,
    })

    return true
}

func (p *Portfolio) GetBalance() decimal.Decimal {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    return p.Balance
}

func (p *Portfolio) GetHoldings() map[string]decimal.Decimal {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    copyHoldings := make(map[string]decimal.Decimal)
    for k, v := range p.Holdings {
        copyHoldings[k] = v
    }
    return copyHoldings
}

func (p *Portfolio) GetTransactionLog() []Transaction {
    p.mutex.Lock()
    defer p.mutex.Unlock()
    return p.TransactionLog
}
