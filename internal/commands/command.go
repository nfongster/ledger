package commands

import (
	"strings"
)

type Command struct {
	Name        string
	Description string
	callback    func(...string)
}

func (c Command) Execute(args ...string) {
	c.callback(args...)
}

func GetCommand(command string) (*Command, bool) {
	cmd := strings.TrimSpace(strings.ToLower(command))
	switch cmd {
	case "quit":
		return &Command{
			Name:        "quit",
			Description: "Exit the program.",
			callback:    quit,
		}, true

	case "transactions":
		return &Command{
			Name:        "transactions",
			Description: "Get a list of all available transactions.",
			callback:    transactions,
		}, true

	default:
		return nil, false
	}
}
