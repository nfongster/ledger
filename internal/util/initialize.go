package util

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/nfongster/ledger/internal/database"
)

func LoadConfig(filePath string) (*Config, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	var config Config
	decoder.Decode(&config)
	return &config, err
}

func ResetDatabase(q *database.Queries) {
	// Truncate will permanently delete rows (i.e., the removal cannot be rolled back) and reset the primary key.
	if err := q.TruncateAllTables(context.Background()); err != nil {
		log.Fatal(err)
	}
}

func AddTransaction(q *database.Queries, description, category string, date time.Time, amount float64) {
	category_id, err := q.GetOrCreateCategory(context.Background(), category)
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        date,
		Description: description,
		Amount:      amount,
		CategoryID:  category_id,
	})
	if err != nil {
		log.Fatal(err)
	}
}

// TODO: Refactor
func SeedDatabase(q *database.Queries) {
	// --- Set Transactions ---
	groceries_id, err := q.GetOrCreateCategory(context.Background(), "Groceries")
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Lembas",
		Amount:      9.99,
		CategoryID:  groceries_id,
	})
	if err != nil {
		log.Fatal(err)
	}

	_, err = q.GetOrCreateCategory(context.Background(), "Groceries")
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		Description: "Old Toby",
		Amount:      49.99,
		CategoryID:  groceries_id,
	})
	if err != nil {
		log.Fatal(err)
	}

	eatingout_id, err := q.GetOrCreateCategory(context.Background(), "Eating Out")
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC),
		Description: "Shultzy's",
		Amount:      25.00,
		CategoryID:  eatingout_id,
		Notes: sql.NullString{
			String: "Bratz n beerz with the boyz",
			Valid:  true},
	})
	if err != nil {
		log.Fatal(err)
	}

	// --- Set Budgets ---
	_, err = q.CreateBudget(context.Background(), database.CreateBudgetParams{
		TargetAmount: 50.00,
		TimePeriod:   database.PeriodWeekly,
		StartDate:    time.Date(2025, time.October, 13, 0, 0, 0, 0, time.UTC),
		Name:         "Groceries",
	})
	if err != nil {
		log.Fatal(err)
	}
}
