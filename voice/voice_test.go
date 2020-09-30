package voice_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/mocks/mock_database"
	mock_voice "github.com/jamestjw/lyrical/mocks/mock_voice"
	"github.com/jamestjw/lyrical/voice"
	"github.com/stretchr/testify/assert"
)

func disconnectAllVoiceConnectionsSetup(ctrl *gomock.Controller) map[string]voice.Connection {
	voiceConnectionMap := make(map[string]voice.Connection)

	for i := 0; i < 5; i++ {
		guildID := fmt.Sprintf("Guild %v", i)
		voiceConnection := mock_voice.NewMockConnection(ctrl)
		voiceConnection.EXPECT().Disconnect().Times(1).Return(nil)
		voiceConnection.EXPECT().GetGuildID().Times(1).Return(guildID)
		voiceConnectionMap[string(i)] = voiceConnection

		channel := mock_voice.NewMockChannel(ctrl)
		channel.EXPECT().RemoveNowPlaying().Times(1)
		voice.ActiveVoiceChannels[guildID] = channel
	}
	return voiceConnectionMap
}

func TestDisconnectAllVoiceConnections(t *testing.T) {
	cleanActiveVoiceChannels()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	voiceConnectionMap := disconnectAllVoiceConnectionsSetup(ctrl)
	mockSession := mock_voice.NewMockConnectable(ctrl)
	mockSession.EXPECT().GetVoiceConnections().Times(1).Return(voiceConnectionMap)

	voice.DisconnectAllVoiceConnections(mockSession)
}

func TestJoinVoiceChannel(t *testing.T) {
	cleanActiveVoiceChannels()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConnection := mock_voice.NewMockConnection(ctrl)

	mockSession := mock_voice.NewMockConnectable(ctrl)
	mockSession.EXPECT().JoinVoiceChannel("guildID", "channelID").Times(1).Return(mockConnection, nil)

	mockDatastore := mock_database.NewMockDatastore(ctrl)
	mockDatastore.EXPECT().GetRandomSongs(gomock.Any()).Return([]database.Song{{Name: "New Song", YoutubeID: "youtubeID"}})

	database.DS = mockDatastore

	voice.JoinVoiceChannel(mockSession, "guildID", "channelID")

	assert.NotNil(t, voice.ActiveVoiceChannels["guildID"], "channel should be initiated")
}

func TestAddSongThatAlreadyExists(t *testing.T) {
	cleanActiveVoiceChannels()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDatastore := mock_database.NewMockDatastore(ctrl)
	mockDatastore.EXPECT().SongExists("youtubeID").Times(1).Return("Song Name", true)
	mockDatastore.EXPECT().GetRandomSongs(gomock.Any()).Return([]database.Song{{Name: "New Song", YoutubeID: "youtubeID"}})

	database.DS = mockDatastore
	voice.AddSong("youtubeID", "guildID")

	if assert.Equal(t, voice.ActiveVoiceChannels["guildID"].FetchPlaylist().IsEmpty(), false) {
		assert.Equal(t, voice.ActiveVoiceChannels["guildID"].FetchPlaylist().First().YoutubeID, "youtubeID")
	}

	assert.Equal(t, voice.ActiveVoiceChannels["guildID"].GetNext().YoutubeID, "youtubeID")
}

func TestAddSongThatDoesNotExistYet(t *testing.T) {
	cleanActiveVoiceChannels()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockDownloader := mock_voice.NewMockDownloader(ctrl)
	mockDownloader.EXPECT().Download("youtubeID").Times(1).Return("New Song", nil)

	mockDatastore := mock_database.NewMockDatastore(ctrl)
	mockDatastore.EXPECT().SongExists("youtubeID").Times(1).Return("", false)
	mockDatastore.EXPECT().AddSongToDB("New Song", "youtubeID").Times(1).Return(nil)
	mockDatastore.EXPECT().GetRandomSongs(gomock.Any()).Return([]database.Song{{Name: "New Song", YoutubeID: "youtubeID"}})

	database.DS = mockDatastore
	voice.Dl = mockDownloader

	voice.AddSong("youtubeID", "guildID")

	if assert.Equal(t, voice.ActiveVoiceChannels["guildID"].FetchPlaylist().IsEmpty(), false) {
		assert.Equal(t, voice.ActiveVoiceChannels["guildID"].FetchPlaylist().First().YoutubeID, "youtubeID")
	}

	assert.Equal(t, voice.ActiveVoiceChannels["guildID"].GetNext().YoutubeID, "youtubeID")
}

func cleanActiveVoiceChannels() {
	voice.ActiveVoiceChannels = voice.NewActiveVoiceChannels()
}
