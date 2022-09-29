package webcam

import (
	"bytes"
	"errors"
	"github.com/fsnotify/fsnotify"
	"os"
	"os/exec"
	"syscall"

	"github.com/DerZombiiie/fotos/fotos"

	"fmt"
	"io"
	"log"
	"sync"
)

type Webcam struct {
	process *exec.Cmd
	Stdin   io.Writer
	Stdout  *bytes.Buffer
	watcher *fsnotify.Watcher

	frames <-chan []byte

	tmpDir  string
	tmpFile string
}

func (w *Webcam) ReadNextFrame() []byte {
	return <-w.frames
}

func (w *Webcam) Ready() {
}

func Unwrap(err error) error {
	u := errors.Unwrap(err)

	if u == nil {
		return err
	} else {
		return Unwrap(u)
	}
}

func (w *Webcam) TakePicture() (*fotos.Image, error) {
	w.process.Process.Signal(syscall.SIGUSR1)

	// read image:
	bytes := w.ReadNextFrame()

	return fotos.ImageFromBytes(bytes), nil
}

type PrefixWriter struct {
	W      io.Writer
	Prefix string

	notInline bool

	mu sync.Mutex
}

// TODO: make better:
func (pw *PrefixWriter) Write(b []byte) (n int, err error) {
	pw.mu.Lock()
	defer pw.mu.Unlock()

	n, err = pw.W.Write(append([]byte(pw.Prefix), b...))
	if n-len(pw.Prefix) <= 0 {
		return 0, err
	} else {
		return n - len(pw.Prefix), err
	}
}

func init() {
	fotos.Runner(func() {
		tmpDir := os.TempDir()
		tmpFile := "tempimage.jpeg"

		cmd := exec.Command("libcamera-still", "-t0", "-s", "-o"+tmpDir+"/"+tmpFile)

		//cmd.Env = append(cmd.Env, "DISPLAY="+os.Getenv("DISPLAY"))

		stdin := &bytes.Buffer{}
		cmd.Stdin = stdin
		cmd.Stderr = &PrefixWriter{W: os.Stderr, Prefix: "[LIBCAMERA-STILL] "}

		go cmd.Run()

		// place watchers on file:
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal("failed to iniitialize fsnotify watcher")
		}

		watcher.Add(tmpDir + "/" + tmpFile)

		ch := make(chan []byte)

		go func(ch chan<- []byte, w *fsnotify.Watcher) {
			for {
				select {
				case event := <-watcher.Events:
					fmt.Printf("event: %s \n", event.String())

					bytes, err := os.ReadFile(tmpDir + "/" + tmpFile)
					if err != nil {
						log.Printf("Error reading tmpfile (event: %s): %s \n", event.String(), err)
					}

					ch <- bytes

				case err := <-watcher.Errors:
					log.Printf("Error in watcher: %s \n", err)
				}
			}
		}(ch, watcher)

		fotos.RegisterCam("libcamera", &Webcam{
			process: cmd,
			Stdin:   stdin,
			watcher: watcher,

			frames: ch,

			tmpDir:  tmpDir,
			tmpFile: tmpFile,
		})
	})
}
