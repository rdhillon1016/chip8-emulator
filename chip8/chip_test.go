package chip8

import "testing"

const (
	testingCycleSleepTime = 0
)

func TestFetchOutOfBoundsInstruction(t *testing.T) {
	defer func() { _ = recover() }()

	chip := NewChip([]byte{0x00E0})
	chip.programCounter = 4095
	chip.ExecuteCycle()

	t.Errorf("did not panic")
}

func Test00E0(t *testing.T) {
	chip := NewChip([]byte{0x00, 0xE0})
	var expectedDisplay [64][32]bool
	for i, v := range chip.Pixels {
		for j := range v {
			chip.Pixels[i][j] = true
			expectedDisplay[i][j] = false
		}
	}
	chip.ExecuteCycle()

	for i, v := range chip.Pixels {
		for j := range v {
			if chip.Pixels[i][j] != expectedDisplay[i][j] {
				t.Errorf("Screen didn't clear %t, %t", chip.Pixels[i][j], expectedDisplay[i][j])
			}
		}
	}
}

func Test00EE(t *testing.T) {
	chip := NewChip([]byte{0x22, 0x02, 0x00, 0xEE})
	chip.ExecuteCycle()
	chip.ExecuteCycle()

	if chip.stackPointer != 0 || chip.programCounter != 0x202 {
		t.Error("Subroutine return failed")
	}
}

func Test1NNN(t *testing.T) {
	chip := NewChip([]byte{0x1E, 0xEE})
	chip.ExecuteCycle()

	if chip.programCounter != 0xEEE {
		t.Error("Program counter didn't jump to correct address")
	}
}

func Test2NNN(t *testing.T) {
	chip := NewChip([]byte{0x2E, 0xEE})
	originalProgramCounter := chip.programCounter
	chip.ExecuteCycle()

	if chip.stackPointer != 1 || chip.stack[0] != originalProgramCounter+2 || chip.programCounter != 0xEEE {
		t.Error("Subroutine call failed")
	}
}

func Test3XNN(t *testing.T) {
	chip := NewChip([]byte{0x31, 0x45, 0x00, 0x00, 0x31, 0x46})
	chip.generalRegisters[1] = 0x45
	chip.ExecuteCycle()

	if chip.programCounter != 0x204 {
		t.Error("Instruction skip failed")
	}

	chip.ExecuteCycle()

	if chip.programCounter != 0x206 {
		t.Error("PC didn't increment correctly")
	}
}

func Test4XNN(t *testing.T) {
	chip := NewChip([]byte{0x41, 0x46, 0x00, 0x00, 0x41, 0x45})
	chip.generalRegisters[1] = 0x45
	chip.ExecuteCycle()

	if chip.programCounter != 0x204 {
		t.Error("Instruction skip failed")
	}

	chip.ExecuteCycle()

	if chip.programCounter != 0x206 {
		t.Error("PC didn't increment correctly")
	}
}

func Test5XY0(t *testing.T) {
	chip := NewChip([]byte{0x50, 0x10, 0x00, 0x00, 0x50, 0x20})
	chip.generalRegisters[0] = 0x45
	chip.generalRegisters[1] = 0x45
	chip.ExecuteCycle()

	if chip.programCounter != 0x204 {
		t.Error("Instruction skip failed")
	}

	chip.ExecuteCycle()

	if chip.programCounter != 0x206 {
		t.Error("PC didn't increment correctly")
	}
}

func Test6XNN(t *testing.T) {
	chip := NewChip([]byte{0x60, 0x11})
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x11 {
		t.Error("Register load failed")
	}
}

func Test7XNN(t *testing.T) {
	chip := NewChip([]byte{0x70, 0xEE, 0x70, 0xEE})
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE || chip.generalRegisters[flagRegisterIndex] == 1 {
		t.Error("Register load failed")
	}

	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xDC || chip.generalRegisters[flagRegisterIndex] == 1 {
		t.Error("Register load failed")
	}
}

func Test8XY0(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x10})
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE {
		t.Error("Register load failed")
	}
}

func Test8XY1(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x11})
	chip.generalRegisters[1] = 0xEE
	chip.generalRegisters[flagRegisterIndex] = 1
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE || chip.generalRegisters[flagRegisterIndex] != 0 {
		t.Error("Register load failed")
	}
}

func Test8XY2(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x12})
	chip.generalRegisters[1] = 0xEE
	chip.generalRegisters[flagRegisterIndex] = 1
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x00 || chip.generalRegisters[flagRegisterIndex] != 0 {
		t.Error("Register load failed")
	}
}

func Test8XY3(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x13})
	chip.generalRegisters[1] = 0xEE
	chip.generalRegisters[flagRegisterIndex] = 1
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE || chip.generalRegisters[flagRegisterIndex] != 0 {
		t.Error("Register load failed")
	}
}

func Test8XY4(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x14, 0x80, 0x14})
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE && chip.generalRegisters[flagRegisterIndex] != 0 {
		t.Error("Register load failed")
	}

	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xDC && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test8XY5(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x15, 0x80, 0x25})
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE && chip.generalRegisters[flagRegisterIndex] != 0 {
		t.Error("Register load failed")
	}

	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test8XY6(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x16})
	chip.generalRegisters[1] = 0x3
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x1 && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test8XY7(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x17, 0x82, 0x17})
	chip.generalRegisters[0] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[1] != 0xEE && chip.generalRegisters[flagRegisterIndex] != 0 {
		t.Error("Register load failed")
	}

	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test8XYE(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x1E})
	chip.generalRegisters[1] = 0x81
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x2 && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test9XY0(t *testing.T) {
	chip := NewChip([]byte{0x90, 0x20, 0x00, 0x00, 0x90, 0x10})
	chip.generalRegisters[0] = 0x45
	chip.generalRegisters[1] = 0x45
	chip.ExecuteCycle()

	if chip.programCounter != 0x204 {
		t.Error("Instruction skip failed")
	}

	chip.ExecuteCycle()

	if chip.programCounter != 0x206 {
		t.Error("PC didn't increment correctly")
	}
}

func TestANNN(t *testing.T) {
	chip := NewChip([]byte{0xAE, 0xEE})
	chip.ExecuteCycle()

	if chip.indexRegister != 0xEEE {
		t.Error("Index register load failed")
	}
}

func TestBNNN(t *testing.T) {
	chip := NewChip([]byte{0xBE, 0xED})
	chip.generalRegisters[0] = 0x1
	chip.ExecuteCycle()

	if chip.programCounter != 0xEEE {
		t.Error("Jump with offset failed")
	}
}

func TestDXYN(t *testing.T) {
	chip := NewChip([]byte{0xD0, 0x15})
	xCord := byte(3)
	yCord := byte(4)

	for j := byte(0); j < 5; j++ {
		chip.Pixels[xCord][yCord+j] = true
	}

	chip.generalRegisters[0] = xCord
	chip.generalRegisters[1] = yCord
	chip.indexRegister = 0x202
	for i := uint16(0); i < 5; i++ {
		chip.memory[chip.indexRegister+i] = 0x81
	}
	chip.ExecuteCycle()

	for j := byte(0); j < 5; j++ {
		for i := byte(0); i < 8; i++ {
			if i == 7 && !chip.Pixels[xCord+i][yCord+j] {
				t.Error("Draw failed")
			}
			if i != 7 && chip.Pixels[xCord+i][yCord+j] {
				t.Error("Draw failed")
			}
		}
	}

	if chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Flag register wasn't set")
	}
}

func TestDXYNWrap(t *testing.T) {
	chip := NewChip([]byte{0xD0, 0x11})
	xCord := byte(64)
	yCord := byte(32)

	chip.generalRegisters[0] = xCord
	chip.generalRegisters[1] = yCord
	chip.indexRegister = 0x202
	chip.memory[chip.indexRegister] = 0xFF
	chip.ExecuteCycle()

	for xCord := 0; xCord < 8; xCord++ {
		if !chip.Pixels[xCord][0] {
			t.Error("Draw failed")
		}
	}
}

func TestDXYNTruncate(t *testing.T) {
	chip := NewChip([]byte{0xD0, 0x11})
	xCord := byte(63)
	yCord := byte(31)

	chip.generalRegisters[0] = xCord
	chip.generalRegisters[1] = yCord
	chip.indexRegister = 0x202
	chip.memory[chip.indexRegister] = 0xFF
	chip.ExecuteCycle()

	if !chip.Pixels[xCord][yCord] {
		t.Error("Draw failed")
	}

	// Check around it
	if chip.Pixels[xCord-1][yCord] || chip.Pixels[xCord][yCord-1] {
		t.Error("Pixels were set that shouldn't have been")
	}

	// Make sure it didn't wrap
	if chip.Pixels[0][0] {
		t.Error("Pixels were set that shouldn't have been")
	}
}

func TestEX9E(t *testing.T) {
	chip := NewChip([]byte{0xE0, 0x9E, 0x00, 0x00, 0xE1, 0x9E})
	chip.keys[0] = true
	chip.generalRegisters[1] = 1
	chip.ExecuteCycle()

	if chip.programCounter != 0x204 {
		t.Error("Skip failed")
	}

	chip.ExecuteCycle()

	if chip.programCounter != 0x206 {
		t.Error("PC increment failed")
	}
}

func TestEXA1(t *testing.T) {
	chip := NewChip([]byte{0xE0, 0xA1, 0x00, 0x00, 0xE1, 0xA1})
	chip.keys[1] = true
	chip.generalRegisters[1] = 1
	chip.ExecuteCycle()

	if chip.programCounter != 0x204 {
		t.Error("Skip failed")
	}

	chip.ExecuteCycle()

	if chip.programCounter != 0x206 {
		t.Error("PC increment failed")
	}
}

func TestFX07(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x07})
	chip.delayTimerValue = 5
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 5 {
		t.Error("Timer value to register failed")
	}
}

func TestFX0A(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x0A})
	chip.ExecuteCycle()

	if chip.programCounter != 0x200 {
		t.Error("Program did not halt while waiting for key press")
	}

	chip.keys[5] = true
	chip.ExecuteCycle()

	if chip.programCounter == 0x202 {
		t.Error("Program is not waiting for key release")
	}

	chip.keys[5] = false
	chip.ExecuteCycle()

	if chip.programCounter != 0x202 {
		t.Error("Program did not continue after key release")
	}
}

func TestFX15(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x15})
	chip.generalRegisters[0] = 5
	chip.ExecuteCycle()

	if chip.delayTimerValue != 5 {
		t.Error("Register to delay timer value failed")
	}
}

func TestFX18(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x18})
	chip.generalRegisters[0] = 5
	chip.ExecuteCycle()

	if chip.SoundTimerValue != 5 {
		t.Error("Register to sound timer value failed")
	}
}

func TestFX1E(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x1E})
	chip.indexRegister = 0xFF
	chip.generalRegisters[0] = 0xFF
	chip.ExecuteCycle()

	if chip.indexRegister != 510 {
		t.Error("Index register not added to properly")
	}
}

func TestFX29(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x29})
	chip.generalRegisters[0] = 0x2
	chip.ExecuteCycle()

	if chip.indexRegister != 0x5A {
		t.Error("Incorrect sprite location was written to index register")
	}
}

func TestFX33(t *testing.T) {
	chip := NewChip([]byte{0xF0, 0x33})
	chip.generalRegisters[0] = 173
	chip.indexRegister = 0x202
	chip.ExecuteCycle()

	if chip.memory[0x202] != 1 || chip.memory[0x202+1] != 7 || chip.memory[0x202+2] != 3 {
		t.Error("Incorrect binary-coded decimal representation")
	}
}

func TestFX55(t *testing.T) {
	chip := NewChip([]byte{0xF5, 0x55})
	for i := range chip.generalRegisters {
		chip.generalRegisters[i] = 0xDE
	}
	chip.indexRegister = 0x202
	chip.ExecuteCycle()

	if chip.indexRegister == 0x202 {
		t.Error("index register was not modified")
	}

	for i := 0; i < 6; i++ {
		if chip.memory[chip.indexRegister-uint16(i)-1] != 0xDE {
			t.Error("Error on register dump")
			break
		}
	}

	if chip.memory[chip.indexRegister] != 0x0 {
		t.Error("Memory location was modified")
	}
}

func TestFX65(t *testing.T) {
	chip := NewChip([]byte{0xF5, 0x65})
	chip.indexRegister = 0x202
	for i := 0; i < 6; i++ {
		chip.memory[chip.indexRegister+uint16(i)] = 0xDE
	}
	chip.ExecuteCycle()

	if chip.indexRegister == 0x202 {
		t.Error("index register was not modified")
	}

	for i := 0; i < 6; i++ {
		if chip.generalRegisters[i] != 0xDE {
			t.Error("Error on register load")
			break
		}
	}
}
