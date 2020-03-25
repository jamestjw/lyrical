package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/matcher"
)

func dummyMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
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

func joinVoiceChannelRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	matched, channelName, err := matcher.Match(matcher.JoinChannelRequestRe, m.Content, "!join-voice", "channel-name")

	if !matched {
		return
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
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
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, channel.ID))
		log.Printf("Joining Guild ID: %s ChannelID: %v \n", m.GuildID, channel.ID)

		vc := joinVoiceChannel(s, m.GuildID, channel.ID)
		nextSong := activeVoiceChannels.channelMap[m.GuildID].Next
		if nextSong == nil {
			s.ChannelMessageSend(m.ChannelID, "Playlist is still empty.")
		} else {
			go playMusic(vc, nextSong)
			s.ChannelMessageSend(m.ChannelID, "Starting music... 🎵")
		}
	}
}

func leaveVoiceChannelRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// TODO: Leave voice channel of current guild only.
	if m.Content == "!leave-voice" {
		s.ChannelMessageSend(m.ChannelID, "Leaving voice channel 👋🏼")
		err := disconnectAllVoiceConnections(s)

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
	}
}

func addToPlaylistRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	matched, youtubeID, err := matcher.Match(matcher.AddPlaylistRequestRe, m.Content, "!add-playlist", "youtube-id")

	if !matched {
		return
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Adding to playlist 😉")
	title, err := addSong(youtubeID, m.GuildID)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Your song **%s** was added 👍", title))
}

func helpRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!help" {
		s.ChannelMessageSend(m.ChannelID, help.Message())
	}
}

func playMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!play-music" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel yet.")
		} else {
			thisVoiceChannel := activeVoiceChannels.channelMap[vc.GuildID]
			if thisVoiceChannel.MusicActive {
				s.ChannelMessageSend(m.ChannelID, "I am already playing music 😁")
			} else {
				if thisVoiceChannel.Next == nil {
					s.ChannelMessageSend(m.ChannelID, "Playlist is still empty.")
				} else {
					go playMusic(vc, thisVoiceChannel.Next)
					s.ChannelMessageSend(m.ChannelID, "Starting music... 🎵")
				}
			}
		}
	}
}

func stopMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!stop-music" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel.")
		} else {
			if activeVoiceChannels.channelMap[vc.GuildID].MusicActive {
				activeVoiceChannels.channelMap[vc.GuildID].AbortChannel <- "stop"
				s.ChannelMessageSend(m.ChannelID, "OK, Shutting up now...")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Well I am not playing any music currently 🤔")
			}
		}
	}
}

func nowPlayingRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!now-playing" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel.")
		} else {
			if activeVoiceChannels.channelMap[vc.GuildID].MusicActive {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Now playing: **%s**", activeVoiceChannels.channelMap[vc.GuildID].GetNowPlayingName()))
			} else {
				s.ChannelMessageSend(m.ChannelID, "Well I am not playing any music currently 🤔")
			}
		}
	}
}
