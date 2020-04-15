package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/jamestjw/lyrical/database"
	"github.com/jamestjw/lyrical/voice"
)

func main() {
	initialiseApplication()

	dg, err := NewSession(config.DiscordToken)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}
	defer dg.CloseConnection()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(heartbeatHandlerFunc)
	dg.AddHandler(defaultMux.Route)

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
