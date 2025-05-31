package cmd

import (
	"ahtapot/scylla"
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
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
		if false {
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
		}

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

			_, err := f.WriteString(mapToSingleQuoteString(tags))
								// if err != nil {
								// 	panic(err)
								// }

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
			fmt.Println("ERRRRRRRRRRRRRRRRR")
		}

		if len(columns) == 0 {
			fmt.Println("tablo bulunamadı veya sütun yok")
		}

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
			case "static":
				line := fmt.Sprintf("    %s %s %s,", col.Name, col.Type, strings.ToUpper(col.Kind))
				otherCols = append(otherCols, line)
			default: // "regular"
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

		// KEY kısmınlarını oluştur
		var pk, cl string
		if len(clusteringKeys) > 0 {
			pk = fmt.Sprintf("PRIMARY KEY ((%s), %s)", strings.Join(partitionKeys, ", "), strings.Join(clusteringKeys, ", "))
			cl += fmt.Sprintf("CLUSTERING ORDER BY (%s)", strings.Join(clusteringLine, ", "))
		} else {
			pk = fmt.Sprintf("PRIMARY KEY (%s)", strings.Join(partitionKeys, ", "))
		}

		// DONT ASK ME ABOUT HOW I WROTE THIS PART... ONLY GOD KNOWS HOW THIS ONE WORKS
		query = "SELECT * FROM system_schema.tables WHERE keyspace_name = ? AND table_name = ?"
		iter = session.Query(query, cfg.Keyspace, cfg.Table).Iter()
		var tableOpts []string

		row := make(map[string]any)
		for iter.MapScan(row) {
			for key, value := range row {
				if key == "flags" || key == "keyspace_name" || key == "table_name" || key == "id" {
					continue
				}

				if reflect.TypeOf(value).Kind() == reflect.Map {
					elementOfValue := reflect.TypeOf(value).Elem()
					if elementOfValue.Kind() == reflect.Slice || elementOfValue.Kind() == reflect.Array {
						if elementOfValue.Elem().Kind() == reflect.Uint8 {
							assertedMap, _ := value.(map[string][]byte)
							for k, v := range assertedMap {
								parsedString := parseExtensionsColumn(v)
								tableOpts = append(tableOpts, fmt.Sprintf("%s = %s", k, mapToSingleQuoteString(parsedString)))
							}
						}
					} else {
						assertedValue := value.(map[string]string)
						tableOpts = append(tableOpts, fmt.Sprintf("%s = %s", key, mapToSingleQuoteString(assertedValue)))
					}
				} else {
					if str, ok := value.(string); ok {
						value = fmt.Sprintf("'%s'", str)
					}
					tableOpts = append(tableOpts, fmt.Sprintf("%s = %v", key, value))
				}
			}
			row = map[string]any{} // rowu boşalt
		}

		sort.Strings(tableOpts)

		for i, val := range tableOpts {
			if i == 0 && len(cl) == 0 {
				tableOpts[i] = "     " + val
			} else {
				tableOpts[i] = "     AND " + val
			}
		}

		// 4. Script'i birleştir
		lines := append([]string{fmt.Sprintf("CREATE TABLE %s.%s (", cfg.Keyspace, cfg.Table)}, otherCols...)
		lines = append(lines, fmt.Sprintf("    %s", pk), ") WITH "+cl)
		lines = append(lines, tableOpts...)

		var finalScript = strings.Join(lines, "\n") + ";"
		fmt.Println(finalScript)

	},
}

func parseExtensionsColumn(data []byte) map[string]string {
	buf := bytes.NewReader(data)
	result := make(map[string]string)

	var numElements int32
	binary.Read(buf, binary.LittleEndian, &numElements)

	for i := 0; i < int(numElements); i++ {
		var keyLen int32
		binary.Read(buf, binary.LittleEndian, &keyLen)

		keyBytes := make([]byte, keyLen)
		buf.Read(keyBytes)
		key := string(keyBytes)

		var valLen int32
		binary.Read(buf, binary.LittleEndian, &valLen)

		valBytes := make([]byte, valLen)
		buf.Read(valBytes)
		val := string(valBytes)

		result[key] = val
	}

	return result
}

func mapToSingleQuoteString(m map[string]string) string {
	var parts []string
	for k, v := range m {
		parts = append(parts, fmt.Sprintf("'%s': '%s'", k, v))
	}
	return "{" + strings.Join(parts, ", ") + "}"
}

func byPosition(columns []Column) func(i, j int) bool {
	return func(i, j int) bool {
		return columns[i].Position < columns[j].Position
	}
}

func init() {
	rootCmd.AddCommand(backupCmd)
}
