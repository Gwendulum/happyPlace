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

type sessionState int

const (
	filepickerState sessionState = iota
	mixerState
)

type rootModel struct {
	state        sessionState
	filepicker   tea.Model
	mixer        tea.Model
	selectedFile string
	engine       *AudioEngine
}

func initialRootModel(fp filePickerModel, mm mixerModel, ae *AudioEngine) rootModel {
	return rootModel{
		state:      filepickerState,
		filepicker: fp,
		mixer:      mm,
		engine:     ae,
	}
}

type startInModeMsg struct{}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case startInModeMsg:
		m.state = filepickerState
		fmt.Printf("did this run?")
		return m, nil

	case fileSelectedMsg:
		m.selectedFile = msg.Path
		m.state = mixerState
		m.engine.Play(m.selectedFile)
	}

	switch m.state {
	case filepickerState:
		m.filepicker, cmd = m.filepicker.Update(msg)
		cmds = append(cmds, cmd)

	case mixerState:
		m.mixer, cmd = m.mixer.Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m rootModel) View() tea.View {
	switch m.state {
	case filepickerState:
		return m.filepicker.View()
	case mixerState:
		return m.mixer.View()
	}
	return tea.NewView("you're in root mode. What are you doing here?")
}

type filePickerModel struct {
	cursor int
	file   string
	engine *AudioEngine
}

func (fp filePickerModel) Init() tea.Cmd {
	return nil
}

func (fp filePickerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return fp, tea.Quit
		case "down", "j":
			if fp.cursor < len(files)-1 {
				fp.cursor++
			}
		case "up", "k":
			if fp.cursor > 0 {
				fp.cursor--
			}
		case "enter":
			fp.file = files[fp.cursor]
			return fp, func() tea.Msg {
				return fileSelectedMsg{Path: fp.file}
			}
		}

	}
	return fp, nil
}

func (fp filePickerModel) View() tea.View {
	s := strings.Builder{}
	s.WriteString("choose a file\n\n")

	for i := range files {

		s.WriteString(" ")

		if fp.cursor == i {
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

type mixerModel struct {
	cursor int
}

func (mm mixerModel) Init() tea.Cmd {
	return nil
}

func (mm mixerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {

	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return mm, tea.Quit
		case "down", "j":
			if mm.cursor < len(files)-1 {
				mm.cursor++
			}
		case "up", "k":
			if mm.cursor > 0 {
				mm.cursor--
			}
		}

	}
	return mm, nil
}

func (mm mixerModel) View() tea.View {

	s := strings.Builder{}
	s.WriteString("\nmixer mode\n")
	s.WriteString("\npress 'q' to quit \n")
	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
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

type fileSelectedMsg struct {
	Path string
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
	fp := filePickerModel{}
	mm := mixerModel{}
	p := tea.NewProgram(initialRootModel(fp, mm, &ae))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error here %v", err)
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
