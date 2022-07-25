package chip8

import "io"

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
	n, err := r.data.Read(instr)

	if n == 0 {
		return nil
	}

	if err != nil {
		panic(err)
	}

	return &instr
}
