package commands

import (
	"fmt"
	"io"
	"net/http"
)

func transactions(args ...string) {
	fmt.Println("Getting all transactions...")

	resp, err := http.Get("http://localhost:8080/api/transactions")
	if err != nil {
		fmt.Printf("Error getting transactions: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	fmt.Printf("Data retrieved!\n%s\n", string(body))
}
