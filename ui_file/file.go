package neopixel

import (
	"github.com/rilendorf/fotos/fotos"

	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type File struct {
	sync.RWMutex

	Msg, Status string
	File        string
}

func NewFile(f string) (cw *File, err error) {
	return &File{File: f}, nil
}

func (cw *File) WriteFile() {
	log.Print("Updating status file.")

	err := os.WriteFile(
		cw.File,
		[]byte(fmt.Sprintf("MSG: %s| STATUS: %s", cw.Msg, cw.Status)),
		0777,
	)

	if err != nil {
		log.Printf("Failed to write *File to '%s': %s", cw.File, err)
	}
}

func (cw *File) ShowImage(*fotos.Image) {
	// cant show
}

func (cw *File) Countdown(duration time.Duration) {
	// dont care
}

func (cw *File) ShowMsg(msg string) {
	cw.Lock()
	defer cw.Unlock()

	cw.Msg = msg
	cw.WriteFile()
}

func (cw *File) SetStatus(msg string) {
	cw.Lock()
	defer cw.Unlock()

	cw.Status = msg
	cw.WriteFile()
}

func (cw *File) Run() {
	<-make(chan struct{})
}

func init() {
	f, err := NewFile("/tmp/ui_file")
	if err != nil {
		log.Printf("Error creating FileUI: %s", err)
	}

	f.WriteFile()

	fotos.RegisterUI("file", f)
}
