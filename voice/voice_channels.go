package voice

import "github.com/jamestjw/lyrical/playlist"

type voiceChannels struct {
	// Maps GuildID to voiceChannel
	ChannelMap map[string]Channel
}

type voiceChannel struct {
	NowPlaying   *playlist.Song
	Next         *playlist.Song
	AbortChannel chan string
	MusicActive  bool
}

func NewActiveVoiceChannels() *voiceChannels {
	var vcs voiceChannels
	vcs.ChannelMap = make(map[string]Channel)
	return &vcs
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
