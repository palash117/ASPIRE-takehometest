package db

import (
	"fmt"
	"sync"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Db struct {
	DB *gorm.DB
}

func SetupDb() *Db {
	db, err := gorm.Open(sqlite.Open("twitter.db"), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("err opening connection to sql lite %s", err))
	}
	return &Db{DB: db}
}
func SetupTestDb(t *testing.T) (*gorm.DB, error) {
	t.Helper()
	once := sync.OnceValues(func() (*gorm.DB, error) {
		return gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	})
	return once()

}

type NoEntriesFound struct{}

func (n NoEntriesFound) Error() string {
	return "No entries found"
}
