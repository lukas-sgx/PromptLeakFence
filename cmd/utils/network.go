package utils

import (
	"fmt"
	"os/exec"
)

func SetRedirect(from string, to string, verbose bool) {
	StopRedirect(from, to, false)
	exec.Command("sudo", "iptables", "-t", "nat", "-I", "OUTPUT", "-p", "tcp", "--dport", from, "-m", "owner", "!", "--uid-owner", "0", "-j", "REDIRECT", "--to-ports", to).Run()

	if verbose {
		fmt.Printf("📡 Redirect traffic: %s -> %s\n", from, to)
	}
}

func StopRedirect(from string, to string, verbose bool) {
	exec.Command("sudo", "iptables", "-t", "nat", "-D", "OUTPUT", "-p", "tcp", "--dport", from, "-m", "owner", "!", "--uid-owner", "0", "-j", "REDIRECT", "--to-ports", to).Run()

	if verbose {
		fmt.Printf("⛔ STOP Redirect traffic: %s x %s\n", from, to)
	}
}
