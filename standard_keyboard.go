package chip8

import "io"

// StandardKeyboard implements KeyBoard
// It's just useful for example "terminal.go"
type StandardKeyboard struct {
	input io.Reader
}

type ConfigKeyboard struct {
	Input io.Reader
}

// NewStandardKeyboard is a function that receive a config as param and return a pointer to StandardKeyboard
func NewStandardKeyboard(config *ConfigKeyboard) *StandardKeyboard {
	return &StandardKeyboard{input: config.Input}
}

// KeyDown check the input and translate to CHIP-8
func (sk *StandardKeyboard) KeyDown() Key {
	var buf [1]byte

	n, err := sk.input.Read(buf[:])
	if n == 0 {
		// Return "none key pressed"
		return Key(0xFF)
	}
	if err != nil {
		panic(err)
	}

	switch string(buf[0]) {
	case "0":
		return Key(0x00)
	case "1":
		return Key(0x01)
	case "2":
		return Key(0x02)
	case "3":
		return Key(0x03)
	case "4":
		return Key(0x04)
	case "5":
		return Key(0x05)
	case "6":
		return Key(0x06)
	case "7":
		return Key(0x07)
	case "8":
		return Key(0x08)
	case "9":
		return Key(0x09)
	case "A":
		return Key(0x0A)
	case "B":
		return Key(0x0B)
	case "C":
		return Key(0x0C)
	case "D":
		return Key(0x0D)
	case "E":
		return Key(0x0E)
	case "F":
		return Key(0x0F)
	}

	return Key(0xFF)
}
