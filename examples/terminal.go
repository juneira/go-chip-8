/*
	EXAMPLE OF USAGE ON TERMINAL (LINUX)
*/

package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

// KeyboardInput implements io.Reader
type KeyBoardInput struct{}

func (k *KeyBoardInput) Read(p []byte) (n int, err error) {
	p[0] = []byte("5")[0]

	return 1, nil
}

// FakeSound implements chip8.Sound
type FakeSound struct{}

func (fs *FakeSound) Beep() {
}

func main() {
	filepath := flag.String("file", "", "path of CHIP-8 program")
	flag.Parse()

	if *filepath == "" {
		panic("param 'file' is required")
	}

	f, _ := os.Open(*filepath)
	buf := bufio.NewReader(f)

	memory := chip8.NewStandardMemory(&chip8.ConfigMemory{Rom: buf})

	output := &bytes.Buffer{}
	display := chip8.NewStandardDisplay(&chip8.ConfigDisplay{Output: output})

	keyboard := chip8.NewStandardKeyboard(&chip8.ConfigKeyboard{Input: &KeyBoardInput{}})

	go paintScreen(output)

	cpu := chip8.NewCpu(&chip8.ConfigCpu{
		Display:  display,
		Keyboard: keyboard,
		Sound:    &FakeSound{},
		Memory:   memory,
		PC:       0x200,
	})

	cpu.Start()
}

func paintScreen(screenBuffer *bytes.Buffer) {
	for {
		if len(screenBuffer.Bytes()) > 0 {
			// Clean Screen
			runCmd("clear")

			fmt.Print(screenBuffer.String())
			screenBuffer.Reset()
		}
	}
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}
