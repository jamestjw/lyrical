package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/matcher"
	"github.com/jamestjw/lyrical/searcher"
	"github.com/jamestjw/lyrical/utils"
	"github.com/jamestjw/lyrical/voice"
)

func init() {
	searcher.InitialiseSearchService(config.YoutubeAPIKey)
}

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

	channel, err := utils.FindVoiceChannel(s, m.GuildID, channelName)

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Unable to find channel fo this name in the server.")
		return
	}

	if voice.AlreadyInVoiceChannel(s, m.GuildID) {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I am already in Voice Channel within Guild ID: %s", m.GuildID))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, channel.ID))
		log.Printf("Joining Guild ID: %s ChannelID: %v \n", m.GuildID, channel.ID)

		vc := voice.JoinVoiceChannel(s, m.GuildID, channel.ID)
		nextSong := voice.ActiveVoiceChannels.ChannelMap[m.GuildID].Next
		if nextSong == nil {
			s.ChannelMessageSend(m.ChannelID, "Playlist is still empty.")
		} else {
			go voice.PlayMusic(vc, nextSong)
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
		err := voice.DisconnectAllVoiceConnections(&botSession{s})

		if err != nil {
			s.ChannelMessageSend(m.ChannelID, err.Error())
		}
	}
}

func addToPlaylistRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	matched, query, err := matcher.Match(matcher.AddPlaylistRequestRe, m.Content, "!add-playlist", "youtube-id")

	if !matched {
		return
	}

	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	youtubeID, err := searcher.GetVideoID(query)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	s.ChannelMessageSend(m.ChannelID, "Adding to playlist 😉")

	title, err := voice.AddSong(youtubeID, m.GuildID)
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
			thisVoiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
			if thisVoiceChannel.MusicActive {
				s.ChannelMessageSend(m.ChannelID, "I am already playing music 😁")
			} else {
				if thisVoiceChannel.Next == nil {
					s.ChannelMessageSend(m.ChannelID, "Playlist is currently empty.")
				} else {
					go voice.PlayMusic(vc, thisVoiceChannel.Next)
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
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel. 😔")
		} else {
			voiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
			if voiceChannel.MusicActive {
				voiceChannel.AbortChannel <- "stop"
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
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel. 😔")
		} else {
			if voice.ActiveVoiceChannels.ChannelMap[vc.GuildID].MusicActive {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Now playing: **%s**", voice.ActiveVoiceChannels.ChannelMap[vc.GuildID].GetNowPlayingName()))
			} else {
				s.ChannelMessageSend(m.ChannelID, "Well I am not playing any music currently 🤔")
			}
		}
	}
}

func skipMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!skip-music" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel yet. 😔")
		} else {
			thisVoiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
			if thisVoiceChannel.MusicActive {
				thisVoiceChannel.AbortChannel <- "stop"
				s.ChannelMessageSend(m.ChannelID, "Skipping song... ❌")
				if thisVoiceChannel.Next != nil {
					go voice.PlayMusic(vc, thisVoiceChannel.Next)
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "Well I am not playing any music currently 🤔")
			}
		}
	}
}
