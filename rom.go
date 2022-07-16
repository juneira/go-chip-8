package chip8

import "io"

type Instruction []byte

type Rom struct {
	data io.Reader
}

// NewRom is a function that receive a "data" as param and return a pointer to Rom
func NewRom(data io.Reader) *Rom {
	return &Rom{data: data}
}

// NextInstruction return the next instruction
func (r *Rom) NextInstruction() *Instruction {
	instr := make(Instruction, 2)
	r.data.Read(instr)

	return &instr
}
