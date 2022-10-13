package main

import (
	"github.com/mitchellh/go-ps"
	"path"

	"log"
	"os"
	"syscall"
)

func getPid() int {
	l, err := ps.Processes()
	if err != nil {
		log.Printf("Error: %s \n", err)
		return -1
	}

	for _, p := range l {
		_, f := path.Split(p.Executable())

		pid := p.Pid()

		if f == "fotos" {
			return pid
		}
	}

	return -1
}

func main() {
	pid := getPid()

	if pid < 0 {
		log.Fatal("Couldn't get PID")
	}

	log.Printf("Found Process, PID: %d \n", pid)

	pr, err := os.FindProcess(pid)
	if err != nil {
		log.Fatal(err)
	}

	err = pr.Signal(syscall.SIGUSR1)
	if err != nil {
		log.Fatal(err)
	}

}
