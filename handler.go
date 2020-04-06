package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/matcher"
	"github.com/jamestjw/lyrical/voice"
)

var defaultMux = &Mux{}

func init() {
	defaultMux = NewMux()
	defaultMux.RegisterHandler(matcher.JoinChannelRequestMatcher, joinVoiceChannelRequest)
	defaultMux.RegisterHandler(matcher.LeaveVoiceMatcher, leaveVoiceChannelRequest)
	defaultMux.RegisterHandler(matcher.AddPlaylistRequestMatcher, addToPlaylistRequest)
	defaultMux.RegisterHandler(matcher.PlayMusicMatcher, playMusicRequest)
	defaultMux.RegisterHandler(matcher.StopMusicMatcher, stopMusicRequest)
	defaultMux.RegisterHandler(matcher.SkipMusicMatcher, skipMusicRequest)
	defaultMux.RegisterHandler(matcher.NowPlayingMatcher, nowPlayingRequest)
	defaultMux.RegisterHandler(matcher.HelpMatcher, helpRequest)
	defaultMux.RegisterHandler(matcher.UpNextMatcher, upNextRequest)
}

func heartbeatHandlerFunc(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
		s.ChannelMessageSend(m.ChannelID, "I am alive!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
		s.ChannelMessageSend(m.ChannelID, "I am alive!")
	}
}

func joinVoiceChannelRequest(event Event, channelName string) {
	event.SendMessage(fmt.Sprintf("Connecting to channel name: %s", channelName))
	channelID, err := event.FindVoiceChannel(channelName)

	if err != nil {
		event.SendMessage("Unable to find channel fo this name in the server.")
		return
	}

	if voice.AlreadyInVoiceChannel(event.GetSession(), event.GetGuildID()) {
		event.SendMessage(fmt.Sprintf("I am already in Voice Channel within Guild ID: %s", event.GetGuildID()))
	} else {
		event.SendMessage(fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v", event.GetGuildID(), channelID))
		log.Printf("Joining Guild ID: %s ChannelID: %v \n", event.GetGuildID(), channelID)

		vc := voice.JoinVoiceChannel(event.GetSession(), event.GetGuildID(), channelID)
		thisChannel := voice.ActiveVoiceChannels[event.GetGuildID()]

		if !thisChannel.ExistsNext() {
			event.SendMessage("Playlist is still empty.")
		} else {
			voice.PlayMusic(vc.GetAudioInputChannel(), event.GetGuildID(), thisChannel, true)
			event.SendMessage("Starting music... 🎵")
		}
	}
}

func leaveVoiceChannelRequest(event Event, _ string) {
	vc, connected := event.GetVoiceConnection()

	if connected {
		event.SendMessage("Leaving voice channel 👋🏼")
		vc.Disconnect()
	} else {
		event.SendMessage("I am not in a voice channel.")
	}
}

func addToPlaylistRequest(event Event, query string) {
	youtubeID, err := searchService.GetVideoID(query)
	if err != nil {
		event.SendMessage(err.Error())
		return
	}

	event.SendMessage("Adding to playlist 😉")

	title, err := voice.AddSong(youtubeID, event.GetGuildID())
	if err != nil {
		event.SendMessage(err.Error())
		return
	}
	event.SendMessage(fmt.Sprintf("Your song **%s** was added 👍", title))

	vc, connected := event.GetVoiceConnection()
	if connected {
		thisVoiceChannel := voice.ActiveVoiceChannels[vc.GetGuildID()]

		if !thisVoiceChannel.IsPlayingMusic() && thisVoiceChannel.ExistsNext() {
			go voice.PlayMusic(vc.GetAudioInputChannel(), event.GetGuildID(), thisVoiceChannel, true)
			event.SendMessage("Playing next song in the playlist... 🎵")
		}
	}
}

func helpRequest(event Event, _ string) {
	event.SendMessage(help.Message())
}

func playMusicRequest(event Event, _ string) {
	vc, connected := event.GetVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel yet.")
	} else {
		thisVoiceChannel := voice.ActiveVoiceChannels[vc.GetGuildID()]
		if thisVoiceChannel.IsPlayingMusic() {
			event.SendMessage("I am already playing music 😁")
		} else {
			if !thisVoiceChannel.ExistsNext() && !thisVoiceChannel.ExistsBackupNext() {
				event.SendMessage("Playlist is currently empty.")
				return
			}
			voice.PlayMusic(vc.GetAudioInputChannel(), event.GetGuildID(), thisVoiceChannel, thisVoiceChannel.ExistsNext())
			event.SendMessage("Starting music... 🎵")
		}
	}
}

func stopMusicRequest(event Event, _ string) {
	vc, connected := event.GetVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel. 😔")
	} else {
		voiceChannel := voice.ActiveVoiceChannelForGuild(vc.GetGuildID())
		if voiceChannel.IsPlayingMusic() {
			voiceChannel.StopMusic()
			event.SendMessage("OK, Shutting up now...")
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}

func nowPlayingRequest(event Event, _ string) {
	vconn, connected := event.GetVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel. 😔")
	} else {
		vchan := voice.ActiveVoiceChannels[vconn.GetGuildID()]
		if vchan.IsPlayingMusic() {
			event.SendMessage(fmt.Sprintf("Now playing: **%s**", vchan.GetNowPlayingName()))
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}

func skipMusicRequest(event Event, _ string) {
	vc, connected := event.GetVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel yet. 😔")
	} else {
		thisVoiceChannel := voice.ActiveVoiceChannelForGuild(vc.GetGuildID())
		if thisVoiceChannel.IsPlayingMusic() {
			thisVoiceChannel.StopMusic()
			event.SendMessage("Skipping song... ❌")
			if thisVoiceChannel.ExistsNext() {
				go voice.PlayMusic(vc.GetAudioInputChannel(), event.GetGuildID(), thisVoiceChannel, true)
			} else if thisVoiceChannel.ExistsBackupNext() {
				go voice.PlayMusic(vc.GetAudioInputChannel(), event.GetGuildID(), thisVoiceChannel, false)
			}
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}

func upNextRequest(event Event, _ string) {
	thisVoiceChannel, exists := voice.ActiveVoiceChannels[event.GetGuildID()]
	if !exists {
		event.SendMessage("Playlist is currently empty.")
		return
	}

	nextSongs, hasSongs := thisVoiceChannel.GetNextSongs()

	if !hasSongs {
		event.SendMessage("Playlist is currently empty.")
		return
	}

	songMessages := []string{"Coming Up Next:"}

	for i, song := range nextSongs {
		songMessages = append(songMessages, fmt.Sprintf("%v. %s", i+1, song.Name))
	}

	event.SendMessage(strings.Join(songMessages, "\n"))
}
