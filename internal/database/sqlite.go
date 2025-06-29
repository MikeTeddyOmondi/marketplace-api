package database

import (
	"fmt"
	"time"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SQLiteDB struct {
	config *config.DatabaseConfig
	db     *gorm.DB
}

func NewSQLiteDB(cfg *config.DatabaseConfig) *SQLiteDB {
	return &SQLiteDB{config: cfg}
}

func (s *SQLiteDB) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(s.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(s.config.ConnectionPool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(s.config.ConnectionPool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(s.config.ConnectionPool.ConnMaxLifetime) * time.Second)

	s.db = db
	return db, nil
}

func (s *SQLiteDB) GetDSN() string {
	return s.config.SQLite.Path
}

func (s *SQLiteDB) Close() error {
	if s.db != nil {
		sqlDB, err := s.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
