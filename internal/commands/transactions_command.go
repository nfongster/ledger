package commands

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func transactions(args ...string) {
	fmt.Println("Getting all transactions...")

	url := "http://localhost:8080/api/transactions"
	if len(args) > 0 {
		id := strings.TrimRight(args[0], "\n")
		if _, err := strconv.Atoi(id); err != nil {
			fmt.Println("Argument was invalid!  Please supply an integer ID.")
			return
		}
		url += "/" + id
	}

	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error getting transaction(s): %v\n", err)
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

	fmt.Printf("Data successfully retrieved!\n%s\n", string(body))
}
