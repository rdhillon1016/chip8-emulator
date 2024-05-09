package main

import (
	"os"
	"time"
	"flag"
	"github.com/rdhillon1016/chip8-emulator/chip8"
	"github.com/rdhillon1016/chip8-emulator/io/pixeladapter"
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
	filePath := flag.String("filePath", "./roms/Pong.ch8", "Location of ROM file (default is ./roms/Pong.ch8)")
	flag.Bool("usePixelEngine", false, "Use the github.com/gopxl/pixel engine instead of ebitengine (default)")

	flag.Parse()

	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		panic("Unable to read game file")
	}

	pixeladapter.Run(chip8.NewChip(fileBytes), executionRateHz)
}
