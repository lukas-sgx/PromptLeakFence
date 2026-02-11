package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"syscall"

	"github.com/elazarl/goproxy"
	"github.com/lukas-sgx/PromptLeakFence/cmd/utils"
	"github.com/spf13/cobra"
)

var (
	listenAddr string
	targetAddr string
	policyFile string
	verbose    bool
)

func excludeWord(pattern *utils.FilePolicy, content string) string {
	replace := "[INTERNAL_PROMPT_REDACTED]"

	for _, p := range pattern.Policy.Exclude {
		content = strings.ReplaceAll(strings.ToLower(content), p, replace)
	}

	newContent := strings.Split(content, " ")
	for i, word := range newContent {
		if strings.Contains(string(word), strings.ToLower(replace)) {
			newContent[i] = replace
		}
	}
	return strings.Join(newContent, " ")
}

func excludeInjection(pattern *utils.FilePolicy, content string) string {
	replace := "[PROMPT_INJECTION_DETECTED]"

	for _, p := range pattern.Policy.InjectionPatterns {
		if !strings.HasPrefix(p, "(?i)") {
			p = "(?i)" + p
		}

		re, err := regexp.Compile(p)
		if err != nil {
			if verbose {
				fmt.Printf("Invalid regex pattern: %s - %v\n", p, err)
			}
			continue
		}
		content = re.ReplaceAllString(content, replace)
	}
	return content
}

func contentControl(content string) string {
	pattern := utils.ReadPolicy(policyFile)

	content = excludeWord(pattern, content)
	content = excludeInjection(pattern, content)
	return content
}

func rewriteRequest(req **http.Request) {
	var data map[string]interface{}

	if (*req).Body == nil {
		return
	}

	bodyBytes, err := io.ReadAll((*req).Body)
	if err != nil {
		return
	}

	if err := json.Unmarshal(bodyBytes, &data); err != nil {
		if verbose {
			fmt.Printf("Failed to parse JSON: %v\n", err)
		}
		(*req).Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return
	}

	if msgs, ok := data["messages"].([]interface{}); ok {
		newMessages := []interface{}{}
		for _, m := range msgs {
			msg, ok := m.(map[string]interface{})
			if !ok {
				newMessages = append(newMessages, m)
				continue
			}

			role, _ := msg["role"].(string)
			content, _ := msg["content"].(string)

			if role == "assistant" && (content == "[object Object]" || content == "") {
				continue
			}

			msg["content"] = contentControl(content)

			newMessages = append(newMessages, msg)
		}
		data["messages"] = newMessages
	}

	bodyBytes, _ = json.Marshal(data)
	(*req).Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if verbose {
		fmt.Println((*req).Body)
	}

	(*req).ContentLength = int64(len(bodyBytes))
	(*req).Header.Set("Content-Length", strconv.Itoa(len(bodyBytes)))
}

func proxyUp(port string, portLLM string) {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = verbose

	target, _ := url.Parse("http://127.0.0.1:" + portLLM)
	reverseProxy := httputil.NewSingleHostReverseProxy(target)

	proxy.NonproxyHandler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		rewriteRequest(&req)
		req.URL.Scheme = "http"
		req.URL.Host = target.Host
		req.Host = target.Host
		reverseProxy.ServeHTTP(w, req)
	})

	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, proxy))
}

func proxyListener(model map[string]string) {
	utils.SetRedirect(model[targetAddr], listenAddr, verbose)

	if verbose {
		fmt.Printf("📁 Policy: %s | Verbose: %v\n", policyFile, verbose)
		fmt.Printf("🚀 Proxy start: %s → %s (%s)\n", listenAddr, targetAddr, model[targetAddr])
	}
	proxyUp(listenAddr, model[targetAddr])
}

var proxyCmd = &cobra.Command{
	Use:     "proxy",
	Short:   "Transparent LLM proxy avec leak protection",
	Long:    `Intercepte et scanne tous les prompts LLM en temps réel`,
	GroupID: "security",
	Run: func(cmd *cobra.Command, args []string) {
		model := map[string]string{
			"ollama":    "11434",
			"llama.cpp": "8080",
			"lmstudio":  "1234",
			"oobabooga": "7860",
			"openwebui": "3000",
			"copilot":   "5000",
			"gemini":    "8080",
			"claude":    "8080",
		}

		utils.CheckRoot()
		utils.CheckTarget(model, targetAddr)

		sigc := make(chan os.Signal, 1)
		signal.Notify(sigc,
			syscall.SIGHUP,
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGQUIT)
		go func() {
			s := <-sigc

			switch s {
			default:
				utils.StopRedirect(model[targetAddr], listenAddr, verbose)
				os.Exit(0)
			}
		}()

		proxyListener(model)
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
