package main

import (
	"os"
	"time"

	"github.com/gopxl/pixel/pixelgl"
	"github.com/rdhillon1016/chip8-emulator/chip8"
	"github.com/rdhillon1016/chip8-emulator/io"
)

const (
	executionRateHz = 700
	timerRateHz     = 60
)

type Game struct {

}

func run() {
	args := os.Args[1:]

	fileBytes, err := os.ReadFile(args[0])
	if err != nil {
		panic("Unable to read game file")
	}

	display := io.NewDisplay()

	chip := chip8.NewChip(fileBytes)

	/* Note that the tickers are unnecessary when their corresponding values
	are 0, and thus can sometimes be wasteful. However, since they only
	tick at a rate of 60Hz, this is a fine tradeoff for now. A better
	solution may be to pause the ticker when its corresponding value is 0. */
	delayTicker := time.NewTicker(time.Second / timerRateHz)
	soundTicker := time.NewTicker(time.Second / timerRateHz)

	for !display.Window.Closed() {
		chip.SetKeys(display.GetKeyPresses())
		var screenUpdated bool
		select {
		case <-delayTicker.C:
			chip.DecrementDelayTimer()
			screenUpdated = chip.ExecuteCycle()
		case <-soundTicker.C:
			chip.DecrementSoundTimer()
			screenUpdated = chip.ExecuteCycle()
		default:
			screenUpdated = chip.ExecuteCycle()
		}
		if screenUpdated {
			display.UpdateScreen(chip.Pixels)
		} else {
			display.Window.UpdateInput()
		}
		time.Sleep(time.Second / executionRateHz)
	}
}

func main() {
	pixelgl.Run(run)
}
