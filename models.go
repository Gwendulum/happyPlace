package main

import (
	tea "charm.land/bubbletea/v2"
	_ "github.com/hajimehoshi/oto/v2"

	"strings"
)

type sessionState int

const (
	filepickerState sessionState = iota
	mixerState
)

type rootModel struct {
	state      sessionState
	filepicker tea.Model
	mixer      tea.Model
}

func initialRootModel(fp filePickerModel, mm mixerModel) rootModel {
	return rootModel{
		state:      mixerState,
		filepicker: fp,
		mixer:      mm,
	}
}

func (m rootModel) Init() tea.Cmd {
	return nil
}

func (m rootModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case fileSelectedMsg:
		m.state = mixerState
		return m, func() tea.Msg {
			return createStreamMsg{Path: msg.Path}
		}
	case loadFileMsg:
		m.state = filepickerState
		return m, nil
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
				return fileSelectedMsg{Path: files[fp.cursor]}
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
	cursor   int
	selected int
	engine   *AudioEngine
}

func (mm mixerModel) Init() tea.Cmd {
	return nil
}

func (mm mixerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case createStreamMsg:
		mm.engine.createStream(msg.Path, mm.selected)
		mm.engine.Play()
		return mm, nil
	case tea.KeyPressMsg:

		switch msg.String() {
		case "ctrl+c", "q":
			return mm, tea.Quit
		case "down", "j":
			if mm.cursor < len(mm.engine.streamers)-1 {
				mm.cursor++
			}
		case "up", "k":
			if mm.cursor > 0 {
				mm.cursor--
			}
		case "enter":

			mm.selected = mm.cursor
			if mm.engine.streamers[mm.selected].loaded {
				mm.engine.streamers[mm.selected].Init()
				return mm, nil
			}
			return mm, func() tea.Msg {
				return loadFileMsg{}
			}
		case "s":
			if !mm.engine.initialized {
				return mm, nil
			}
			if !mm.engine.streamers[mm.cursor].loaded {
				return mm, nil
			}

			if mm.engine.streamers[mm.cursor].Paused {
				mm.engine.streamers[mm.cursor].Play()
			} else {
				mm.engine.streamers[mm.cursor].Stop()
			}
			return mm, nil

		}
	}
	return mm, nil
}

func (mm mixerModel) View() tea.View {

	s := strings.Builder{}
	s.WriteString("\nmixer mode\n")
	for i := 0; i < 4; i++ {
		if mm.cursor == i {
			s.WriteString("• ")
		}
		if !mm.engine.streamers[i].loaded {
			s.WriteString("[add file]\n")
		} else {
			s.WriteString(mm.engine.streamers[i].filepath)
			s.WriteString("\n")
		}
	}
	s.WriteString("\npress 'q' to quit \n")
	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}

type fileSelectedMsg struct {
	Path string
}

type loadFileMsg struct {
}

type createStreamMsg struct {
	Path string
}
type audioErrorMsg error
