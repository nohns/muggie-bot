package discord

import (
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/nohns/muggie-bot"
)

func NewDiscordSession(token string) (*discordgo.Session, error) {

	s, err := discordgo.New(fmt.Sprintf("Bot %s", token))
	if err != nil {
		return nil, fmt.Errorf("could not create discord session: %v", err)
	}
	s.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAllWithoutPrivileged)

	// Open a websocket connection to Discord and begin listening.
	s.Identify.Intents = discordgo.IntentsGuildMessages
	err = s.Open()
	if err != nil {
		return nil, fmt.Errorf("error opening discord connection: %v", err)
	}

	s.UpdateGameStatus(0, "Pentanque")
	return s, nil
}

type discordProvider struct {
	s *discordgo.Session
}

func NewDiscordProvider(s *discordgo.Session) *discordProvider {
	return &discordProvider{
		s: s,
	}
}

// Register handler when message comes in on any channel
func (dp *discordProvider) OnMsg(h muggie.MsgHandlerFunc) error {
	dp.s.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if err := h(&message{MessageCreate: m}); err != nil {
			log.Printf("error occured while handling discord message: %v", err)
		}
	})

	return nil
}

func (dp *discordProvider) ReplyTo(msg muggie.Message, content string) error {
	dm, ok := msg.(*message)
	if !ok {
		return fmt.Errorf("reply message given must be from discord")
	}

	// Send reply via Discord API
	_, err := dp.s.ChannelMessageSendReply(dm.ChannelID, content, dm.Reference())
	if err != nil {
		return fmt.Errorf("could not send discord message reply: %v", err)
	}

	return nil
}

// Discord message that implements muggie.Message interface
type message struct {
	*discordgo.MessageCreate
}

func (m *message) Content() string {
	return m.MessageCreate.Content
}

func (m *message) ID() string {
	return m.MessageCreate.ID
}
