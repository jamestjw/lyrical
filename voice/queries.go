package voice

import (
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jinzhu/gorm"
)

var DB Database

type SongDatabase struct {
	connection *gorm.DB
}

func ConnectToDatabase() {
	DB = SongDatabase{database.InitialiseDatabase()}
}

// AddSongToDB adds song details to the database
func (db SongDatabase) AddSongToDB(name string, youtubeID string) error {
	song := &database.Song{Name: name, YoutubeID: youtubeID}

	db.connection.Create(song)
	return nil
}

// SongExists checks if a given youtubeID corresponds to a song in the database
func (db SongDatabase) SongExists(youtubeID string) (name string, exists bool) {
	var song database.Song
	db.connection.Where(&database.Song{YoutubeID: youtubeID}).First(&song)
	if song != (database.Song{}) {
		exists = true
		name = song.Name
	}
	return
}

// LoadPlaylist will load a playlist from the database.
func (db SongDatabase) LoadPlaylist(p *playlist.Playlist) {
	var songs []database.Song
	db.connection.Order("random()").Limit(10).Find(&songs)

	for _, song := range songs {
		p.AddSong(song.Name, song.YoutubeID)
	}

	p.QueueNext(p.First())
}
