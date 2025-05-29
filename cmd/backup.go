package cmd

import (
	"ahtapot/scylla"
	"fmt"
	"strings"

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

		session, _ := scylla.CreateSession(scylla.CreateConfig(cfg))
		defer session.Close()

		// var query = session.Query("select * from system.clients limit 2")

		// // var query = session.Query("select * from alternator_zaimoglu2.zaimoglu2 LIMIT 2")

		// if rows, err := query.Iter().SliceMap(); err == nil {
		// 	for _, row := range rows {
		// 		fmt.Printf("%v\n", row)
		// 		fmt.Println(row["driver_name"])
		// 	}
		// } else {
		// 	panic("Query error: " + err.Error())
		// }

		colQuery := `SELECT column_name FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?`
		iter := session.Query(colQuery, cfg.Keyspace, cfg.Table).Iter()

		var columns []string
		var col string
		for iter.Scan(&col) {
			columns = append(columns, col)
		}
		if err := iter.Close(); err != nil {
			fmt.Errorf("column fetch error: %w", err)
		}
		if len(columns) == 0 {
			fmt.Errorf("no columns found")
		}

		// Ana SELECT sorgusu
		query := fmt.Sprintf("SELECT %s FROM %s.%s limit 2", strings.Join(columns, ", "), cfg.Keyspace, cfg.Table)
		dataIter := session.Query(query).PageSize(500).Iter()

		// Satır satır oku
		row := make(map[string]interface{})
		for dataIter.MapScan(row) {
			fmt.Println(row)
			fmt.Println("Yeni satır:")
			for _, c := range columns {
				fmt.Printf("  %s: %v\n", c, row[c])
			}
			row = map[string]interface{}{} // temizle
		}

		if err := dataIter.Close(); err != nil {
			fmt.Errorf("veri okuma hatası: %w", err)
		}

	},
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
