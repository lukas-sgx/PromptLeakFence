package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Println("Prompt Leak Fence is a tool to detect and prevent prompt leaks in your code.")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("    plf [flags]")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("    -h, --help -> show this pannel")
}

func checkFlags() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "-h" || args[1] == "--help" {
			help()
		}
	} else {
		help()
	}
}

func main() {
	checkFlags()
}
