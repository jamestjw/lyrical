package main

import (
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/searcher"
)

var searchService Searcher

func initialiseApplication() {
	help.InitialiseHelpText()
	database.InitialiseDatabase()
	loadConfig()
	searchService = searcher.InitialiseSearchService(config.YoutubeAPIKey)
}
