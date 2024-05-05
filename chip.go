package main

type Chip struct {
	memory []byte
	display int
	programCounter uint16
	indexRegister uint16
	stack []byte
	delayTimer int
	soundTimer int
	generalRegisters [16]byte
}