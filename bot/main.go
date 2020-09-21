package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"

	"github.com/dangdennis/crossing/bot/consumers"
	"github.com/dangdennis/crossing/bot/db"
	"github.com/dangdennis/crossing/bot/env"
)

func main() {
	client := db.Client()

	defer func() {
		err := client.Disconnect()
		if err != nil {
			panic(err)
		}
	}()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + env.GetDiscordBotKey())
	if err != nil {
		log.Panic("failed to create Discord session", err)
	}

	defer func() {
		_ = dg.Close()
	}()

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(consumers.MessageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
