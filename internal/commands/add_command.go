package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Define a client-side transaction
type transaction struct {
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Notes       string    `json:"notes"`
}

func add(args ...string) {
	if len(args) < 5 {
		fmt.Println("Insufficient number of args was supplied!")
		return
	}

	dateStr, desc, amountStr, cat, notes := args[0], args[1], args[2], args[3], strings.TrimRight(args[4], "\n")
	date, err := parseDate(dateStr)
	if err != nil {
		fmt.Printf("error parsing date: %v\n", err)
		return
	}

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		fmt.Printf("error parsing amount: %v\n", err)
		return
	}

	url := "http://localhost:8080/api/transactions"
	t := transaction{
		Date:        date,
		Description: desc,
		Amount:      amount,
		Category:    cat,
		Notes:       notes,
	}

	jsonData, err := json.Marshal(t)
	if err != nil {
		fmt.Printf("error marshalling json: %v\n", err)
		return
	}

	reqBody := bytes.NewBuffer(jsonData)
	resp, err := http.Post(url, "application/json", reqBody)
	if err != nil {
		fmt.Printf("error in post request: %v\n", err)
		return
	}
	defer resp.Body.Close()
}

func parseDate(dateStr string) (time.Time, error) {
	const layout = "2006/01/02"
	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}
