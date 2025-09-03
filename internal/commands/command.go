package commands

import (
	"strings"
)

type Command struct {
	Description string
	Args        map[string]string
	callback    func(...string)
}

func (c Command) Execute(args ...string) {
	c.callback(args...)
}

type Registry struct {
	commandMap map[string]Command
}

func InitializeCommandRegistry() *Registry {
	r := Registry{}
	r.commandMap = map[string]Command{
		"quit": {
			Description: "Exit the program.",
			callback:    quit,
		},
		"help": {
			Description: "List all available commands and their arguments (if applicable).",
			callback:    help(&r),
		},
		"transactions": {
			Description: "Get a list of all available transactions.",
			Args: map[string]string{
				"id": "(Optional) Primary ID of a single transaction.  Supply this argument to get a single transaction.",
			},
			callback: transactions,
		},
		"add": {
			Description: "Add a new entry to the ledger.  All args are required.",
			Args: map[string]string{
				"date":        "Transaction date (YYYY/MM/DD).",
				"description": "Transaction description.",
				"amount":      "Transaction amount (USD).  Must be a floating point value.",
				"category":    "Transaction category.",
				"notes":       "Notes pertaining to the transaction.",
			},
			callback: add,
		},
	}
	return &r
}

func (r *Registry) GetCommand(command string) (*Command, bool) {
	command = strings.TrimSpace(strings.ToLower(command))
	cmd, exists := r.commandMap[command]
	return &cmd, exists
}
