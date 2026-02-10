package utils

import (
	"fmt"
	"os/exec"
)

func SetRedirect(from string, to string, verbose bool) {

	exec.Command("sudo", "iptables", "-t", "nat", "-F").Run()
	exec.Command("sudo", "iptables", "-t", "nat", "-I", "OUTPUT", "-p", "tcp", "--dport", from, "-m", "owner", "!", "--uid-owner", "0", "-j", "REDIRECT", "--to-ports", to).Run()

	if verbose {
		fmt.Printf("📡 Redirect traffic: %s -> %s\n", from, to)
	}
}
