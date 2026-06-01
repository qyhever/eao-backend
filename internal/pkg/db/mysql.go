package db

import (
	"database/sql"
	"errors"
	"strings"
	"time"

	"eao/internal/config"

	_ "github.com/go-sql-driver/mysql"
)

var ErrMySQLDSNRequired = errors.New("mysql dsn is required")

func OpenMySQLFromConfig(cfg *config.Config) (*sql.DB, error) {
	if cfg == nil {
		return nil, nil
	}

	dsn := config.BuildMySQLDSN(cfg)
	if strings.TrimSpace(dsn) == "" {
		return nil, nil
	}

	return OpenMySQL(dsn)
}

func OpenMySQL(dsn string) (*sql.DB, error) {
	if strings.TrimSpace(dsn) == "" {
		return nil, ErrMySQLDSNRequired
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(10)

	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}
