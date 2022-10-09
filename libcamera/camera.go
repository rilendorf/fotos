package webcam

import (
	"github.com/DerZombiiie/go-libcamera-example/collector"

	"bytes"
	"errors"

	"log"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/DerZombiiie/fotos/fotos"
)

type Status uint8

const (
	Idle Status = iota
	Ready
	WaitingFrame
	Sleep
)

type Webcam struct {
	process *exec.Cmd
	status  Status

	frames chan []byte

	tmpDir  string
	tmpFile string
}

func (w *Webcam) ReadNextFrame() []byte {
	return <-w.frames
}

func (w *Webcam) Ready() {
	if w.status == Idle {
		w.status = Ready
	} else {
		log.Println("Ready pressed too often in a short timespan, ignoring")
	}
}

var (
	ErrInvalidStatus = errors.New("Invalid Status while taking Picture")
)

func (w *Webcam) TakePicture() (*fotos.Image, error) {
	if w.status != Ready {
		log.Printf("Trying to take picture with status != Ready, status: %s \n", w.status.String())
		return nil, ErrInvalidStatus
	}

	w.status = WaitingFrame

	w.process.Process.Signal(syscall.SIGUSR1)

	// read image:
	bytes := w.ReadNextFrame()
	w.status = Sleep
	go func() {
		time.Sleep(time.Millisecond * 100)

		w.status = Idle
	}()

	return fotos.ImageFromBytes(bytes), nil
}

func init() {
	fotos.Runner(func() {
		cmd := exec.Command("libcamera-still", "-t0", "-s", "-o-", "-ejpg", "--vflip", "-f")

		ch := make(chan []byte)

		cmd.Stdout = collector.MakeCollector(time.Millisecond*250, func(buf *bytes.Buffer) {
			ch <- buf.Bytes()
		})

		cmd.Stderr = os.Stderr

		cmd.Start()

		fotos.RegisterCam("libcamera", &Webcam{
			process: cmd,
			frames:  ch,
		})
	})
}
