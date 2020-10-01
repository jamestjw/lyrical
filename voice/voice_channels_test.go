package voice_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	mock_models "github.com/jamestjw/lyrical/mocks/mock_models"
	"github.com/jamestjw/lyrical/models"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/voice"
	"github.com/stretchr/testify/assert"
)

func TestExistsNext(t *testing.T) {
	vc := voice.NewVoiceChannel()

	assert.False(t, vc.ExistsNext(), "no songs exists initially")

	vc.Playlist.QueueNext(&playlist.Song{
		Name:      "Song Name",
		YoutubeID: "Youtube ID",
	})

	assert.True(t, vc.ExistsNext(), "song exists now")
}

func TestExistsBackupNextWithEmptyDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vc := voice.NewVoiceChannel()

	mockDS := mock_models.NewMockDatastore(ctrl)
	mockDS.EXPECT().GetRandomSongs(gomock.Any()).Return(make([]models.Song, 0))
	models.DS = mockDS

	assert.False(t, vc.ExistsBackupNext(), "no songs exists initially")

	vc.BackupPlaylist.QueueNext(&playlist.Song{
		Name:      "Song Name",
		YoutubeID: "Youtube ID",
	})

	assert.True(t, vc.ExistsBackupNext(), "song exists now")
}

func TestExistsBackupNextWithPopulatedDB(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	vc := voice.NewVoiceChannel()

	mockDS := mock_models.NewMockDatastore(ctrl)
	mockDS.EXPECT().GetRandomSongs(gomock.Any()).Return([]models.Song{{}})
	models.DS = mockDS

	assert.True(t, vc.ExistsBackupNext(), "songs populated from DB")
}

func TestGetNowPlayingName(t *testing.T) {
	nowPlayingSong := &playlist.Song{
		Name: "NP Song",
	}
	vc := voice.NewVoiceChannel()
	vc.Playlist.SetNowPlaying(nowPlayingSong)

	assert.Equal(t, "NP Song", vc.GetNowPlayingName())
}

func TestGetAbortChannel(t *testing.T) {
	vc := voice.NewVoiceChannel()

	assert.IsType(t, make(chan string), vc.GetAbortChannel())
}

func TestFetchPlaylist(t *testing.T) {
	vc := voice.NewVoiceChannel()

	assert.Equal(t, vc.Playlist, vc.FetchPlaylist())
}

func TestSetNowPlaying(t *testing.T) {
	nextNextSong := &playlist.Song{}
	nextSong := &playlist.Song{Name: "nextSongName", Next: nextNextSong}

	vc := voice.NewVoiceChannel()
	vc.SetNowPlaying(nextSong)

	assert.Equal(t, "nextSongName", vc.GetNowPlayingName())
	assert.Equal(t, nextNextSong, vc.Playlist.GetNext())
}

func TestRemoveNowPlaying(t *testing.T) {
	vc := voice.NewVoiceChannel()
	vc.Playlist.SetNowPlaying(&playlist.Song{})

	assert.True(t, vc.Playlist.IsPlayingMusic(), "should initially has now playing")

	vc.RemoveNowPlaying()
	assert.False(t, vc.Playlist.IsPlayingMusic(), "shoould no longer have now playing")
}

func TestStopMusic(t *testing.T) {
	vc := voice.NewVoiceChannel()
	s := &playlist.Song{}
	vc.SetNowPlaying(s)

	assert.Empty(t, vc.AbortChannel, "should be initially empty")
	assert.True(t, vc.Playlist.IsPlayingMusic())
	vc.StopMusic()
	assert.NotEmpty(t, vc.AbortChannel, "should be not empty")
	assert.False(t, vc.Playlist.IsPlayingMusic())
}

func TestIsPlayingMusic(t *testing.T) {
	vc := voice.NewVoiceChannel()

	assert.False(t, vc.IsPlayingMusic(), "should be initially not playing")
	vc.Playlist.SetNowPlaying(&playlist.Song{})
	assert.True(t, vc.IsPlayingMusic(), "should be playing music now")
}

func TestGetNext(t *testing.T) {
	nextNextSong := &playlist.Song{}
	nextSong := &playlist.Song{Name: "nextSongName", Next: nextNextSong}

	vc := voice.NewVoiceChannel()
	vc.Playlist.QueueNext(nextSong)

	assert.Equal(t, nextSong, vc.GetNext())
	assert.Equal(t, "nextSongName", vc.GetNowPlayingName())
	assert.Equal(t, nextNextSong, vc.Playlist.GetNext())
}

func TestGetBackupNext(t *testing.T) {
	nextNextSong := &playlist.Song{}
	nextSong := &playlist.Song{Name: "nextSongName", Next: nextNextSong}

	vc := voice.NewVoiceChannel()
	vc.BackupPlaylist.QueueNext(nextSong)

	assert.Equal(t, nextSong, vc.GetBackupNext())
	assert.Equal(t, "nextSongName", vc.GetNowPlayingName())
	assert.Empty(t, vc.BackupPlaylist.GetNext())
}

func TestGetNextSongs(t *testing.T) {
	vc := voice.NewVoiceChannel()

	_, exists := vc.GetNextSongs()
	assert.False(t, exists, "should have no songs at the start")

	songs := []*playlist.Song{}

	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Song %v", i+1)
		youtubeID := fmt.Sprintf("ID %v", i+1)
		song := vc.Playlist.AddSong(name, youtubeID)
		songs = append(songs, song)
	}

	vc.Playlist.QueueNext(songs[0])

	nextSongs, exists := vc.GetNextSongs()
	assert.True(t, exists, "should now have songs")
	assert.Equal(t, songs, nextSongs, "should have songs in the order they were added")
}

func TestGetNextBackupSongs(t *testing.T) {
	vc := voice.NewVoiceChannel()

	_, exists := vc.GetNextSongs()
	assert.False(t, exists, "should have no songs at the start")

	songs := []*playlist.Song{}

	for i := 0; i < 3; i++ {
		name := fmt.Sprintf("Song %v", i+1)
		youtubeID := fmt.Sprintf("ID %v", i+1)
		song := vc.BackupPlaylist.AddSong(name, youtubeID)
		songs = append(songs, song)
	}

	vc.BackupPlaylist.QueueNext(songs[0])

	nextSongs, exists := vc.GetNextBackupSongs()
	assert.True(t, exists, "should now have songs")
	assert.Equal(t, songs, nextSongs, "should have songs in the order they were added")
}
