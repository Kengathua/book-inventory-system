package main

import (
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/Kengathua/book-inventory-system/pkg/infrastructure/database/gorm"
	"github.com/Kengathua/book-inventory-system/pkg/models"
)

var (
	dbClient = &gorm.DBClient{}
)

func init() {
	dbClient = gorm.NewPGDBClient()
}

func main() {
	dbClient.DB.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp" WITH SCHEMA public;`)
	dbClient.DB.AutoMigrate(&models.Token{})
	dbClient.DB.AutoMigrate(&models.User{})
	dbClient.DB.AutoMigrate(&models.Author{})
	dbClient.DB.AutoMigrate(&models.Librarian{})
	dbClient.DB.AutoMigrate(&models.StoreKeeper{})
	dbClient.DB.AutoMigrate(&models.Book{})
}
