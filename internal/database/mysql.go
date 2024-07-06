package database

import (
	"fmt"

	"github.com/ditatompel/xmr-remote-nodes/internal/config"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// DB holds the database
type DB struct{ *sqlx.DB }

// database instance
var defaultDB = &DB{}

// connect sets the db client of database using configuration
func (db *DB) connect(cfg *config.DB) (err error) {
	if defaultDB.DB != nil {
		return nil // reuse existing connection if available
	}

	dbURI := fmt.Sprintf("%s:%s@(%s:%d)/%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
	)

	db.DB, err = sqlx.Connect("mysql", dbURI)
	if err != nil {
		return err
	}

	// Try to ping database.
	if err := db.Ping(); err != nil {
		defer db.Close() // close database connection
		return fmt.Errorf("can't sent ping to database, %w", err)
	}

	return nil
}

// GetDB returns db instance
func GetDB() *DB {
	return defaultDB
}

// ConnectDB sets the db client of database using default configuration
func ConnectDB() error {
	return defaultDB.connect(config.DBCfg())
}
