package cmd

import (
	"fmt"
	"log"

	"ahtapot/scylla"

	"github.com/spf13/cobra"
)

var backupCmd = &cobra.Command{
	Use:   "backup",
	Short: "Runs backup operation",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Backup başlıyor...")
		fmt.Printf("Host: %s\n", cfg.Host)
		fmt.Printf("User: %s\n", cfg.User)
		fmt.Printf("Password: %s\n", cfg.Password)
		fmt.Printf("Keyspace: %s\n", cfg.Keyspace)
		fmt.Printf("Table: %s\n", cfg.Table)
		fmt.Printf("Port: %s\n", cfg.Port)
		fmt.Printf("Consistency: %s\n", cfg.Consistency)
		fmt.Printf("Format: %s\n", cfg.Format)

		scyllaConfig := scylla.ConnectionConfig{
			Hosts:       cfg.Host,
			Port:        cfg.Port,
			Username:    cfg.User,
			Password:    cfg.Password,
			Keyspace:    cfg.Keyspace,
			Consistency: cfg.Consistency,
		}

		session, err := scylla.CreateSession(scyllaConfig)
		if err != nil {
			log.Fatalf("HATA: %v", err)
		}
		defer session.Close()

		log.Println("Scylla bağlantısı başarılı!")
	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
