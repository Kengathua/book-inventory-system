package gorm

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/Kengathua/book-inventory-system/pkg/config"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPGDBClient() *DBClient {
	conf, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load config")
	}

	dbUrl := conf.DatabaseURL
	db, err := connectToPostgresDB(dbUrl, conf.Environment)
	if err != nil {
		panic(err)
	}

	return &DBClient{DB: db}
}

func connectToPostgresDB(dbUrl string, enviroment string) (*gorm.DB, error) {
	_ = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Disable color
		},
	)

	DB, err := gorm.Open(postgresDriver.Open(dbUrl), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, errors.New("failed to connect database")
	}

	log.Printf("Successfully connected to db\n")
	return DB, nil
}
