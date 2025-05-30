package cmd

import (
	"ahtapot/scylla"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type Column struct {
	Name            string
	ClusteringOrder string
	Kind            string
	Position        int
	Type            string
}

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
		fmt.Println()

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
		/*
			colQuery := `SELECT column_name FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?`
			iter := session.Query(colQuery, cfg.Keyspace, cfg.Table).Iter()

			var columns []string
			var col string
			for iter.Scan(&col) {
				col = `"` + col + `"`
				columns = append(columns, col)
			}

			if err := iter.Close(); err != nil {
				fmt.Println("column fetch error: %w", err)
			}

			if len(columns) == 0 {
				fmt.Println("no columns found")
			}

			sort.Strings(columns) // SORT COLUMNS

			// Ana SELECT sorgusu
			query := fmt.Sprintf("SELECT %s FROM %s.%s limit 2", strings.Join(columns, ", "), cfg.Keyspace, cfg.Table)
			dataIter := session.Query(query).PageSize(500).Iter()

			fmt.Println(query)

			// f, _ := os.Create("dump.txt")
			// defer f.Close()

			// Satır satır oku
			row := make(map[string]any)
			for dataIter.MapScan(row) {

				fmt.Println(row)
				// fmt.Println("Yeni satır:")
				// for _, c := range columns {
				// 	fmt.Printf("  %s: %v\n", c, row[c])
				// }

				for k, v := range row {
					fmt.Printf("Alan: %s, Tipi: %T\n", k, v)
				}
				row = map[string]any{} // rowu boşalt
			}

			if err := dataIter.Close(); err != nil {
				fmt.Println("veri okuma hatası: %w", err)
			}
		*/

		query := "SELECT column_name, clustering_order, kind, position, type FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ? order by column_name desc"
		iter := session.Query(query, cfg.Keyspace, cfg.Table).Iter()

		var columns []Column
		var colName, clusteringOrder, kind, position, colType string

		for iter.Scan(&colName, &clusteringOrder, &kind, &position, &colType) {
			colName = `"` + colName + `"`
			pos, _ := strconv.Atoi(position)
			columns = append(columns, Column{Name: colName, ClusteringOrder: clusteringOrder, Kind: kind, Position: pos, Type: colType})
		}
		if err := iter.Close(); err != nil {
			fmt.Errorf("ERRRRRRRRRRRRRRRRR")
		}

		if len(columns) == 0 {
			fmt.Errorf("tablo bulunamadı veya sütun yok")
		}

		// fmt.Println(columns)

		// STATIC ANAHTAR KELIMESI EKLENECEK

		// 2. Sütunları türüne göre ayır
		var partitionColumns, clusteringColumns []Column
		var otherCols []string

		for _, col := range columns {
			line := fmt.Sprintf("    %s %s,", col.Name, col.Type)
			switch col.Kind {
			case "partition_key":
				partitionColumns = append(partitionColumns, col)
				otherCols = append(otherCols, line)
			case "clustering":
				clusteringColumns = append(clusteringColumns, col)
				otherCols = append(otherCols, line)
			default: // "regular", "static"
				otherCols = append(otherCols, line)
			}
		}

		sort.Slice(partitionColumns, byPosition(partitionColumns))

		sort.Slice(clusteringColumns, byPosition(clusteringColumns))

		var partitionKeys, clusteringKeys []string

		for _, col := range partitionColumns {
			partitionKeys = append(partitionKeys, col.Name)
		}

		var clusteringLine []string

		for _, col := range clusteringColumns {
			clusteringKeys = append(clusteringKeys, col.Name)
			clusteringLine = append(clusteringLine, col.Name+" "+col.ClusteringOrder)
		}

		// 3. PRIMARY KEY kısmını oluştur
		var pk, cl string
		if len(clusteringKeys) > 0 {
			pk = fmt.Sprintf("PRIMARY KEY ((%s), %s)", strings.Join(partitionKeys, ", "), strings.Join(clusteringKeys, ", "))
			cl = fmt.Sprintf("WITH CLUSTERING ORDER BY (%s)", strings.Join(clusteringLine, ", "))
		} else {
			pk = fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(partitionKeys, ", "))
		}

		// 4. Script'i birleştir
		lines := append([]string{fmt.Sprintf("CREATE TABLE %s.%s (", cfg.Keyspace, cfg.Table)}, otherCols...)
		lines = append(lines, fmt.Sprintf("    %s", pk), ") "+cl)

		var finalScript = strings.Join(lines, "\n")
		fmt.Println(finalScript)

	},
}

func byPosition(columns []Column) func(i, j int) bool {
	return func(i, j int) bool {
		return columns[i].Position < columns[j].Position
	}
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
