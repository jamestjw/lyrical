package main

import (
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/searcher"
	"github.com/jamestjw/lyrical/voice"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
)

var searchService Searcher
var DB *gorm.DB

func initialiseApplication() {
	help.InitialiseHelpText()
	DB = database.InitialiseDatabase("production")
	voice.ConnectToDatabase(DB)
	loadConfig()
	searchService = searcher.InitialiseYoutubeService(config.YoutubeAPIKey)
}

func shutdownApplication(dg voice.Connectable) {
	log.Println("Received signal to terminate, cleaning up...")
	// Cleanly close down the Discord session.
	voice.DisconnectAllVoiceConnections(dg.(voice.Connectable))
	err := DB.Close()
	if err != nil {
		log.Error("Error closing database connection")
		return
	}

	log.Info("Exit successful!")
}
