package main

import (
	"encoding/json"
	"os"

	log "github.com/sirupsen/logrus"
)

var config Config

// Config has all application level configurations
type Config struct {
	DiscordToken        string `yaml:"discordToken"`
	YoutubeAPIKey       string `yaml:"youtubeApiKey"`
	UpNextMaxSongsCount int    `yaml:"upNextMaxSongsCount"`
}

func loadConfig() {
	file, err := os.Open("configs/config.json")
	if err != nil {
		log.Fatalln("Cannot open config file", err)
	}
	decoder := json.NewDecoder(file)
	config = Config{}
	err = decoder.Decode(&config)
	if err != nil {
		log.Fatalln("Cannot get configuration from file", err)
	}
}
