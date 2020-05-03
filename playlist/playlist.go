package playlist

import (
	"sync"

	"github.com/jamestjw/lyrical/utils"
)

// Song is a type that contains information about a song saved in the bot
type Song struct {
	YoutubeID string
	Name      string
	Next      *Song
}

// Playlist contains all the songs available for the guild
type Playlist struct {
	Songs      []*Song
	nowPlaying *Song
	next       *Song

	m sync.Mutex
}

// IsEmpty is a method on a playlist to check if it is empty
func (p *Playlist) IsEmpty() bool {
	p.m.Lock()
	defer p.m.Unlock()

	return p.isEmpty()
}

func (p *Playlist) isEmpty() bool {
	return len(p.Songs) == 0
}

// First returns first song in the playlist
func (p *Playlist) First() *Song {
	p.m.Lock()
	defer p.m.Unlock()

	return p.Songs[0]
}

// Last returns last song in the playlist
func (p *Playlist) Last() *Song {
	p.m.Lock()
	defer p.m.Unlock()

	return p.last()
}

func (p *Playlist) last() *Song {
	// p.m.Lock()
	// defer p.m.Unlock()

	return p.Songs[len(p.Songs)-1]
}

// IsPlayingMusic returns whether or not now playing is populated
func (p *Playlist) IsPlayingMusic() bool {
	p.m.Lock()
	defer p.m.Unlock()
	return p.nowPlaying != nil
}

// NowPlayingName returns name of the currently playing song
func (p *Playlist) NowPlayingName() string {
	p.m.Lock()
	defer p.m.Unlock()
	return p.nowPlaying.Name
}

// SetNowPlaying sets the currently playing song
func (p *Playlist) SetNowPlaying(s *Song) {
	p.m.Lock()
	defer p.m.Unlock()
	p.nowPlaying = s
	utils.LogInfo("", utils.KvForEvent("set-now-playing", utils.KVs("name", s.Name)))
}

// RemoveNowPlaying removes the currently playing song
func (p *Playlist) RemoveNowPlaying() {
	p.m.Lock()
	defer p.m.Unlock()
	p.nowPlaying = nil
}

// QueueNext sets a song as next in the playlist
func (p *Playlist) QueueNext(s *Song) {
	p.m.Lock()
	defer p.m.Unlock()
	p.next = s
}

// GetNext returns the next song in the playlist
func (p *Playlist) GetNext() *Song {
	p.m.Lock()
	defer p.m.Unlock()

	return p.next
}

// AddSong adds a song with this youtubeID to a playlist
func (p *Playlist) AddSong(songName string, youtubeID string) *Song {
	p.m.Lock()
	defer p.m.Unlock()

	newSong := &Song{
		YoutubeID: youtubeID,
		Name:      songName,
	}
	if !p.isEmpty() {
		p.last().Next = newSong
	}
	p.Songs = append(p.Songs, newSong)
	return newSong
}

// GetNextSongs returns the next songs that will be played in the playlist
func (p *Playlist) GetNextSongs() []*Song {
	p.m.Lock()
	defer p.m.Unlock()

	var songs []*Song
	nextSong := p.next

	for nextSong != nil {
		songs = append(songs, nextSong)
		nextSong = nextSong.Next
	}

	return songs
}
