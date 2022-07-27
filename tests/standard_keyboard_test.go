package chip8

import (
	"fmt"
	"strconv"
	"testing"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

func TestStandardKeyboard_KeyBoard(t *testing.T) {
	keys := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F"}
	for _, key := range keys {
		t.Run(fmt.Sprintf("when key is %s", key), func(t *testing.T) {
			b := []byte(key)[0]

			m := &MockInputKeyboard{key: key}
			sk := chip8.NewStandardKeyboard(&chip8.ConfigKeyboard{Input: m})

			expected := chip8.Key((b - 65) + 0x0A)
			if b < byte(65) {
				value, _ := strconv.Atoi(key)
				expected = chip8.Key(value)
			}

			result := sk.KeyDown()

			if result != expected {
				t.Errorf("result: 0x%X, expected: 0x%X", result, expected)
			}
		})

		t.Run("when key is empty", func(t *testing.T) {
			m := &MockInputKeyboard{key: ""}
			sk := chip8.NewStandardKeyboard(&chip8.ConfigKeyboard{Input: m})

			expected := chip8.Key(0xFF)
			result := sk.KeyDown()

			if result != expected {
				t.Errorf("result: 0x%X, expected: 0x%X", result, expected)
			}
		})
	}
}

type MockInputKeyboard struct {
	key string
}

func (mi *MockInputKeyboard) Read(p []byte) (n int, err error) {
	if mi.key != "" {
		p[0] = []byte(mi.key)[0]
		return 1, nil
	}

	return 0, nil
}
