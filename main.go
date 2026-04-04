package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/v2"
	_ "github.com/hajimehoshi/oto/v2"

	"os"
)

var files = []string{"test1.wav", "testshort.wav"}

func main() {
	mixer := beep.Mixer{}

	vol := effects.Volume{
		Streamer: &mixer,
		Base:     10.0,
		Volume:   0.0,
		Silent:   false,
	}

	ctrl := beep.Ctrl{
		Streamer: &vol,
	}
	ae := AudioEngine{
		mixer:      &mixer,
		volumeCtrl: &vol,
		pauseCtrl:  &ctrl,
	}
	fp := filePickerModel{}
	mm := mixerModel{
		stream: make([]streamData, 4),
	}
	p := tea.NewProgram(initialRootModel(fp, mm, &ae))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error here %v", err)
		os.Exit(1)
	}

}
