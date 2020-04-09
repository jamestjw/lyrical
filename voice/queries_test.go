package voice_test

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/voice"
	"github.com/stretchr/testify/assert"
)

func setup() {
	voice.ConnectToDatabase("test")
}

func teardown(c *sql.DB) {
	statement, err := c.Prepare("DROP TABLE IF EXISTS songs")

	if err != nil {
		log.Fatal(err)
	}

	defer statement.Close()
	statement.Exec()
}

func cleanSongs() {
	voice.DB.(voice.SongDatabase).Connection.Delete(database.Song{})
}

func TestMain(m *testing.M) {
	setup()

	os.Exit(m.Run())

	// teardown(c)
}
func TestAddSongToDatabase(t *testing.T) {
	cleanSongs()

	name, exists := voice.DB.SongExists("Youtube ID")
	assert.False(t, exists)

	err := voice.DB.AddSongToDB("Song Name", "Youtube ID")
	assert.Nil(t, err)

	name, exists = voice.DB.SongExists("Youtube ID")
	assert.True(t, exists)
	assert.Equal(t, name, "Song Name")
}

func TestLoadPlaylist(t *testing.T) {
	cleanSongs()

	voice.DB.AddSongToDB("Song Name 1", "Youtube ID 1")
	voice.DB.AddSongToDB("Song Name 2", "Youtube ID 2")

	p := &playlist.Playlist{}
	voice.DB.LoadPlaylist(p)

	assert.Contains(t, []string{p.First().Name, p.First().Next.Name}, "Song Name 1")
	assert.Contains(t, []string{p.First().Name, p.First().Next.Name}, "Song Name 2")
	assert.Contains(t, p.GetNext().Name, "Song Name")
}
