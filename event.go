package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/voice"
)

// DiscordEvent will be passed into handler functions and contains all
// that is necessary to respond accordingly.
type DiscordEvent struct {
	session *discordgo.Session
	message *discordgo.MessageCreate
}

// SendMessage sends a message to the channel within
// the guild that invoked this event.
func (e DiscordEvent) SendMessage(message string) {
	e.session.ChannelMessageSend(e.message.ChannelID, message)
}

// FindVoiceChannel tries to find a voice channel with this channel name
// within the guild.
func (e DiscordEvent) FindVoiceChannel(channelName string) (channelID string, err error) {
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

// GetSession returns a Connectable belonging to the guild of this event
func (e DiscordEvent) GetSession() voice.Connectable {
	return botSession{e.session}
}

// GetGuildID returns the guild ID of the guild in which this event
// was invoked.
func (e DiscordEvent) GetGuildID() string {
	return e.message.GuildID
}

// GetVoiceConnection returns a Connection that the bot is connected to
// in the guild in which this event was fired.
func (e DiscordEvent) GetVoiceConnection() (voice.Connection, bool) {
	vc, connected := e.session.VoiceConnections[e.GetGuildID()]
	return voice.DGVoiceConnection{Connection: vc}, connected
}
