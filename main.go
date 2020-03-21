package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

const guildID = "690751445607251988"
const voiceChannelID = "690751445607251992"

type voiceChannels struct {
	channelMap map[*discordgo.VoiceConnection](chan string)
}

func newActiveVoiceChannels() *voiceChannels {
	var vcs voiceChannels
	vcs.channelMap = make(map[*discordgo.VoiceConnection](chan string), 1)
	return &vcs
}

var activeVoiceChannels *voiceChannels

func init() {
	activeVoiceChannels = newActiveVoiceChannels()
}

func main() {
	dg, err := discordgo.New("Bot " + config.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	defer dg.Close()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)
	dg.AddHandler(joinChannelRequest)
	dg.AddHandler(stopMusicRequest)
	dg.AddHandler(playMusicRequest)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
		log.Println("Received signal to terminate, cleaning up...")
		// Cleanly close down the Discord session.
		disconnectAllVoiceConnections(dg)
		log.Println("Exit successful!")
		return
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	// If the message is "ping" reply with "Pong!"
	if m.Content == "!ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "!pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

func joinChannelRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "!join-voice" {
		if alreadyInVoiceChannel(s, m.GuildID) {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I am already in Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, voiceChannelID))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Joining Voice Channel: Guild ID: %s ChannelID: %v \n", m.GuildID, voiceChannelID))
			log.Printf("Joining Guild ID: %s ChannelID: %v \n", m.GuildID, voiceChannelID)
			vc := joinVoiceChannel(s, m.GuildID, voiceChannelID)
			playMusic(vc)
		}
	}
}

func joinVoiceChannel(s *discordgo.Session, guildID string, voiceChannelID string) *discordgo.VoiceConnection {
	vc, err := s.ChannelVoiceJoin(guildID, voiceChannelID, false, false)
	if err != nil {
		log.Fatal(err)
	}
	activeVoiceChannels.channelMap[vc] = make(chan string)
	return vc
}

func playMusicRequest(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "!play-music" {
		vc, connected := s.VoiceConnections[guildID]
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
		vc, connected := s.VoiceConnections[guildID]
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
