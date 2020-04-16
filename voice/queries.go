package voice

import (
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jinzhu/gorm"
)

var DB Database

type SongDatabase struct {
	Connection *gorm.DB
}

func ConnectToDatabase(db *gorm.DB) {
	DB = SongDatabase{db}
}

// AddSongToDB adds song details to the database
func (db SongDatabase) AddSongToDB(name string, youtubeID string) error {
	song := &database.Song{Name: name, YoutubeID: youtubeID}

	db.Connection.Create(song)
	return nil
}

// SongExists checks if a given youtubeID corresponds to a song in the database
func (db SongDatabase) SongExists(youtubeID string) (name string, exists bool) {
	var song database.Song
	db.Connection.Where(&database.Song{YoutubeID: youtubeID}).First(&song)
	if song != (database.Song{}) {
		exists = true
		name = song.Name
	}
	return
}

// LoadPlaylist will load a playlist from the database.
func (db SongDatabase) LoadPlaylist(p *playlist.Playlist) {
	var songs []database.Song
	db.Connection.Order("random()").Limit(20).Find(&songs)

	for i, song := range songs {
		newSong := p.AddSong(song.Name, song.YoutubeID)
		if i == 0 && p.GetNext() == nil {
			p.QueueNext(newSong)
		}
	}
}
