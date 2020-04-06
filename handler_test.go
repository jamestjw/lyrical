package main

import (
	"fmt"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	mock_main "github.com/jamestjw/lyrical/mocks/mock_main"
	mock_voice "github.com/jamestjw/lyrical/mocks/mock_voice"
	"github.com/jamestjw/lyrical/playlist"
	"github.com/jamestjw/lyrical/voice"
	"github.com/stretchr/testify/mock"
)

func TestHelpRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().SendMessage(gomock.Any())

	helpRequest(mockEvent, "")
}

func TestLeaveVoiceChannelRequestWhenConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().Disconnect().Times(1)

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().SendMessage("Leaving voice channel 👋🏼")

	leaveVoiceChannelRequest(mockEvent, "")
}

func TestLeaveVoiceChannelRequestWhenNotConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(nil, false)
	mockEvent.EXPECT().SendMessage("I am not in a voice channel.")

	leaveVoiceChannelRequest(mockEvent, "")
}

func TestNowPlayingRequestWhileNotConnected(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(nil, false)
	mockEvent.EXPECT().SendMessage("Hey I dont remember being invited to a voice channel. 😔")

	nowPlayingRequest(mockEvent, "")
}

func TestNowPlayingRequestWhileConnectedAndPlayingMusic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().IsPlayingMusic().Return(true)
	mockChannel.EXPECT().GetNowPlayingName().Return("current song name")
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().GetGuildID().Return("guildID")

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().SendMessage("Now playing: **current song name**")

	nowPlayingRequest(mockEvent, "")
}

func TestNowPlayingRequestWhileConnectedAndNotPlayingMusic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().IsPlayingMusic().Return(false)
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().GetGuildID().Return("guildID")

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().SendMessage("Well I am not playing any music currently 🤔")

	nowPlayingRequest(mockEvent, "")
}

func TestSkipMusicRequestWhileConnectedAndPlayingMusic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().IsPlayingMusic().Return(true)
	mockChannel.EXPECT().StopMusic()
	mockChannel.EXPECT().ExistsNext().Return(false)
	mockChannel.EXPECT().ExistsBackupNext().Return(false)
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().GetGuildID().Return("guildID")

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().SendMessage("Skipping song... ❌")

	skipMusicRequest(mockEvent, "")
}

func TestStopMusicRequestWhileConnectedAndPlayingMusic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().IsPlayingMusic().Return(true)
	mockChannel.EXPECT().StopMusic()
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().GetGuildID().Return("guildID")

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().SendMessage("OK, Shutting up now...")

	stopMusicRequest(mockEvent, "")
}

func TestPlayMusicRequestWhileConnectedAndPlayingMusic(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	audiochan := make(chan []byte)

	// Waiting is required because we call a goroutine.
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().IsPlayingMusic().Return(false)
	mockChannel.EXPECT().ExistsNext().Return(true)
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockPlayer := mock_voice.NewMockMusicPlayer(ctrl)
	mockPlayer.EXPECT().PlayMusic(audiochan, "guildID", mockChannel, true).Do(func(chan []byte, string, voice.Channel, bool) { wg.Done() })

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().GetGuildID().Return("guildID")
	mockConnection.EXPECT().GetAudioInputChannel().Return(audiochan)

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().GetGuildID().Return("guildID")
	mockEvent.EXPECT().SendMessage("Starting music... 🎵")

	voice.DefaultMusicPlayer = mockPlayer

	playMusicRequest(mockEvent, "")
}

// AddToPlaylist
type mockSearchService struct {
	mock.Mock
}

func (s *mockSearchService) GetVideoID(id string) (string, error) {
	args := s.Called(id)
	return args.String(0), args.Error(1)
}

type mockMusicDownloader struct {
	mock.Mock
}

func (m *mockMusicDownloader) Download(query string) (title string, err error) {
	args := m.Called(query)
	return args.String(0), args.Error(1)
}

func TestAddToPlaylist(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSS := new(mockSearchService)
	mockSS.On("GetVideoID", "song name").Return("video id", nil)

	mockDB := mock_voice.NewMockDatabase(ctrl)
	mockDB.EXPECT().AddSongToDB("song name", "video id").Return(nil)
	mockDB.EXPECT().SongExists("video id").Return("", false)

	mockDl := new(mockMusicDownloader)
	mockDl.On("Download", "video id").Return("song name", nil)

	// Waiting is required because we call a goroutine.
	var wg sync.WaitGroup
	wg.Add(1)
	defer wg.Wait()

	audiochan := make(chan []byte)

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().FetchPlaylist().Return(&playlist.Playlist{})
	gomock.InOrder(
		mockChannel.EXPECT().ExistsNext().Return(false),
		mockChannel.EXPECT().ExistsNext().Return(true),
	)
	mockChannel.EXPECT().SetNext(gomock.Any())
	mockChannel.EXPECT().IsPlayingMusic().Return(false)

	mockPlayer := mock_voice.NewMockMusicPlayer(ctrl)
	mockPlayer.EXPECT().PlayMusic(audiochan, "guildID", mockChannel, true).Do(func(chan []byte, string, voice.Channel, bool) { wg.Done() })

	mockConnection := mock_voice.NewMockConnection(ctrl)
	mockConnection.EXPECT().GetGuildID().Return("guildID")
	mockConnection.EXPECT().GetAudioInputChannel().Return(audiochan)

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetVoiceConnection().Return(mockConnection, true)
	mockEvent.EXPECT().GetGuildID().AnyTimes().Return("guildID")
	gomock.InOrder(
		mockEvent.EXPECT().SendMessage("Adding to playlist 😉"),
		mockEvent.EXPECT().SendMessage("Your song **song name** was added 👍"),
		mockEvent.EXPECT().SendMessage("Playing next song in the playlist... 🎵"),
	)

	searchService = mockSS
	voice.DB = mockDB
	voice.Dl = mockDl
	voice.DefaultMusicPlayer = mockPlayer
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	addToPlaylistRequest(mockEvent, "song name")
}

func TestJoinChannelRequest(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	connMap := make(map[string]voice.Connection)

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().ExistsNext().Return(false)
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockConnection := mock_voice.NewMockConnection(ctrl)

	mockConnectable := mock_voice.NewMockConnectable(ctrl)
	mockConnectable.EXPECT().GetVoiceConnections().Return(connMap)
	mockConnectable.EXPECT().JoinVoiceChannel("guildID", "channel-id").Return(mockConnection, nil)

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().FindVoiceChannel("channel-name").Return("channel-id", nil)
	mockEvent.EXPECT().GetSession().AnyTimes().Return(mockConnectable)
	mockEvent.EXPECT().GetGuildID().AnyTimes().Return("guildID")
	gomock.InOrder(
		mockEvent.EXPECT().SendMessage("Connecting to channel name: channel-name"),
		mockEvent.EXPECT().SendMessage("Joining Voice Channel: Guild ID: guildID ChannelID: channel-id"),
		mockEvent.EXPECT().SendMessage("Playlist is still empty."),
	)
	joinVoiceChannelRequest(mockEvent, "channel-name")
}

func TestUpNextRequest(t *testing.T) {
	var songs []*playlist.Song

	for i := 0; i < 2; i++ {
		songName := fmt.Sprintf("Song %v", i)
		song := &playlist.Song{
			YoutubeID: "",
			Name:      songName,
			Next:      nil,
		}
		songs = append(songs, song)
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockChannel := mock_voice.NewMockChannel(ctrl)
	mockChannel.EXPECT().GetNextSongs().Return(songs, true)
	voice.ActiveVoiceChannels["guildID"] = mockChannel

	mockEvent := mock_main.NewMockEvent(ctrl)
	mockEvent.EXPECT().GetGuildID().AnyTimes().Return("guildID")
	gomock.InOrder(
		mockEvent.EXPECT().SendMessage("Coming Up Next:\n1. Song 0\n2. Song 1"),
	)
	upNextRequest(mockEvent, "")
}
