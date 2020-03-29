package voice

import (
	"errors"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/ytmp3"
	"github.com/jonas747/dca"
)

// ActiveVoiceChannels is a global struct containing information
// on the active voice channels belonging to each guild.
var ActiveVoiceChannels *voiceChannels

func init() {
	ActiveVoiceChannels = NewActiveVoiceChannels()
}

// AlreadyInVoiceChannel checks if a Connectable object is currently
// connected to a voice channel that belongs to a particular guild.
func AlreadyInVoiceChannel(s Connectable, guildID string) bool {
	vcs := s.GetVoiceConnections()
	_, connected := vcs[guildID]
	return connected
}

// DisconnectAllVoiceConnections will disconnect all voice channels belonging
// to a Connectable object. It will also remove the actively playing music
// status and the NowPlaying song of each voice channel in the global
// ActiveVoiceChannels object.
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

func maybeSetNext(guildID string, s *playlist.Song) {
	if _, exists := ActiveVoiceChannels.ChannelMap[guildID]; !exists {
		initialiseVoiceChannelForGuild(guildID)
	}

	vc := ActiveVoiceChannels.ChannelMap[guildID]
	if vc.GetNext() == nil {
		vc.SetNext(s)
	}
}

// PlayMusic plays a given song into an input Audio Channel that belongs to a guild
// with guildID. The given song will be set as the currently playing song of the guild and
// the voice channel of the guild will be marked as active..It will also automatically play
// the next song if there is one.
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

// JoinVoiceChannel invokes the JoinVoiceChannel method of a Connectable object.
// It also initialises an entry in the ChannelMap of the global ActiveVoiceChannels
// struct.
func JoinVoiceChannel(s Connectable, guildID string, voiceChannelID string) Connection {
	vc, err := s.JoinVoiceChannel(guildID, voiceChannelID)
	if err != nil {
		log.Fatal(err)
	}

	if _, exists := ActiveVoiceChannels.ChannelMap[guildID]; !exists {
		initialiseVoiceChannelForGuild(guildID)
	}
	return vc
}

func initialiseVoiceChannelForGuild(guildID string) {
	vcd := &voiceChannel{AbortChannel: make(chan string, 1)}
	ActiveVoiceChannels.ChannelMap[guildID] = vcd
}

// AddSong will download a song based on the youtubeID for the guild
// with guildID if it has not already been downloaded. It will
// also add a database entry of it and add it to the playlist
// of the guild.
func AddSong(youtubeID string, guildID string) (title string, err error) {
	title, exists := DB.SongExists(youtubeID)

	if !exists {
		title, err = Dl.Download(youtubeID)
		if err != nil {
			err = fmt.Errorf("Error adding the song %s 🤨: %s", youtubeID, err.Error())
		}

		dbErr := DB.AddSongToDB(title, youtubeID)
		if dbErr != nil {
			log.Printf("Error writing song ID %s to the database: %s", youtubeID, dbErr)
		}
	}

	newSong := playlist.LyricalPlaylist.AddSong(title, youtubeID)
	maybeSetNext(guildID, newSong)
	return
}
