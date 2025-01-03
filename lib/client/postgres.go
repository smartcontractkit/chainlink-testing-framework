package client

import (
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	// import for side effect of sql packages
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/lib/k8s/environment"
)

// PostgresConnector sqlx postgres connector
type PostgresConnector struct {
	*sqlx.DB
	Cfg *PostgresConfig
}

// PostgresConfig postgres config for connection string
type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresConnector creates new sqlx postgres connector
func NewPostgresConnector(cfg *PostgresConfig) (*PostgresConnector, error) {
	var sslmode string
	if cfg.SSLMode == "" {
		sslmode = "sslmode=disable"
	} else {
		sslmode = fmt.Sprintf("sslmode=%s", cfg.SSLMode)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s %s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, sslmode)
	log.Debug().Str("ConnectionString", psqlInfo).Msg("Connecting")
	db, err := sqlx.Connect("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		// nolint
		db.Close()
		return nil, err
	}
	log.Debug().Str("ConnectionString", psqlInfo).Msg("Connected")
	return &PostgresConnector{DB: db, Cfg: cfg}, nil
}

// ConnectDB establishes a connection to a PostgreSQL database using the provided environment settings.
// It returns a PostgresConnector instance or an error if the connection fails.
func ConnectDB(nodeNum int, e *environment.Environment) (*PostgresConnector, error) {
	spl := strings.Split(e.URLs["chainlink_db"][nodeNum], ":")
	port := spl[len(spl)-1]
	db, err := NewPostgresConnector(&PostgresConfig{
		Host:     "127.0.0.1",
		Port:     port,
		User:     "postgres",
		Password: "postgres",
		DBName:   "chainlink",
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}
