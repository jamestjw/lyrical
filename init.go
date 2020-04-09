package main

import (
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/searcher"
	"github.com/jamestjw/lyrical/voice"
)

var searchService Searcher

func initialiseApplication() {
	help.InitialiseHelpText()
	voice.ConnectToDatabase()
	loadConfig()
	searchService = searcher.InitialiseSearchService(config.YoutubeAPIKey)
}
