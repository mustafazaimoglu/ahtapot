package cmd

import (
	"ahtapot/common"

	"github.com/spf13/cobra"
)

var cfg common.FlagConfig

var rootCmd = &cobra.Command{
	Use:   "ahtapot",
	Short: "Backup/Restore Tool for ScyllaDB",
}

func RootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.PersistentFlags().Bool("help", false, "Help Page for Ahtapot")
	rootCmd.PersistentFlags().StringVarP(&cfg.Host, "host", "h", "127.0.0.1", "Host/Hosts")
	rootCmd.PersistentFlags().StringVarP(&cfg.User, "user", "U", "cassandra", "User")
	rootCmd.PersistentFlags().StringVarP(&cfg.Password, "password", "P", "cassandra", "Password")
	rootCmd.PersistentFlags().StringVarP(&cfg.Keyspace, "keyspace", "k", "system", "Keyspace")
	rootCmd.PersistentFlags().StringVarP(&cfg.Table, "table", "t", "clients", "Table")
	rootCmd.PersistentFlags().StringVarP(&cfg.Port, "port", "p", "9042", "Port")
	rootCmd.PersistentFlags().StringVarP(&cfg.Consistency, "consistency", "c", "ONE", "Consistency")
	rootCmd.PersistentFlags().StringVarP(&cfg.Format, "format", "f", "csv", "Format (custom,csv)")
}
