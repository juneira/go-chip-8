package chip8

import (
	"fmt"
	"io"
)

type Register [0x10]byte

type Cpu struct {
	register Register
	output   io.Writer
}

// NewCpu is a function that receive a "output" as param and return a pointer to Cpu
func NewCpu(register Register, output io.Writer) *Cpu {
	return &Cpu{register: register, output: output}
}

// Log is a function that write values of registers to "output" of Cpu
func (c *Cpu) Log() {
	for i := 0; i < len(c.register); i++ {
		str := fmt.Sprintf("register[%d] = %x\n", i, c.register[i])
		c.output.Write([]byte(str))
	}
}

// Process is a function that process a register
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

	instrType, instrSubtype, err := instr.GetTypeAndSubType()
	if err != nil {
		return err
	}

	switch instrType {
	case InstructionType(0x08):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x8XY0(x, y)
		case InstructionSubType(0x04):
			c.process0x8XY4(x, y)
		case InstructionSubType(0x05):
			c.process0x8XY5(x, y)
		case InstructionSubType(0x07):
			c.process0x8XY7(x, y)
		}
	}

	return nil
}

func (c *Cpu) process0x8XY0(x, y byte) {
	c.register[x] = c.register[y]
}

func (c *Cpu) process0x8XY4(x, y byte) {
	val := c.register[x] + c.register[y]

	if int(c.register[x])+int(c.register[y]) > 0xFF {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
}

func (c *Cpu) process0x8XY5(x, y byte) {
	val := c.register[x] - c.register[y]

	if int(c.register[x])-int(c.register[y]) >= 0x0 {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
}

func (c *Cpu) process0x8XY7(x, y byte) {
	val := c.register[y] - c.register[x]

	if int(c.register[y])-int(c.register[x]) >= 0x0 {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
}
