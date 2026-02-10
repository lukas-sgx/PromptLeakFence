package utils

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func CheckRoot() {
	cmd := exec.Command("id", "-u")
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	uid, err := strconv.Atoi(string(output[:len(output)-1]))
	if err != nil {
		log.Fatal(err)
	}
	if uid != 0 {
		fmt.Println("This program must be run as root! (sudo)")
		os.Exit(1)
	}
}

func CheckTarget(model map[string]string, targetAddr string) {
	if _, exists := model[targetAddr]; !exists {
		fmt.Printf("Target '%s' not found\n", targetAddr)
		fmt.Println("Available targets:")
		for name := range model {
			fmt.Printf("    %s\n", name)
		}
		os.Exit(1)
	}
}
