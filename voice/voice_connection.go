package voice

import (
	"github.com/bwmarrin/discordgo"
)

type DGVoiceConnection struct {
	Connection *discordgo.VoiceConnection
}

func (vc DGVoiceConnection) Disconnect() error {
	return vc.Connection.Disconnect()
}

func (vc DGVoiceConnection) GetGuildID() string {
	return vc.Connection.GuildID
}
