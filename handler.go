package main

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	joinChannelRequestRe = regexp.MustCompile(`^!join-voice\s?(.*)$`)
)

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		s.ChannelMessageSend(m.ChannelID, "Try running !pong ;)")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
		s.ChannelMessageSend(m.ChannelID, "Try running !ping ;)")
	}
}

func joinChannelRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	matched := joinChannelRequestRe.FindStringSubmatch(m.Content)

	if matched == nil {
		return
	}

	channelName := strings.TrimSpace(matched[1])
	if channelName == "" {
		fmt.Println("!join-voice expects an argument <channel-name> ")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Connecting to channel name: %s", channelName))

	channel, err := findVoiceChannel(s, m.GuildID, channelName)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find channel fo this name in the server.")
		return
	}

	if alreadyInVoiceChannel(s, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I am already in Voice Channel within Guild ID: %s", m.GuildID))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, config.VoiceChannelID))
		log.Printf("Joining Guild ID: %s ChannelID: %v \n", m.GuildID, channel.ID)
		vc := joinVoiceChannel(s, m.GuildID, channel.ID)
		playMusic(vc)
	}
}
