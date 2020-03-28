package voice_test

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	mock_voice "github.com/jamestjw/lyrical/mocks/mock_voice"
	"github.com/jamestjw/lyrical/voice"
	"github.com/stretchr/testify/mock"
)

type voiceConnectionMock struct {
	mock.Mock
}

func (m *voiceConnectionMock) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func disconnectAllVoiceConnectionsSetup(ctrl *gomock.Controller) map[string]voice.Connection {
	voice.ActiveVoiceChannels = voice.NewActiveVoiceChannels()
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
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := mock_voice.NewMockConnectable(ctrl)
	voiceConnectionMap := disconnectAllVoiceConnectionsSetup(ctrl)
	mockSession.EXPECT().GetVoiceConnections().Times(1).Return(voiceConnectionMap)

	voice.DisconnectAllVoiceConnections(mockSession)
}
