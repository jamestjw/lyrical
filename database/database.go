package database

import (
	"database/sql"
	"log"

	// Load Sqlite driver
	_ "github.com/mattn/go-sqlite3"
)

// Connection is the connection to our Sqlite database
var Connection *sql.DB

func init() {
	var err error
	Connection, err = sql.Open("sqlite3", "db/discordbot.db")
	if err != nil {
		log.Fatal(err)
	}

	statement, _ := Connection.Prepare("CREATE TABLE IF NOT EXISTS songs (id INTEGER PRIMARY KEY, youtube_id TEXT, name TEXT)")
	statement.Exec()
}
