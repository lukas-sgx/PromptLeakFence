package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "plf",
	Short: "PromptLeakFence CLI",
	Long:  "The first bi-directional LLM prompt firewall",
}

var logLevel string

func Execute() error {
	rootCmd.AddGroup(&cobra.Group{
		ID:    "security",
		Title: "Security Commands:",
	})
	rootCmd.AddCommand(NewProxyCmd())
	rootCmd.AddCommand(NewAuditCmd())

	rootCmd.PersistentFlags().StringVarP(&logLevel, "log-level", "L", "info", "Log level")

	if err := rootCmd.Execute(); err != nil {
		return err
	}
	return nil
}

func init() {
	rootCmd.SetHelpFunc(func(c *cobra.Command, a []string) {
		fmt.Println(`
        ╔══════════════════════════════════════╗
        ║    PromptLeakFence v0.1.0            ║
        ║    First LLM Prompt Firewall         ║
        ╚══════════════════════════════════════╝
        `)
		c.Usage()
	})
}
