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
	channelMap map[*discordgo.VoiceConnection](chan string)
}

func newActiveVoiceChannels() *voiceChannels {
	var vcs voiceChannels
	vcs.channelMap = make(map[*discordgo.VoiceConnection](chan string), 1)
	return &vcs
}

func playMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!play-music" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel yet.")
		} else {
			go playMusic(vc)
			s.ChannelMessageSend(m.ChannelID, "Starting music... 👍")
		}
	}
}

func stopMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!stop-music" {
		vc, connected := s.VoiceConnections[m.GuildID]
		if !connected {
			s.ChannelMessageSend(m.ChannelID, "Hey I dont remember being invited to a voice channel.")
		} else {
			activeVoiceChannels.channelMap[vc] <- "stop"
			s.ChannelMessageSend(m.ChannelID, "OK, Shutting up now...")
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
		case <-activeVoiceChannels.channelMap[vc]:
			return
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			log.Println("Unable to send audio..")
			return
		}
	}
}
