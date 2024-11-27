package storage

import (
	_ "embed"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	_ "modernc.org/sqlite"
)

// Can be set to file::memory:?cache=shared
const dataSourceName string = "file:sqlite.db?cache=shared&mode=rwc"

func NewDbStorage() *gorm.DB {
	db, err := gorm.Open(sqlite.New(sqlite.Config{
		DSN:        dataSourceName,
		DriverName: "sqlite",
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func DBMigrate(db *gorm.DB) {
	if err := db.AutoMigrate(&Note{}); err != nil {
		panic("Failed to run migrations:" + err.Error())
	}
}

func SeedData(db *gorm.DB) {
	var count int64
	db.Model(&Note{}).Count(&count)
	if count == 0 {
		db.Create(&NotesSeed)
	}
}
