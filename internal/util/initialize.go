package util

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/nfongster/ledger/internal/database"
)

func GetDbConnectionString() string {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port)
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

func AddBudget(q *database.Queries, category string, amount float64, period database.Period, startDate time.Time) {
	_, err := q.CreateBudget(context.Background(), database.CreateBudgetParams{
		TargetAmount: amount,
		TimePeriod:   period,
		StartDate:    startDate,
		Name:         category,
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
