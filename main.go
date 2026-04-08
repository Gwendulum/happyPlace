package main

import (
	tea "charm.land/bubbletea/v2"
	"fmt"
	"github.com/gopxl/beep/v2"
	_ "github.com/hajimehoshi/oto/v2"

	"os"
)

var files = []string{"test1.wav", "testshort.wav"}

func main() {
	mixer := beep.Mixer{}

	ae := AudioEngine{
		mixer:     &mixer,
		streamers: make([]*StreamPlayer, 4),
	}

	for i := range ae.streamers {
		ae.streamers[i] = &StreamPlayer{}
	}
	fp := filePickerModel{}
	mm := mixerModel{
		engine: &ae,
	}
	p := tea.NewProgram(initialRootModel(fp, mm))
	_, err := p.Run()
	if err != nil {
		fmt.Printf("error here %v\n", err)
		os.Exit(1)
	}

}
