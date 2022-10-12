package gpio

import (
	"github.com/DerZombiiie/fotos/fotos"

	"github.com/stianeikeland/go-rpio/v4"

	"log"
	"time"
)

type Pin struct {
	rpio.Pin
}

var pin Pin

func init() {
	fotos.Runner(func() {
		err := rpio.Open()
		if err != nil {
			log.Printf("Error, cant open raspi/GPIO: %s \nmaybe you aren't on rpi? \n", err)
			return
		}

		pin = Pin{rpio.Pin(16)}

		pin.Output()

		fotos.RegisterUI("flash", pin)
	})
}

func (p Pin) Run()                   {}
func (p Pin) SetStatus(string)       {}
func (p Pin) ShowMsg(string)         {}
func (p Pin) ShowImage(*fotos.Image) {}

func (p Pin) Countdown(d time.Duration) {
	log.Println("flash countdown!")
	p.Write(rpio.High)
	time.Sleep(d)

	go func() {
		time.Sleep(time.Second * 4)

		log.Println("flash countdown done")
		pin.Write(rpio.Low)
	}()
}
