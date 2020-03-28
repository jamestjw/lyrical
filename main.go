package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/voice"
)

func init() {
	database.InitialiseDatabase()
}

func main() {
	dg, err := NewSession(config.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	defer dg.CloseConnection()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(dummyMessageCreate)
	dg.AddHandler(joinVoiceChannelRequest)
	dg.AddHandler(leaveVoiceChannelRequest)
	dg.AddHandler(stopMusicRequest)
	dg.AddHandler(playMusicRequest)
	dg.AddHandler(skipMusicRequest)
	dg.AddHandler(addToPlaylistRequest)
	dg.AddHandler(nowPlayingRequest)
	dg.AddHandler(helpRequest)

	// Open a websocket connection to Discord and begin listening.
	err = dg.ListenAndServe()
	if err != nil {
		log.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running. Press CTRL-C to exit.")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	select {
	case <-sc:
		log.Println("Received signal to terminate, cleaning up...")
		// Cleanly close down the Discord session.
		voice.DisconnectAllVoiceConnections(dg.(voice.Connectable))
		database.Connection.Close()
		log.Println("Exit successful!")
		return
	}
}
