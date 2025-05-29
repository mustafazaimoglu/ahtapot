package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Runs restore operation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Restore başlıyor...")
		fmt.Printf("Host: %s\n", cfg.Host)
		fmt.Printf("User: %s\n", cfg.User)
		fmt.Printf("Password: %s\n", cfg.Password)
		fmt.Printf("Keyspace: %s\n", cfg.Keyspace)
		fmt.Printf("Table: %s\n", cfg.Table)
		fmt.Printf("Port: %s\n", cfg.Port)
		fmt.Printf("Consistency: %s\n", cfg.Consistency)
		fmt.Printf("Format: %s\n", cfg.Format)

		// TODO: Restore işlemini burada yazabilirsin
	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
}
