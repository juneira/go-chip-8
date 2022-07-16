package chip8_test

import (
	"bytes"
	"fmt"
	"testing"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

func TestCpu_Log(t *testing.T) {
	expectedRegister := chip8.Register{0x02, 0xB3, 0xAA, 0xFF, 0x00, 0xF0, 0xBC}
	output := &bytes.Buffer{}
	cpu := chip8.NewCpu(expectedRegister, output)

	cpu.Log()

	result := output.Bytes()
	expected := registersToStr(expectedRegister)

	if string(result) != string(expected) {
		t.Errorf("result: %s, expected: %s", result, expected)
	}
}

func TestCpu_Process(t *testing.T) {
	t.Run("instruction 0x8XY4", func(t *testing.T) {
		instr := chip8.Instruction{0x80, 0x14}

		t.Run("when carry does not occur", func(t *testing.T) {
			register := chip8.Register{0xF0, 0x0F}
			output := &bytes.Buffer{}
			cpu := chip8.NewCpu(register, output)

			cpu.Process(instr)
			cpu.Log()

			expectedRegister := chip8.Register{0xFF, 0x0F}
			expectedRegister[0xF] = 0x0

			result := output.Bytes()
			expected := registersToStr(expectedRegister)

			if string(result) != string(expected) {
				t.Errorf("result: %s, expected: %s", result, expected)
			}
		})

		t.Run("when carry occurs", func(t *testing.T) {
			register := chip8.Register{0xFF, 0x0FF}
			output := &bytes.Buffer{}
			cpu := chip8.NewCpu(register, output)

			cpu.Process(instr)
			cpu.Log()

			expectedRegister := chip8.Register{0xFF, 0xFF}
			expectedRegister[0xF] = 0x01

			result := output.Bytes()
			expected := registersToStr(expectedRegister)

			if string(result) != string(expected) {
				t.Errorf("result: %s, expected: %s", result, expected)
			}
		})
	})
}

func registersToStr(registers chip8.Register) []byte {
	str := ""

	for i := 0; i < len(registers); i++ {
		str += fmt.Sprintf("register[%d] = %x\n", i, registers[i])
	}

	return []byte(str)
}
