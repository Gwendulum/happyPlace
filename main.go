package main

import (
	"bufio"
	"fmt"
	"os"
	"time"

	_ "charm.land/bubbletea/v2"
	"github.com/gopxl/beep/effects"
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	_ "github.com/hajimehoshi/oto/v2"
)

func main() {
	s1 := StreamPlayer{}
	s1.loadFile("test1.wav")

	s2 := StreamPlayer{}
	s2.loadFile("testshort.wav")
	mixer := beep.Mixer{}

	ctrl1 := effects.Volume{
		Streamer: &mixer,
		Base:     10.0,
		Volume:   -0.5,
		Silent:   false,
	}
	fmt.Println(time.Second * 15)
	speaker.Init(s1.Format.SampleRate, s1.Format.SampleRate.N(time.Second/10))
	speaker.Play(&ctrl1)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Printf(">")
		if !scanner.Scan() {
			break
		}
		input := scanner.Text()

		if input == "add" {
			s := createStream("testshort.wav")
			mixer.Add(&s)
			input = scanner.Text()
		}
		/*TODO: add function for creating new streams from file and adding that file to a registry.
		Registry will ensure there are no duplicate streams(though supports multiple streams from the same file), and sets a max number of streams running at a time.
		*/

		if input == "stop" {
			//TODO: add "stop stream" functionality that stops and removes streams from .Mixer and the registry.
		}
		if input == "volume" {
			//TODO: add "smooth volume" for each stream that waits for next sample. Will later be slider.
			ctrl1.Volume += .1
		}
		if input == "osc" {
			//TODO: add oscillating volume for each stream.
		}
	}
}

func createStream(filepath string) StreamPlayer {
	s := StreamPlayer{}
	s.loadFile(filepath)
	return s
}
