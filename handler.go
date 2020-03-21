package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
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
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func joinChannelRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!join-voice" {
		if alreadyInVoiceChannel(s, m.GuildID) {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I am already in Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, config.VoiceChannelID))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, config.VoiceChannelID))
			log.Printf("Joining Guild ID: %s ChannelID: %v \n", m.GuildID, config.VoiceChannelID)
			vc := joinVoiceChannel(s, m.GuildID, config.VoiceChannelID)
			playMusic(vc)
		}
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID string, voiceChannelID string) *discordgo.VoiceConnection {
	vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, false)
	if err != nil {
		log.Fatal(err)
	}
	activeVoiceChannels.channelMap[vc] = make(chan string)
	return vc
}
