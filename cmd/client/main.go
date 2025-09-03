package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/nfongster/ledger/internal/commands"
)

type Transaction struct {
	ID          int       `json:"id"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	Category    string    `json:"category"`
	Notes       string    `json:"notes"`
}

func main() {
	fmt.Println("Welcome to ledger!")

	reader := bufio.NewReader(os.Stdin)
	registry := commands.InitializeCommandRegistry()

	for {
		fmt.Print("> ")

		str, _ := reader.ReadString('\n')
		input := strings.Split(str, " ")

		cmd, args := input[0], input[1:]
		command, exists := registry.GetCommand(cmd)
		if !exists {
			fmt.Println("Command does not exist!")
		} else {
			command.Execute(args...)
		}
	}
}
