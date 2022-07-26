package chip8

import (
	"io"
)

const white = "□"
const black = "■"

// StandardDisplay implements interface Display
type StandardDisplay struct {
	output io.Writer
	screen [31][63]byte
}

type ConfigDisplay struct {
	Output io.Writer
}

// NewStandardDisplay is a function that receive a config as param and return a pointer to StandardDisplay
func NewStandardDisplay(config *ConfigDisplay) *StandardDisplay {
	return &StandardDisplay{output: config.Output}
}

// Flush is a function that paint the screen with information of attribute "screen"
func (sd *StandardDisplay) Flush() {
	buf := ""
	for i := 0; i < 31; i++ {
		for j := 0; j < 63; j++ {
			if sd.screen[i][j] == 1 {
				buf += black
			} else {
				buf += white
			}
		}
		buf += "\n"
	}
	sd.output.Write([]byte(buf))
}

// Clear sets all pixels to 0
func (sd *StandardDisplay) Clear() {
	for i := 0; i < 31; i++ {
		for j := 0; j < 63; j++ {
			sd.screen[i][j] = 0
		}
	}
}

// Draw draws a sprint on position xDisplay and yDisplay
func (sd *StandardDisplay) Draw(xDisplay, yDisplay, sprite byte) bool {
	collision := false

	for bitIdx := byte(0); bitIdx < 8; bitIdx++ {
		newPixel := (sprite & (1 << (7 - bitIdx))) >> (7 - bitIdx)
		oldPixel := sd.screen[yDisplay][xDisplay+bitIdx]

		if newPixel != 0 && oldPixel != 0 {
			collision = true
		}

		sd.screen[yDisplay][xDisplay+bitIdx] = newPixel ^ oldPixel
	}

	return collision
}
