package chip8_test

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

func TestStandardMemory_Log(t *testing.T) {
	mem, log := newMemory()

	mem.Log()

	memExpected := initialMemory()

	expected := memToStr(memExpected[:])
	result := log.Bytes()

	if string(result) != string(expected) {
		t.Errorf("result: %s\nexpected: %s\n", string(result), string(expected))
	}
}

func TestStandardMemory_Save(t *testing.T) {
	mem, log := newMemory()

	mem.Save([]byte{0x1F, 0x2F}, 0x100)
	mem.Log()

	memExpected := initialMemory()
	memExpected[0x100] = 0x1F
	memExpected[0x101] = 0x2F

	expected := memToStr(memExpected[:])
	result := log.Bytes()
	if string(result) != string(expected) {
		t.Errorf("result: %s\nexpected: %s\n", string(result), string(expected))
	}
}

func TestStandardMemory_SaveBCD(t *testing.T) {
	mem, log := newMemory()

	memExpected := initialMemory()
	memExpected[0x100] = 1
	memExpected[0x101] = 5
	memExpected[0x102] = 2

	vx := byte(152)
	i := uint16(0x100)

	mem.SaveBCD(vx, i)
	mem.Log()

	expected := memToStr(memExpected[:])
	result := log.Bytes()
	if string(result) != string(expected) {
		t.Errorf("result: %s\nexpected: %s\n", string(result), string(expected))
	}
}

func TestStandardMemory_Load(t *testing.T) {
	mem, _ := newMemory()

	expected := []byte{0x60, 0x20}

	i := uint16(0x6)
	result := []byte{0x0, 0x0}
	mem.Load(result, i)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v\nexpected: %v\n", result, expected)
	}
}

func TestStandardMemory_LoadInstruction(t *testing.T) {
	mem, _ := newMemory()

	expected := chip8.Instruction{0x60, 0x20}

	pc := uint16(0x6)
	result := mem.LoadInstruction(pc)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v\nexpected: %v\n", result, expected)
	}
}

func TestStandardMemory_LoadChar(t *testing.T) {
	mem, _ := newMemory()

	expected := uint16(0x23)

	vx := byte(0x07)
	result := mem.LoadChar(vx)

	if result != expected {
		t.Errorf("result: %v\nexpected: %v\n", result, expected)
	}
}

func TestStandardMemory_LoadSprite(t *testing.T) {
	mem, _ := newMemory()

	expected := byte(0xF0)

	i := uint16(0x11)
	result := mem.LoadSprite(i)

	if result != expected {
		t.Errorf("result: %v\nexpected: %v\n", result, expected)
	}
}

func newMemory() (*chip8.StandardMemory, *bytes.Buffer) {
	log := &bytes.Buffer{}
	mem := chip8.NewStandardMemory(&chip8.ConfigMemory{Rom: &MockRom{}, Log: log})
	return mem, log
}

func initialMemory() [0xFFF]byte {
	return [0xFFF]byte{
		// 0
		0xF0, 0x90, 0x90, 0x90, 0xF0,
		// 1
		0x20, 0x60, 0x20, 0x20, 0x70,
		// 2
		0xF0, 0x10, 0xF0, 0x80, 0xF0,
		// 3
		0xF0, 0x10, 0xF0, 0x10, 0xF0,
		// 4
		0x90, 0x90, 0xF0, 0x10, 0x10,
		// 5
		0xF0, 0x80, 0xF0, 0x10, 0xF0,
		// 6
		0xF0, 0x80, 0xF0, 0x90, 0xF0,
		// 7
		0xF0, 0x10, 0x20, 0x40, 0x40,
		// 8
		0xF0, 0x90, 0xF0, 0x90, 0xF0,
		// 9
		0xF0, 0x90, 0xF0, 0x10, 0xF0,
		// A
		0xF0, 0x90, 0xF0, 0x90, 0x90,
		// B
		0xE0, 0x90, 0xE0, 0x90, 0xE0,
		// C
		0xF0, 0x80, 0x80, 0x80, 0xF0,
		// D
		0xE0, 0x90, 0x90, 0x90, 0xE0,
		// E
		0xF0, 0x80, 0xF0, 0x80, 0xF0,
		// F
		0xF0, 0x80, 0xF0, 0x80, 0x80,
	}
}

func memToStr(mem []byte) []byte {
	return []byte(fmt.Sprintf("memory: %v\n", mem))
}
