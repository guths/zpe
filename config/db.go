package config

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type dbManager struct {
	db   *gorm.DB
	once sync.Once
}

var DBManager = &dbManager{}

func (manager *dbManager) GetDB() *gorm.DB {
	manager.once.Do(func() {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			AppConfig.DBUsername,
			AppConfig.DBPassword,
			AppConfig.DBHost,
			AppConfig.DBPort,
			AppConfig.DBDatabase,
		)

		db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to connect to MySQL database: %v", err)
		}

		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get database instance: %v", err)
		}
		if err := sqlDB.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}

		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(100)
		sqlDB.SetConnMaxLifetime(0)

		manager.db = db
	})

	return manager.db
}
