// Code generated by MockGen. DO NOT EDIT.
// Source: infrastructures/discord.go

// Package mocks is a generated GoMock package.
package mocks

import (
	discordgo "github.com/bwmarrin/discordgo"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockDiscordSession is a mock of DiscordSession interface
type MockDiscordSession struct {
	ctrl     *gomock.Controller
	recorder *MockDiscordSessionMockRecorder
}

// MockDiscordSessionMockRecorder is the mock recorder for MockDiscordSession
type MockDiscordSessionMockRecorder struct {
	mock *MockDiscordSession
}

// NewMockDiscordSession creates a new mock instance
func NewMockDiscordSession(ctrl *gomock.Controller) *MockDiscordSession {
	mock := &MockDiscordSession{ctrl: ctrl}
	mock.recorder = &MockDiscordSessionMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDiscordSession) EXPECT() *MockDiscordSessionMockRecorder {
	return m.recorder
}

// AddHandler mocks base method
func (m *MockDiscordSession) AddHandler(handler interface{}) func() {
	ret := m.ctrl.Call(m, "AddHandler", handler)
	ret0, _ := ret[0].(func())
	return ret0
}

// AddHandler indicates an expected call of AddHandler
func (mr *MockDiscordSessionMockRecorder) AddHandler(handler interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddHandler", reflect.TypeOf((*MockDiscordSession)(nil).AddHandler), handler)
}

// Channel mocks base method
func (m *MockDiscordSession) Channel(arg0 string) (*discordgo.Channel, error) {
	ret := m.ctrl.Call(m, "Channel", arg0)
	ret0, _ := ret[0].(*discordgo.Channel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Channel indicates an expected call of Channel
func (mr *MockDiscordSessionMockRecorder) Channel(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Channel", reflect.TypeOf((*MockDiscordSession)(nil).Channel), arg0)
}

// ChannelMessageSend mocks base method
func (m *MockDiscordSession) ChannelMessageSend(arg0, arg1 string) (*discordgo.Message, error) {
	ret := m.ctrl.Call(m, "ChannelMessageSend", arg0, arg1)
	ret0, _ := ret[0].(*discordgo.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ChannelMessageSend indicates an expected call of ChannelMessageSend
func (mr *MockDiscordSessionMockRecorder) ChannelMessageSend(arg0, arg1 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ChannelMessageSend", reflect.TypeOf((*MockDiscordSession)(nil).ChannelMessageSend), arg0, arg1)
}

// Close mocks base method
func (m *MockDiscordSession) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockDiscordSessionMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDiscordSession)(nil).Close))
}

// Guild mocks base method
func (m *MockDiscordSession) Guild(arg0 string) (*discordgo.Guild, error) {
	ret := m.ctrl.Call(m, "Guild", arg0)
	ret0, _ := ret[0].(*discordgo.Guild)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Guild indicates an expected call of Guild
func (mr *MockDiscordSessionMockRecorder) Guild(arg0 interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Guild", reflect.TypeOf((*MockDiscordSession)(nil).Guild), arg0)
}

// Open mocks base method
func (m *MockDiscordSession) Open() error {
	ret := m.ctrl.Call(m, "Open")
	ret0, _ := ret[0].(error)
	return ret0
}

// Open indicates an expected call of Open
func (mr *MockDiscordSessionMockRecorder) Open() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Open", reflect.TypeOf((*MockDiscordSession)(nil).Open))
}