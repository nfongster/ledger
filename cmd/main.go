package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nfongster/ledger/internal/commands"
)

func main() {
	fmt.Println("Welcome to ledger!")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		str, _ := reader.ReadString('\n')
		input := strings.Split(str, " ")

		cmd, args := input[0], input[1:]
		command, exists := commands.GetCommand(cmd)
		if !exists {
			fmt.Println("Command does not exist!")
		} else {
			command.Execute(args...)
		}
	}
}
