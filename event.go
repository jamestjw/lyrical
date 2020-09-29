package main

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"

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
func (e DiscordEvent) SendMessage(message string) *discordgo.Message {
	m, err := e.session.ChannelMessageSend(e.message.ChannelID, message)
	if err != nil {
		log.Fatal(err)
	}
	return m
}

// SendMessageWithMentions sends a message to the channel within the guild that invoked
// the event while mentioning a list of users present.
// message: Message to send
// userIDs: List of users to mention
func (e DiscordEvent) SendMessageWithMentions(message string) *discordgo.Message {
	// m, err := e.session.ChannelMessageSend(e.message.ChannelID, message)
	message = fmt.Sprintf("<@james412> %s", message)
	m, err := e.session.ChannelMessageSendComplex(e.message.ChannelID, &discordgo.MessageSend{
		Content: message,
		AllowedMentions: &discordgo.MessageAllowedMentions{
			Parse: []discordgo.AllowedMentionType{discordgo.AllowedMentionTypeUsers},
		},
	})
	if err != nil {
		log.Fatal(err)
	}
	return m
}

// SendQuotedMessage sends a message to the channel within
// the guild that invoked this event with an added quote.
func (e DiscordEvent) SendQuotedMessage(quote string, message string) *discordgo.Message {
	quotedMessage := fmt.Sprintf(`>>> %s`, quote)
	e.SendMessage(quotedMessage)
	return e.SendMessage(message)
}

// React will add a reaction from the bot to the message that triggered
// this event.
func (e DiscordEvent) React(emoji string) {
	e.ReactToMessage(emoji, e.message.ID)
}

// ReactToMessage will react to a particular message in the channel that triggered this event
func (e DiscordEvent) ReactToMessage(emoji string, messageID string) {
	err := e.session.MessageReactionAdd(e.message.ChannelID, messageID, emoji)
	if err != nil {
		log.Fatal(err)
	}
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

func (e DiscordEvent) GetChannelID() string {
	return e.message.ChannelID
}

func (e DiscordEvent) GetMessageByMessageID(messageID string) (*discordgo.Message, error) {
	m, err := e.session.ChannelMessage(e.GetChannelID(), messageID)
	return m, err
}

// GetReactionsFromMessage will fetch a list of IDs of users that reacted
// to a particular message.
func (e DiscordEvent) GetReactionsFromMessage(messageID string) (map[string][]string, error) {
	var reactions map[string][]string
	m, err := e.GetMessageByMessageID(messageID)

	if err != nil {
		return make(map[string][]string), err
	}
	for _, reaction := range m.Reactions {
		var userIDs []string

		users, err := e.session.MessageReactions(e.message.ChannelID, e.message.ID, reaction.Emoji.ID, 100, "", "")

		if err != nil {
			return make(map[string][]string), fmt.Errorf("unable to fetch users that reacted to the message")
		}

		for _, user := range users {
			userIDs = append(userIDs, user.ID)
		}
		reactions[reaction.Emoji.Name] = userIDs
	}

	return reactions, nil
}
