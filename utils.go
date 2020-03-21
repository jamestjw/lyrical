package main

import (
	"errors"

	"github.com/bwmarrin/discordgo"
)

// Finds a voice channel based on a given name within a Guild.
func findVoiceChannel(s *discordgo.Session, guildID string, channelName string) (*discordgo.Channel, error) {
	channels, err := s.GuildChannels(guildID)
	if err != nil {
		return nil, err
	}

	for _, channel := range channels {
		if channel.Bitrate == 0 || channel.ParentID == "" {
			continue
		}

		if channel.Name == channelName {
			return channel, nil
		}
	}

	return nil, errors.New("unable to find channel with this name")
}
