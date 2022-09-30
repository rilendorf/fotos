package gpio

import (
	"github.com/DerZombiiie/fotos/fotos"

	"github.com/stianeikeland/go-rpio/v4"

	"log"
	"time"
)

func init() {
	fotos.Runner(func() {
		err := rpio.Open()
		if err != nil {
			log.Printf("Error, cant open raspi/GPIO: %s \nmaybe you aren't on RASPI? \n", err)
			return
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
