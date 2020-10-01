// Code generated by MockGen. DO NOT EDIT.
// Source: voice/interfaces.go

// Package mock_voice is a generated GoMock package.
package mock_voice

import (
	gomock "github.com/golang/mock/gomock"
	playlist "github.com/jamestjw/lyrical/playlist"
	voice "github.com/jamestjw/lyrical/voice"
	reflect "reflect"
)

// MockConnectable is a mock of Connectable interface
type MockConnectable struct {
	ctrl     *gomock.Controller
	recorder *MockConnectableMockRecorder
}

// MockConnectableMockRecorder is the mock recorder for MockConnectable
type MockConnectableMockRecorder struct {
	mock *MockConnectable
}

// NewMockConnectable creates a new mock instance
func NewMockConnectable(ctrl *gomock.Controller) *MockConnectable {
	mock := &MockConnectable{ctrl: ctrl}
	mock.recorder = &MockConnectableMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConnectable) EXPECT() *MockConnectableMockRecorder {
	return m.recorder
}

// GetVoiceConnections mocks base method
func (m *MockConnectable) GetVoiceConnections() map[string]voice.Connection {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVoiceConnections")
	ret0, _ := ret[0].(map[string]voice.Connection)
	return ret0
}

// GetVoiceConnections indicates an expected call of GetVoiceConnections
func (mr *MockConnectableMockRecorder) GetVoiceConnections() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVoiceConnections", reflect.TypeOf((*MockConnectable)(nil).GetVoiceConnections))
}

// JoinVoiceChannel mocks base method
func (m *MockConnectable) JoinVoiceChannel(guildID, voiceChannelID string) (voice.Connection, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JoinVoiceChannel", guildID, voiceChannelID)
	ret0, _ := ret[0].(voice.Connection)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// JoinVoiceChannel indicates an expected call of JoinVoiceChannel
func (mr *MockConnectableMockRecorder) JoinVoiceChannel(guildID, voiceChannelID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JoinVoiceChannel", reflect.TypeOf((*MockConnectable)(nil).JoinVoiceChannel), guildID, voiceChannelID)
}

// MockConnection is a mock of Connection interface
type MockConnection struct {
	ctrl     *gomock.Controller
	recorder *MockConnectionMockRecorder
}

// MockConnectionMockRecorder is the mock recorder for MockConnection
type MockConnectionMockRecorder struct {
	mock *MockConnection
}

// NewMockConnection creates a new mock instance
func NewMockConnection(ctrl *gomock.Controller) *MockConnection {
	mock := &MockConnection{ctrl: ctrl}
	mock.recorder = &MockConnectionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockConnection) EXPECT() *MockConnectionMockRecorder {
	return m.recorder
}

// Disconnect mocks base method
func (m *MockConnection) Disconnect() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Disconnect")
	ret0, _ := ret[0].(error)
	return ret0
}

// Disconnect indicates an expected call of Disconnect
func (mr *MockConnectionMockRecorder) Disconnect() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Disconnect", reflect.TypeOf((*MockConnection)(nil).Disconnect))
}

// GetGuildID mocks base method
func (m *MockConnection) GetGuildID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGuildID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetGuildID indicates an expected call of GetGuildID
func (mr *MockConnectionMockRecorder) GetGuildID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGuildID", reflect.TypeOf((*MockConnection)(nil).GetGuildID))
}

// GetAudioInputChannel mocks base method
func (m *MockConnection) GetAudioInputChannel() chan []byte {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAudioInputChannel")
	ret0, _ := ret[0].(chan []byte)
	return ret0
}

// GetAudioInputChannel indicates an expected call of GetAudioInputChannel
func (mr *MockConnectionMockRecorder) GetAudioInputChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAudioInputChannel", reflect.TypeOf((*MockConnection)(nil).GetAudioInputChannel))
}

// MockChannel is a mock of Channel interface
type MockChannel struct {
	ctrl     *gomock.Controller
	recorder *MockChannelMockRecorder
}

// MockChannelMockRecorder is the mock recorder for MockChannel
type MockChannelMockRecorder struct {
	mock *MockChannel
}

// NewMockChannel creates a new mock instance
func NewMockChannel(ctrl *gomock.Controller) *MockChannel {
	mock := &MockChannel{ctrl: ctrl}
	mock.recorder = &MockChannelMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockChannel) EXPECT() *MockChannelMockRecorder {
	return m.recorder
}

// RemoveNowPlaying mocks base method
func (m *MockChannel) RemoveNowPlaying() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RemoveNowPlaying")
}

// RemoveNowPlaying indicates an expected call of RemoveNowPlaying
func (mr *MockChannelMockRecorder) RemoveNowPlaying() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveNowPlaying", reflect.TypeOf((*MockChannel)(nil).RemoveNowPlaying))
}

// GetNext mocks base method
func (m *MockChannel) GetNext() *playlist.Song {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNext")
	ret0, _ := ret[0].(*playlist.Song)
	return ret0
}

// GetNext indicates an expected call of GetNext
func (mr *MockChannelMockRecorder) GetNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNext", reflect.TypeOf((*MockChannel)(nil).GetNext))
}

// SetNext mocks base method
func (m *MockChannel) SetNext(arg0 *playlist.Song) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNext", arg0)
}

// SetNext indicates an expected call of SetNext
func (mr *MockChannelMockRecorder) SetNext(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNext", reflect.TypeOf((*MockChannel)(nil).SetNext), arg0)
}

// SetNowPlaying mocks base method
func (m *MockChannel) SetNowPlaying(s *playlist.Song) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetNowPlaying", s)
}

// SetNowPlaying indicates an expected call of SetNowPlaying
func (mr *MockChannelMockRecorder) SetNowPlaying(s interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetNowPlaying", reflect.TypeOf((*MockChannel)(nil).SetNowPlaying), s)
}

// GetAbortChannel mocks base method
func (m *MockChannel) GetAbortChannel() chan string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAbortChannel")
	ret0, _ := ret[0].(chan string)
	return ret0
}

// GetAbortChannel indicates an expected call of GetAbortChannel
func (mr *MockChannelMockRecorder) GetAbortChannel() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAbortChannel", reflect.TypeOf((*MockChannel)(nil).GetAbortChannel))
}

// IsPlayingMusic mocks base method
func (m *MockChannel) IsPlayingMusic() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsPlayingMusic")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsPlayingMusic indicates an expected call of IsPlayingMusic
func (mr *MockChannelMockRecorder) IsPlayingMusic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsPlayingMusic", reflect.TypeOf((*MockChannel)(nil).IsPlayingMusic))
}

// GetNowPlayingName mocks base method
func (m *MockChannel) GetNowPlayingName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNowPlayingName")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetNowPlayingName indicates an expected call of GetNowPlayingName
func (mr *MockChannelMockRecorder) GetNowPlayingName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNowPlayingName", reflect.TypeOf((*MockChannel)(nil).GetNowPlayingName))
}

// StopMusic mocks base method
func (m *MockChannel) StopMusic() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "StopMusic")
}

// StopMusic indicates an expected call of StopMusic
func (mr *MockChannelMockRecorder) StopMusic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopMusic", reflect.TypeOf((*MockChannel)(nil).StopMusic))
}

// FetchPlaylist mocks base method
func (m *MockChannel) FetchPlaylist() *playlist.Playlist {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FetchPlaylist")
	ret0, _ := ret[0].(*playlist.Playlist)
	return ret0
}

// FetchPlaylist indicates an expected call of FetchPlaylist
func (mr *MockChannelMockRecorder) FetchPlaylist() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FetchPlaylist", reflect.TypeOf((*MockChannel)(nil).FetchPlaylist))
}

// GetNextSongs mocks base method
func (m *MockChannel) GetNextSongs() ([]*playlist.Song, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextSongs")
	ret0, _ := ret[0].([]*playlist.Song)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetNextSongs indicates an expected call of GetNextSongs
func (mr *MockChannelMockRecorder) GetNextSongs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextSongs", reflect.TypeOf((*MockChannel)(nil).GetNextSongs))
}

// ExistsNext mocks base method
func (m *MockChannel) ExistsNext() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsNext")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ExistsNext indicates an expected call of ExistsNext
func (mr *MockChannelMockRecorder) ExistsNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsNext", reflect.TypeOf((*MockChannel)(nil).ExistsNext))
}

// ExistsBackupNext mocks base method
func (m *MockChannel) ExistsBackupNext() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExistsBackupNext")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ExistsBackupNext indicates an expected call of ExistsBackupNext
func (mr *MockChannelMockRecorder) ExistsBackupNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExistsBackupNext", reflect.TypeOf((*MockChannel)(nil).ExistsBackupNext))
}

// GetBackupNext mocks base method
func (m *MockChannel) GetBackupNext() *playlist.Song {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBackupNext")
	ret0, _ := ret[0].(*playlist.Song)
	return ret0
}

// GetBackupNext indicates an expected call of GetBackupNext
func (mr *MockChannelMockRecorder) GetBackupNext() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBackupNext", reflect.TypeOf((*MockChannel)(nil).GetBackupNext))
}

// GetNextBackupSongs mocks base method
func (m *MockChannel) GetNextBackupSongs() ([]*playlist.Song, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNextBackupSongs")
	ret0, _ := ret[0].([]*playlist.Song)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetNextBackupSongs indicates an expected call of GetNextBackupSongs
func (mr *MockChannelMockRecorder) GetNextBackupSongs() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNextBackupSongs", reflect.TypeOf((*MockChannel)(nil).GetNextBackupSongs))
}

// MockDownloader is a mock of Downloader interface
type MockDownloader struct {
	ctrl     *gomock.Controller
	recorder *MockDownloaderMockRecorder
}

// MockDownloaderMockRecorder is the mock recorder for MockDownloader
type MockDownloaderMockRecorder struct {
	mock *MockDownloader
}

// NewMockDownloader creates a new mock instance
func NewMockDownloader(ctrl *gomock.Controller) *MockDownloader {
	mock := &MockDownloader{ctrl: ctrl}
	mock.recorder = &MockDownloaderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDownloader) EXPECT() *MockDownloaderMockRecorder {
	return m.recorder
}

// Download mocks base method
func (m *MockDownloader) Download(query string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Download", query)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Download indicates an expected call of Download
func (mr *MockDownloaderMockRecorder) Download(query interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Download", reflect.TypeOf((*MockDownloader)(nil).Download), query)
}

// MockMusicPlayer is a mock of MusicPlayer interface
type MockMusicPlayer struct {
	ctrl     *gomock.Controller
	recorder *MockMusicPlayerMockRecorder
}

// MockMusicPlayerMockRecorder is the mock recorder for MockMusicPlayer
type MockMusicPlayerMockRecorder struct {
	mock *MockMusicPlayer
}

// NewMockMusicPlayer creates a new mock instance
func NewMockMusicPlayer(ctrl *gomock.Controller) *MockMusicPlayer {
	mock := &MockMusicPlayer{ctrl: ctrl}
	mock.recorder = &MockMusicPlayerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockMusicPlayer) EXPECT() *MockMusicPlayerMockRecorder {
	return m.recorder
}

// PlayMusic mocks base method
func (m *MockMusicPlayer) PlayMusic(input chan []byte, guildID string, song voice.Channel, main bool) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "PlayMusic", input, guildID, song, main)
}

// PlayMusic indicates an expected call of PlayMusic
func (mr *MockMusicPlayerMockRecorder) PlayMusic(input, guildID, song, main interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlayMusic", reflect.TypeOf((*MockMusicPlayer)(nil).PlayMusic), input, guildID, song, main)
}
