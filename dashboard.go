package main

import (
    "encoding/json"
    "log"
    "net/http"
)

type DashboardMetrics struct {
    TotalSOL      string  `json:"total_sol"`
    TotalValueSOL string  `json:"total_value_sol"`
    ProfitLossSOL string  `json:"profit_loss_sol"`
    ProfitLossPct float64 `json:"profit_loss_pct"`
}

func (mm *MonitoringModule) ServeDashboard(w http.ResponseWriter, r *http.Request) {
    metrics := mm.CollectMetrics()

    dashboard := DashboardMetrics{
        TotalSOL:      metrics.TotalSOL.String(),
        TotalValueSOL: metrics.TotalValue.String(),
        ProfitLossSOL: metrics.ProfitLossSOL.String(),
        ProfitLossPct: metrics.ProfitLossPct.InexactFloat64(),
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(dashboard)
}

func InitializeDashboard(monitoring *MonitoringModule) {
    http.HandleFunc("/dashboard", monitoring.ServeDashboard)
    go func() {
        log.Fatal(http.ListenAndServe(":8080", nil))
    }()
}
