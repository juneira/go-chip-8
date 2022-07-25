package chip8

import (
	"fmt"
	"io"
)

const romAddressOffset = 0x200
const fontAddressOffset = 0x0

// Sprite of fonts
var fonts = []byte{
	// 0
	0xFA, 0x90, 0x90, 0x90, 0xFA,
	// 1
	0x20, 0x60, 0x20, 0x20, 0x70,
	// 2
	0xF0, 0x10, 0xF0, 0x80, 0xF0,
	// 3
	0xF0, 0x10, 0xF0, 0x10, 0xF0,
	// 4
	0x90, 0x90, 0xF0, 0x10, 0x10,
	// 5
	0xF0, 0x80, 0xF0, 0x10, 0xF0,
	// 6
	0xF0, 0x80, 0xF0, 0x90, 0xF0,
	// 7
	0xF0, 0x10, 0x20, 0x40, 0x40,
	// 8
	0xF0, 0x90, 0xF0, 0x90, 0xF0,
	// 9
	0xF0, 0x90, 0xF0, 0x10, 0xF0,
	// A
	0xF0, 0x90, 0xF0, 0x90, 0x90,
	// B
	0xE0, 0x90, 0xE0, 0x90, 0xE0,
	// C
	0xF0, 0x80, 0x80, 0x80, 0xF0,
	// D
	0xE0, 0x90, 0x90, 0x90, 0xE0,
	// E
	0xF0, 0x80, 0xF0, 0x80, 0xF0,
	// F
	0xF0, 0x80, 0xF0, 0x80, 0x80,
}

// StardardMemory implements interface Memory
type StardardMemory struct {
	mem [0xFFF]byte
	log io.Writer
}

type ConfigMemory struct {
	Rom *Rom
	Log io.Writer
}

// NewRom is a function that receive a "data" as param and return a pointer to Rom
func NewStardardMemory(config *ConfigMemory) *StardardMemory {
	sm := &StardardMemory{log: config.Log}
	sm.loadFonts()
	sm.loadGame(config.Rom)

	return sm
}

func (sm *StardardMemory) loadFonts() {
	for i := 0; i < len(fonts); i++ {
		sm.mem[fontAddressOffset+i] = fonts[i]
	}
}

func (sm *StardardMemory) loadGame(rom *Rom) {
	i := 0
	instr := rom.NextInstruction()

	for instr != nil {
		sm.mem[romAddressOffset+i] = []byte(*instr)[0]
		i++
		sm.mem[romAddressOffset+i] = []byte(*instr)[1]
		i++

		instr = rom.NextInstruction()
	}
}

// Log writes values of memory to "log" of Memory
func (sm *StardardMemory) Log() {
	sm.log.Write([]byte(fmt.Sprintf("memory: %v\n", sm.mem)))
}
