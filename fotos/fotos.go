package fotos

import (
	"log"
	"sync"
)

func TakePicture() {
	log.Println("[fotos] TakePicture")

	ui := GetUI()
	cam := GetCam()

	cam.Ready()

	ui.Countdown(3)

	img, err := cam.TakePicture()
	if err != nil {
		ui.ShowMsg(Msgs("error.takeimg") + err.Error())
		return
	}

	ui.ShowImage(img)

	err = save(img)

	if err != nil {
		ui.ShowMsg(Msgs("error.saveimg") + err.Error())
		return
	}
}

type Camera interface {
	Ready() // gets called a few seconds before picture is taken
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

type UserInterface interface {
	Countdown(i int)  // should wait i seconds and display countdown
	ShowImage(*Image) // should down a image as result
	ShowMsg(string)   // should show a message string
	SetStatus(string) // should set a (almost) always visible status string

	Run() // start ui
}

var (
	userInterfaces   = make(map[string]UserInterface)
	userInterfacesMu sync.RWMutex
)

func RegisterUI(name string, ui UserInterface) {
	userInterfacesMu.Lock()
	defer userInterfacesMu.Unlock()

	userInterfaces[name] = ui
}

// returns a registerd UI panics if none is present
func anyUI() UserInterface {
	userInterfacesMu.RLock()
	defer userInterfacesMu.RUnlock()

	for _, v := range userInterfaces {
		return v
	}

	Panic("no UserInterface registerd")
	return nil
}

func getUI(s string) UserInterface {
	userInterfacesMu.RLock()
	defer userInterfacesMu.RUnlock()

	return userInterfaces[s]
}

func GetUI() UserInterface {
	var ui UserInterface

	if u, ok := Conf()["ui"]; ok && uiExists(u) {
		ui = getUI(u)
	} else {
		ui = anyUI()
	}

	return ui
}

func uiExists(ui string) bool {
	userInterfacesMu.RLock()
	defer userInterfacesMu.RUnlock()

	_, ok := userInterfaces[ui]

	return ok
}
