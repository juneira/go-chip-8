package chip8_test

import (
	"testing"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

type instrTestCase struct {
	instr           chip8.Instruction
	expected        uint16
	isExpectedError bool
}

func TestInstruction_GetX(t *testing.T) {
	tests := []instrTestCase{
		{chip8.Instruction{0x31, 0x02}, 0x01, false},
		{chip8.Instruction{0x42, 0x03}, 0x02, false},
		{chip8.Instruction{0x53, 0x02}, 0x03, false},
		{chip8.Instruction{0xFF, 0x06}, 0x0F, false},
		{chip8.Instruction{0xFF, 0x06, 0x06}, 0x0, true},
		{chip8.Instruction{0xFF}, 0x0, true},
		{chip8.Instruction{}, 0x0, true},
	}

	for _, test := range tests {
		result, err := test.instr.GetX()

		checkResult(t, uint16(result), err, test)
	}
}

func TestInstruction_GetY(t *testing.T) {
	tests := []instrTestCase{
		{chip8.Instruction{0x52, 0x14}, 0x01, false},
		{chip8.Instruction{0x83, 0x22}, 0x02, false},
		{chip8.Instruction{0x91, 0x30}, 0x03, false},
		{chip8.Instruction{0xDE, 0xFA}, 0x0F, false},
		{chip8.Instruction{0xFF, 0x06, 0x06}, 0x0, true},
		{chip8.Instruction{0xFF}, 0x0, true},
		{chip8.Instruction{}, 0x0, true},
	}

	for _, test := range tests {
		result, err := test.instr.GetY()

		checkResult(t, uint16(result), err, test)
	}
}

func TestInstruction_GetN(t *testing.T) {
	tests := []instrTestCase{
		{chip8.Instruction{0xD1, 0xF1}, 0x01, false},
		{chip8.Instruction{0xD2, 0xA2}, 0x02, false},
		{chip8.Instruction{0xDA, 0x03}, 0x03, false},
		{chip8.Instruction{0xDE, 0x2F}, 0x0F, false},
		{chip8.Instruction{0xFF, 0x05, 0x65}, 0x0, true},
		{chip8.Instruction{0xFF}, 0x0, true},
		{chip8.Instruction{}, 0x0, true},
	}

	for _, test := range tests {
		result, err := test.instr.GetN()

		checkResult(t, uint16(result), err, test)
	}
}

func TestInstruction_GetNN(t *testing.T) {
	tests := []instrTestCase{
		{chip8.Instruction{0x3A, 0x10}, 0x10, false},
		{chip8.Instruction{0x4E, 0x20}, 0x20, false},
		{chip8.Instruction{0x60, 0xAB}, 0xAB, false},
		{chip8.Instruction{0x71, 0xFF}, 0xFF, false},
		{chip8.Instruction{0xFF, 0x06, 0x65}, 0x0, true},
		{chip8.Instruction{0xFF}, 0x0, true},
		{chip8.Instruction{}, 0x0, true},
	}

	for _, test := range tests {
		result, err := test.instr.GetNN()

		checkResult(t, uint16(result), err, test)
	}
}

func TestInstruction_GetNNN(t *testing.T) {
	tests := []instrTestCase{
		{chip8.Instruction{0x1A, 0x10}, 0xA10, false},
		{chip8.Instruction{0x2E, 0x20}, 0xE20, false},
		{chip8.Instruction{0xA0, 0xAB}, 0x0AB, false},
		{chip8.Instruction{0xB1, 0xFF}, 0x1FF, false},
		{chip8.Instruction{0xFF, 0x06, 0x65}, 0x0, true},
		{chip8.Instruction{0xFF}, 0x0, true},
		{chip8.Instruction{}, 0x0, true},
	}

	for _, test := range tests {
		result, err := test.instr.GetNNN()

		checkResult(t, result, err, test)
	}
}

func checkResult(t *testing.T, result uint16, err error, test instrTestCase) {
	t.Helper()

	if err != nil && !test.isExpectedError {
		t.Fatalf("error not expected: %s", err.Error())
	}

	if err == nil && test.isExpectedError {
		t.Fatal("error is expected but doesn't ocorrs")
	}

	if result != test.expected {
		t.Errorf("result: 0x%X, expected: 0x%X", result, test.expected)
	}
}
