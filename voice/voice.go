package voice

import (
	"fmt"
	"io"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/jamestjw/lyrical/models"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/ytmp3"
	"github.com/jonas747/dca"
)

// ActiveVoiceChannels is a global struct containing information
// on the active voice channels belonging to each guild.
var ActiveVoiceChannels map[string]Channel

// DefaultMusicPlayer handles playing music to a given channel.
var DefaultMusicPlayer MusicPlayer

func init() {
	ActiveVoiceChannels = NewActiveVoiceChannels()
	DefaultMusicPlayer = &defaultMusicPlayer{}
}

type defaultMusicPlayer struct{}

// ActiveVoiceChannelForGuild will create or fetch the channel for this guild
func ActiveVoiceChannelForGuild(guildID string) Channel {
	if c, exists := ActiveVoiceChannels[guildID]; exists {
		return c
	}

	initialiseVoiceChannelForGuildIfNotExists(guildID)
	return ActiveVoiceChannels[guildID]
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
		ActiveVoiceChannels[channel.GetGuildID()].RemoveNowPlaying()
	}
	log.Println("Disconnected from voice channel...")
	return nil
}

func maybeSetNext(guildID string, s *playlist.Song) {
	vc := ActiveVoiceChannelForGuild(guildID)
	if !vc.ExistsNext() {
		vc.SetNext(s)
	}
}

// PlayMusic plays a given song into an input Audio Channel that belongs to a guild
// with guildID. The given song will be set as the currently playing song of the guild and
// the voice channel of the guild will be marked as active. It will also automatically play
// the next song if there is one.
func (d *defaultMusicPlayer) PlayMusic(input chan []byte, guildID string, vc Channel, mainPlaylist bool) {
	if !vc.ExistsNext() && mainPlaylist {
		panic("Song does not exist in playlist but PlayMusic was called.")
	} else if !vc.ExistsBackupNext() && !mainPlaylist {
		panic("Song does not exist in backup playlist but PlayMusic was called.")
	}

	var song *playlist.Song
	var aborted bool

	if mainPlaylist {
		song = vc.GetNext()
	} else {
		song = vc.GetBackupNext()
	}

	defer func() {
		if aborted {
			return
		}

		vc.RemoveNowPlaying()

		if vc.ExistsNext() {
			go d.PlayMusic(input, guildID, vc, true)
		} else if vc.ExistsBackupNext() {
			go d.PlayMusic(input, guildID, vc, false)
		}
	}()

	encodeSession, err := dca.EncodeFile(ytmp3.PathToAudio(song.YoutubeID), dca.StdEncodeOptions)
	if err != nil {
		log.Print(err)
		return
	}
	defer encodeSession.Cleanup()

	decoder := dca.NewDecoder(encodeSession)
	abortChannel := ActiveVoiceChannels[guildID].GetAbortChannel()

	for {
		frame, err := decoder.OpusFrame()
		if err != nil {
			if err != io.EOF {
				log.Print(err)
				return
			}
			break
		}

		// Do something with the frame, in this example were sending it to discord
		select {
		case input <- frame:
		case <-abortChannel:
			aborted = true
			return
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			log.Println("TIMEOUT: Unable to send audio..")
			return
		}
	}

	return
}

// JoinVoiceChannel invokes the JoinVoiceChannel method of a Connectable object.
// It also initialises an entry in the ChannelMap of the global ActiveVoiceChannels
// struct.
func JoinVoiceChannel(s Connectable, guildID string, voiceChannelID string) Connection {
	vc, err := s.JoinVoiceChannel(guildID, voiceChannelID)
	if err != nil {
		log.Fatal(err)
	}

	initialiseVoiceChannelForGuildIfNotExists(guildID)

	return vc
}

func initialiseVoiceChannelForGuildIfNotExists(guildID string) {
	if _, exists := ActiveVoiceChannels[guildID]; exists {
		return
	}

	backupPlaylist := &playlist.Playlist{}
	LoadPlaylist(backupPlaylist)

	vc := &voiceChannel{
		AbortChannel:   make(chan string, 1),
		Playlist:       &playlist.Playlist{},
		BackupPlaylist: backupPlaylist,
	}

	ActiveVoiceChannels[guildID] = vc
}

// AddSong will download a song based on the youtubeID for the guild
// with guildID if it has not already been downloaded. It will
// also add a models entry of it and add it to the playlist
// of the guild.
func AddSong(youtubeID string, guildID string) (title string, err error) {
	title, exists := models.DS.SongExists(youtubeID)

	if !exists {
		title, err = Dl.Download(youtubeID)
		if err != nil {
			err = fmt.Errorf("Error adding the song %s 🤨: %s", youtubeID, err.Error())
			return
		}

		dbErr := models.DS.AddSongToDB(title, youtubeID)
		if dbErr != nil {
			log.Printf("Error writing song ID %s to the models: %s", youtubeID, dbErr)
		}
	}

	initialiseVoiceChannelForGuildIfNotExists(guildID)
	guildPlaylist := ActiveVoiceChannels[guildID].FetchPlaylist()
	newSong := guildPlaylist.AddSong(title, youtubeID)
	maybeSetNext(guildID, newSong)
	return
}

// PlayMusic plays the next song in a given Channel
func PlayMusic(input chan []byte, guildID string, vc Channel, main bool) {
	go DefaultMusicPlayer.PlayMusic(input, guildID, vc, main)
}
