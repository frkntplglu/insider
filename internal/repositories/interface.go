package repositories

import "gorm.io/gorm"

type Database interface {
	GetConnection() *gorm.DB
}
