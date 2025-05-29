package util

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

type BackupConfig struct {
	Session   *gocql.Session
	Keyspace  string
	Table     string
	OutputDir string
}

// Meta bilgisi için basit yapı
type MetaInfo struct {
	Rows      int       `json:"rows"`
	Timestamp time.Time `json:"timestamp"`
	Duration  string    `json:"duration"`
}

func DumpTable(cfg BackupConfig) error {

	// Output klasörü oluştur
	// if err := os.MkdirAll(cfg.OutputDir, 0755); err != nil {
	// 	return fmt.Errorf("output klasörü oluşturulamadı: %w", err)
	// }

	// Şemayı al
	schemaQuery := `SELECT column_name, kind, type FROM system_schema.columns WHERE keyspace_name = ? AND table_name = ?`
	iter := cfg.Session.Query(schemaQuery, cfg.Keyspace, cfg.Table).Iter()

	var cols []string
	colTypes := make(map[string]string)

	var name, kind, typ string
	for iter.Scan(&name, &kind, &typ) {
		if kind == "partition_key" || kind == "clustering" || kind == "regular" {
			cols = append(cols, name)
			colTypes[name] = typ
		}
	}
	if err := iter.Close(); err != nil {
		return err
	}

	// schema.cql oluştur
	schemaFile, err := os.Create(cfg.OutputDir + "/schema.cql")
	if err != nil {
		return err
	}
	defer schemaFile.Close()

	schemaParts := []string{}
	for _, col := range cols {
		schemaParts = append(schemaParts, fmt.Sprintf("%s %s", col, colTypes[col]))
	}
	schemaLine := fmt.Sprintf("CREATE TABLE %s.%s (\n  %s\n);", cfg.Keyspace, cfg.Table, strings.Join(schemaParts, ",\n  "))
	schemaFile.WriteString(schemaLine)

	// data.dump oluştur
	dataFile, err := os.Create(cfg.OutputDir + "/data.dump")
	if err != nil {
		return err
	}
	defer dataFile.Close()

	log.Printf("Backup tamamlandı.\n")
	return nil
}
