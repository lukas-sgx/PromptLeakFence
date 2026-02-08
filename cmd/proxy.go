package cmd

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"strconv"

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

    exec.Command("sudo", "iptables", "-t", "nat", "-F").Run()
    exec.Command("sudo", "iptables", "-t", "nat", "-I", "OUTPUT", "-p", "tcp", "--dport", from, "-m", "owner", "!", "--uid-owner", "0", "-j", "REDIRECT", "--to-ports", to).Run()
	
	if verbose {
    	fmt.Printf("✅ Interception active sur ports %s -> %s\n", from, to)
	}
}

func proxyUp(port string, portLLM string) {
    proxy := goproxy.NewProxyHttpServer()
    proxy.Verbose = verbose

    target, _ := url.Parse("http://127.0.0.1:" + portLLM)
    reverseProxy := httputil.NewSingleHostReverseProxy(target)

    proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
        req.URL.Scheme = "http"
        req.URL.Host = target.Host
        req.Host = target.Host
        reverseProxy.ServeHTTP(w, req)
    })

    log.Fatal(http.ListenAndServe("0.0.0.0:"+port, proxy))
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
	
	if verbose {
		fmt.Printf("🚀 Proxy démarré: %s → %s (%s)\n", listenAddr, targetAddr, modelPort[targetAddr])
	}
	proxyUp(listenAddr, modelPort[targetAddr])
	
	if verbose {
		fmt.Printf("📁 Policy: %s | Verbose: %v\n", policyFile, verbose)
	}
}

func checkRoot() {
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

var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Short:   "Transparent LLM proxy avec leak protection",
	Long:    `Intercepte et scanne tous les prompts LLM en temps réel`,
	GroupID: "security",
	Run: func(cmd *cobra.Command, args []string) {
		checkRoot()
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
