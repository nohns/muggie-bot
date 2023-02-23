package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
	"github.com/nohns/muggie-bot/pkg/discordmp3"
)

var botCmds = []discordgo.ApplicationCommand{
	{
		Name:        "jazzroll",
		Description: "Join a voice channel and play jazz",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "target-user",
				Description: "The user whose voice channel should be jazz'ed. Leave empty for yourself",
			},
		},
	},
	{
		Name:        "duvebarattack",
		Description: "Join a voice channel and play suprise 'du ve bar'",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "target-user",
				Description: "The user whose voice channel should be 'du ve bar'ed. Leave empty for yourself",
			},
		},
	},
	{
		Name:        "nuerdenher",
		Description: "Join a voice channel and play John",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionUser,
				Name:        "target-user",
				Description: "A hungry user. Leave empty for yourself",
			},
		},
	},
}

func registerCmds(s *discordgo.Session, appID string) ([]*discordgo.ApplicationCommand, error) {
	var registeredCmds []*discordgo.ApplicationCommand
	for _, cmd := range botCmds {
		c, err := s.ApplicationCommandCreate(appID, "", &cmd)
		if err != nil {
			return nil, err
		}
		registeredCmds = append(registeredCmds, c)
	}
	return registeredCmds, nil
}

// Get user from command options. Always takes first option if available
func getUserOpt(s *discordgo.Session, i *discordgo.InteractionCreate) *discordgo.User {
	if len(i.ApplicationCommandData().Options) == 0 {
		return nil
	}

	return i.ApplicationCommandData().Options[0].UserValue(s)
}

type interactionRespErr struct {
	msg   string
	inner error
}

func (e interactionRespErr) Error() string {
	return e.inner.Error()
}

func (e interactionRespErr) Unwrap() error {
	return e.inner
}

func forCmd(cmd string, h func(s *discordgo.Session, i *discordgo.InteractionCreate) error) func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		// Only accept cmd name given
		if i.ApplicationCommandData().Name != cmd {
			return
		}

		// Call handler, and log error and respond to interaction if any message is set on error
		err := h(s, i)
		var e *interactionRespErr
		if errors.As(err, &e) && e.msg != "" {
			ierr := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: e.msg,
				},
			})
			if ierr != nil {
				log.Printf("error responding to interaction: %v", ierr)
			}
		}
		if err != nil {
			log.Printf("error handling command '%s': %v", cmd, err)
		}
	}
}

func handlePizzaburgerCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Open pizzaburger commercial mp3 file
	mp3file, err := readPwdFile("pizzaburgeren.mp3")
	if err != nil {
		return &interactionRespErr{
			msg:   "A ka sgu æt find æ mappe mæ æ musik å æ server",
			inner: fmt.Errorf("error occured while getting working directory: %v", err),
		}
	}

	// Get target user, either from command options or from interaction author as fallback
	target := getUserOpt(s, i)
	if target == nil {
		target = i.Member.User
	}

	// Join target user's voice channel
	p, err := joinAndPlayAudio(s, target, mp3file)
	if err != nil {
		return err
	}

	// Tell user that NU ER DEN HER
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Nu er æn her %s", target.Mention()),
		},
	})

	// Wait for audio to finish playing or stop playing
	p.WaitForEnd()
	if err := p.Close(); err != nil {
		return err
	}

	return nil
}

func handleDuvebarCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Open duvebar mp3 file
	mp3file, err := readPwdFile("duvebar.mp3")
	if err != nil {
		return &interactionRespErr{
			msg:   "A ka sgu æt find æ mappe mæ æ musik å æ server",
			inner: fmt.Errorf("error occured while getting working directory: %v", err),
		}
	}

	// Get target user, either from command options or from interaction author as fallback
	target := getUserOpt(s, i)
	if target == nil {
		target = i.Member.User
	}

	// Join target user's voice channel
	p, err := joinAndPlayAudio(s, target, mp3file)
	if err != nil {
		return err
	}

	// Tell user that they ve bar
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s du ve bar'", target.Mention()),
		},
	})

	// Wait for audio to finish playing or stop playing
	p.WaitForEnd()
	if err := p.Close(); err != nil {
		return err
	}

	return nil
}

func handleJazzrollCmd(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	// Open jazz mp3 file
	mp3file, err := readPwdFile("moveonup.mp3")
	if err != nil {
		return &interactionRespErr{
			msg:   "A ka sgu æt find æ mappe mæ æ musik å æ server",
			inner: fmt.Errorf("error occured while getting working directory: %v", err),
		}
	}

	// Get target user, either from command options or from interaction author as fallback
	target := getUserOpt(s, i)
	if target == nil {
		target = i.Member.User
	}

	// Join target user's voice channel
	p, err := joinAndPlayAudio(s, target, mp3file)
	if err != nil {
		return err
	}

	// Tell user that jazz is on the way
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("A wil ind og spil noe jazz for %s", target.Mention()),
		},
	})

	// Wait for audio to finish playing or stop playing
	p.WaitForEnd()
	if err := p.Close(); err != nil {
		return err
	}

	return nil
}

// Read a file from current working directory. The directory from where the executable is started from
func readPwdFile(filename string) (io.Reader, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filepath.Join(dir, filename), os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func joinAndPlayAudio(s *discordgo.Session, target *discordgo.User, mp3file io.Reader) (*discordmp3.Player, error) {

	// Get voice channel and guild of user
	var guidID string
	var channelID string
	for _, g := range s.State.Guilds {
		for _, v := range g.VoiceStates {
			if v.UserID == target.ID {
				// get channel
				channel, err := s.Channel(v.ChannelID)
				if err != nil {
					return nil, &interactionRespErr{msg: "Ejj fordævden da... A ka sgu æt find æ bruger si kanal", inner: err}
				}
				// Only allow voice channels
				if channel.Type != discordgo.ChannelTypeGuildVoice {
					break
				}

				guidID = channel.GuildID
				channelID = channel.ID
				break
			}
		}
	}

	// If user is not in a voice channel, return
	if channelID == "" {
		return nil, &interactionRespErr{msg: "Æ bruger er æt i en tele kanal", inner: nil}
	}

	// Jazz roll the channel
	vc, err := s.ChannelVoiceJoin(guidID, channelID, false, true)
	if err != nil {
		return nil, &interactionRespErr{
			msg:   "For sivan da! A ka sgu æt join æ kanal",
			inner: fmt.Errorf("error occured while joining voice channel: %v", err),
		}
	}

	err = vc.Speaking(true)
	if err != nil {
		return nil, &interactionRespErr{
			msg:   "Herre jemini..! Nøj gik gal da A sgu til og snak",
			inner: fmt.Errorf("error occured while setting speaking: %v", err),
		}
	}

	p := discordmp3.NewPlayer(vc)
	if err = p.Play(mp3file); err != nil {
		return nil, &interactionRespErr{
			msg:   "A ka sgu æt spil æ musik å æ server",
			inner: fmt.Errorf("error occured while playing mp3 file: %v", err),
		}
	}

	return p, nil
}
