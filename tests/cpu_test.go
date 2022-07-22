package chip8_test

import (
	"bytes"
	"fmt"
	"math/rand"
	"testing"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

func TestCpu_Log(t *testing.T) {
	expectedRegister := chip8.Register{0x02, 0xB3, 0xAA, 0xFF, 0x00, 0xF0, 0xBC}
	expectedStack := chip8.Stack{0x02, 0xB3, 0xAA, 0xFF, 0x00, 0xF0, 0xBC}
	iExpected := uint16(0x2)
	pcExpected := uint16(0x5)
	dtExpected := byte(0x6)
	stExpected := byte(0x7)
	spExpected := byte(0x5)

	log := &bytes.Buffer{}
	cpu := chip8.NewCpu(&chip8.ConfigCpu{
		Keyboard: nil,
		Stack:    expectedStack,
		Log:      log,
		Register: expectedRegister,
		I:        iExpected,
		PC:       pcExpected,
		SP:       spExpected,
		DT:       dtExpected,
		ST:       stExpected,
	})

	cpu.Log()

	result := log.Bytes()
	expected := cpuToStr(pcExpected, iExpected, expectedRegister, expectedStack, spExpected, dtExpected, stExpected)

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
	stack            chip8.Stack
	keyPressed       chip8.Key
	sp               byte
	expectedRegister chip8.Register
	expectedStack    chip8.Stack
	spExpected       byte
	dtExpected       byte
	stExpected       byte
	pcExpected       uint16
	iExpected        uint16
	flag             bool
}

func TestCpu_Process(t *testing.T) {
	tests := []cpuTestCase{
		{
			describe: "instruction 0x00EE",
			instr:    chip8.Instruction{0x00, 0xEE},
			contexts: []cpuTestCaseContext{
				{
					context:          "when SP is one",
					register:         chip8.Register{},
					stack:            chip8.Stack{0xFF},
					expectedRegister: chip8.Register{},
					expectedStack:    chip8.Stack{0xFF},
					sp:               0x1,
					spExpected:       0x0,
					pcExpected:       0xFF,
				},
			},
		},
		{
			describe: "instruction 0x1NNN",
			instr:    chip8.Instruction{0x1F, 0x5A},
			contexts: []cpuTestCaseContext{
				{
					context:          "when PC has different value of NNN",
					register:         chip8.Register{},
					expectedRegister: chip8.Register{},
					pcExpected:       0xF5A,
				},
			},
		},
		{
			describe: "instruction 0x2NNN",
			instr:    chip8.Instruction{0x2F, 0x90},
			contexts: []cpuTestCaseContext{
				{
					context:          "when SP is zero",
					register:         chip8.Register{},
					expectedRegister: chip8.Register{},
					expectedStack:    chip8.Stack{0x01},
					pcExpected:       0xF90,
					spExpected:       0x01,
				},
			},
		},
		{
			describe: "instruction 0x3XNN",
			instr:    chip8.Instruction{0x30, 0x5A},
			contexts: []cpuTestCaseContext{
				{
					context:          "when value of V[X] is equal to NN",
					register:         chip8.Register{0x5A},
					expectedRegister: chip8.Register{0x5A},
					pcExpected:       0x2,
				},
				{
					context:          "when value of V[X] is different to NN",
					register:         chip8.Register{0x4F},
					expectedRegister: chip8.Register{0x4F},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x4XNN",
			instr:    chip8.Instruction{0x40, 0x5A},
			contexts: []cpuTestCaseContext{
				{
					context:          "when value of V[X] is equal to NN",
					register:         chip8.Register{0x5A},
					expectedRegister: chip8.Register{0x5A},
					pcExpected:       0x1,
				},
				{
					context:          "when value of V[X] is different to NN",
					register:         chip8.Register{0x4F},
					expectedRegister: chip8.Register{0x4F},
					pcExpected:       0x2,
				},
			},
		},
		{
			describe: "instruction 0x5XY0",
			instr:    chip8.Instruction{0x50, 0x10},
			contexts: []cpuTestCaseContext{
				{
					context:          "when value of V[X] is equal to V[Y]",
					register:         chip8.Register{0x5A, 0x5A},
					expectedRegister: chip8.Register{0x5A, 0x5A},
					pcExpected:       0x2,
				},
				{
					context:          "when value of V[X] is different to NN",
					register:         chip8.Register{0x5A, 0x4F},
					expectedRegister: chip8.Register{0x5A, 0x4F},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x6XNN",
			instr:    chip8.Instruction{0x61, 0x50},
			contexts: []cpuTestCaseContext{
				{
					context:          "when X has different value of NN",
					register:         chip8.Register{0x0F, 0xAB, 0xCD},
					expectedRegister: chip8.Register{0x0F, 0x50, 0xCD},
					pcExpected:       0x1,
				},
				{
					context:          "when X has equals value of NN",
					register:         chip8.Register{0x0F, 0x50, 0xAB},
					expectedRegister: chip8.Register{0x0F, 0x50, 0xAB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x7XNN",
			instr:    chip8.Instruction{0x71, 0x50},
			contexts: []cpuTestCaseContext{
				{
					context:          "when X+NN is less than (0xFF + 1)",
					register:         chip8.Register{0x0F, 0xAB, 0xCD},
					expectedRegister: chip8.Register{0x0F, 0xFB, 0xCD},
					pcExpected:       0x1,
				},
				{
					context:          "when X+NN is greater than (0xFF + 1)",
					register:         chip8.Register{0x0F, 0xFB, 0xAB},
					expectedRegister: chip8.Register{0x0F, 0x4B, 0xAB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY0",
			instr:    chip8.Instruction{0x81, 0x20},
			contexts: []cpuTestCaseContext{
				{
					context:          "when X has different value of Y",
					register:         chip8.Register{0x0F, 0xAB, 0xCD},
					expectedRegister: chip8.Register{0x0F, 0xCD, 0xCD},
					pcExpected:       0x1,
				},
				{
					context:          "when X has equals value of Y",
					register:         chip8.Register{0x0F, 0xAB, 0xAB},
					expectedRegister: chip8.Register{0x0F, 0xAB, 0xAB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY1",
			instr:    chip8.Instruction{0x80, 0x21},
			contexts: []cpuTestCaseContext{
				{
					context:          "when X has different value of Y",
					register:         chip8.Register{0xAB, 0x0F, 0xCD},
					expectedRegister: chip8.Register{0xEF, 0x0F, 0xCD},
					pcExpected:       0x1,
				},
				{
					context:          "when X has equals value of Y",
					register:         chip8.Register{0xAB, 0x0F, 0xAB},
					expectedRegister: chip8.Register{0xAB, 0x0F, 0xAB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY2",
			instr:    chip8.Instruction{0x80, 0x22},
			contexts: []cpuTestCaseContext{
				{
					context:          "when X has different value of Y",
					register:         chip8.Register{0xAB, 0x0F, 0xCD},
					expectedRegister: chip8.Register{0x89, 0x0F, 0xCD},
					pcExpected:       0x1,
				},
				{
					context:          "when X has equals value of Y",
					register:         chip8.Register{0xAB, 0x0F, 0xAB},
					expectedRegister: chip8.Register{0xAB, 0x0F, 0xAB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY3",
			instr:    chip8.Instruction{0x80, 0x23},
			contexts: []cpuTestCaseContext{
				{
					context:          "when X has different value of Y",
					register:         chip8.Register{0xAB, 0x0F, 0xCD},
					expectedRegister: chip8.Register{0x66, 0x0F, 0xCD},
					pcExpected:       0x1,
				},
				{
					context:          "when X has equals value of Y",
					register:         chip8.Register{0xAB, 0x0F, 0xAB},
					expectedRegister: chip8.Register{0x0, 0x0F, 0xAB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY4",
			instr:    chip8.Instruction{0x80, 0x14},
			contexts: []cpuTestCaseContext{
				{
					context:          "when carry does not occur",
					register:         chip8.Register{0xF0, 0x0F},
					expectedRegister: chip8.Register{0xFF, 0x0F},
					pcExpected:       0x1,
				},
				{
					context:          "when carry occurs",
					register:         chip8.Register{0xFF, 0xFF},
					expectedRegister: chip8.Register{0xFE, 0xFF},
					flag:             true,
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY5",
			instr:    chip8.Instruction{0x80, 0x15},
			contexts: []cpuTestCaseContext{
				{
					context:          "when borrow does not occur",
					register:         chip8.Register{0xFF, 0x0F},
					expectedRegister: chip8.Register{0xF0, 0x0F},
					flag:             true,
					pcExpected:       0x1,
				},
				{
					context:          "when borrow occurs",
					register:         chip8.Register{0x0F, 0xFF},
					expectedRegister: chip8.Register{0x10, 0xFF},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY6",
			instr:    chip8.Instruction{0x80, 0x16},
			contexts: []cpuTestCaseContext{
				{
					context:          "when the old least significant bit is 0",
					register:         chip8.Register{0xAE, 0x03},
					expectedRegister: chip8.Register{0x15, 0x03},
					pcExpected:       0x1,
				},
				{
					context:          "when the old least significant bit is 1",
					register:         chip8.Register{0xAF, 0x03},
					expectedRegister: chip8.Register{0x15, 0x03},
					flag:             true,
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XY7",
			instr:    chip8.Instruction{0x80, 0x17},
			contexts: []cpuTestCaseContext{
				{
					context:          "when borrow does not occur",
					register:         chip8.Register{0x0F, 0xFF},
					expectedRegister: chip8.Register{0xF0, 0xFF},
					flag:             true,
					pcExpected:       0x1,
				},
				{
					context:          "when borrow occurs",
					register:         chip8.Register{0xFF, 0x0F},
					expectedRegister: chip8.Register{0x10, 0x0F},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x8XYE",
			instr:    chip8.Instruction{0x80, 0x1E},
			contexts: []cpuTestCaseContext{
				{
					context:          "when the old most significant bit is 0",
					register:         chip8.Register{0x1E, 0x03},
					expectedRegister: chip8.Register{0xF0, 0x03},
					pcExpected:       0x1,
				},
				{
					context:          "when the old most significant bit is 1",
					register:         chip8.Register{0xAE, 0x03},
					expectedRegister: chip8.Register{0x70, 0x03},
					flag:             true,
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0x9XY0",
			instr:    chip8.Instruction{0x90, 0x10},
			contexts: []cpuTestCaseContext{
				{
					context:          "when value of V[X] is equal to V[Y]",
					register:         chip8.Register{0x5A, 0x5A},
					expectedRegister: chip8.Register{0x5A, 0x5A},
					pcExpected:       0x1,
				},
				{
					context:          "when value of V[X] is different to NN",
					register:         chip8.Register{0x5A, 0x4F},
					expectedRegister: chip8.Register{0x5A, 0x4F},
					pcExpected:       0x2,
				},
			},
		},
		{
			describe: "instruction 0xANNN",
			instr:    chip8.Instruction{0xAA, 0xBC},
			contexts: []cpuTestCaseContext{
				{
					context:          "when value of I is different to NNN",
					register:         chip8.Register{0x5A, 0x5A},
					expectedRegister: chip8.Register{0x5A, 0x5A},
					pcExpected:       0x1,
					iExpected:        0xABC,
				},
			},
		},
		{
			describe: "instruction 0xBNNN",
			instr:    chip8.Instruction{0xBA, 0xBC},
			contexts: []cpuTestCaseContext{
				{
					context:          "when NNN+V0 is less or equal than 0xFFF)",
					register:         chip8.Register{0x4A, 0xAB, 0xCD},
					expectedRegister: chip8.Register{0x4A, 0xAB, 0xCD},
					pcExpected:       0xB06,
				},
			},
		},
		{
			describe: "instruction 0xCXNN",
			instr:    chip8.Instruction{0xC0, 0x5E},
			contexts: []cpuTestCaseContext{
				{
					context:          "when random number is 0x0E",
					register:         chip8.Register{0xFA, 0xBB},
					expectedRegister: chip8.Register{0x0E, 0xBB},
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0xEX9E",
			instr:    chip8.Instruction{0xE1, 0x9E},
			contexts: []cpuTestCaseContext{
				{
					context:          "when key pressed is different the V[X]",
					register:         chip8.Register{0xFA, 0x0A},
					expectedRegister: chip8.Register{0xFA, 0x0A},
					keyPressed:       0x05,
					pcExpected:       0x1,
				},
				{
					context:          "when key pressed is equals the V[X]",
					register:         chip8.Register{0xFA, 0x0A},
					expectedRegister: chip8.Register{0xFA, 0x0A},
					keyPressed:       0x0A,
					pcExpected:       0x2,
				},
			},
		},
		{
			describe: "instruction 0xEXA1",
			instr:    chip8.Instruction{0xE1, 0xA1},
			contexts: []cpuTestCaseContext{
				{
					context:          "when key pressed is different the V[X]",
					register:         chip8.Register{0xFA, 0x0A},
					expectedRegister: chip8.Register{0xFA, 0x0A},
					keyPressed:       0x05,
					pcExpected:       0x2,
				},
				{
					context:          "when key pressed is equals the V[X]",
					register:         chip8.Register{0xFA, 0x0A},
					expectedRegister: chip8.Register{0xFA, 0x0A},
					keyPressed:       0x0A,
					pcExpected:       0x1,
				},
			},
		},
		{
			describe: "instruction 0xFX07",
			instr:    chip8.Instruction{0xF1, 0x07},
			contexts: []cpuTestCaseContext{
				{
					context:          "when V[X] is different to DT",
					register:         chip8.Register{0xFA, 0xBB},
					expectedRegister: chip8.Register{0xFA, 0x00},
					pcExpected:       0x1,
					dtExpected:       0x0,
				},
			},
		},
		{
			describe: "instruction 0xFX15",
			instr:    chip8.Instruction{0xF1, 0x15},
			contexts: []cpuTestCaseContext{
				{
					context:          "when DT is different to V[X]",
					register:         chip8.Register{0xFA, 0xBB},
					expectedRegister: chip8.Register{0xFA, 0xBB},
					pcExpected:       0x1,
					dtExpected:       0xBB,
				},
			},
		},
		{
			describe: "instruction 0xFX18",
			instr:    chip8.Instruction{0xF1, 0x18},
			contexts: []cpuTestCaseContext{
				{
					context:          "when DT is different to V[X]",
					register:         chip8.Register{0xFA, 0xBB},
					expectedRegister: chip8.Register{0xFA, 0xBB},
					pcExpected:       0x1,
					stExpected:       0xBB,
				},
			},
		},
		{
			describe: "instruction 0xFX1E",
			instr:    chip8.Instruction{0xF0, 0x1E},
			contexts: []cpuTestCaseContext{
				{
					context:          "when i + V[X] is less than 0xFF",
					register:         chip8.Register{0xFA, 0xBB},
					expectedRegister: chip8.Register{0xFA, 0xBB},
					pcExpected:       0x1,
					iExpected:        0xFA,
				},
			},
		},
	}

	rand.Seed(51153153)

	// Set register 0xF when flag is true
	setFlags(tests)

	for _, test := range tests {
		t.Run(test.describe, func(t *testing.T) {
			for _, context := range test.contexts {
				t.Run(context.context, func(t *testing.T) {
					log := &bytes.Buffer{}

					cpu := chip8.NewCpu(&chip8.ConfigCpu{
						Keyboard: MockKeyBoard{Key: context.keyPressed},
						Stack:    context.stack,
						Log:      log,
						Register: context.register,
						I:        0x0,
						PC:       0x0,
						SP:       context.sp,
						DT:       0x0,
						ST:       0x0,
					})

					err := cpu.Process(test.instr)
					if err != nil {
						t.Fatalf("error not expected: %s", err.Error())
					}

					cpu.Log()

					result := log.Bytes()
					expected := cpuToStr(context.pcExpected, context.iExpected, context.expectedRegister,
						context.expectedStack, context.spExpected, context.dtExpected, context.stExpected)

					if string(result) != string(expected) {
						t.Errorf("result: %s, expected: %s", result, expected)
					}
				})
			}
		})
	}
}

func setFlags(tests []cpuTestCase) {
	for i, test := range tests {
		for j, context := range test.contexts {
			if context.flag {
				tests[i].contexts[j].expectedRegister[0xF] = 0x1
			}
		}
	}
}

func cpuToStr(pc, i uint16, register chip8.Register, stack chip8.Stack, sp, dt, st byte) []byte {
	str := fmt.Sprintf("pc = %x\nsp = %x\ndt = %x\nst = %x\ni = %x\nstack = %v\n", pc, sp, dt, st, i, stack)

	for i := 0; i < len(register); i++ {
		str += fmt.Sprintf("register[%d] = %x\n", i, register[i])
	}

	return []byte(str)
}

type MockKeyBoard struct {
	Key chip8.Key
}

func (m MockKeyBoard) KeyDown() chip8.Key {
	return m.Key
}
