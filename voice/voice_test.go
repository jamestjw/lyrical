package voice_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
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
		voice.ActiveVoiceChannels.ChannelMap[guildID] = channel
	}
	return voiceConnectionMap
}

func TestDisconnectAllVoiceConnections(t *testing.T) {
	cleanActiveVoiceChannels()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := mock_voice.NewMockConnectable(ctrl)
	voiceConnectionMap := disconnectAllVoiceConnectionsSetup(ctrl)
	mockSession.EXPECT().GetVoiceConnections().Times(1).Return(voiceConnectionMap)

	voice.DisconnectAllVoiceConnections(mockSession)
}

func TestJoinVoiceChannel(t *testing.T) {
	cleanActiveVoiceChannels()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := mock_voice.NewMockConnectable(ctrl)
	mockConnection := mock_voice.NewMockConnection(ctrl)

	mockSession.EXPECT().JoinVoiceChannel("guildID", "channelID").Times(1).Return(mockConnection, nil)

	voice.JoinVoiceChannel(mockSession, "guildID", "channelID")

	assert.NotNil(t, voice.ActiveVoiceChannels.ChannelMap["guildID"], "channel should be initiated")
}

func cleanActiveVoiceChannels() {
	voice.ActiveVoiceChannels = voice.NewActiveVoiceChannels()
}
