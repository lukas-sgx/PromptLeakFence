package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	// "os"
	"os/exec"

	"github.com/elazarl/goproxy"
	"github.com/spf13/cobra"
)

var (
	listenAddr string
	targetAddr string
	policyFile string
	verbose    bool
)

func setRedirect(from string, to string) {
	exec.Command("sudo", "iptables", "-t", "nat", "-A", "PREROUTING", "-p", "tcp", "--dport", from, "-j", "REDIRECT", "--to-ports", to).Run()
}

func proxyUp(port string, portLLM string) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	target, _ := url.Parse("http://127.0.0.1:" + portLLM)
	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.Host = target.Host
		reverseProxy.ServeHTTP(w, req)
	})

	log.Fatal(http.ListenAndServe(":"+port, proxy))
}

func proxyListener() {
	modelPort := map[string]string{
		"ollama":    "11434",
		"llama.cpp": "8080",
		"lmstudio":  "1234",
		"oobabooga": "7860",
		"openwebui": "3000",
	}

	setRedirect(modelPort[targetAddr], listenAddr)
	fmt.Printf("🚀 Proxy démarré: %s → %s (%s)\n", listenAddr, targetAddr, modelPort[targetAddr])
	proxyUp(listenAddr, modelPort[targetAddr])
	fmt.Printf("📁 Policy: %s | Verbose: %v\n", policyFile, verbose)
}

var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Short:   "Transparent LLM proxy avec leak protection",
	Long:    `Intercepte et scanne tous les prompts LLM en temps réel`,
	GroupID: "security",
	Run: func(cmd *cobra.Command, args []string) {
		proxyListener()
	},
}

func NewProxyCmd() *cobra.Command {
	proxyCmd.Flags().StringVarP(&listenAddr, "listen", "l", "8080", "Listen address")
	proxyCmd.Flags().StringVarP(&targetAddr, "target", "t", "", "Target LLM service")
	proxyCmd.Flags().StringVarP(&policyFile, "policy", "p", "configs/policy.yaml", "Policy file")
	proxyCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logs")
	proxyCmd.MarkFlagRequired("target")
	return proxyCmd
}
