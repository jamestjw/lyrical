package voice

import (
	"sync"

	"github.com/jamestjw/lyrical/playlist"
)

type voiceChannel struct {
	m              sync.Mutex
	AbortChannel   chan string
	Playlist       *playlist.Playlist
	BackupPlaylist *playlist.Playlist
}

// NewActiveVoiceChannels is a factory method to create voice channels map
func NewActiveVoiceChannels() map[string]Channel {
	vcs := make(map[string]Channel)
	return vcs
}

func NewVoiceChannel() *voiceChannel {
	return &voiceChannel{
		AbortChannel:   make(chan string, 1),
		Playlist:       &playlist.Playlist{},
		BackupPlaylist: &playlist.Playlist{},
	}
}

func (vc *voiceChannel) GetNowPlayingName() string {
	return vc.Playlist.NowPlayingName()
}

func (vc *voiceChannel) GetNext() *playlist.Song {
	s := vc.Playlist.GetNext()
	vc.setNowPlaying(s)
	return s
}

func (vc *voiceChannel) GetBackupNext() *playlist.Song {
	s := vc.BackupPlaylist.GetNext()
	vc.Playlist.SetNowPlaying(s)
	vc.BackupPlaylist.QueueNext(s.Next)
	return s
}

func (vc *voiceChannel) ExistsBackupNext() bool {
	// If BackupPlaylist has been completely used up, load again from DB.
	if vc.BackupPlaylist.GetNext() == nil {
		DB.LoadPlaylist(vc.BackupPlaylist)
	}

	return vc.BackupPlaylist.GetNext() != nil
}
func (vc *voiceChannel) ExistsNext() bool {
	return vc.Playlist.GetNext() != nil
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
	vc.setNowPlaying(s)
}

func (vc *voiceChannel) setNowPlaying(s *playlist.Song) {
	vc.Playlist.SetNowPlaying(s)
	vc.Playlist.QueueNext(s.Next)
}

func (vc *voiceChannel) RemoveNowPlaying() {
	vc.Playlist.RemoveNowPlaying()
}

func (vc *voiceChannel) FetchPlaylist() *playlist.Playlist {
	return vc.Playlist
}

func (vc *voiceChannel) GetNextSongs() (nextSongs []*playlist.Song, exists bool) {
	nextSongs = vc.Playlist.GetNextSongs()
	exists = len(nextSongs) > 0
	return
}

func (vc *voiceChannel) GetNextBackupSongs() (nextSongs []*playlist.Song, exists bool) {
	nextSongs = vc.BackupPlaylist.GetNextSongs()
	exists = len(nextSongs) > 0
	return
}
