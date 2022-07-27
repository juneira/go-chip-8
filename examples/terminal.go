/*
	EXAMPLE OF USAGE ON TERMINAL (LINUX)
*/

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"time"

	chip8 "github.com/MarceloMPJR/go-chip-8"
)

func main() {
	f, _ := os.Open("../roms/space_invaders.ch8")
	buf := bufio.NewReader(f)

	rom := chip8.NewRom(buf)
	memory := chip8.NewStandardMemory(&chip8.ConfigMemory{Rom: rom})

	output := &bytes.Buffer{}
	display := chip8.NewStandardDisplay(&chip8.ConfigDisplay{Output: output})

	keyboard := &chip8.StandardKeyboard{}

	go paintScreen(output)

	cpu := chip8.NewCpu(&chip8.ConfigCpu{
		Display:  display,
		Keyboard: keyboard,
		Memory:   memory,
		PC:       0x200,
	})

	cpu.Start()
	keyboard.Start()

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
