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
	dbClient.DB.Migrator().DropTable(&models.Token{})
	dbClient.DB.Migrator().DropTable(&models.User{})
	dbClient.DB.Migrator().DropTable(&models.Author{})
	dbClient.DB.Migrator().DropTable(&models.Librarian{})
	dbClient.DB.Migrator().DropTable(&models.StoreKeeper{})
	dbClient.DB.Migrator().DropTable(&models.Book{})
}
