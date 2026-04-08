package main

import (
	"fmt"
	"time"

	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	_ "github.com/hajimehoshi/oto/v2"

	"sync"
)

func (ae *AudioEngine) createStream(filepath string, streamID int) *StreamPlayer {
	s := StreamPlayer{
		filepath: filepath,
	}
	s.loadFile()
	s.loaded = true
	s.Paused = true
	if !ae.initialized {

		ae.Format = s.Format
	}

	ae.mixer.Add(&s)
	ae.streamers[streamID] = &s
	return &s
}

type AudioEngine struct {
	mixer       *beep.Mixer
	streamers   []*StreamPlayer
	initialized bool
	mu          sync.RWMutex
	Format      beep.Format
}

func (ae *AudioEngine) Play() error {
	if !ae.initialized {
		ae.initialized = true
		speaker.Init(ae.Format.SampleRate, ae.Format.SampleRate.N(time.Second/10))
		speaker.Play(ae.mixer)
		fmt.Printf("\nspeaker init\n")

	}
	return nil
}

func (ae *AudioEngine) Stop(streamID int) error {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	ae.streamers[streamID].Paused = true
	err := ae.streamers[streamID].CurrentStream.Seek(0)
	if err != nil {
		return audioErrorMsg(err)
	}
	return nil
}
