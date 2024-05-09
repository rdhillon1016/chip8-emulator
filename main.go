package main

import (
	"flag"
	"os"

	"github.com/rdhillon1016/chip8-emulator/chip8"
	"github.com/rdhillon1016/chip8-emulator/io"
)

func main() {
	filePath := flag.String("filePath", "./roms/Tetris.ch8", "Location of ROM file (default is ./roms/Tetris.ch8)")
	executionRateHz := flag.Int("executionRate", 700, "Execution rate of the chip in Hz (default is 700)")

	flag.Parse()

	fileBytes, err := os.ReadFile(*filePath)
	if err != nil {
		panic("Unable to read game file")
	}

	io.Run(chip8.NewChip(fileBytes), *executionRateHz)
}
