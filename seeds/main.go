package main

import (
	"flag"
	"fmt"

	"github.com/guths/zpe/config"
	"github.com/guths/zpe/models"
	"github.com/guths/zpe/seeds/seeders"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func init() {
	config.InitializeAppConfig()
}

func main() {
	tableFlag := flag.String("table", "all_table", "specify the table")

	flag.Parse()

	table := *tableFlag

	var db *gorm.DB
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.AppConfig.DBUsername,
		config.AppConfig.DBPassword,
		config.AppConfig.DBHost,
		config.AppConfig.DBPort,
		config.AppConfig.DBDatabase,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	db.AutoMigrate(&models.Role{}, &models.User{}, &models.UserRole{})

	if err != nil {
		panic("error connecting db to run seeders")
	}

	switch {
	case table == "user_role":
		seeders.UserRoleSeeder(db)
	}
}
