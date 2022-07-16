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

type cpuTestCase struct {
	describe string
	instr    chip8.Instruction
	contexts []cpuTestCaseContext
}

type cpuTestCaseContext struct {
	context          string
	register         chip8.Register
	expectedRegister chip8.Register
}

func TestCpu_Process(t *testing.T) {
	tests := []cpuTestCase{
		{
			describe: "instruction 0x8XY4",
			instr:    chip8.Instruction{0x80, 0x14},
			contexts: []cpuTestCaseContext{
				{
					context:          "when carry does not occur",
					register:         chip8.Register{0xF0, 0x0F},
					expectedRegister: chip8.Register{0xFF, 0x0F},
				},
				{
					context:          "when carry occurs",
					register:         chip8.Register{0xFF, 0xFF},
					expectedRegister: chip8.Register{0xFE, 0xFF},
				},
			},
		},
	}

	// set flags
	tests[0].contexts[1].expectedRegister[0xF] = 1

	for _, test := range tests {
		t.Run(test.describe, func(t *testing.T) {
			instr := chip8.Instruction{0x80, 0x14}

			for _, context := range test.contexts {
				t.Run(context.context, func(t *testing.T) {
					output := &bytes.Buffer{}
					cpu := chip8.NewCpu(context.register, output)

					err := cpu.Process(instr)
					if err != nil {
						t.Fatalf("error not expected: %s", err.Error())
					}

					cpu.Log()

					result := output.Bytes()
					expected := registersToStr(context.expectedRegister)

					if string(result) != string(expected) {
						t.Errorf("result: %s, expected: %s", result, expected)
					}
				})
			}
		})
	}
}

func registersToStr(registers chip8.Register) []byte {
	str := ""

	for i := 0; i < len(registers); i++ {
		str += fmt.Sprintf("register[%d] = %x\n", i, registers[i])
	}

	return []byte(str)
}
