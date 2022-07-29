package chip8

import (
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Register [0x10]byte
type Stack [0x10]uint16

type Cpu struct {
	display  Display
	keyboard Keyboard
	sound    Sound
	memory   Memory
	register Register
	stack    Stack
	log      io.Writer
	sp       byte
	dt       byte
	st       byte
	pc       uint16
	i        uint16
}

type ConfigCpu struct {
	// Externals devices
	Display  Display
	Keyboard Keyboard
	Sound    Sound
	Memory   Memory
	Stack    Stack
	Log      io.Writer

	// Registers
	Register Register
	I        uint16

	// Pseudo Registers
	PC uint16
	SP byte
	DT byte
	ST byte
}

// NewCpu receives params and return a pointer to Cpu
func NewCpu(config *ConfigCpu) *Cpu {
	return &Cpu{
		display:  config.Display,
		keyboard: config.Keyboard,
		sound:    config.Sound,
		memory:   config.Memory,
		register: config.Register,
		stack:    config.Stack,
		log:      config.Log,
		pc:       config.PC,
		i:        config.I,
		sp:       config.SP,
		dt:       config.DT,
		st:       config.ST,
	}
}

// Start initialize decrementation of DT and ST registers
func (c *Cpu) Start() {
	go func() {
		for {
			if c.dt > 0 {
				c.dt--
			}

			if c.st > 0 {
				c.st--
				if c.st == 0 {
					c.sound.Beep()
				}
			}

			time.Sleep(16 * time.Millisecond)
		}
	}()
}

// Log writes values of registers to "log" of Cpu
func (c *Cpu) Log() {
	str := fmt.Sprintf("pc = %x\nsp = %x\ndt = %x\nst = %x\ni = %x\nstack = %v\n", c.pc, c.sp, c.dt, c.st, c.i, c.stack)

	for i := 0; i < len(c.register); i++ {
		str += fmt.Sprintf("register[%d] = %x\n", i, c.register[i])
	}

	c.log.Write([]byte(str))
}

// Process is a function that process a instruction
func (c *Cpu) Process(instr Instruction) error {
	if err := c.handle(instr); err != nil {
		return err
	}

	return nil
}

func (c *Cpu) NextInstruction() uint16 {
	return c.pc
}

func (c *Cpu) handle(instr Instruction) error {
	x, err := instr.GetX()
	if err != nil {
		return err
	}

	y, err := instr.GetY()
	if err != nil {
		return err
	}

	n, err := instr.GetN()
	if err != nil {
		return err
	}

	nn, err := instr.GetNN()
	if err != nil {
		return err
	}

	nnn, err := instr.GetNNN()
	if err != nil {
		return err
	}

	instrType, instrSubtype, err := instr.GetTypeAndSubType()
	if err != nil {
		return err
	}

	switch instrType {
	case InstructionType(0x00):
		switch nn {
		case 0xE0:
			c.process0x00E0()
		case 0xEE:
			c.process0x00EE()
		default:
			panic("invalid instruction")
		}
	case InstructionType(0x01):
		c.process0x1NNN(nnn)
	case InstructionType(0x02):
		c.process0x2NNN(nnn)
	case InstructionType(0x03):
		c.process0x3XNN(x, nn)
	case InstructionType(0x04):
		c.process0x4XNN(x, nn)
	case InstructionType(0x05):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x5XY0(x, y)
		default:
			panic("invalid instruction")
		}
	case InstructionType(0x06):
		c.process0x6XNN(x, nn)
	case InstructionType(0x07):
		c.process0x7XNN(x, nn)
	case InstructionType(0x08):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x8XY0(x, y)
		case InstructionSubType(0x01):
			c.process0x8XY1(x, y)
		case InstructionSubType(0x02):
			c.process0x8XY2(x, y)
		case InstructionSubType(0x03):
			c.process0x8XY3(x, y)
		case InstructionSubType(0x04):
			c.process0x8XY4(x, y)
		case InstructionSubType(0x05):
			c.process0x8XY5(x, y)
		case InstructionSubType(0x06):
			c.process0x8XY6(x, y)
		case InstructionSubType(0x07):
			c.process0x8XY7(x, y)
		case InstructionSubType(0x0E):
			c.process0x8XYE(x, y)
		default:
			panic("invalid instruction")
		}
	case InstructionType(0x09):
		switch instrSubtype {
		case InstructionSubType(0x00):
			c.process0x9XY0(x, y)
		default:
			panic("invalid instruction")
		}
	case InstructionType(0x0A):
		c.process0xANNN(nnn)
	case InstructionType(0x0B):
		c.process0xBNNN(nnn)
	case InstructionType(0x0C):
		c.process0xCXNN(x, nn)
	case InstructionType(0x0D):
		c.process0xDXYN(x, y, n)
	case InstructionType(0x0E):
		switch nn {
		case 0x9E:
			c.process0xEX9E(x)
		case 0xA1:
			c.process0xEXA1(x)
		default:
			panic("invalid instruction")
		}
	case InstructionType(0x0F):
		switch nn {
		case 0x07:
			c.process0xFX07(x)
		case 0x0A:
			c.process0xFX0A(x)
		case 0x15:
			c.process0xFX15(x)
		case 0x18:
			c.process0xFX18(x)
		case 0x1E:
			c.process0xFX1E(x)
		case 0x29:
			c.process0xFX29(x)
		case 0x33:
			c.process0xFX33(x)
		case 0x55:
			c.process0xFX55(x)
		case 0x65:
			c.process0xFX65(x)
		default:
			panic("invalid instruction")
		}
	default:
		panic("invalid instruction")
	}

	return nil
}

func (c *Cpu) process0x00E0() {
	c.display.Clear()
	c.display.Flush()
	c.pc += 2
}

func (c *Cpu) process0x00EE() {
	c.sp--
	c.pc = c.stack[c.sp]
}

func (c *Cpu) process0x1NNN(nnn uint16) {
	c.pc = nnn
}

func (c *Cpu) process0x2NNN(nnn uint16) {
	c.stack[c.sp] = c.pc + 0x2
	c.sp++
	c.pc = nnn
}

func (c *Cpu) process0x3XNN(x, nn byte) {
	if c.register[x] != nn {
		c.pc += 2
		return
	}
	c.pc += 4
}

func (c *Cpu) process0x4XNN(x, nn byte) {
	if c.register[x] == nn {
		c.pc += 2
		return
	}
	c.pc += 4
}

func (c *Cpu) process0x5XY0(x, y byte) {
	if c.register[x] != c.register[y] {
		c.pc += 2
		return
	}
	c.pc += 4
}

func (c *Cpu) process0x6XNN(x, nn byte) {
	c.register[x] = nn
	c.pc += 2
}

func (c *Cpu) process0x7XNN(x, nn byte) {
	c.register[x] += nn
	c.pc += 2
}

func (c *Cpu) process0x8XY0(x, y byte) {
	c.register[x] = c.register[y]
	c.pc += 2
}

func (c *Cpu) process0x8XY1(x, y byte) {
	c.register[x] |= c.register[y]
	c.pc += 2
}

func (c *Cpu) process0x8XY2(x, y byte) {
	c.register[x] &= c.register[y]
	c.pc += 2
}

func (c *Cpu) process0x8XY3(x, y byte) {
	c.register[x] ^= c.register[y]
	c.pc += 2
}

func (c *Cpu) process0x8XY4(x, y byte) {
	val := c.register[x] + c.register[y]

	c.register[0xF] = 0x00
	if int(c.register[x])+int(c.register[y]) > 0xFF {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
	c.pc += 2
}

func (c *Cpu) process0x8XY5(x, y byte) {
	val := c.register[x] - c.register[y]

	c.register[0xF] = 0x00
	if int(c.register[x])-int(c.register[y]) >= 0x0 {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
	c.pc += 2
}

func (c *Cpu) process0x8XY6(x, y byte) {
	c.register[0xF] = c.register[x] & 0x01
	c.register[x] >>= 1
	c.pc += 2
}

func (c *Cpu) process0x8XY7(x, y byte) {
	val := c.register[y] - c.register[x]

	c.register[0xF] = 0x00
	if int(c.register[y])-int(c.register[x]) >= 0x0 {
		c.register[0xF] = 0x01
	}

	c.register[x] = val
	c.pc += 2
}

func (c *Cpu) process0x8XYE(x, y byte) {
	c.register[0xF] = (c.register[x] >> 7)
	c.register[x] <<= 1
	c.pc += 2
}

func (c *Cpu) process0x9XY0(x, y byte) {
	if c.register[x] == c.register[y] {
		c.pc += 2
		return
	}
	c.pc += 4
}

func (c *Cpu) process0xANNN(nnn uint16) {
	c.i = nnn
	c.pc += 2
}

func (c *Cpu) process0xBNNN(nnn uint16) {
	c.pc = nnn + uint16(c.register[0])
}

func (c *Cpu) process0xCXNN(x, nn byte) {
	c.register[x] = byte(rand.Intn(0xFF)) & nn
	c.pc += 2
}

func (c *Cpu) process0xDXYN(x, y, n byte) {
	xDisplay, yDisplay := c.register[x], c.register[y]
	colission := false
	for i := 0; i < int(n); i++ {
		sprite := c.memory.LoadSprite(c.i + uint16(i))
		colission = c.display.Draw(xDisplay, yDisplay+byte(i), sprite) || colission
	}

	c.register[0xF] = 0x00
	if colission {
		c.register[0xF] = 0x01
	}

	c.display.Flush()
	c.pc += 2
}

func (c *Cpu) process0xEX9E(x byte) {
	if Key(c.register[x]) != c.keyboard.KeyDown() {
		c.pc += 2
		return
	}
	c.pc += 4
}

func (c *Cpu) process0xEXA1(x byte) {
	if Key(c.register[x]) == c.keyboard.KeyDown() {
		c.pc += 2
		return
	}
	c.pc += 4
}

func (c *Cpu) process0xFX07(x byte) {
	c.register[x] = c.dt
	c.pc += 2
}

func (c *Cpu) process0xFX0A(x byte) {
	// Wait for the key press
	key := c.keyboard.KeyDown()
	for key == Key(0xFF) {
		key = c.keyboard.KeyDown()
	}

	c.register[x] = byte(key)
	c.pc += 2
}

func (c *Cpu) process0xFX15(x byte) {
	c.dt = c.register[x]
	c.pc += 2
}

func (c *Cpu) process0xFX18(x byte) {
	c.st = c.register[x]
	c.pc += 2
}

func (c *Cpu) process0xFX1E(x byte) {
	val := int(c.i) + int(c.register[x])
	c.register[0xF] = 0x00
	if val > 0xFFF {
		c.register[0xF] = 0x01
	}

	c.i += uint16(c.register[x])
	c.pc += 2
}

func (c *Cpu) process0xFX29(x byte) {
	c.i = c.memory.LoadChar(c.register[x])
	c.pc += 2
}

func (c *Cpu) process0xFX33(x byte) {
	c.memory.SaveBCD(c.register[x], c.i)
	c.pc += 2
}

func (c *Cpu) process0xFX55(x byte) {
	c.memory.Save(c.register[0:x+1], c.i)
	c.pc += 2
}

func (c *Cpu) process0xFX65(x byte) {
	c.memory.Load(c.register[0:x+1], c.i)
	c.pc += 2
}
