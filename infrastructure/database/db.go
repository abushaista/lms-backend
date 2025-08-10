package database

import (
	"fmt"
	"log"

	"github.com/abushaista/lms-backend/infrastructure/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewGormDB(cfg *config.Config) (*gorm.DB, error) {
	user := cfg.DBUser
	pass := cfg.DBPass
	host := cfg.DBHost
	port := cfg.DBPort
	name := cfg.DBName

	if user == "" {
		log.Fatal("DB env not set")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	return gorm.Open(mysql.Open(dsn), &gorm.Config{})

}
