package chip8

import (
	"fmt"
)

type StandardKeyboard struct {
	key   Key
	count int
}

func (sk *StandardKeyboard) Start() {
	fmt.Println("START")
}

func (sk *StandardKeyboard) KeyDown() Key {
	sk.count++
	if sk.count%5 == 0 {
		sk.key = 0x4
	} else {
		sk.key = 0x5
	}

	return sk.key
}
