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
	if os.Getenv("DATABASE_URL") != "" {
		return os.Getenv("DATABASE_URL")
	}

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

func AddTransaction(q *database.Queries, description, category, notes string, date time.Time, amount float64) {
	category_id, err := q.GetOrCreateCategory(context.Background(), category)
	if err != nil {
		log.Fatal(err)
	}
	_, err = q.CreateTransaction(context.Background(), database.CreateTransactionParams{
		Date:        date,
		Description: description,
		Amount:      amount,
		CategoryID:  category_id,
		Notes: sql.NullString{
			String: notes,
			Valid:  notes != ""},
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

func SeedDatabase(q *database.Queries) {
	AddTransaction(q,
		"Lembas Bread",
		"Groceries",
		"",
		time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		8.00)

	AddTransaction(q,
		"Salted Pork",
		"Groceries",
		"",
		time.Date(2025, time.October, 16, 0, 0, 0, 0, time.UTC),
		15.00)

	AddTransaction(q,
		"The Green Dragon",
		"Eating Out",
		"beerz with the boyz",
		time.Date(2025, time.October, 17, 0, 0, 0, 0, time.UTC),
		40.00)

	AddBudget(q,
		"Groceries",
		50.00,
		database.PeriodWeekly,
		time.Date(2025, time.October, 13, 0, 0, 0, 0, time.UTC))
}
