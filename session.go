package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jamestjw/lyrical/voice"
)

// NewSession creates a new session connected to discord
func NewSession(discordToken string) (s Session, err error) {
	dg, err := discordgo.New("Bot " + discordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	s = botSession{dg}
	return
}

func (s botSession) GetVoiceConnections() map[string]voice.Connection {
	vcMap := make(map[string]voice.Connection)
	for key, value := range s.session.VoiceConnections {
		vcMap[key] = voice.DGVoiceConnection{Connection: value}
	}
	return vcMap
}

func (s botSession) CloseConnection() error {
	return s.session.Close()
}

func (s botSession) ListenAndServe() error {
	return s.session.Open()
}

func (s botSession) AddHandler(handler interface{}) func() {
	return s.session.AddHandler(handler)
}

func (s botSession) JoinVoiceChannel(guildID string, voiceChannelID string) (voice.Connection, error) {
	vc, err := s.session.ChannelVoiceJoin(guildID, voiceChannelID, false, false)
	return voice.DGVoiceConnection{Connection: vc}, err
}
