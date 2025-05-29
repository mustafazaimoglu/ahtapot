package scylla

import (
	"ahtapot/common"
	"fmt"
	"strings"
	"time"

	"github.com/gocql/gocql"
)

type ConnectionConfig struct {
	Hosts       []string
	Port        int
	Username    string
	Password    string
	Keyspace    string
	Consistency string // e.g. "ONE", "QUORUM"
}

// CreateSession returns an authenticated ScyllaDB session
func CreateSession(cfg ConnectionConfig) (*gocql.Session, error) {
	cluster := gocql.NewCluster(cfg.Hosts...)
	cluster.Port = cfg.Port
	cluster.Keyspace = cfg.Keyspace
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: cfg.Username,
		Password: cfg.Password,
	}
	cluster.Consistency = parseConsistency(cfg.Consistency)
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("connection error: %w", err)
	}
	return session, nil
}

func parsePort(portStr string) int {
	var port int
	_, err := fmt.Sscanf(portStr, "%d", &port)
	if err != nil || port <= 0 {
		fmt.Println("Port setted to 9042")
		return 9042
	}
	return port
}

func parseConsistency(level string) gocql.Consistency {
	level = strings.ToUpper(level)
	switch level {
	case "QUORUM":
		return gocql.Quorum
	case "LOCAL_QUORUM":
		return gocql.LocalQuorum
	case "ALL":
		return gocql.All
	case "ANY":
		return gocql.Any
	case "ONE":
		return gocql.One
	default:
		fmt.Println("Consisteny unknown! setted to ONE")
		return gocql.One
	}
}

func CreateConfig(cfg common.FlagConfig) ConnectionConfig {
	scyllaConfig := ConnectionConfig{
		Hosts:       strings.Split(cfg.Host, ","),
		Port:        parsePort(cfg.Port),
		Username:    cfg.User,
		Password:    cfg.Password,
		Keyspace:    cfg.Keyspace,
		Consistency: cfg.Consistency,
	}

	return scyllaConfig
}
