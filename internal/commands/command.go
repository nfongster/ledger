package commands

import (
	"strings"
)

type Command struct {
	Description string
	callback    func(...string)
}

func (c Command) Execute(args ...string) {
	c.callback(args...)
}

var commandMap map[string]Command = map[string]Command{
	"quit": {
		Description: "Exit the program.",
		callback:    quit,
	},
	"transactions": {
		Description: "Get a list of all available transactions, or supply an ID as an argument to get a single transaction.",
		callback:    transactions,
	},
	"add": {
		Description: "Add a new entry to the ledger.  Args: date, desc, amount, cat, notes",
		callback:    add,
	},
}

func GetCommand(command string) (*Command, bool) {
	command = strings.TrimSpace(strings.ToLower(command))
	cmd, exists := commandMap[command]
	return &cmd, exists
}
