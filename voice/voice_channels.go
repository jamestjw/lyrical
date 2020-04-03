package voice

import "github.com/jamestjw/lyrical/playlist"

type voiceChannel struct {
	AbortChannel chan string
	Playlist     *playlist.Playlist
}

// NewActiveVoiceChannels is a factory method to create voice channels map
func NewActiveVoiceChannels() map[string]Channel {
	vcs := make(map[string]Channel)
	return vcs
}

func (vc *voiceChannel) GetNowPlayingName() string {
	return vc.Playlist.NowPlayingName()
}

func (vc *voiceChannel) GetNext() *playlist.Song {
	return vc.Playlist.GetNext()
}

func (vc *voiceChannel) SetNext(s *playlist.Song) {
	vc.Playlist.QueueNext(s)
}

func (vc *voiceChannel) GetAbortChannel() chan string {
	return vc.AbortChannel
}

func (vc *voiceChannel) IsPlayingMusic() bool {
	return vc.Playlist.IsPlayingMusic()
}

func (vc *voiceChannel) StopMusic() {
	vc.AbortChannel <- "stop"
}

func (vc *voiceChannel) SetNowPlaying(s *playlist.Song) {
	vc.Playlist.SetNowPlaying(s)
	vc.Playlist.QueueNext(s.Next)
}

func (vc *voiceChannel) RemoveNowPlaying() {
	vc.Playlist.RemoveNowPlaying()
}

func (vc *voiceChannel) FetchPlaylist() *playlist.Playlist {
	return vc.Playlist
}
