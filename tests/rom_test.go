package rom_test

import (
	"bytes"
	"reflect"
	"testing"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

func TestRom_NextInstruction(t *testing.T) {
	buf := chip8.Instruction("MARCEL")

	data := bytes.NewBuffer(buf)

	rom := chip8.NewRom(data)

	for i := 0; i < len(buf); i += 2 {
		expected := chip8.Instruction(buf[i : i+2])

		checkInstruction(t, rom.NextInstruction(), &expected)
	}
}

func checkInstruction(t *testing.T, result, expected *chip8.Instruction) {
	t.Helper()

	if result == nil {
		t.Fatal("result is nil")
	}

	if !reflect.DeepEqual(*result, *expected) {
		t.Errorf("result %v, expected %v", *result, *expected)
	}
}
