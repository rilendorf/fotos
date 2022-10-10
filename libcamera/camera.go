package webcam

import (
	"github.com/fsnotify/fsnotify"
	"github.com/mitchellh/go-ps"
	"path"

	"log"
	"os"
	"syscall"
	"time"

	"github.com/DerZombiiie/fotos/fotos"
)

type Camera struct {
	takeImage chan struct{}
	recImage  chan []byte
}

func (c *Camera) CameraThread() {
	c.takeImage = make(chan struct{})
	c.recImage = make(chan []byte)

	l, err := ps.Processes()
	if err != nil {
		log.Printf("Error: %s \n", err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Printf("Error: %s \n", err)
		return
	}

	watcher.Add("image.jpg")

	for _, p := range l {
		_, f := path.Split(p.Executable())

		pid := p.Pid()

		if f == "libcamera-still" {
			log.Printf("Found libcamera-still process (PID: %d) \n", pid)

			if err != nil {
				log.Printf("Error Finding Process (PID: %d): %s \n", pid, err)
				return
			}

			for {
				log.Println("waiting for channel takeImage")
				<-c.takeImage

				log.Printf("sending signal to PID %d \n", pid)

				pr, err := os.FindProcess(pid)
				if err != nil {
					log.Fatal(err)
				}

				pr.Signal(syscall.SIGUSR1)

			loop:
				for {
					log.Println("Waiting for events")
					c := time.NewTimer(time.Second * 3).C

					select {
					case <-c:
						log.Println("Timeout waiting for events!")
						break loop

					case e := <-watcher.Events:

						log.Printf("Event: %s, op: %s \n", e, e.Op)

						switch e.Op {
						case fsnotify.Remove:
						case fsnotify.Rename:
						case fsnotify.Chmod:
							continue

						case fsnotify.Write:
							fallthrough
						case fsnotify.Create:
							time.Sleep(time.Second * 2)

							break loop
						}
					case e := <-watcher.Errors:
						log.Printf("watcher: error: %s \n", e)
					}
				}

				file := "tmpimage.jpg"

				log.Println("Renaming...")
				err = os.Rename("image.jpg", file)
				if err != nil {
					log.Printf("Error renaming: %s \n", err)
				}
				log.Println("Done.")

				b, err := os.ReadFile(file)
				if err != nil {
					c.recImage <- []byte{}
					log.Printf("error reading file %s: %s \n", file, err)
				}

				c.recImage <- b

				// timeout
				time.Sleep(time.Second * 2)
				log.Println("done taking image, returning to:")
			}
		}
	}

	log.Fatal("No 'libcamera-still' process found!")
}

func (c *Camera) Ready(time.Duration) {}
func (c *Camera) TakePicture() (*fotos.Image, error) {
	c.takeImage <- struct{}{}

	// read image:
	bytes := <-c.recImage

	return fotos.ImageFromBytes(bytes), nil
}

func init() {
	fotos.Runner(func() {
		c := &Camera{}

		go c.CameraThread()

		fotos.RegisterCam("libcamera", c)
	})
}
