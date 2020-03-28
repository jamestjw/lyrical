package voice

import (
	"testing"

	"github.com/golang/mock/gomock"
	mock_interfaces "github.com/jamestjw/lyrical/mock_main"
	"github.com/stretchr/testify/mock"
)

type voiceConnectionMock struct {
	mock.Mock
}

func (m *voiceConnectionMock) Disconnect() error {
	args := m.Called()
	return args.Error(0)
}

func disconnectAllVoiceConnectionsSetup() map[string]*voiceConnectionMock {
	voiceConnectionMap := make(map[string]*voiceConnectionMock)
	for i := 0; i < 5; i++ {
		voiceConnection := new(voiceConnectionMock)
		voiceConnection.On("Disconnect").Return(nil)
		voiceConnectionMap[string(i)] = voiceConnection
	}
	return voiceConnectionMap
}

func TestDisconnectAllVoiceConnections(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSession := mock_interfaces.NewMockSession(ctrl)
	voiceConnectionMap := disconnectAllVoiceConnectionsSetup()
	mockSession.EXPECT().GetVoiceConnections().Times(1).Return(voiceConnectionMap)

	DisconnectAllVoiceConnections(mockSession)
}
