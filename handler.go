package main

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/help"
	"github.com/jamestjw/lyrical/matcher"
	lyrical_poll "github.com/jamestjw/lyrical/poll"
	"github.com/jamestjw/lyrical/utils"
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
	defaultMux.RegisterHandler(matcher.VoteMatcher, newPollRequest)
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
	utils.LogInfo("request to join", utils.KvForHandler(event.GetGuildID(), "joinVoiceChannelRequest", utils.KVs("channelName", channelName)))

	event.SendMessage(fmt.Sprintf("Connecting to channel name: %s", channelName))
	channelID, err := event.FindVoiceChannel(channelName)

	if err != nil {
		event.SendMessage("Unable to find channel fo this name in the server.")
		return
	}

	if voice.AlreadyInVoiceChannel(event.GetSession(), event.GetGuildID()) {
		event.SendMessage(fmt.Sprintf("I am already in Voice Channel within Guild ID: %s", event.GetGuildID()))
	} else {
		vc := voice.JoinVoiceChannel(event.GetSession(), event.GetGuildID(), channelID)
		utils.LogInfo("joined channel", utils.KvForHandler(event.GetGuildID(), "joinVoiceChannelRequest", utils.KVs("channelID", channelID)))

		thisChannel := voice.ActiveVoiceChannels[event.GetGuildID()]

		if !thisChannel.ExistsNext() && !thisChannel.ExistsBackupNext() {
			event.SendMessage("Playlist is still empty.")
		} else {
			voice.PlayMusic(vc.GetAudioInputChannel(), event.GetGuildID(), thisChannel, thisChannel.ExistsNext())
		}
	}
}

func leaveVoiceChannelRequest(event Event, _ string) {
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "leaveVoiceChannelRequest", nil))

	vc, connected := event.GetVoiceConnection()

	if connected {
		voiceChannel := voice.ActiveVoiceChannelForGuild(event.GetGuildID())
		if voiceChannel.IsPlayingMusic() {
			voiceChannel.StopMusic()
		}
		vc.Disconnect()
		event.SendMessage("Left voice channel 👋🏼")
	} else {
		event.SendMessage("I am not in a voice channel.")
	}
}

func addToPlaylistRequest(event Event, query string) {
	utils.LogInfo(query, utils.KvForHandler(event.GetGuildID(), "addToPlaylistRequest", nil))

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
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "helpRequest", nil))

	event.SendMessage(help.Message())
}

func playMusicRequest(event Event, _ string) {
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "playMusicRequest", nil))

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
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "stopMusicRequest", nil))

	vc, connected := event.GetVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel. 😔")
	} else {
		voiceChannel := voice.ActiveVoiceChannelForGuild(vc.GetGuildID())
		if voiceChannel.IsPlayingMusic() {
			voiceChannel.StopMusic()
		} else {
			event.SendMessage("Well I am not playing any music currently 🤔")
		}
	}
}

func nowPlayingRequest(event Event, _ string) {
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "nowPlayingRequest", nil))

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
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "skipMusicRequest", nil))

	vc, connected := event.GetVoiceConnection()
	if !connected {
		event.SendMessage("Hey I dont remember being invited to a voice channel yet. 😔")
	} else {
		thisVoiceChannel := voice.ActiveVoiceChannelForGuild(vc.GetGuildID())
		if thisVoiceChannel.IsPlayingMusic() {
			thisVoiceChannel.StopMusic()
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
	utils.LogInfo("", utils.KvForHandler(event.GetGuildID(), "upNextRequest", nil))

	thisVoiceChannel := voice.ActiveVoiceChannelForGuild(event.GetGuildID())

	nextSongs, hasSongs := thisVoiceChannel.GetNextSongs()
	nextBackupSongs, hasBackupSongs := thisVoiceChannel.GetNextBackupSongs()

	if !hasSongs && !hasBackupSongs {
		event.SendMessage("Playlist is currently empty.")
		return
	}

	allSongs := utils.LimitSongsArrayLengths(nextSongs, nextBackupSongs, config.UpNextMaxSongsCount)
	message := utils.FormatNowPlayingText(allSongs, "Coming Up Next:")
	event.SendMessage(message)
}

func newPollRequest(event Event, pollParams string) {
	utils.LogInfo(pollParams, utils.KvForHandler(event.GetGuildID(), "newPollRequest", nil))

	p, err := lyrical_poll.FromStringParams(pollParams)

	if err != nil {
		event.SendMessage(err.Error())
		return
	}

	pollMessageContents, pollEmojis := p.GeneratePollMessage()
	sentMessage := event.SendMessage(pollMessageContents)
	for _, emoji := range pollEmojis {
		event.ReactToMessage(emoji, sentMessage.ID)
	}

	defer func() {
		time.Sleep(p.GetDuration())
		finalMsg, err := event.GetMessageByMessageID(sentMessage.ID)

		if err != nil {
			utils.LogInfo(err.Error(), utils.KvForHandler(event.GetGuildID(), "newPollRequest", nil))
			event.SendMessage("Unable to find the poll, was the message deleted? :eyes:")
			return
		}

		counts := utils.ExtractEmojiCounts(finalMsg.Reactions)
		p.AddResult(counts)

		event.SendQuotedMessage(pollMessageContents, p.GetVerdict())
	}()
}
