package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("i3status")
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to open stdout pipe: %s", err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to start i3status: %s", err)
	}

	s := bufio.NewScanner(pipe)

	for s.Scan() {
		c,err := os.ReadFile("/tmp/ui_file")
		if err != nil {
			log.Printf("Failed to read file from ui_file: %s", err)

			c = []byte("*no ui_file*")
		}
		
		fmt.Printf("%s | %s\n", c, s.Text())
	}
}
