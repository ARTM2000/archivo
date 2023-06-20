package archive

import (
	"fmt"
	"log"

	"github.com/ARTM2000/archive1/internal/archive/auth"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DBConfig struct {
	DBHost    string
	DBPort    int
	DBUser    string
	DBPass    string
	DBName    string
	DBZone    string
	DBSSLMode bool
}

func NewDBConnection(dbc DBConfig) *gorm.DB {
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
		auth.UserSchema{},
	)

	return db
}
