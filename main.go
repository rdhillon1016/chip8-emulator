package main

import (
	"os"
	"time"

	"github.com/rdhillon1016/chip8-emulator/chip8"
)

const (
	executionRateHz   = 700
	timerRateHz       = 60
)

func main() {
	args := os.Args[1:]

	fileBytes, err := os.ReadFile(args[0])
	if err != nil {
		panic("Unable to read game file")
	}
	chip := chip8.NewChip(fileBytes, time.Second / executionRateHz)
	
	/* Note that the tickers are unnecessary when their corresponding values
	are 0, and thus can sometimes be wasteful. However, since they only
	tick at a rate of 60Hz, this is a fine tradeoff for now. A better
	solution may be to pause the ticker when its corresponding value is 0. */
	delayTicker := time.NewTicker(time.Second / timerRateHz)
	soundTicker := time.NewTicker(time.Second / timerRateHz)

	for {
		select {
		case <-delayTicker.C:
			chip.DecrementDelayTimer()
			chip.ExecuteCycle()
		case <-soundTicker.C:
			chip.DecrementSoundTimer()
			chip.ExecuteCycle()
		default:
			chip.ExecuteCycle()
		}
	}
}