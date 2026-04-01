package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	_ "github.com/hajimehoshi/oto/v2"
	"os"
	"strings"
	"sync"
	"time"
)

var files = []string{"test1.wav", "testshort.wav"}

type model struct {
	cursor int
	file   string
	engine *AudioEngine
}

func initialModel(ae *AudioEngine) model {
	return model{
		engine: ae,
	}
}

type audioErrorMsg error

func playFileCmd(ae *AudioEngine, path string) tea.Cmd {
	return func() tea.Msg {
		err := ae.Play(path)
		if err != nil {
			return audioErrorMsg(err)
		}
		return "success"
	}
}

func (m model) Init() tea.Cmd {
	return nil
}
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "down", "j":
			if m.cursor < len(files)-1 {
				m.cursor++
			}
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "enter":
			m.file = files[m.cursor]
			return m, playFileCmd(m.engine, m.file)

		}

	}
	return m, nil
}

func (m model) View() tea.View {
	s := strings.Builder{}
	s.WriteString("choose a file\n\n")

	for i := range files {

		s.WriteString(" ")

		if m.cursor == i {
			s.WriteString("• ")
		}
		s.WriteString(files[i])
		s.WriteString("\n")
	}

	s.WriteString("\npress 'q' to quit \n")
	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}

func main() {
	mixer := beep.Mixer{}

	ctrl1 := effects.Volume{
		Streamer: &mixer,
		Base:     10.0,
		Volume:   0.0,
		Silent:   false,
	}

	ae := AudioEngine{
		mixer:      &mixer,
		volumeCtrl: &ctrl1,
	}
	/*TODO: add function for creating new streams from file and adding that file to a registry.
	Registry will ensure there are no duplicate streams(though supports multiple streams from the same file), and sets a max number of streams running at a time.
	*/

	//TODO: add "smooth volume" for each stream that waits for next sample. Will later be slider.
	//TODO: add oscillating volume for each stream.
	//TODO: add "stop stream" functionality that stops and removes streams from .Mixer and the registry.
	p := tea.NewProgram(initialModel(&ae))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error %v", err)
		os.Exit(1)
	}

}

func createStream(filepath string) StreamPlayer {
	s := StreamPlayer{}
	s.loadFile(filepath)
	return s
}

type AudioEngine struct {
	mixer       *beep.Mixer
	Active      map[string]*StreamPlayer
	initialized bool
	volumeCtrl  *effects.Volume
	mu          sync.RWMutex
}

func (ae *AudioEngine) Play(filepath string) error {
	ae.mu.Lock()
	defer ae.mu.Unlock()

	s := createStream(filepath)

	if !ae.initialized {
		speaker.Init(s.Format.SampleRate, s.Format.SampleRate.N(time.Second/10))
		speaker.Play(ae.volumeCtrl)
		ae.initialized = true
	}

	ae.mixer.Add(&s)
	return nil
}
