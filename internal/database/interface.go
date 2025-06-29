package database

import (
	"fmt"
	"log"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"

	"gorm.io/gorm"
)

type Database interface {
	Connect() (*gorm.DB, error)
	GetDSN() string
	Close() error
}

type Manager struct {
	Db     *gorm.DB
	Config *config.DatabaseConfig
}

func NewManager(cfg *config.DatabaseConfig) (Database, error) {
	log.Printf("Database config received: driver=%s", cfg.Driver)
	var db Database

	switch cfg.Driver {
	case "sqlite":
		db = NewSQLiteDB(cfg)
	case "postgres":
		db = NewPostgresDB(cfg)
	case "mysql":
		db = NewMySQLDB(cfg)
		log.Fatalf("MySQL support is not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	return db, nil
}
