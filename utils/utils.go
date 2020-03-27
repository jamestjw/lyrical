package utils

import (
	"errors"
	"time"

	"github.com/bwmarrin/discordgo"
)

// FindVoiceChannel tries to find a voice channel based on a
// given name within a Guild.
func FindVoiceChannel(s *discordgo.Session, guildID string, channelName string) (*discordgo.Channel, error) {
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

// VideoDurationValid parses the duration of a YouTube video
// and checks if it valid
func VideoDurationValid(videoDuration time.Duration) (err error) {
	if videoDuration.Minutes() > 10 {
		err = errors.New("video is more than 10 minutes long")
	}
	return
}
