package libcamera

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

const file = "tmpimage.jpg"

type Libcamera struct {
	takeImage chan struct{}
	recImage  chan []byte
}

func (c *Libcamera) CameraThread() {
	c.takeImage = make(chan struct{})
	c.recImage = make(chan []byte)

	l, err := ps.Processes()
	if err != nil {
		log.Fatalf("[libcamera] Error aquiring process list: %s \n", err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("[libcamera] Error creating Wacher: %s \n", err)
		return
	}

	watcher.Add("image.jpg")

	for _, p := range l {
		_, f := path.Split(p.Executable())

		pid := p.Pid()

		if f == "libcamera-still" {
			log.Printf("[libcamera] Found libcamera-still process (PID: %d) \n", pid)

			if err != nil {
				log.Fatalf("[libcamera] Error Finding Process (PID: %d): %s \n", pid, err)
				return
			}

			for {
				log.Println("[libcamera] waiting for channel takeImage")
				<-c.takeImage

				log.Printf("[libcamera] sending signal SIGUSR1 to PID %d \n", pid)

				pr, err := os.FindProcess(pid)
				if err != nil {
					log.Fatal(err)
				}

				pr.Signal(syscall.SIGUSR1)

			loop:
				for {
					log.Println("[libcamera] Waiting for events")
					c := time.NewTimer(time.Second * 3).C

					select {
					case <-c:
						log.Println("[libcamera] Timeout waiting for events!")
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
						log.Printf("[libcamera] watcher: error: %s \n", e)
					}
				}

				log.Println("Renaming...")
				err = os.Rename("image.jpg", file)
				if err != nil {
					log.Printf("[libcamera] Error renaming: %s \n", err)
				}

				b, err := os.ReadFile(file)
				if err != nil {
					log.Printf("[libcamera] error reading file %s: %s \n", file, err)
				}

				c.recImage <- b

				// sleep, to let the camera reset
				time.Sleep(time.Second * 2)
			}
		}
	}

	log.Fatal("[libcamera] No 'libcamera-still' process found!")
}

func (c *Libcamera) Ready(time.Duration) {}
func (c *Libcamera) TakePicture() (*fotos.Image, error) {
	c.takeImage <- struct{}{}

	// read image:
	bytes := <-c.recImage
	if bytes == nil || len(bytes) == 0 {
		return nil, ErrNoBytes
	}

	return fotos.ImageFromBytes(bytes), nil
}

func init() {
	fotos.Runner(func() {
		c := &Libcamera{}

		// Own camera thread because i want to make sure only one picture is taken at a time
		go c.CameraThread()

		fotos.RegisterCam("libcamera", c)
	})
}
