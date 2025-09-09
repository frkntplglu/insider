package database

import (
	"fmt"

	"github.com/frkntplglu/insider/pkg/logger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

type Database struct {
	Db *gorm.DB
}

func NewDatabase(cfg DatabaseConfig) *Database {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Istanbul", cfg.Host, cfg.User, cfg.Password, cfg.Database, cfg.Port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		logger.Fatal("Error while connecting db : %v", err)
	}

	return &Database{Db: db}

}

func (d *Database) GetConnection() *gorm.DB {
	return d.Db
}

func (d *Database) Close() {
	sqlDB, err := d.Db.DB()
	if err != nil {
		logger.Fatal("Error while getting db instance : %v", err)
	}
	sqlDB.Close()
}
