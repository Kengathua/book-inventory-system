package gorm

import (
	"testing"

	"gorm.io/gorm"
)

var (
	tb       testing.TB
	dbClient = &DBClient{}
)

type DBClient struct {
	DB *gorm.DB
}
