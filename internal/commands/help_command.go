package commands

import "fmt"

func help(r *Registry, _ ...string) func(args ...string) {
	return func(args ...string) {
		for name, cmd := range r.commandMap {
			fmt.Printf("%s: %s\n", name, cmd.Description)
			for arg, desc := range cmd.Args {
				fmt.Printf("- %s: %s\n", arg, desc)
			}
		}
	}
}
