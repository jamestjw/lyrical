package main

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/matcher"
	"github.com/jamestjw/lyrical/searcher"
	"github.com/jamestjw/lyrical/voice"
)

var defaultMux = &Mux{}

func init() {
	searcher.InitialiseSearchService(config.YoutubeAPIKey)

	defaultMux = NewMux()
	defaultMux.RegisterHandler(matcher.JoinChannelRequestMatcher, joinVoiceChannelRequest)
	defaultMux.RegisterHandler(matcher.LeaveVoiceMatcher, leaveVoiceChannelRequest)
	defaultMux.RegisterHandler(matcher.AddPlaylistRequestMatcher, addToPlaylistRequest)
	defaultMux.RegisterHandler(matcher.PlayMusicMatcher, playMusicRequest)
	defaultMux.RegisterHandler(matcher.StopMusicMatcher, stopMusicRequest)
	defaultMux.RegisterHandler(matcher.SkipMusicMatcher, skipMusicRequest)
	defaultMux.RegisterHandler(matcher.NowPlayingMatcher, nowPlayingRequest)
	defaultMux.RegisterHandler(matcher.HelpMatcher, helpRequest)
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

	if voice.AlreadyInVoiceChannel(event.getSession(), event.getGuildID()) {
		event.SendMessage(fmt.Sprintf("I am already in Voice Channel within Guild ID: %s", event.getGuildID()))
	} else {
		event.SendMessage(fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v \n", event.getGuildID(), channelID))
		log.Printf("Joining Guild ID: %s ChannelID: %v \n", event.getGuildID(), channelID)

		vc := voice.JoinVoiceChannel(event.getSession(), event.getGuildID(), channelID)
		nextSong := voice.ActiveVoiceChannels.ChannelMap[event.getGuildID()].GetNext()
		if nextSong == nil {
			event.SendMessage("Playlist is still empty.")
		} else {
			go voice.PlayMusic(vc.GetAudioInputChannel(), event.getGuildID(), nextSong)
			event.SendMessage("Starting music... 🎵")
		}
	}
}

func leaveVoiceChannelRequest(event Event, _ string) {
	// TODO: Leave voice channel of current guild only.
	event.SendMessage("Leaving voice channel 👋🏼")
	err := voice.DisconnectAllVoiceConnections(event.getSession())

	if err != nil {
		event.SendMessage(err.Error())
	}
}

func addToPlaylistRequest(event Event, query string) {
	youtubeID, err := searcher.GetVideoID(query)
	if err != nil {
		event.SendMessage(err.Error())
		return
	}

	event.SendMessage("Adding to playlist 😉")

	title, err := voice.AddSong(youtubeID, event.getGuildID())
	if err != nil {
		event.SendMessage(err.Error())
		return
	}
	event.SendMessage(fmt.Sprintf("Your song **%s** was added 👍", title))

	vc, connected := event.getVoiceConnection()
	if connected {
		thisVoiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
		if !thisVoiceChannel.IsPlayingMusic() && thisVoiceChannel.GetNext() != nil {
			go voice.PlayMusic(vc.OpusSend, event.getGuildID(), thisVoiceChannel.GetNext())
			event.SendMessage("Playing next song in the playlist... 🎵")
		}
	}
}

func helpRequest(event Event, _ string) {
	event.SendMessage(help.Message())
}

func playMusicRequest(event Event, _ string) {
	vc, connected := event.getVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel yet.")
	} else {
		thisVoiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
		if thisVoiceChannel.IsPlayingMusic() {
			event.SendMessage("I am already playing music 😁")
		} else {
			if thisVoiceChannel.GetNext() == nil {
				event.SendMessage("Playlist is currently empty.")
			} else {
				go voice.PlayMusic(vc.OpusSend, event.getGuildID(), thisVoiceChannel.GetNext())
				event.SendMessage("Starting music... 🎵")
			}
		}
	}
}

func stopMusicRequest(event Event, _ string) {
	vc, connected := event.getVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel. 😔")
	} else {
		voiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
		if voiceChannel.IsPlayingMusic() {
			voiceChannel.StopMusic()
			event.SendMessage("OK, Shutting up now...")
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}

func nowPlayingRequest(event Event, _ string) {
	vc, connected := event.getVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel. 😔")
	} else {
		if voice.ActiveVoiceChannels.ChannelMap[vc.GuildID].IsPlayingMusic() {
			event.SendMessage(fmt.Sprintf("Now playing: **%s**", voice.ActiveVoiceChannels.ChannelMap[vc.GuildID].GetNowPlayingName()))
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}

func skipMusicRequest(event Event, _ string) {
	vc, connected := event.getVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel yet. 😔")
	} else {
		thisVoiceChannel := voice.ActiveVoiceChannels.ChannelMap[vc.GuildID]
		if thisVoiceChannel.IsPlayingMusic() {
			thisVoiceChannel.StopMusic()
			event.SendMessage("Skipping song... ❌")
			if thisVoiceChannel.GetNext() != nil {
				go voice.PlayMusic(vc.OpusSend, event.getGuildID(), thisVoiceChannel.GetNext())
			}
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}
