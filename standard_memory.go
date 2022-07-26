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

// StandardMemory implements interface Memory
type StandardMemory struct {
	mem [0xFFF]byte
	log io.Writer
}

type ConfigMemory struct {
	Rom *Rom
	Log io.Writer
}

// NewRom is a function that receive a "data" as param and return a pointer to Rom
func NewStandardMemory(config *ConfigMemory) *StandardMemory {
	sm := &StandardMemory{log: config.Log}
	sm.loadFonts()
	sm.loadGame(config.Rom)

	return sm
}

func (sm *StandardMemory) loadFonts() {
	for i := 0; i < len(fonts); i++ {
		sm.mem[fontAddressOffset+i] = fonts[i]
	}
}

func (sm *StandardMemory) loadGame(rom *Rom) {
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
func (sm *StandardMemory) Log() {
	sm.log.Write([]byte(fmt.Sprintf("memory: %v\n", sm.mem)))
}

// SaveBCD convert vx byte to decimal and save each digit on I, I+1 and I+2
func (sm *StandardMemory) SaveBCD(vx byte, i uint16) {
	sm.mem[i] = vx / 100
	sm.mem[i+1] = (vx % 100) / 10
	sm.mem[i+2] = (vx % 10)
}

// Save saves the registers on memory starting on register I
func (sm *StandardMemory) Save(register []byte, i uint16) {
	for idx, reg := range register {
		sm.mem[int(i)+idx] = reg
	}
}

// Load loads to the register from of memory starting on register I
func (sm *StandardMemory) Load(register []byte, i uint16) {
	for idx := 0; idx < len(register); idx++ {
		register[idx] = sm.mem[int(i)+idx]
	}
}

// LoadInstruction returns the instruction addressed by register PC
func (sm *StandardMemory) LoadInstruction(pc uint16) Instruction {
	return Instruction{sm.mem[pc], sm.mem[pc+1]}
}

// LoadChar returns the address to char VX
func (sm *StandardMemory) LoadChar(vx byte) uint16 {
	if vx > 0xF {
		return 0
	}

	return uint16(vx * 5)
}

// LoadSprit returns the sprite on position I
func (sm *StandardMemory) LoadSprite(i uint16) byte {
	return sm.mem[i]
}
