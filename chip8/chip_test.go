package chip8

import "testing"

const (
	testingCycleSleepTime = 0
)

func TestFetchOutOfBoundsInstruction(t *testing.T) {
	defer func() { _ = recover() }()

	chip := NewChip([]byte{0x00E0}, testingCycleSleepTime)
	chip.programCounter = 4095
	chip.ExecuteCycle()

	t.Errorf("did not panic")
}

func Test00E0(t *testing.T) {
	chip := NewChip([]byte{0x00, 0xE0}, testingCycleSleepTime)
	var expectedDisplay [64][32]bool
	for i, v := range chip.display {
		for j := range v {
			chip.display[i][j] = true
			expectedDisplay[i][j] = false
		}
	}
	chip.ExecuteCycle()

	for i, v := range chip.display {
		for j := range v {
			if chip.display[i][j] != expectedDisplay[i][j] {
				t.Errorf("Screen didn't clear %t, %t", chip.display[i][j], expectedDisplay[i][j])
			}
		}
	}
}

func Test00EE(t *testing.T) {
	chip := NewChip([]byte{0x22, 0x02, 0x00, 0xEE}, testingCycleSleepTime)
	chip.ExecuteCycle()
	chip.ExecuteCycle()

	if chip.stackPointer != 0 || chip.programCounter != 0x202 {
		t.Error("Subroutine return failed")
	}
}

func Test1NNN(t *testing.T) {
	chip := NewChip([]byte{0x1E, 0xEE}, testingCycleSleepTime)
	chip.ExecuteCycle()

	if chip.programCounter != 0xEEE {
		t.Error("Program counter didn't jump to correct address")
	}
}

func Test2NNN(t *testing.T) {
	chip := NewChip([]byte{0x2E, 0xEE}, testingCycleSleepTime)
	originalProgramCounter := chip.programCounter
	chip.ExecuteCycle()

	if chip.stackPointer != 1 || chip.stack[0] != originalProgramCounter+2 || chip.programCounter != 0xEEE {
		t.Error("Subroutine call failed")
	}
}

func Test3XNN(t *testing.T) {
	chip := NewChip([]byte{0x31, 0x45, 0x00, 0x00, 0x31, 0x46}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x41, 0x46, 0x00, 0x00, 0x41, 0x45}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x50, 0x10, 0x00, 0x00, 0x50, 0x20}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x60, 0x11}, testingCycleSleepTime)
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x11 {
		t.Error("Register load failed")
	}
}

func Test7XNN(t *testing.T) {
	chip := NewChip([]byte{0x70, 0xEE, 0x70, 0xEE}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x80, 0x10}, testingCycleSleepTime)
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE {
		t.Error("Register load failed")
	}
}

func Test8XY1(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x11}, testingCycleSleepTime)
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE {
		t.Error("Register load failed")
	}
}

func Test8XY2(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x12}, testingCycleSleepTime)
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x00 {
		t.Error("Register load failed")
	}
}

func Test8XY3(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x13}, testingCycleSleepTime)
	chip.generalRegisters[1] = 0xEE
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0xEE {
		t.Error("Register load failed")
	}
}

func Test8XY4(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x14, 0x80, 0x14}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x80, 0x15, 0x80, 0x25}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x80, 0x16}, testingCycleSleepTime)
	chip.generalRegisters[0] = 0x3
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x1 && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test8XY7(t *testing.T) {
	chip := NewChip([]byte{0x80, 0x17, 0x82, 0x17}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0x80, 0x1E}, testingCycleSleepTime)
	chip.generalRegisters[0] = 0x81
	chip.ExecuteCycle()

	if chip.generalRegisters[0] != 0x2 && chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Register load failed")
	}
}

func Test9XY0(t *testing.T) {
	chip := NewChip([]byte{0x90, 0x20, 0x00, 0x00, 0x90, 0x10}, testingCycleSleepTime)
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
	chip := NewChip([]byte{0xAE, 0xEE}, testingCycleSleepTime)
	chip.ExecuteCycle()

	if chip.indexRegister != 0xEEE {
		t.Error("Index register load failed")
	}
}

func TestBNNN(t *testing.T) {
	chip := NewChip([]byte{0xBE, 0xED}, testingCycleSleepTime)
	chip.generalRegisters[0] = 0x1
	chip.ExecuteCycle()

	if chip.programCounter != 0xEEE {
		t.Error("Jump with offset failed")
	}
}

func TestDXYN(t *testing.T) {
	chip := NewChip([]byte{0xD0, 0x15}, testingCycleSleepTime)
	xCord := byte(3)
	yCord := byte(4)

	for j := byte(0); j < 5; j++ {
		chip.display[xCord][yCord+j] = true
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
			if i == 7 && !chip.display[xCord+i][yCord+j] {
				t.Error("Draw failed")
			}
			if i != 7 && chip.display[xCord+i][yCord+j] {
				t.Error("Draw failed")
			}
		}
	}

	if chip.generalRegisters[flagRegisterIndex] != 1 {
		t.Error("Flag register wasn't set")
	}
}

func TestDXYNWrap(t *testing.T) {
	chip := NewChip([]byte{0xD0, 0x11}, testingCycleSleepTime)
	xCord := byte(64)
	yCord := byte(32)

	chip.generalRegisters[0] = xCord
	chip.generalRegisters[1] = yCord
	chip.indexRegister = 0x202
	chip.memory[chip.indexRegister] = 0xFF
	chip.ExecuteCycle()

	for xCord := 0; xCord < 8; xCord++ {
		if !chip.display[xCord][0] {
			t.Error("Draw failed")
		}
	}
}

func TestDXYNTruncate(t *testing.T) {
	chip := NewChip([]byte{0xD0, 0x11}, testingCycleSleepTime)
	xCord := byte(63)
	yCord := byte(31)

	chip.generalRegisters[0] = xCord
	chip.generalRegisters[1] = yCord
	chip.indexRegister = 0x202
	chip.memory[chip.indexRegister] = 0xFF
	chip.ExecuteCycle()

	if !chip.display[xCord][yCord] {
		t.Error("Draw failed")
	}

	// Check around it
	if chip.display[xCord - 1][yCord] || chip.display[xCord][yCord - 1] {
		t.Error("Pixels were set that shouldn't have been")
	}

	// Make sure it didn't wrap
	if chip.display[0][0] {
		t.Error("Pixels were set that shouldn't have been")
	}
}