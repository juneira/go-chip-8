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
func (c *Cpu) Process(instr Instruction) {
	x := 0x0F & instr[0]
	y := instr[1] >> 4

	val := c.register[x] + c.register[y]

	if int(c.register[x])+int(c.register[y]) > 0xFF {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
}
