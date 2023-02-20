package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nohns/muggie-bot"
	"github.com/nohns/muggie-bot/discord"
)

// Message matches that triggers Muggie response
var muggieMatches = []string{"det er ik fordi", "det ik fordi", "det er ikke fordi"}

type bot struct {
	msgProvider muggie.MessageProvider
	msgReplier  muggie.MessageReplier
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

	// Return instantiated bot
	return &bot{
		msgProvider: dp,
		msgReplier:  dp,
	}, nil
}

func (b *bot) run() error {
	// Register handler for responding with Muggie quote
	err := b.msgProvider.OnMsg(b.respondWithMuggie)
	if err != nil {
		return fmt.Errorf("could not register respond with muggie handler: %v", err)
	}

	// Keep app running, until SIGTERM etc...
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
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
