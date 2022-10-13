package fotos

import (
	"sync"
	"time"
)

type UserInterface interface {
	Countdown(i time.Duration) // should wait i seconds and display countdown
	ShowImage(*Image)          // should down a image as result
	ShowMsg(string)            // should show a message string
	SetStatus(string)          // should set a (almost) always visible status string

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

func Countdown(t time.Duration) {
	userInterfacesMu.RLock()
	defer userInterfacesMu.RUnlock()

	wg := sync.WaitGroup{}

	for _, ui := range userInterfaces {
		wg.Add(1)
		go func(ui UserInterface) {
			ui.Countdown(t)
			wg.Done()
		}(ui)
	}

	wg.Wait()
}

func ShowMsg(msg string) {
	userInterfacesMu.RLock()
	defer userInterfacesMu.RUnlock()

	wg := sync.WaitGroup{}

	for _, ui := range userInterfaces {
		wg.Add(1)
		go func(ui UserInterface) {
			ui.ShowMsg(msg)
			wg.Done()
		}(ui)
	}

	wg.Wait()
}

func ShowImage(i *Image) {
	userInterfacesMu.RLock()
	defer userInterfacesMu.RUnlock()

	wg := sync.WaitGroup{}

	for _, ui := range userInterfaces {
		wg.Add(1)
		go func(ui UserInterface) {
			ui.ShowImage(i)
			wg.Done()
		}(ui)
	}

	wg.Wait()

}

func RunUI() {
	userInterfacesMu.RLock()

	var u = make(map[string]UserInterface)
	for k, v := range userInterfaces {
		u[k] = v
	}

	userInterfacesMu.RUnlock()

	var wg sync.WaitGroup

	for _, u := range u {
		wg.Add(1)
		go func(u UserInterface) { u.Run(); wg.Done() }(u)
	}

	wg.Wait()
}
