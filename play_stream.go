package main

import (
	"log"
	"os"
	"time"

	_ "charm.land/bubbletea/v2"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/wav"
	_ "github.com/hajimehoshi/oto/v2"
)

type StreamPlayer struct {
	CurrentStream beep.StreamSeekCloser
	NextStream    beep.StreamSeekCloser
	Crossfade     beep.Streamer
	Format        beep.Format
}

func (s *StreamPlayer) loadFile(filepath string) {

	f1, err := os.Open(filepath)
	if err != nil {
		log.Fatal(err)
	}

	s1, format, err := wav.Decode(f1)
	if err != nil {
		log.Fatal(err)
	}
	f2, _ := os.Open(filepath)

	s2, _, err := wav.Decode(f2)
	if err != nil {
		log.Fatal(err)
	}

	s.CurrentStream = s1

	s.NextStream = s2

	s.Format = format
}

func (s *StreamPlayer) Stream(samples [][2]float64) (n int, ok bool) {
	currentPosition := s.CurrentStream.Position()
	totalLength := s.CurrentStream.Len()
	fadeTime := s.Format.SampleRate.N(time.Second * 15)
	if currentPosition <= totalLength-fadeTime {
	}
	if currentPosition > totalLength-fadeTime && s.Crossfade == nil {
		log.Printf("init crossfade")
		s.Crossfade = beep.Take(
			fadeTime,
			beep.Mix(
				effects.Transition(s.CurrentStream, fadeTime, 1.0, 0.0, effects.TransitionEqualPower),
				effects.Transition(s.NextStream, fadeTime, 0.0, 1.0, effects.TransitionEqualPower),
			))
	}
	if s.Crossfade != nil {
		n, ok = s.Crossfade.Stream(samples)
		if !ok {
			log.Printf("hand-off")
			s.CurrentStream.Seek(0)
			tmpStream := s.CurrentStream
			s.CurrentStream = s.NextStream
			s.NextStream = tmpStream
			s.Crossfade = nil

			if n < len(samples) {
				nextN, nextOk := s.CurrentStream.Stream(samples[n:])
				return n + nextN, nextOk
			}
			return s.CurrentStream.Stream(samples)
		}
		return n, true
	}
	return s.CurrentStream.Stream(samples)
}

func (s *StreamPlayer) Err() error {
	return nil
}
