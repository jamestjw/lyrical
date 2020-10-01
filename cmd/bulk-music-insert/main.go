package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/jamestjw/lyrical/models"
	"github.com/jamestjw/lyrical/searcher"
	"github.com/jamestjw/lyrical/voice"
)

// Loading config, move this somewhere else? Or make the main config module shareable
var config bulkInsertConfig
var searchService *searcher.YoutubeService

type bulkInsertConfig struct {
	YoutubeAPIKey string `yaml:"youtubeApiKey"`
}

func loadConfig() {
	file, err := os.Open("configs/config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = bulkInsertConfig{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
	log.Print("Loaded configuration file.")
}

// Main function
func main() {
	playlistIDPtr := flag.String("playlist-id", "", "ID of a Youtube playlist")
	maxItemsPtr := flag.Int64("max-items", 20, "Max number of videos to take from this playlist")
	flag.Parse()

	videoIDs, err := searchService.GetVideoIDs(*playlistIDPtr, *maxItemsPtr)

	if err != nil {
		log.Fatalln(err)
	}

	addSongs(videoIDs)
}

func init() {
	loadConfig()
	setupLogger()
	searchService = searcher.InitialiseYoutubeService(config.YoutubeAPIKey)
	models.InitialiseDatabase("production")
}

func addSongs(videoIDs []string) {
	for _, videoID := range videoIDs {
		title, err := addSong(videoID)
		if err != nil {
			log.Error(err)
		} else {
			log.Printf("Successfully downloaded song with title: %s", title)
		}
	}
}

func addSong(youtubeID string) (title string, err error) {
	title, exists := models.DS.SongExists(youtubeID)

	if exists {
		err = fmt.Errorf("The song %s (%s) already exists in the models", title, youtubeID)
	} else {
		title, err = voice.Dl.Download(youtubeID)
		if err != nil {
			err = fmt.Errorf("Error adding the song %s 🤨: %s", youtubeID, err.Error())
			return
		}

		dbErr := models.DS.AddSongToDB(title, youtubeID)
		if dbErr != nil {
			log.Printf("Error writing song ID %s to the models: %s", youtubeID, dbErr)
		}
	}
	return
}

func setupLogger() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)

	logFile, err := os.OpenFile("log/bulk-music-insert.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, logFile)
		log.SetOutput(mw)
	} else {
		log.Info("Failed to log to file, using default stderr")
	}
}
