package webcam

import (
	"github.com/DerZombiiie/go-libcamera-example/collector"

	"bytes"
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/DerZombiiie/fotos/fotos"
)

type Webcam struct {
	process *exec.Cmd

	frames chan []byte

	tmpDir  string
	tmpFile string
}

func (w *Webcam) ReadNextFrame() []byte {
	return <-w.frames
}

func (w *Webcam) Ready() {
}

func (w *Webcam) TakePicture() (*fotos.Image, error) {
	w.process.Process.Signal(syscall.SIGUSR1)

	// read image:
	bytes := w.ReadNextFrame()

	return fotos.ImageFromBytes(bytes), nil
}

func init() {
	fotos.Runner(func() {
		cmd := exec.Command("libcamera-still", "-t0", "-s", "-o-", "-ejpg")

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
