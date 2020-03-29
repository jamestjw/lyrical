package playlist

// Song is a type that contains information about a song saved in the bot
type Song struct {
	YoutubeID string
	Name      string
	Next      *Song
}

// Playlist contains all the songs available for the guild
type Playlist struct {
	Songs []*Song
}

var LyricalPlaylist = &Playlist{}

// IsEmpty is a method on a playlist to check if it is empty
func (p *Playlist) IsEmpty() bool {
	return len(p.Songs) == 0
}

// First returns first song in the playlist
func (p *Playlist) First() *Song {
	return p.Songs[0]
}

// Last returns last song in the playlist
func (p *Playlist) Last() *Song {
	return p.Songs[len(p.Songs)-1]
}

// AddSong adds a song with this youtubeID to a playlist
func (p *Playlist) AddSong(songName string, youtubeID string) *Song {
	newSong := &Song{
		YoutubeID: youtubeID,
		Name:      songName,
	}
	if !p.IsEmpty() {
		p.Last().Next = newSong
	}
	p.Songs = append(p.Songs, newSong)
	return newSong
}
