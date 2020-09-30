package voice

import (
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/playlist"
)

// LoadPlaylist will load a playlist from the database.
func LoadPlaylist(p *playlist.Playlist) {
	songs := database.DS.GetRandomSongs(20)
	for i, song := range songs {
		newSong := p.AddSong(song.Name, song.YoutubeID)
		if i == 0 && p.GetNext() == nil {
			p.QueueNext(newSong)
		}
	}
}
