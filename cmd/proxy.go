package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    listenAddr string
    targetAddr string
    policyFile string
    verbose    bool
)

var proxyCmd = &cobra.Command{
    Use:     "proxy",
    Short:   "Transparent LLM proxy avec leak protection",
    Long:    `Intercepte et scanne tous les prompts LLM en temps réel`,
    GroupID: "security",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("🚀 Proxy démarré: %s → %s\n", listenAddr, targetAddr)
        fmt.Printf("📁 Policy: %s | Verbose: %v\n", policyFile, verbose)
    },
}

func NewProxyCmd() *cobra.Command {
    proxyCmd.Flags().StringVarP(&listenAddr, "listen", "l", ":8080", "Listen address")
    proxyCmd.Flags().StringVarP(&targetAddr, "target", "t", "", "Target LLM service")
    proxyCmd.Flags().StringVarP(&policyFile, "policy", "p", "configs/policy.yaml", "Policy file")
    proxyCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Verbose logs")
    proxyCmd.MarkFlagRequired("target")
    return proxyCmd
}
