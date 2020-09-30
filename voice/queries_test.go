package voice_test

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/mocks/mock_database"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/voice"
	"github.com/stretchr/testify/assert"
)

func setup() {
	// voice.ConnectToDatabase(database.InitialiseDatabase("test"))
}

func cleanSongs() {
	// voice.DB.(voice.SongDatabase).Connection.Delete(database.Song{})
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
	cleanSongs()
}

func TestLoadPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	songs := []database.Song{
		database.Song{Name: "Song Name 1"},
		database.Song{Name: "Song Name 2"},
	}

	mockDS := mock_database.NewMockDatastore(ctrl)
	mockDS.EXPECT().GetRandomSongs(20).Return(songs)
	database.DS = mockDS

	p := &playlist.Playlist{}
	voice.LoadPlaylist(p)

	assert.Contains(t, []string{p.First().Name, p.First().Next.Name}, "Song Name 1")
	assert.Contains(t, []string{p.First().Name, p.First().Next.Name}, "Song Name 2")
	assert.Contains(t, p.GetNext().Name, "Song Name")
}
