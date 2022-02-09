package client

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	// import for side effect of sql packages
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
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
}

// NewPostgresConnector creates new sqlx postgres connector
func NewPostgresConnector(cfg *PostgresConfig) (*PostgresConnector, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
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
