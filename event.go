package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Event will be passed into handler functions and contains all
// that is necessary to respond accordingly.
type Event struct {
	session *discordgo.Session
	message *discordgo.MessageCreate
}

// SendMessage sends a message to the channel within
// the guild that invoked this event.
func (e Event) SendMessage(message string) {
	e.session.ChannelMessageSend(e.message.ChannelID, message)
}

// FindVoiceChannel tries to find a voice channel with this channel name
// within the guild.
func (e Event) FindVoiceChannel(channelName string) (channelID string, err error) {
	channels, err := e.session.GuildChannels(e.message.GuildID)
	if err != nil {
		return
	}

	for _, channel := range channels {
		if channel.Bitrate == 0 || channel.ParentID == "" {
			continue
		}

		if channel.Name == channelName {
			channelID = channel.ID
			return
		}
	}
	err = errors.New("unable to find channel with this name")
	return
}

func (e Event) getSession() botSession {
	return botSession{e.session}
}

func (e Event) getGuildID() string {
	return e.message.GuildID
}

func (e Event) getVoiceConnection() (vc *discordgo.VoiceConnection, connected bool) {
	vc, connected = e.session.VoiceConnections[e.getGuildID()]
	return
}
