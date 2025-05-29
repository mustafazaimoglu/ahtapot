package config

import "github.com/spf13/pflag"

// Yapı: Flag değerlerini taşıyacak
type Config struct {
	Host        string
	User        string
	Password    string
	Keyspace    string
	Table       string
	Consistency string
	Port        string
	Format      string
}

// ParseFlags: CLI bayraklarını tanımlar ve döner
func ParseFlags() *Config {
	cfg := &Config{}

	pflag.StringVarP(&cfg.Host, "host", "h", "127.0.0.1", "Host/Hosts")
	pflag.StringVarP(&cfg.User, "user", "U", "cassandra", "User")
	pflag.StringVarP(&cfg.Password, "password", "P", "cassandra", "Password")
	pflag.StringVarP(&cfg.Keyspace, "keyspace", "k", "", "Keyspace")
	pflag.StringVarP(&cfg.Table, "table", "t", "", "Table")
	pflag.StringVarP(&cfg.Port, "port", "p", "9042", "Port")
	pflag.StringVarP(&cfg.Consistency, "consistency", "c", "ONE", "Consistency")
	pflag.StringVarP(&cfg.Format, "format", "f", "csv", "Format (csv,json)")

	pflag.Parse()
	return cfg
}
