package commands

import (
	"fmt"
	"os"
)

func quit(args ...string) {
	fmt.Println("Exiting...")
	os.Exit(0)
}
