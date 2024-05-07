package chip8

import (
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"math"
	"time"
)

const (
	flagRegisterIndex       = 15
	memoryStartIndexForFont = 0x50
	memoryStartIndexForGame = 0x200
	displayWidth            = 64
	displayHeight           = 32
)

var font = [80]byte{
	0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
	0x20, 0x60, 0x20, 0x20, 0x70, // 1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
	0x90, 0x90, 0xF0, 0x10, 0x10, // 4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
	0xF0, 0x10, 0x20, 0x40, 0x40, // 7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
	0xF0, 0x90, 0xF0, 0x90, 0x90, // A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
	0xF0, 0x80, 0x80, 0x80, 0xF0, // C
	0xE0, 0x90, 0x90, 0x90, 0xE0, // D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
	0xF0, 0x80, 0xF0, 0x80, 0x80, // F
}

var keysToIndexMap map[rune]uint = map[rune]uint{
	'X': 0x0,
	'1': 0x1,
	'2': 0x2,
	'3': 0x3,
	'Q': 0x4,
	'W': 0x5,
	'E': 0x6,
	'A': 0x7,
	'S': 0x8,
	'D': 0x9,
	'Z': 0xA,
	'C': 0xB,
	'4': 0xC,
	'R': 0xD,
	'F': 0xE,
	'V': 0xF,
}

type Chip struct {
	memory           [4096]byte
	programCounter   uint16
	indexRegister    uint16
	stack            [16]uint16
	stackPointer     int
	delayTimerValue  uint8
	soundTimerValue  uint8
	generalRegisters [16]byte
	display          [displayWidth][displayHeight]bool
	keys             [16]bool
	cycleSleepTime   time.Duration
}

func NewChip(fileBytes []byte, cycleSleepTime time.Duration) *Chip {
	chip := Chip{
		programCounter: 0x200,
		cycleSleepTime: cycleSleepTime,
	}
	chip.loadGameIntoMemory(fileBytes)
	chip.loadFontIntoMemory()
	return &chip
}

func (chip *Chip) ExecuteCycle() bool {
	instruction := chip.fetchInstruction()
	screenUpdated := chip.executeInstruction(instruction)
	time.Sleep(chip.cycleSleepTime)
	return screenUpdated
}

func (chip *Chip) fetchInstruction() uint16 {
	if chip.programCounter+1 > uint16(len(chip.memory)) {
		panic("Instruction is out of bounds")
	}
	currInstruction := binary.BigEndian.Uint16(chip.memory[chip.programCounter : chip.programCounter+2])
	chip.programCounter += 2
	return currInstruction
}

func (chip *Chip) executeInstruction(instruction uint16) bool {
	screenUpdated := false
	firstHexit := (instruction >> 12) & 0xF
	secondHexit := (instruction >> 8) & 0xF
	thirdHexit := (instruction >> 4) & 0xF
	fourthHexit := instruction & 0xF
	secondByteOfInstruction := byte(instruction)
	last12BitsOfInstruction := instruction & 0x0FFF

	switch firstHexit {
	case 0x0:
		switch instruction {
		case 0x00E0:
			for i, v := range chip.display {
				for j := range v {
					chip.display[i][j] = false
				}
			}
			screenUpdated = true
		case 0x00EE:
			chip.stackPointer--
			if chip.stackPointer < 0 {
				panic("Returning from main routine")
			}
			chip.programCounter = chip.stack[chip.stackPointer]
		default:
			panic(fmt.Sprintf("Cannot recognize instruction: %d", instruction))
		}
	case 0x1:
		chip.programCounter = last12BitsOfInstruction
	case 0x2:
		if chip.stackPointer > 15 {
			panic("Stack limit of 16 reached")
		} else {
			chip.stack[chip.stackPointer] = chip.programCounter
			chip.programCounter = last12BitsOfInstruction
			chip.stackPointer++
		}
	case 0x3:
		registerValue := chip.generalRegisters[secondHexit]
		if registerValue == byte(secondByteOfInstruction) {
			chip.programCounter += 2
		}
	case 0x4:
		registerValue := chip.generalRegisters[secondHexit]
		if registerValue != byte(secondByteOfInstruction) {
			chip.programCounter += 2
		}
	case 0x5:
		registerValueOne := chip.generalRegisters[secondHexit]
		registerValueTwo := chip.generalRegisters[thirdHexit]
		if registerValueOne == registerValueTwo {
			chip.programCounter += 2
		}
	case 0x6:
		chip.generalRegisters[secondHexit] = byte(secondByteOfInstruction)
	case 0x7:
		chip.generalRegisters[secondHexit] += byte(secondByteOfInstruction)
	case 0x8:
		switch fourthHexit {
		case 0x0:
			chip.generalRegisters[secondHexit] = chip.generalRegisters[thirdHexit]
		case 0x1:
			chip.generalRegisters[secondHexit] |= chip.generalRegisters[thirdHexit]
		case 0x2:
			chip.generalRegisters[secondHexit] &= chip.generalRegisters[thirdHexit]
		case 0x3:
			chip.generalRegisters[secondHexit] ^= chip.generalRegisters[thirdHexit]
		case 0x4:
			registerValueOne := chip.generalRegisters[secondHexit]
			registerValueTwo := chip.generalRegisters[thirdHexit]
			if registerValueOne > math.MaxUint8-registerValueTwo {
				chip.generalRegisters[flagRegisterIndex] = 1
			} else {
				chip.generalRegisters[flagRegisterIndex] = 0
			}
			chip.generalRegisters[secondHexit] += registerValueTwo
		case 0x5:
			registerValueOne := chip.generalRegisters[secondHexit]
			registerValueTwo := chip.generalRegisters[thirdHexit]
			if registerValueOne >= registerValueTwo {
				chip.generalRegisters[flagRegisterIndex] = 1
			} else {
				chip.generalRegisters[flagRegisterIndex] = 0
			}
			chip.generalRegisters[secondHexit] -= registerValueTwo
		case 0x6:
			registerValue := chip.generalRegisters[secondHexit]
			chip.generalRegisters[flagRegisterIndex] = registerValue & 0x1
			chip.generalRegisters[secondHexit] >>= 1
		case 0x7:
			registerValueOne := chip.generalRegisters[secondHexit]
			registerValueTwo := chip.generalRegisters[thirdHexit]
			if registerValueTwo >= registerValueOne {
				chip.generalRegisters[flagRegisterIndex] = 1
			} else {
				chip.generalRegisters[flagRegisterIndex] = 0
			}
			chip.generalRegisters[secondHexit] = registerValueTwo - registerValueOne
		case 0xE:
			registerValue := chip.generalRegisters[secondHexit]
			chip.generalRegisters[flagRegisterIndex] = registerValue & 0x80
			chip.generalRegisters[secondHexit] <<= 1
		}
	case 0x9:
		registerValueOne := chip.generalRegisters[secondHexit]
		registerValueTwo := chip.generalRegisters[thirdHexit]
		if registerValueOne != registerValueTwo {
			chip.programCounter += 2
		}
	case 0xA:
		chip.indexRegister = last12BitsOfInstruction
	case 0xB:
		chip.programCounter = uint16(chip.generalRegisters[0]) + last12BitsOfInstruction
	case 0xC:
		randomByteSlice := make([]byte, 1)
		_, err := rand.Read(randomByteSlice)
		if err != nil {
			panic("Random number generation failed")
		}
		chip.generalRegisters[secondHexit] = secondByteOfInstruction & randomByteSlice[0]
	case 0xD:
		height := fourthHexit
		startingX := chip.generalRegisters[secondHexit] % displayWidth
		startingY := chip.generalRegisters[thirdHexit] % displayHeight
		for j := byte(0); j < byte(height); j++ {
			currByte := chip.memory[chip.indexRegister+uint16(j)]
			currentY := startingY + j
			if currentY >= displayHeight {
				break
			}
			for i := byte(0); i < 8; i++ {
				currX := startingX + i
				if currX >= displayWidth {
					break
				}
				currPixel := chip.display[currX][currentY]
				var newPixel bool
				if ((currByte >> (8 - i - 1)) & 1) == 1 {
					newPixel = true
				} else {
					newPixel = false
				}
				if currPixel && newPixel {
					chip.generalRegisters[flagRegisterIndex] = 1
				}
				chip.display[currX][currentY] = currPixel != newPixel
			}
		}
		screenUpdated = true
	case 0xE:
		key := chip.generalRegisters[secondHexit]
		switch secondByteOfInstruction {
		case 0x9E:
			if keyPressed(key) {
				chip.programCounter += 2
			}
		case 0xA1:
			if !keyPressed(key) {
				chip.programCounter += 2
			}
		}
	case 0xF:
		switch secondByteOfInstruction {
		case 0x07:
			chip.generalRegisters[secondHexit] = chip.delayTimerValue
		case 0x0A:
			// wait for key press
		case 0x15:
			chip.delayTimerValue = chip.generalRegisters[secondHexit]
		case 0x18:
			chip.soundTimerValue = chip.generalRegisters[secondHexit]
		case 0x1E:
			chip.indexRegister += uint16(chip.generalRegisters[secondHexit])
		case 0x29:
			// something to do with sprites
		case 0x33:
			registerValue := chip.generalRegisters[secondHexit]

			hundredsDigit := (registerValue / 100) % 10
			tensDigit := (registerValue / 10) % 10
			onesDigit := (registerValue / 1) % 10

			chip.memory[chip.indexRegister] = hundredsDigit
			chip.memory[chip.indexRegister+1] = tensDigit
			chip.memory[chip.indexRegister+2] = onesDigit
		case 0x55:
			chip.dumpRegisters(secondHexit)
		case 0x65:
			chip.loadRegisters(secondHexit)
		}
	}
	return screenUpdated
}

func (chip *Chip) dumpRegisters(finalRegisterIndex uint16) {
	for i := 0; i <= int(finalRegisterIndex); i++ {
		chip.memory[int(chip.indexRegister)+i] = chip.generalRegisters[i]
	}
}

func (chip *Chip) loadRegisters(finalRegisterIndex uint16) {
	for i := 0; i <= int(finalRegisterIndex); i++ {
		chip.generalRegisters[i] = chip.memory[int(chip.indexRegister)+i]
	}
}

func (chip *Chip) DecrementDelayTimer() {
	if chip.delayTimerValue != 0 {
		chip.delayTimerValue--
	}
}

func (chip *Chip) DecrementSoundTimer() {
	if chip.soundTimerValue != 0 {
		chip.soundTimerValue--
	}
}

func keyPressed(key byte) bool {
	// TODO
	return true
}

func (chip *Chip) SetKeys(keyState [16]bool) {
	chip.keys = keyState
}

func (chip *Chip) loadGameIntoMemory(fileBytes []byte) {
	copy(chip.memory[memoryStartIndexForGame:memoryStartIndexForGame+len(fileBytes)], fileBytes)
}

func (chip *Chip) loadFontIntoMemory() {
	copy(chip.memory[memoryStartIndexForFont:memoryStartIndexForFont+len(font)], font[:])
}
