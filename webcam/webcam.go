package webcam

import (
	"github.com/blackjack/webcam"

	"github.com/rilendorf/fotos/fotos"

	"fmt"
	"log"
	"sort"
	"sync"
)

type Webcam struct {
	*webcam.Webcam
}

func (w Webcam) Ready() {
	w.StartStreaming()

	go w.TrashFrames()
}

var trashFrames chan struct{}
var trashFramesMu sync.RWMutex

func (w Webcam) TrashFrames() {
	trashFramesMu.Lock()
	trashFrames = make(chan struct{})
	trashFramesMu.Unlock()

	for {
		select {
		case <-trashFrames:
			return
		default:
			w.ReadFrame()
		}
	}
}

func (w Webcam) TakePicture() (*fotos.Image, error) {
	w.StartStreaming()
	defer w.StopStreaming()

	trashFramesMu.RLock()
	close(trashFrames)
	trashFramesMu.RUnlock()

	err := w.WaitForFrame(2)
	if err != nil {
		log.Println("[webcam.TakePicture] ", err)
		return nil, err
	}

	frame, err := w.ReadFrame()
	if err != nil {
		log.Println("[webcam.TakePicture] ", err)
		return nil, err
	}

	newimage := make([]byte, len(frame))
	for i := 0; i < len(frame); i++ {
		newimage[i] = frame[i]
	}

	return fotos.ImageFromBytes(newimage), err
}

func init() {
	fotos.Runner(func() {
		cam, err := webcam.Open(fotos.Device())
		if err != nil {
			log.Fatal("[webcam.Init] ", err)
		}

		format_desc := cam.GetSupportedFormats()
		var format webcam.PixelFormat
		var found bool
		for f, desc := range format_desc {
			if desc == "Motion-JPEG" {
				format = f
				found = true
			}
		}

		if !found {
			log.Fatal("[webcam.Init] Format \"Motion-JPEG\" not available on device\n", fotos.Device())
		}

		// Supported Frame Sizes
		frames := Sizes(cam.GetSupportedFrameSizes(format))
		sort.Sort(frames)

		f, w, h, err := cam.SetImageFormat(format, frames[0].MaxWidth, frames[0].MaxHeight)
		if err != nil {
			log.Fatal("[webcam.Init] " + err.Error())
		} else {
			log.Println(fmt.Sprintf("[webcam.Init] Resulting image format: %s (%dx%d)", format_desc[f], w, h))
		}

		fotos.RegisterCam("webcam", Webcam{cam})
	})
}

type Sizes []webcam.FrameSize

func (s Sizes) Len() int {
	return len(s)
}

func (s Sizes) Less(i, j int) bool {
	return s[i].MaxHeight*s[i].MaxWidth > s[j].MaxHeight*s[j].MaxWidth
}

func (s Sizes) Swap(i, j int) {
	tmp := s[i]

	s[i] = s[j]
	s[j] = tmp
}
