package database

import (
	"fmt"
	"time"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresDB struct {
	config *config.DatabaseConfig
	db     *gorm.DB
}

func NewPostgresDB(cfg *config.DatabaseConfig) *PostgresDB {
	return &PostgresDB{config: cfg}
}

func (p *PostgresDB) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(p.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(p.config.ConnectionPool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(p.config.ConnectionPool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(p.config.ConnectionPool.ConnMaxLifetime) * time.Second)

	p.db = db
	return db, nil
}

func (p *PostgresDB) GetDSN() string {
	cfg := p.config.Postgres
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.Host, cfg.User, cfg.Password, cfg.DBName, cfg.Port, cfg.SSLMode, cfg.TimeZone)
}

func (p *PostgresDB) Close() error {
	if p.db != nil {
		sqlDB, err := p.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
