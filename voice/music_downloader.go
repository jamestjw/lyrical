package voice

import (
	"github.com/jamestjw/lyrical/ytmp3"
)

var Dl Downloader

func init() {
	Dl = &MusicDownloader{downloadMusicFunction: ytmp3.Download}
}

type MusicDownloader struct {
	downloadMusicFunction func(string) (string, error)
}

func (d *MusicDownloader) Download(query string) (title string, err error) {
	title, err = d.downloadMusicFunction(query)
	return
}
