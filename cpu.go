package chip8

import (
	"fmt"
	"io"
	"math/rand"
)

type Register [0x10]byte
type Stack [0x10]uint16

type Cpu struct {
	register Register
	stack    Stack
	output   io.Writer
	sp       byte
	pc       uint16
	i        uint16
}

// NewCpu is a function that receive a "output" as param and return a pointer to Cpu
func NewCpu(register Register, stack Stack, output io.Writer, pc, i uint16, sp byte) *Cpu {
	return &Cpu{register: register, stack: stack, output: output, pc: pc, i: i, sp: sp}
}

// Log is a function that write values of registers to "output" of Cpu
func (c *Cpu) Log() {
	str := fmt.Sprintf("pc = %x\nsp = %x\ni = %x\nstack = %v\n", c.pc, c.sp, c.i, c.stack)

	for i := 0; i < len(c.register); i++ {
		str += fmt.Sprintf("register[%d] = %x\n", i, c.register[i])
	}

	c.output.Write([]byte(str))
}

// Process is a function that process a instruction
func (c *Cpu) Process(instr Instruction) error {
	if err := c.handle(instr); err != nil {
		return err
	}

	return nil
}

func (c *Cpu) handle(instr Instruction) error {
	x, err := instr.GetX()
	if err != nil {
		return err
	}

	y, err := instr.GetY()
	if err != nil {
		return err
	}

	nn, err := instr.GetNN()
	if err != nil {
		return err
	}

	nnn, err := instr.GetNNN()
	if err != nil {
		return err
	}

	instrType, instrSubtype, err := instr.GetTypeAndSubType()
	if err != nil {
		return err
	}

	switch instrType {
	case InstructionType(0x01):
		c.process0x1NNN(nnn)
	case InstructionType(0x03):
		c.process0x3XNN(x, nn)
	case InstructionType(0x04):
		c.process0x4XNN(x, nn)
	case InstructionType(0x05):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x5XY0(x, y)
		}
	case InstructionType(0x06):
		c.process0x6XNN(x, nn)
	case InstructionType(0x07):
		c.process0x7XNN(x, nn)
	case InstructionType(0x08):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x8XY0(x, y)
		case InstructionSubType(0x01):
			c.process0x8XY1(x, y)
		case InstructionSubType(0x02):
			c.process0x8XY2(x, y)
		case InstructionSubType(0x03):
			c.process0x8XY3(x, y)
		case InstructionSubType(0x04):
			c.process0x8XY4(x, y)
		case InstructionSubType(0x05):
			c.process0x8XY5(x, y)
		case InstructionSubType(0x06):
			c.process0x8XY6(x, y)
		case InstructionSubType(0x07):
			c.process0x8XY7(x, y)
		case InstructionSubType(0x0E):
			c.process0x8XYE(x, y)
		}
	case InstructionType(0x09):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x9XY0(x, y)
		}
	case InstructionType(0x0A):
		c.process0xANNN(nnn)
	case InstructionType(0x0B):
		c.process0xBNNN(nnn)
	case InstructionType(0x0C):
		c.process0xCXNN(x, nn)
	case InstructionType(0x0F):
		switch nn {
		case 0x1E:
			c.process0xFX1E(x)
		}
	}

	return nil
}

func (c *Cpu) process0x1NNN(nnn uint16) {
	c.pc = nnn
}

func (c *Cpu) process0x3XNN(x, nn byte) {
	if c.register[x] != nn {
		c.pc++
		return
	}
	c.pc += 2
}

func (c *Cpu) process0x4XNN(x, nn byte) {
	if c.register[x] == nn {
		c.pc++
		return
	}
	c.pc += 2
}

func (c *Cpu) process0x5XY0(x, y byte) {
	if c.register[x] != c.register[y] {
		c.pc++
		return
	}
	c.pc += 2
}

func (c *Cpu) process0x6XNN(x, nn byte) {
	c.register[x] = nn
	c.pc++
}

func (c *Cpu) process0x7XNN(x, nn byte) {
	c.register[x] += nn
	c.pc++
}

func (c *Cpu) process0x8XY0(x, y byte) {
	c.register[x] = c.register[y]
	c.pc++
}

func (c *Cpu) process0x8XY1(x, y byte) {
	c.register[x] |= c.register[y]
	c.pc++
}

func (c *Cpu) process0x8XY2(x, y byte) {
	c.register[x] &= c.register[y]
	c.pc++
}

func (c *Cpu) process0x8XY3(x, y byte) {
	c.register[x] ^= c.register[y]
	c.pc++
}

func (c *Cpu) process0x8XY4(x, y byte) {
	val := c.register[x] + c.register[y]

	if int(c.register[x])+int(c.register[y]) > 0xFF {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
	c.pc++
}

func (c *Cpu) process0x8XY5(x, y byte) {
	val := c.register[x] - c.register[y]

	if int(c.register[x])-int(c.register[y]) >= 0x0 {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
	c.pc++
}

func (c *Cpu) process0x8XY6(x, y byte) {
	c.register[0xF] = c.register[x] % 2
	c.register[x] >>= c.register[y]
	c.pc++
}

func (c *Cpu) process0x8XY7(x, y byte) {
	val := c.register[y] - c.register[x]

	if int(c.register[y])-int(c.register[x]) >= 0x0 {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
	c.pc++
}

func (c *Cpu) process0x8XYE(x, y byte) {
	c.register[0xF] = (c.register[x] >> 7)
	c.register[x] <<= c.register[y]
	c.pc++
}

func (c *Cpu) process0x9XY0(x, y byte) {
	if c.register[x] == c.register[y] {
		c.pc++
		return
	}
	c.pc += 2
}

func (c *Cpu) process0xANNN(nnn uint16) {
	c.i = nnn
	c.pc++
}

func (c *Cpu) process0xBNNN(nnn uint16) {
	c.pc = nnn + uint16(c.register[0])
}

func (c *Cpu) process0xCXNN(x, nn byte) {
	c.register[x] = byte(rand.Intn(0xFF)) & nn
	c.pc++
}

func (c *Cpu) process0xFX1E(x byte) {
	c.i += uint16(c.register[x])
	c.pc++
}
