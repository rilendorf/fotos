package fotos

import (
	"log"
	"sync"
	"time"
)

func TakePicture() {
	log.Println("[fotos] TakePicture")

	cam := GetCam()

	cam.Ready(CountdownTime())

	Countdown(CountdownTime())

	img, err := cam.TakePicture()
	if err != nil {
		ShowMsg(Msgs("error.takeimg") + err.Error())
		return
	}

	ShowImage(img)

	err = save(img)

	if err != nil {
		ShowMsg(Msgs("error.saveimg") + err.Error())
		return
	}
}

type Camera interface {
	Ready(in time.Duration) // gets called a few seconds before picture is taken
	TakePicture() (*Image, error)
}

var (
	cameras   = make(map[string]Camera)
	camerasMu sync.RWMutex
)

func RegisterCam(name string, c Camera) {
	camerasMu.Lock()
	defer camerasMu.Unlock()

	cameras[name] = c
}

func GetCam() Camera {
	c, ok := Conf()["cam"]
	if !ok {
		return anyCam()
	}

	camerasMu.RLock()
	defer camerasMu.RUnlock()

	cam, ok := cameras[c]
	if !ok {
		log.Println("camera dosnt exist using anycam")
		return anyCam()
	}

	return cam
}

func anyCam() Camera {
	camerasMu.RLock()
	defer camerasMu.RUnlock()

	for _, c := range cameras {
		return c
	}

	log.Fatal("No camera module")
	return nil
}
