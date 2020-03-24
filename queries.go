package main

import (
	"log"
	"time"

	"github.com/jamestjw/lyrical/database"
)

// AddSongToDB adds song details to the database
func AddSongToDB(s *Song) error {
	statement, err := database.Connection.Prepare("INSERT INTO songs (youtube_id, name, created_at) VALUES (?, ?, ?)")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = statement.Exec(s.YoutubeID, s.Name, time.Now().Format("2006-01-02 15:04:05.000"))
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
