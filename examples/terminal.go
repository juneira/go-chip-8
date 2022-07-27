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
	"time"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

// KeyboardInput implements io.Reader
type KeyBoardInput struct {
}

func (k *KeyBoardInput) Read(p []byte) (n int, err error) {
	p[0] = []byte("5")[0]

	return 1, nil
}

func main() {
	filepath := flag.String("file", "", "path of CHIP-8 program")
	flag.Parse()

	if *filepath == "" {
		panic("param 'file' is required")
	}

	f, _ := os.Open(*filepath)
	buf := bufio.NewReader(f)

	rom := chip8.NewRom(buf)
	memory := chip8.NewStandardMemory(&chip8.ConfigMemory{Rom: rom})

	output := &bytes.Buffer{}
	display := chip8.NewStandardDisplay(&chip8.ConfigDisplay{Output: output})

	keyboard := chip8.NewStandardKeyboard(&chip8.ConfigKeyboard{Input: &KeyBoardInput{}})

	go paintScreen(output)

	cpu := chip8.NewCpu(&chip8.ConfigCpu{
		Display:  display,
		Keyboard: keyboard,
		Memory:   memory,
		PC:       0x200,
	})

	cpu.Start()

	for {
		time.Sleep(5 * time.Millisecond)
		pc := cpu.NextInstruction()
		instr := memory.LoadInstruction(pc)

		k := [0x1FF]byte{}
		memory.Load(k[:], 0)

		cpu.Process(instr)
	}
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
