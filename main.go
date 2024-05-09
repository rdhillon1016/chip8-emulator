package main

import (
	"os"
	"flag"
	"github.com/rdhillon1016/chip8-emulator/chip8"
	"github.com/rdhillon1016/chip8-emulator/io/pixeladapter"
	"github.com/rdhillon1016/chip8-emulator/io/ebitadapter"
)

func main() {
	filePath := flag.String("filePath", "./roms/Pong.ch8", "Location of ROM file (default is ./roms/Pong.ch8)")
	usePixelEngine := flag.Bool("usePixelEngine", false, "Use the github.com/gopxl/pixel engine instead of ebitengine (default)")
	executionRateHz := flag.Int("executionRate", 700, "Execution rate of the chip in Hz (default is 700)")

	flag.Parse()

	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		panic("Unable to read game file")
	}

	if (*usePixelEngine) {
		pixeladapter.Run(chip8.NewChip(fileBytes), *executionRateHz)
	} else {
		ebitadapter.Run(chip8.NewChip(fileBytes), *executionRateHz)
	}
	ebitadapter.Run(chip8.NewChip(fileBytes), *executionRateHz)
}
