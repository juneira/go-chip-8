package chip8_test

import chip8 "github.com/MarceloMPJR/go-chip-8"

type MockDisplay struct {
	clearCount int
	drawCount  int
}

func (md *MockDisplay) Clear() {
	md.clearCount++
}

func (md *MockDisplay) Draw(xDisplay, yDisplay, sprite byte) bool {
	md.drawCount++

	return false
}

func (md *MockDisplay) Flush() {
}

type MockKeyBoard struct {
	Key chip8.Key
}

func (mk MockKeyBoard) KeyDown() chip8.Key {
	return mk.Key
}

type MockMemory struct {
	saveCount     int
	saveBCDCount  int
	loadCount     int
	loadCharCount int
}

func (mm *MockMemory) Save(register []byte, i uint16) {
	mm.saveCount++
}

func (mm *MockMemory) SaveBCD(vx byte, i uint16) {
	mm.saveBCDCount++
}

func (mm *MockMemory) Load(register []byte, i uint16) {
	mm.loadCount++
}

func (mm *MockMemory) LoadInstruction(pc uint16) chip8.Instruction {
	return chip8.Instruction{0x0, 0x0}
}

func (mm *MockMemory) LoadChar(vx byte) uint16 {
	mm.loadCharCount++
	return uint16(vx) + 0x2
}

func (mm *MockMemory) LoadSprite(i uint16) byte {
	return 0x00
}
