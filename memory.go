package chip8

type Memory interface {
	/*
		Save should saves the register on the memory starting from I on memory
	*/
	Save(register []Register, i uint16)

	/*
		Load should loads the register on the memory starting from I on memory
	*/
	Load(register []Register, i uint16)

	/*
		LoadChar should return the address of char referring to vx
	*/
	LoadChar(vx byte) uint16
}
