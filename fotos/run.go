package fotos

import (
	"log"
	"sync"
)

func Run() {
	runOnce.Do(run)
}

var runOnce sync.Once

func run() {
	println("run")

	err := readConfig()
	if err != nil {
		log.Fatal("Error reading config " + err.Error())
	}

	runnersMu.RLock()
	for _, r := range runners {
		go r()
	}
	runnersMu.RUnlock()

	// start ui
	RunUI()

	println("exit")
}

var (
	runners   []func()
	runnersMu sync.RWMutex
)

func Runner(r func()) {
	runnersMu.Lock()
	defer runnersMu.Unlock()

	runners = append(runners, r)
}
