package voice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/ytmp3"
	"github.com/jonas747/dca"
)

var ActiveVoiceChannels *voiceChannels

func init() {
	ActiveVoiceChannels = NewActiveVoiceChannels()
}

func AlreadyInVoiceChannel(s *discordgo.Session, guildID string) bool {
	_, connected := s.VoiceConnections[guildID]
	return connected
}

func DisconnectAllVoiceConnections(s Connectable) error {
	for _, channel := range s.GetVoiceConnections() {
		err := channel.Disconnect()
		if err != nil {
			return err
		}
		log.Println("Disconnected from voice channel...")
		ActiveVoiceChannels.ChannelMap[channel.GetGuildID()].RemoveNowPlaying()
	}
	return nil
}

func MaybeSetNext(guildID string, s *playlist.Song) {
	if _, exists := ActiveVoiceChannels.ChannelMap[guildID]; !exists {
		InitialiseVoiceChannelForGuild(guildID)
	}

	vc := ActiveVoiceChannels.ChannelMap[guildID]
	if vc.GetNext() == nil {
		vc.SetNext(s)
	}
}

func (vc *voiceChannel) SetNowPlaying(s *playlist.Song) {
	vc.MusicActive = true
	vc.NowPlaying = s
	vc.Next = s.Next
}

func (vc *voiceChannel) RemoveNowPlaying() {
	vc.MusicActive = false
	vc.NowPlaying = nil
}

func PlayMusic(input chan []byte, guildID string, song *playlist.Song) error {
	encodeSession, err := dca.EncodeFile(ytmp3.PathToAudio(song.YoutubeID), dca.StdEncodeOptions)
	if err != nil {
		log.Print(err)
		return errors.New("unable to open this song")
	}
	defer encodeSession.Cleanup()

	decoder := dca.NewDecoder(encodeSession)

	ActiveVoiceChannels.ChannelMap[guildID].SetNowPlaying(song)
	defer ActiveVoiceChannels.ChannelMap[guildID].RemoveNowPlaying()

	abortChannel := ActiveVoiceChannels.ChannelMap[guildID].GetAbortChannel()

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
		case input <- frame:
		case <-abortChannel:
			return nil
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			log.Println("TIMEOUT: Unable to send audio..")
			return nil
		}
	}

	// Being able to get here means that audio clip has ended
	if song.Next != nil {
		go PlayMusic(input, guildID, song.Next)
	}
	return nil
}

func JoinVoiceChannel(s Connectable, guildID string, voiceChannelID string) Connection {
	vc, err := s.JoinVoiceChannel(guildID, voiceChannelID)
	if err != nil {
		log.Fatal(err)
	}

	if _, exists := ActiveVoiceChannels.ChannelMap[guildID]; !exists {
		InitialiseVoiceChannelForGuild(guildID)
	}
	return vc
}

func InitialiseVoiceChannelForGuild(guildID string) {
	vcd := &voiceChannel{AbortChannel: make(chan string, 1)}
	ActiveVoiceChannels.ChannelMap[guildID] = vcd
}

func DownloadByYoutubeID(youtubeID string) (title string, err error) {
	title, err = ytmp3.Download(youtubeID)
	return
}

func AddSong(youtubeID string, guildID string) (title string, err error) {
	title, exists := database.SongExists(youtubeID)

	if !exists {
		title, err = DownloadByYoutubeID(youtubeID)
		if err != nil {
			err = fmt.Errorf("Error adding the song %s 🤨: %s", youtubeID, err.Error())
		}

		dbErr := database.AddSongToDB(title, youtubeID)
		if dbErr != nil {
			log.Printf("Error writing song ID %s to the database: %s", youtubeID, dbErr)
		}
	}

	newSong := playlist.LyricalPlaylist.AddSongWithYoutubeID(title, youtubeID)
	MaybeSetNext(guildID, newSong)
	return
}
