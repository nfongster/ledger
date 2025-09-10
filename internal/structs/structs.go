package data

import (
	"encoding/json"
	"os"
	"time"

	"github.com/nfongster/ledger/internal/database"
)

type State struct {
	Database *database.Queries
}

type Config struct {
	DbConnectionString string
}

func LoadConfig() (*Config, error) {
	file, err := os.Open("config.json")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	decoder.Decode(&config)
	return &config, err
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
	CategoryId   int       `json:"category_id"`
}
