package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/voice"
)

// Session that represents current session of discord bot
type Session interface {
	CloseConnection() error
	ListenAndServe() error
	AddHandler(interface{}) func()
}

// Event is an interface for a discord message event
type Event interface {
	SendMessage(message string)
	FindVoiceChannel(channelName string) (channelID string, err error)
	getSession() voice.Connectable
	getGuildID() string
	getVoiceConnection() (voice.Connection, bool)
}

type botSession struct {
	session *discordgo.Session
}
