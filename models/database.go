package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var DS Datastore

type DB struct {
	*gorm.DB
}

func InitialiseDatabase(env string) {
	DbEnvMap := map[string]string{
		"production": "db/discordbot.db",
		"test":       "../db/test.db",
	}

	db, err := gorm.Open("sqlite3", DbEnvMap[env])

	if err != nil {
		panic("failed to connect models")
	}

	// Migrate the schema
	db.AutoMigrate(&Song{})

	DS = &DB{db}
}

func (db *DB) Close() error {
	return db.Close()
}
