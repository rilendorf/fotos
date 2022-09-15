package fotos

import (
	"sync"
)

var (
	imageUpdaters   []func(string, *Image)
	imageUpdatersMu sync.RWMutex
)

func addImage(name string, img *Image) {
	imageUpdatersMu.RLock()
	defer imageUpdatersMu.RUnlock()

	for _, f := range imageUpdaters {
		f(name, img)
	}
}

func ImageUpdater(f func(string, *Image)) {
	imageUpdatersMu.Lock()
	defer imageUpdatersMu.Unlock()

	imageUpdaters = append(imageUpdaters, f)
}
