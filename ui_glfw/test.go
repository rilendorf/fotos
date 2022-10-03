package ui

import (
	"github.com/DerZombiiie/fotos/fotos"
	"github.com/DerZombiiie/fotos/ui_glfw/textimage"

	"fmt"
	"image"
	"math"
	"sync"
	"time"
)

const finalText = "cheese"

type ui_glfw struct {
	cache   map[string]*image.RGBA
	cacheMu sync.RWMutex

	width, height int
}

func (ui *ui_glfw) GetImage(text string) *image.RGBA {
	ui.cacheMu.Lock()
	defer ui.cacheMu.Unlock()

	if img, ok := ui.cache[text]; ok {
		return img
	}

	img := textimage.GenerateImage(ui.width, ui.height, text)
	ui.cache[text] = img

	return img
}

func (ui *ui_glfw) SetStatus(string)       {}
func (ui *ui_glfw) ShowMsg(string)         {}
func (ui *ui_glfw) ShowImage(*fotos.Image) {}

func (ui *ui_glfw) Countdown(d time.Duration) {
	countdown := int(math.Floor(d.Seconds()))

	if countdown <= 0 {
		countdown = 3
	}
	fmt.Printf("Counting down from %d \n", countdown)

	for {
		img := ui.GetImage(fmt.Sprintf("%d", countdown))

		time.Sleep(time.Second)
		countdown--

		if countdown < 0 {
			break
		}

		imageUpdates <- img
	}

	//time.Sleep(time.Second)
	imageUpdates <- ui.GetImage("cheese")
	time.Sleep(time.Second / 4)

	go func() {
		time.Sleep(time.Second / 4 * 3)

		imageUpdates <- ui.GetImage("")
	}()
}

func init() {
	ui := &ui_glfw{
		cache: make(map[string]*image.RGBA),

		width:  720,
		height: 480,
	}

	fotos.RegisterUI("glfw", ui)
}
