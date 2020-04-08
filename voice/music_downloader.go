package voice

import (
	"github.com/jamestjw/lyrical/ytmp3"
)

// Dl handles all downloading of music given a search query
var Dl Downloader

func init() {
	Dl = &MusicDownloader{downloadMusicFunction: ytmp3.Download}
}

// MusicDownloader is a struct that contains a methid that is capable of
// of downloading a song based on a query
type MusicDownloader struct {
	downloadMusicFunction func(string) (string, error)
}

// Download accepts a query string and downloads a songs based on it,
// this is the entry point that should be used to download a song.
func (d *MusicDownloader) Download(query string) (title string, err error) {
	title, err = d.downloadMusicFunction(query)
	return
}
