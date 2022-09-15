package gpio

import (
	"github.com/DerZombiiie/fotos/fotos"

	"github.com/stianeikeland/go-rpio/v4"

	"time"
)

func init() {
	fotos.Runner(func() {
		err := rpio.Open()
		if err != nil {
			panic(err)
		}

		pin := rpio.Pin(17)

		pin.PullUp()
		pin.Input()
		pin.Detect(rpio.RiseEdge)

		for {
			if pin.EdgeDetected() {
				fotos.TakePicture()
			}

			time.Sleep(time.Second / 2)
		}
	})
}
