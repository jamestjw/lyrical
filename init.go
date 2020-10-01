package main

import (
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/models"
	"github.com/jamestjw/lyrical/searcher"
	"github.com/jamestjw/lyrical/voice"
	log "github.com/sirupsen/logrus"
)

var searchService Searcher

func initialiseApplication() {
	help.InitialiseHelpText()
	models.InitialiseDatabase("production")
	loadConfig()
	searchService = searcher.InitialiseYoutubeService(config.YoutubeAPIKey)
}

func shutdownApplication(dg voice.Connectable) {
	log.Println("Received signal to terminate, cleaning up...")
	// Cleanly close down the Discord session.
	voice.DisconnectAllVoiceConnections(dg.(voice.Connectable))
	err := models.DS.Close()
	if err != nil {
		log.Error("Error closing models connection")
		return
	}

	log.Info("Exit successful!")
}
