package main

import (
	"encoding/json"
	"log"
	"os"
)

var config Config

// Config has all application level configurations
type Config struct {
	DiscordToken   string `yaml:"discordToken"`
	GuildID        string `yaml:"guildID"`
	VoiceChannelID string `yaml:"voiceChannelID"`
}

func loadConfig() {
	file, err := os.Open("config.json")
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

func init() {
	loadConfig()
}
