package main

import (
	"errors"
	"fmt"
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
	NowPlaying   *Song
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
		activeVoiceChannels.channelMap[channel].RemoveNowPlaying()
	}
	return nil
}

func (vc *voiceChannel) SetNowPlaying(s *Song) {
	vc.MusicActive = true
	vc.NowPlaying = s
}

func (vc *voiceChannel) RemoveNowPlaying() {
	vc.MusicActive = false
	vc.NowPlaying = nil
}

func playMusic(vc *discordgo.VoiceConnection, song *Song) error {
	encodeSession, err := dca.EncodeFile(ytmp3.PathToAudio(song.YoutubeID), dca.StdEncodeOptions)
	if err != nil {
		log.Print(err)
		return errors.New("unable to open this song")
	}
	defer encodeSession.Cleanup()

	decoder := dca.NewDecoder(encodeSession)

	activeVoiceChannels.channelMap[vc].SetNowPlaying(song)
	defer activeVoiceChannels.channelMap[vc].RemoveNowPlaying()

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Print(err)
				return errors.New("unable to decode this song")
			}
			break
		}

		// Do something with the frame, in this example were sending it to discord
		select {
		case vc.OpusSend <- frame:
		case <-activeVoiceChannels.channelMap[vc].AbortChannel:
			return nil
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			log.Println("TIMEOUT: Unable to send audio..")
			return nil
		}
	}

	// Being able to get here means that audio clip has ended
	if song.Next != nil {
		go playMusic(vc, song.Next)
	}
	return nil
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

func downloadByYoutubeID(youtubeID string) (title string, err error) {
	title, err = ytmp3.Download(youtubeID)
	return
}

func (vc *voiceChannel) GetNowPlayingName() string {
	return vc.NowPlaying.Name
}

func addSong(youtubeID string) (title string, err error) {
	title, err = downloadByYoutubeID(youtubeID)
	if err != nil {
		err = fmt.Errorf("Error adding the song %s 🤨: %s", youtubeID, err.Error())
	}

	newSong := lyricalPlaylist.AddSongWithYoutubeID(title, youtubeID)

	dbErr := AddSongToDB(newSong)
	if dbErr != nil {
		log.Printf("Error writing song ID %s to the database: %s", youtubeID, dbErr)
	}
	return
}
