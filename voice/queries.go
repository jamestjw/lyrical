package voice

import (
	"database/sql"
	"log"
	"time"

	"github.com/jamestjw/lyrical/database"
)

var DB Database

type SongDatabase struct{}

func init() {
	DB = SongDatabase{}
}

// AddSongToDB adds song details to the database
func (SongDatabase) AddSongToDB(name string, youtubeID string) error {
	statement, err := database.Connection.Prepare("INSERT INTO songs (youtube_id, name, created_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Print(err)
		return err
	}
	defer statement.Close()

	_, err = statement.Exec(youtubeID, name, time.Now().Format("2006-01-02 15:04:05.000"))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}

// SongExists checks if a given youtubeID corresponds to a song in the database
func (SongDatabase) SongExists(youtubeID string) (name string, exists bool) {
	err := database.Connection.QueryRow("SELECT name from songs where youtube_id = ? LIMIT 1", youtubeID).Scan(&name)
	if err != nil && err != sql.ErrNoRows {
		log.Fatalf("error checking if row exists %v", err)
	}
	if name != "" {
		exists = true
	}
	return
}
