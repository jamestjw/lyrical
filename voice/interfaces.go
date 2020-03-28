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
}
