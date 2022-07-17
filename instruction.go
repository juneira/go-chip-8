package chip8

import (
	"fmt"
)

type Instruction []byte
type InstructionType byte
type InstructionSubType byte

// GetX returns the parameter X of instruction
func (instr *Instruction) GetX() (byte, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return instr.firstByte() & 0x0F, nil
}

// GetY returns the parameter Y of instruction
func (instr *Instruction) GetY() (byte, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return instr.secondByte() >> 4, nil
}

// GetN returns the parameter N of instruction
func (instr *Instruction) GetN() (byte, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return instr.secondByte() & 0x0F, nil
}

// GetNN returns the parameter NN of instruction
func (instr *Instruction) GetNN() (byte, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return instr.secondByte(), nil
}

// GetNNN returns the parameter NNN of instruction
func (instr *Instruction) GetNNN() (uint16, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return ((uint16(instr.firstByte()) & 0x0F) << 8) | uint16(instr.secondByte()), nil
}

func (instr *Instruction) GetTypeAndSubType() (t InstructionType, s InstructionSubType, err error) {
	t, err = instr.instructionType()
	if err != nil {
		return
	}

	s, err = instr.subInstructionType()
	if err != nil {
		return
	}

	return
}

func (instr *Instruction) instructionType() (InstructionType, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return InstructionType(instr.firstByte() >> 4), nil
}

func (instr *Instruction) subInstructionType() (InstructionSubType, error) {
	if err := instr.validate(); err != nil {
		return 0, err
	}

	return InstructionSubType(instr.secondByte() & 0x0F), nil
}

func (instr *Instruction) firstByte() byte {
	return (*instr)[0]
}

func (instr *Instruction) secondByte() byte {
	return (*instr)[1]
}

func (instr *Instruction) validate() error {
	if len(*instr) != 2 {
		return fmt.Errorf("invalid instruction")
	}

	return nil
}
