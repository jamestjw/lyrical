package voice

import (
	"github.com/bwmarrin/discordgo"
)

// DGVoiceConnection is a wrapper on discordgo.VoiceConnection
// that will handle all operations related to it.
type DGVoiceConnection struct {
	Connection *discordgo.VoiceConnection
}

// Disconnect will disconnect the bot from a particular connection
func (vc DGVoiceConnection) Disconnect() error {
	return vc.Connection.Disconnect()
}

// GetGuildID returns the guild ID that belongs to a particular connection
func (vc DGVoiceConnection) GetGuildID() string {
	return vc.Connection.GuildID
}

// GetAudioInputChannel returns the Opus audio channel that belongs to
// the audio channel of this connection
func (vc DGVoiceConnection) GetAudioInputChannel() chan []byte {
	return vc.Connection.OpusSend
}
