package discordmp3

import (
	"io"
	"log"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

const (
	frameSize  = 4
	sampleRate = 48000
	channels   = 2
)

type Player struct {
	vc      *discordgo.VoiceConnection
	playing bool
	closed  bool
	mu      sync.Mutex
	audio   *dca.EncodeSession
	end     chan struct{}
}

func NewPlayer(vc *discordgo.VoiceConnection) *Player {
	p := &Player{
		vc: vc,
	}
	go p.stream()

	return p
}

func (p *Player) Play(mp3 io.Reader) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	d, err := dca.EncodeMem(mp3, dca.StdEncodeOptions)

	if err != nil {
		return err
	}
	p.audio = d
	p.end = make(chan struct{})
	p.playing = true

	return nil
}

func (p *Player) Pause() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.playing = false
	return nil
}

func (p *Player) WaitForEnd() {
	if p.end == nil {
		return
	}

	<-p.end
}

func (p *Player) Close() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.playing = false
	p.closed = true

	if err := p.vc.Disconnect(); err != nil {
		return err
	}
	p.vc.Close()

	if p.audio != nil {
		p.audio.Cleanup()
	}

	return nil
}

func (p *Player) stream() {

	for {
		p.mu.Lock()

		if p.closed {
			p.end <- struct{}{}
			close(p.end)
			p.mu.Unlock()
			return
		}
		if !p.playing {
			p.mu.Unlock()
			time.Sleep(100 * time.Millisecond)
			continue
		}

		frame, err := p.audio.OpusFrame()
		if err != nil {
			if err != io.EOF {
				// Handle the error
				log.Printf("error reading opus frame from opus encoder: %v", err)
				p.mu.Unlock()
				return
			}

			p.end <- struct{}{}
			close(p.end)
			p.playing = false
			p.mu.Unlock()
			break
		}

		// Do something with the frame, in this example were sending it to discord
		select {
		case p.vc.OpusSend <- frame:
		case <-time.After(time.Second):
			// We haven't been able to send a frame in a second, assume the connection is borked
			log.Printf("error timeout sending opus frame to discord")

			p.mu.Lock()
			p.playing = false
			p.mu.Unlock()
			return
		}

		p.mu.Unlock()
	}
}
