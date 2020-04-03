package voice

import "github.com/jamestjw/lyrical/playlist"

type voiceChannel struct {
	NowPlaying   *playlist.Song
	Next         *playlist.Song
	AbortChannel chan string
	MusicActive  bool
	Playlist     *playlist.Playlist
}

// NewActiveVoiceChannels is a factory method to create voice channels map
func NewActiveVoiceChannels() map[string]Channel {
	vcs := make(map[string]Channel)
	return vcs
}

func (vc *voiceChannel) GetNowPlayingName() string {
	return vc.NowPlaying.Name
}

func (vc *voiceChannel) GetNext() *playlist.Song {
	return vc.Next
}

func (vc *voiceChannel) SetNext(s *playlist.Song) {
	vc.Next = s
}

func (vc *voiceChannel) GetAbortChannel() chan string {
	return vc.AbortChannel
}

func (vc *voiceChannel) IsPlayingMusic() bool {
	return vc.MusicActive
}

func (vc *voiceChannel) StopMusic() {
	vc.AbortChannel <- "stop"
}

func (vc *voiceChannel) SetNowPlaying(s *playlist.Song) {
	vc.MusicActive = true
	vc.NowPlaying = s
	vc.Next = s.Next
}

func (vc *voiceChannel) RemoveNowPlaying() {
	vc.MusicActive = false
	vc.NowPlaying = nil
}

func (vc *voiceChannel) FetchPlaylist() *playlist.Playlist {
	return vc.Playlist
}
