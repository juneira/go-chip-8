package chip8

type Key byte

type Keyboard interface {
	/*
		KeyDown should return a Key mapping the input to any key of CHIP-8 (0x0 - 0xF).
		When none key is down return Key(0xFF)
	*/
	KeyDown() Key
}
