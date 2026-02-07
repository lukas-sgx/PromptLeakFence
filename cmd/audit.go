package cmd

import (
    "fmt"
    "github.com/spf13/cobra"
)

var (
    dbPath string
    port   uint
)

var auditCmd = &cobra.Command{
    Use:     "audit",
    Short:   "Dashboard d'audit des prompts bloqués",
    Long:    `Analyse et visualise les tentatives de fuites détectées`,
    GroupID: "security",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Printf("📊 Audit dashboard sur :%d\n", port)
        fmt.Printf("📁 DB: %s\n", dbPath)
    },
}

func NewAuditCmd() *cobra.Command {
    auditCmd.Flags().StringVarP(&dbPath, "db", "d", "plf.db", "SQLite audit DB")
    auditCmd.Flags().UintVarP(&port, "port", "p", 9090, "Dashboard port")
    return auditCmd
}
