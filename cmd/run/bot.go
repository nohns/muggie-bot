package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/nohns/muggie-bot"
	"github.com/nohns/muggie-bot/discord"
)

// Message matches that triggers Muggie response
var muggieMatches = []string{"det er ik fordi", "det ik fordi", "det er ikke fordi"}

type bot struct {
	msgProvider muggie.MessageProvider
	msgReplier  muggie.MessageReplier
	s           *discordgo.Session
	cmds        []*discordgo.ApplicationCommand
}

// Bootstrap bot with configuration and dependencies
func bootstrap() (*bot, error) {

	// Read config
	conf, err := readConfFromEnv()
	if err != nil {
		return nil, fmt.Errorf("could not read config: %v", err)
	}

	// Connect to discord and create provider for it
	s, err := discord.NewDiscordSession(conf.token)
	if err != nil {
		return nil, fmt.Errorf("could not start discord session: %v", err)
	}
	dp := discord.NewDiscordProvider(s)

	log.Printf("registering bot commands...")
	cmds, err := registerCmds(s, conf.appID)
	if err != nil {
		return nil, fmt.Errorf("could not create command: %v", err)
	}

	// Return instantiated bot
	return &bot{
		msgProvider: dp,
		msgReplier:  dp,
		s:           s,
		cmds:        cmds,
	}, nil
}

func (b *bot) run() error {
	// Register handler for responding with Muggie quote
	err := b.msgProvider.OnMsg(b.respondWithMuggie)
	if err != nil {
		return fmt.Errorf("could not register respond with muggie handler: %v", err)
	}

	// Handle commands
	b.s.AddHandler(forCmd("jazzroll", handleJazzrollCmd))
	b.s.AddHandler(forCmd("duvebarattack", handleDuvebarCmd))
	b.s.AddHandler(forCmd("nuerdenher", handlePizzaburgerCmd))

	log.Printf("bot is now running and awaiting events...")

	// Keep app running, until SIGTERM etc...
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	fmt.Println("")
	log.Printf("unregistering commands...")
	for _, cmd := range b.cmds {
		err := b.s.ApplicationCommandDelete(b.s.State.User.ID, "", cmd.ID)
		if err != nil {
			log.Printf("could not unregister command '%s': %v", cmd.Name, err)
		}
	}
	log.Printf("closing bot connections...")
	b.s.Close()
	return fmt.Errorf("goodbye")
}

// Respond with muggie "Du ve' bar" response to specific messages
func (b *bot) respondWithMuggie(msg muggie.Message) error {
	// If message is not a muggie match, dont reply
	if !containsMuggieMatch(msg.Content()) {
		return nil
	}

	time.Sleep(2 * time.Second)
	b.msgReplier.ReplyTo(msg, "du ve' bar")
	return nil
}

// Test string for potential muggie matches
func containsMuggieMatch(s string) bool {
	for _, mm := range muggieMatches {
		if strings.Contains(strings.ToLower(s), mm) {
			return true
		}
	}

	return false
}

// Print out hello world
