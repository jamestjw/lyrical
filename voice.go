package main

import (
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/ytmp3"
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

func addToPlaylist(youtubeID string) {
	ytmp3.Download(youtubeID)
}
