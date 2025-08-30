package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	fmt.Println("Welcome to ledger!")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		str, _ := reader.ReadString('\n')
		input := strings.Split(str, " ")

		cmd, args := input[0], input[1:]
		command, exists := GetCommand(cmd)
		if !exists {
			fmt.Println("Command does not exist!")
		} else {
			command.Execute(args...)
		}
	}
}

type Command struct {
	Name        string
	Description string
	Callback    func(...string)
}

func (c Command) Execute(args ...string) {
	c.Callback(args...)
}

func GetCommand(command string) (*Command, bool) {
	cmd := strings.TrimSpace(strings.ToLower(command))
	switch cmd {
	case "quit":
		return &Command{
			Name:        "quit",
			Description: "Exit the program.",
			Callback: func(args ...string) {
				fmt.Println("Exiting...")
				os.Exit(0)
			},
		}, true

	default:
		return nil, false
	}
}
