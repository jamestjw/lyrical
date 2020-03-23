package main

import (
	"log"

	"github.com/jamestjw/lyrical/database"
)

// AddSongToDB adds song details to the database
func AddSongToDB(s *Song) error {
	statement, err := database.Connection.Prepare("INSERT INTO songs (youtube_id, name) VALUES (?, ?)")
	if err != nil {
		log.Print(err)
		return err
	}

	_, err = statement.Exec(s.YoutubeID, s.Name)
	if err != nil {
		log.Print(err)
		return err
	}
	return nil
}
