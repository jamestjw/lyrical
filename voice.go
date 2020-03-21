package main

import (
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

var activeVoiceChannels *voiceChannels

func init() {
	activeVoiceChannels = newActiveVoiceChannels()
}

type voiceChannels struct {
	channelMap map[*discordgo.VoiceConnection]*voiceChannel
}

type voiceChannel struct {
	AbortChannel chan string
	MusicActive  bool
}

func newActiveVoiceChannels() *voiceChannels {
	var vcs voiceChannels
	vcs.channelMap = make(map[*discordgo.VoiceConnection]*voiceChannel)
	return &vcs
}

func playMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!play-music" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel yet.")
		} else {
			if activeVoiceChannels.channelMap[vc].MusicActive {
				s.ChannelMessageSend(m.ChannelID, "I am already playing music 😁")
			} else {
				go playMusic(vc)
				s.ChannelMessageSend(m.ChannelID, "Starting music... 👍")
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
			if activeVoiceChannels.channelMap[vc].MusicActive {
				activeVoiceChannels.channelMap[vc].AbortChannel <- "stop"
				s.ChannelMessageSend(m.ChannelID, "OK, Shutting up now...")
			} else {
				s.ChannelMessageSend(m.ChannelID, "Well I am not playing any music currently 🤔")
			}
		}
	}
}

func alreadyInVoiceChannel(s *discordgo.Session, guildID string) bool {
	_, connected := s.VoiceConnections[guildID]
	return connected
}

func disconnectAllVoiceConnections(s *discordgo.Session) error {
	for _, channel := range s.VoiceConnections {
		err := channel.Disconnect()
		if err != nil {
			return err
		}
		log.Println("Disconnected from voice channel...")
	}
	return nil
}

func playMusic(vc *discordgo.VoiceConnection) {
	encodeSession, err := dca.EncodeFile("chopin.mp3", dca.StdEncodeOptions)
	if err != nil {
		log.Fatal(err)
	}
	defer encodeSession.Cleanup()

	decoder := dca.NewDecoder(encodeSession)

	activeVoiceChannels.channelMap[vc].MusicActive = true
	defer func() {
		activeVoiceChannels.channelMap[vc].MusicActive = false
	}()

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Fatal(err)
			}
			break
		}

		// Do something with the frame, in this example were sending it to discord
		select {
		case vc.OpusSend <- frame:
		case <-activeVoiceChannels.channelMap[vc].AbortChannel:
			return
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			log.Println("TIMEOUT: Unable to send audio..")
			return
		}
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID string, voiceChannelID string) *discordgo.VoiceConnection {
	vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, false)
	if err != nil {
		log.Fatal(err)
	}
	vcd := &voiceChannel{AbortChannel: make(chan string, 1)}
	activeVoiceChannels.channelMap[vc] = vcd
	return vc
}
