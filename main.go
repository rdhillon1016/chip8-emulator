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
)

type VirtualHardware struct {
	Chip *chip8.Chip
	DelayTicker *time.Ticker
	SoundTicker *time.Ticker
}

type gameEngine interface {
	Run()
}

func main() {
	args := os.Args[1:]

	fileBytes, err := os.ReadFile(args[0])
	if err != nil {
		panic("Unable to read game file")
	}

	/* Note that the tickers are unnecessary when their corresponding values
	are 0, and thus can sometimes be wasteful. However, since they only
	tick at a rate of 60Hz, this is a fine tradeoff for now. A better
	solution may be to pause the ticker when its corresponding value is 0. */

	pixelgl.Run(func() {
		io.RunWithPixel(chip8.NewChip(fileBytes), executionRateHz)
	})
}
