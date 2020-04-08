package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type Song struct {
	gorm.Model
	YoutubeID string `gorm:"unique;not null"`
	Name      string
}

var Connection *gorm.DB

func InitialiseDatabase() {
	db, err := gorm.Open("sqlite3", "db/discordbot.db")
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Song{})

	Connection = db
}
