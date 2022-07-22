package chip8

type Key byte

type Keyboard interface {
	/*
		KeyDown should return a Key mapping the input to any key of CHIP-8 (0x00 - 0x0F).
		When none key is down return Key(0xFF)
	*/
	KeyDown() Key
}
