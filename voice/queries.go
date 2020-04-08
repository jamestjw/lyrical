package voice

import (
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/playlist"
)

var DB Database

type SongDatabase struct{}

func init() {
	DB = SongDatabase{}
}

// AddSongToDB adds song details to the database
func (SongDatabase) AddSongToDB(name string, youtubeID string) error {
	song := &database.Song{Name: name, YoutubeID: youtubeID}

	database.Connection.Create(song)
	return nil
}

// SongExists checks if a given youtubeID corresponds to a song in the database
func (SongDatabase) SongExists(youtubeID string) (name string, exists bool) {
	var song database.Song
	database.Connection.Where(&database.Song{YoutubeID: youtubeID}).First(&song)
	if song != (database.Song{}) {
		exists = true
		name = song.Name
	}
	return
}

// LoadPlaylist will load a playlist from the database.
func (SongDatabase) LoadPlaylist(p *playlist.Playlist) {
	var songs []database.Song
	database.Connection.Order("random()").Limit(10).Find(&songs)

	for _, song := range songs {
		p.AddSong(song.Name, song.YoutubeID)
	}

	p.QueueNext(p.First())
}
