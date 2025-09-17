package util

import (
	"time"

	"github.com/nfongster/ledger/internal/database"
)

type State struct {
	Database *database.Queries
}

type Config struct {
	DbConnectionString string
}

type TransactionClientParams struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Notes       string    `json:"notes"`
}

type BudgetClientParams struct {
	TargetAmount float64   `json:"target_amount"`
	TimePeriod   string    `json:"time_period"`
	StartDate    time.Time `json:"start_date"`
	Notes        string    `json:"notes"`
	Category     string    `json:"category"`
}
