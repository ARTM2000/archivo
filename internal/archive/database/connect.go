package database

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Config struct {
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBZone    string
	DBSSLMode bool
}

type Manager struct {
	db *gorm.DB
}

func NewManager(dbc Config) Manager {
	sslMode := "disable"
	if dbc.DBSSLMode {
		sslMode = "enable"
	}

	postgresDSN := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		dbc.DBHost,
		dbc.DBUser,
		dbc.DBPass,
		dbc.DBName,
		dbc.DBPort,
		sslMode,
		dbc.DBZone,
	)

	db, err := gorm.Open(postgres.Open(postgresDSN), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Info),
		TranslateError: true,
	})
	if err != nil {
		log.Fatalln("fail to connect database.", err.Error())
	}

	// auto migration.
	// todo: make its safety more
	db.AutoMigrate(
		UserSchema{},
	)

	return Manager{
		db,
	}
}
