// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mock_main is a generated GoMock package.
package mock_main

import (
	discordgo "github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
	voice "github.com/jamestjw/lyrical/voice"
	reflect "reflect"
)

// MockSession is a mock of Session interface.
type MockSession struct {
	ctrl     *gomock.Controller
	recorder *MockSessionMockRecorder
}

// MockSessionMockRecorder is the mock recorder for MockSession.
type MockSessionMockRecorder struct {
	mock *MockSession
}

// NewMockSession creates a new mock instance.
func NewMockSession(ctrl *gomock.Controller) *MockSession {
	mock := &MockSession{ctrl: ctrl}
	mock.recorder = &MockSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSession) EXPECT() *MockSessionMockRecorder {
	return m.recorder
}

// CloseConnection mocks base method.
func (m *MockSession) CloseConnection() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseConnection")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseConnection indicates an expected call of CloseConnection.
func (mr *MockSessionMockRecorder) CloseConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseConnection", reflect.TypeOf((*MockSession)(nil).CloseConnection))
}

// ListenAndServe mocks base method.
func (m *MockSession) ListenAndServe() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListenAndServe")
	ret0, _ := ret[0].(error)
	return ret0
}

// ListenAndServe indicates an expected call of ListenAndServe.
func (mr *MockSessionMockRecorder) ListenAndServe() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListenAndServe", reflect.TypeOf((*MockSession)(nil).ListenAndServe))
}

// AddHandler mocks base method.
func (m *MockSession) AddHandler(arg0 interface{}) func() {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddHandler", arg0)
	ret0, _ := ret[0].(func())
	return ret0
}

// AddHandler indicates an expected call of AddHandler.
func (mr *MockSessionMockRecorder) AddHandler(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHandler", reflect.TypeOf((*MockSession)(nil).AddHandler), arg0)
}

// MockEvent is a mock of Event interface.
type MockEvent struct {
	ctrl     *gomock.Controller
	recorder *MockEventMockRecorder
}

// MockEventMockRecorder is the mock recorder for MockEvent.
type MockEventMockRecorder struct {
	mock *MockEvent
}

// NewMockEvent creates a new mock instance.
func NewMockEvent(ctrl *gomock.Controller) *MockEvent {
	mock := &MockEvent{ctrl: ctrl}
	mock.recorder = &MockEventMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEvent) EXPECT() *MockEventMockRecorder {
	return m.recorder
}

// SendMessage mocks base method.
func (m *MockEvent) SendMessage(message string) *discordgo.Message {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMessage", message)
	ret0, _ := ret[0].(*discordgo.Message)
	return ret0
}

// SendMessage indicates an expected call of SendMessage.
func (mr *MockEventMockRecorder) SendMessage(message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMessage", reflect.TypeOf((*MockEvent)(nil).SendMessage), message)
}

// FindVoiceChannel mocks base method.
func (m *MockEvent) FindVoiceChannel(channelName string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindVoiceChannel", channelName)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindVoiceChannel indicates an expected call of FindVoiceChannel.
func (mr *MockEventMockRecorder) FindVoiceChannel(channelName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindVoiceChannel", reflect.TypeOf((*MockEvent)(nil).FindVoiceChannel), channelName)
}

// GetSession mocks base method.
func (m *MockEvent) GetSession() voice.Connectable {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSession")
	ret0, _ := ret[0].(voice.Connectable)
	return ret0
}

// GetSession indicates an expected call of GetSession.
func (mr *MockEventMockRecorder) GetSession() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSession", reflect.TypeOf((*MockEvent)(nil).GetSession))
}

// GetGuildID mocks base method.
func (m *MockEvent) GetGuildID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGuildID")
	ret0, _ := ret[0].(string)
	return ret0
}

// GetGuildID indicates an expected call of GetGuildID.
func (mr *MockEventMockRecorder) GetGuildID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGuildID", reflect.TypeOf((*MockEvent)(nil).GetGuildID))
}

// GetVoiceConnection mocks base method.
func (m *MockEvent) GetVoiceConnection() (voice.Connection, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVoiceConnection")
	ret0, _ := ret[0].(voice.Connection)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetVoiceConnection indicates an expected call of GetVoiceConnection.
func (mr *MockEventMockRecorder) GetVoiceConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVoiceConnection", reflect.TypeOf((*MockEvent)(nil).GetVoiceConnection))
}

// GetMessageByMessageID mocks base method.
func (m *MockEvent) GetMessageByMessageID(messageID string) (*discordgo.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageByMessageID", messageID)
	ret0, _ := ret[0].(*discordgo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMessageByMessageID indicates an expected call of GetMessageByMessageID.
func (mr *MockEventMockRecorder) GetMessageByMessageID(messageID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageByMessageID", reflect.TypeOf((*MockEvent)(nil).GetMessageByMessageID), messageID)
}

// React mocks base method.
func (m *MockEvent) React(emoji string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "React", emoji)
}

// React indicates an expected call of React.
func (mr *MockEventMockRecorder) React(emoji interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "React", reflect.TypeOf((*MockEvent)(nil).React), emoji)
}

// MockSearcher is a mock of Searcher interface.
type MockSearcher struct {
	ctrl     *gomock.Controller
	recorder *MockSearcherMockRecorder
}

// MockSearcherMockRecorder is the mock recorder for MockSearcher.
type MockSearcherMockRecorder struct {
	mock *MockSearcher
}

// NewMockSearcher creates a new mock instance.
func NewMockSearcher(ctrl *gomock.Controller) *MockSearcher {
	mock := &MockSearcher{ctrl: ctrl}
	mock.recorder = &MockSearcherMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSearcher) EXPECT() *MockSearcherMockRecorder {
	return m.recorder
}

// GetVideoID mocks base method.
func (m *MockSearcher) GetVideoID(arg0 string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVideoID", arg0)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVideoID indicates an expected call of GetVideoID.
func (mr *MockSearcherMockRecorder) GetVideoID(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVideoID", reflect.TypeOf((*MockSearcher)(nil).GetVideoID), arg0)
}
