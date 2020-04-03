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
	return vc.Playlist.NowPlaying.Name
}

func (vc *voiceChannel) GetNext() *playlist.Song {
	return vc.Playlist.Next
}

func (vc *voiceChannel) SetNext(s *playlist.Song) {
	vc.Playlist.Next = s
}

func (vc *voiceChannel) GetAbortChannel() chan string {
	return vc.AbortChannel
}

func (vc *voiceChannel) IsPlayingMusic() bool {
	return vc.Playlist.NowPlaying != nil
}

func (vc *voiceChannel) StopMusic() {
	vc.AbortChannel <- "stop"
}

func (vc *voiceChannel) SetNowPlaying(s *playlist.Song) {
	vc.Playlist.NowPlaying = s
	vc.Playlist.Next = s.Next
}

func (vc *voiceChannel) RemoveNowPlaying() {
	vc.Playlist.NowPlaying = nil
}

func (vc *voiceChannel) FetchPlaylist() *playlist.Playlist {
	return vc.Playlist
}
