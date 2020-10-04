package consumers

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"

	prisma "github.com/dangdennis/crossing/common/db"
	"github.com/dangdennis/crossing/common/repositories/users"
)

// MessageCreate consumes Discord MessageCreate events
func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}

	err := initUser(m)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch m.Content {
	case "!ping":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Pong!")
	case "!pong":
		_, _ = s.ChannelMessageSend(m.ChannelID, "Ping!")
	case "!raid":
		RaidCommand(s, m)
	case "!join":
		JoinCommand(s, m)
	case "!action":
		ActionCommand(s, m)
	case "!help":
		HelpCommand(s, m)
	case "!bomb":
		BombCommand(s, m)
	case "!intro":
		IntroCommand(s, m)
	case "!outro":
		OutroCommand(s, m)
	}
}

// initUser creates a new user, avatar, and wallet
func initUser(m *discordgo.MessageCreate) error {
	// Consider hardening this with an additional cache layer. Check the LRU cache for a discord user id that's recently messaged the channel
	if !strings.HasPrefix(m.Content, "!") {
		return nil
	}

	_, err := users.FindUserByDiscordID(prisma.Client(), m.Author.ID)
	if err == nil {
		fmt.Println("user already exists")
		return nil
	}

	fmt.Println("initializing new user")

	user, err := users.CreateUser(prisma.Client(), users.UserAttrs{DiscordUserID: m.Author.ID})
	if err != nil {
		return err
	}

	_, err = users.CreateAvatar(prisma.Client(), user.ID)
	if err != nil {
		return err
	}

	_, err = users.CreateWallet(prisma.Client(), user.ID)
	if err != nil {
		return err
	}

	return nil
}