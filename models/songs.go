package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

type Song struct {
	gorm.Model
	YoutubeID string `gorm:"unique;not null"`
	Name      string
}

// AddSongToDB adds song details to the models
func (db *DB) AddSongToDB(name string, youtubeID string) error {
	song := Song{Name: name, YoutubeID: youtubeID}

	return db.Create(song).Error
}

// SongExists checks if a given youtubeID corresponds to a song in the models
func (db *DB) SongExists(youtubeID string) (name string, exists bool) {
	var song Song
	err := db.Where(Song{YoutubeID: youtubeID}).First(&song).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		panic(err)
	}

	if song != (Song{}) {
		exists = true
		name = song.Name
	}
	return
}

func (db *DB) GetRandomSongs(limit int) []Song {
	var songs []Song
	db.Order("random()").Limit(limit).Find(&songs)
	return songs
}
