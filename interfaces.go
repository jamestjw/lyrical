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
	SendMessage(message string) *discordgo.Message
	SendMessageWithMentions(message string, userIDs []string) *discordgo.Message
	FindVoiceChannel(channelName string) (channelID string, err error)
	GetSession() voice.Connectable
	GetGuildID() string
	GetVoiceConnection() (voice.Connection, bool)
	GetMessageByMessageID(messageID string) (*discordgo.Message, error)
	React(emoji string)
	ReactToMessage(emoji string, messageID string)
	SendQuotedMessage(quote string, message string) *discordgo.Message
	SendQuotedMessageWithMentions(quote string, message string, userIDs []string) *discordgo.Message
	GetReactionsFromMessage(messageID string) (map[string][]string, error)
	GetUserForBot() (*discordgo.User, error)
}

type Searcher interface {
	GetVideoID(string) (string, error)
}

type botSession struct {
	session *discordgo.Session
}
