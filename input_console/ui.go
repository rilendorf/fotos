package input

import (
	"github.com/DerZombiiie/fotos/fotos"

	"time"
)

func init() {
	fotos.Runner(func() {
		go func() {

			time.Sleep(time.Second * 2)

			//fotos.GetUI().Countdown(2)
			//return

			fotos.TakePicture()
		}()
	})
}
