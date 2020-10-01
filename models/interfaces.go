package models

type Datastore interface {
	AddSongToDB(name string, youtubeID string) error
	SongExists(youtubeID string) (name string, exists bool)
	GetRandomSongs(limit int) []Song
	Close() error
}
