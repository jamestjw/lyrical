package voice

import "github.com/jamestjw/lyrical/playlist"

type Connectable interface {
	GetVoiceConnections() map[string]Connection
	JoinVoiceChannel(guildID string, voiceChannelID string) (Connection, error)
}

type Connection interface {
	Disconnect() (err error)
	GetGuildID() string
	GetAudioInputChannel() chan []byte
}

type Channel interface {
	RemoveNowPlaying()
	GetNext() *playlist.Song
	SetNext(*playlist.Song)
	SetNowPlaying(s *playlist.Song)
	GetAbortChannel() chan string
	IsPlayingMusic() bool
	GetNowPlayingName() string
	StopMusic()
	FetchPlaylist() *playlist.Playlist
}

type Database interface {
	SongExists(youtubeID string) (string, bool)
	AddSongToDB(name string, youtubeID string) error
}

type Downloader interface {
	Download(query string) (title string, err error)
}
