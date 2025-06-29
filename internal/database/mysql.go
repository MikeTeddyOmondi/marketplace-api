package database

import (
	"fmt"
	"time"

	"github.com/MikeTeddyOmondi/marketplace-api/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDB struct {
	config *config.DatabaseConfig
	db     *gorm.DB
}

func NewMySQLDB(cfg *config.DatabaseConfig) *MySQLDB {
	return &MySQLDB{config: cfg}
}

func (m *MySQLDB) Connect() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(m.GetDSN()), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to MySQL: %w", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(m.config.ConnectionPool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(m.config.ConnectionPool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(m.config.ConnectionPool.ConnMaxLifetime) * time.Second)

	m.db = db
	return db, nil
}

func (m *MySQLDB) GetDSN() string {
	cfg := m.config.MySQL
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
		cfg.Charset,
		cfg.ParseTime,
		cfg.Loc,
	)
}

func (m *MySQLDB) Close() error {
	if m.db != nil {
		sqlDB, err := m.db.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
