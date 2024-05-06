package main

import (
	"encoding/binary"
	"math"
	"time"
)

const (
	flagRegisterIndex = 15
	executionRateHz   = 700
	timerRateHz       = 60
)

var font = [...]byte{
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
	memory           []byte
	programCounter   uint16
	indexRegister    uint16
	stack            [16]uint16
	stackPointer     uint16
	delayTimerValue  uint8
	soundTimerValue  uint8
	generalRegisters [16]byte
	display          [64][32]bool
	key              [16]bool
}

func (chip *Chip) Run(filePath string) {
	chip.loadGameIntoMemory(filePath)
	chip.loadFontIntoMemory()

	/* Note that the tickers are unnecessary when their corresponding values
	are 0, and thus can sometimes be wasteful. However, since they only
	tick at a rate of 60Hz, this is a fine tradeoff for now. A better
	solution may be to pause the ticker when its corresponding value is 0. */
	delayTicker := time.NewTicker(time.Second / timerRateHz)
	soundTicker := time.NewTicker(time.Second / timerRateHz)

	for {
		select {
		case <-delayTicker.C:
			if chip.delayTimerValue != 0 {
				chip.delayTimerValue--
			}
			chip.executeCycle()
		case <-soundTicker.C:
			if chip.soundTimerValue != 0 {
				chip.soundTimerValue--
			}
			chip.executeCycle()
		default:
			chip.executeCycle()
		}
	}
}

func (chip *Chip) executeCycle() {
	instruction := chip.fetchInstruction()
	chip.executeInstruction(instruction)
	time.Sleep(time.Second / executionRateHz)
}

func (chip *Chip) fetchInstruction() uint16 {
	currInstruction := binary.BigEndian.Uint16(chip.memory[chip.programCounter : chip.programCounter+2])
	chip.programCounter += 2
	return currInstruction
}

func (chip *Chip) executeInstruction(instruction uint16) {
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
			// Clear screen
		case 0x00EE:
			// Return from subroutine
		default:
			// Calls machine code routine
		}
	case 0x1:
		// jump to address
	case 0x2:
		// calls routine
	case 0x3:
		registerValue := chip.generalRegisters[secondHexit]
		if (registerValue == byte(secondByteOfInstruction)) {
			chip.programCounter += 2
		}
	case 0x4:
		registerValue := chip.generalRegisters[secondHexit]
		if (registerValue != byte(secondByteOfInstruction)) {
			chip.programCounter += 2
		}
	case 0x5:
		registerValueOne := chip.generalRegisters[secondHexit]
		registerValueTwo := chip.generalRegisters[thirdHexit]
		if (registerValueOne == registerValueTwo) {
			chip.programCounter += 2
		}
	case 0x6:
		chip.generalRegisters[secondHexit] = byte(secondByteOfInstruction)
	case 0x7:
		chip.generalRegisters[secondHexit] += byte(secondByteOfInstruction)
	case 0x8:
		switch (fourthHexit) {
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
			if (registerValueOne > math.MaxUint8 - registerValueTwo) {
				chip.generalRegisters[flagRegisterIndex] = 1
			} else {
				chip.generalRegisters[flagRegisterIndex] = 0
			}
			chip.generalRegisters[secondHexit] += registerValueTwo
		case 0x5:
			registerValueOne := chip.generalRegisters[secondHexit]
			registerValueTwo := chip.generalRegisters[thirdHexit]
			if (registerValueOne >= registerValueTwo) {
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
			if (registerValueTwo >= registerValueOne) {
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
		if (registerValueOne != registerValueTwo) {
			chip.programCounter += 2
		}
	case 0xA:
		chip.indexRegister = last12BitsOfInstruction
	case 0xB:
		chip.programCounter = uint16(chip.generalRegisters[0]) + last12BitsOfInstruction
	case 0xC:
		
	case 0xD:
	case 0xE:
	case 0xF:
	}
}

func (chip *Chip) loadGameIntoMemory(filePath string) {}

func (chip *Chip) loadFontIntoMemory() {}
