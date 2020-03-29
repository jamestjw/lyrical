package playlist_test

import (
	"testing"

	"github.com/jamestjw/lyrical/playlist"
	"github.com/stretchr/testify/assert"
)

func TestPlaylistEmpty(t *testing.T) {
	p := playlist.Playlist{}

	assert.Equal(t, len(p.Songs), 0)
	assert.True(t, p.IsEmpty(), "playlist should be empty")

	p.Songs = append(p.Songs, &playlist.Song{})

	assert.False(t, p.IsEmpty(), "playlist should be empty")
}

func TestAddSongWhenEmpty(t *testing.T) {
	p := playlist.Playlist{}

	song := p.AddSong("songName", "youtubeID")

	assert.Equal(t, song.Name, "songName", "should have the right name")
	assert.Equal(t, song.YoutubeID, "youtubeID", "should have the right ID")
	assert.Equal(t, song, p.Songs[0])
}

func TestAddSongWhenNotEmpty(t *testing.T) {
	p := playlist.Playlist{}
	p.Songs = append(p.Songs, &playlist.Song{
		YoutubeID: "YoutubeID 1",
		Name:      "Song 1",
	})

	song := p.AddSong("songName", "youtubeID")

	assert.Equal(t, song.Name, "songName", "should have the right name")
	assert.Equal(t, song.YoutubeID, "youtubeID", "should have the right ID")
	assert.Equal(t, song, p.Songs[1], "last song should be this song")
	assert.Equal(t, song, p.Songs[0].Next, "should set new song as last song of the previous")
}
