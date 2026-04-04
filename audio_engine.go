package main

import (
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	_ "github.com/hajimehoshi/oto/v2"

	"sync"
	"time"
)

func createStream(filepath string) StreamPlayer {
	s := StreamPlayer{}
	s.loadFile(filepath)
	return s
}

type AudioEngine struct {
	mixer       *beep.Mixer
	Active      map[int]*StreamPlayer
	initialized bool
	volumeCtrl  *effects.Volume
	pauseCtrl   *beep.Ctrl
	mu          sync.RWMutex
}

func (ae *AudioEngine) Play(filepath string, streamID int) error {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	s := createStream(filepath)

	if !ae.initialized {
		speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/10))
		speaker.Play(ae.pauseCtrl)
		ae.initialized = true
	}
	ae.mixer.Add(&s)

	return nil
}

func (ae *AudioEngine) Stop(streamID int) error {
	ae.mu.Lock()
	defer ae.mu.Unlock()
	ae.pauseCtrl.Paused = true
	return nil
}
