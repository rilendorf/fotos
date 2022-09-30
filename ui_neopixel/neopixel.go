package neopixel

import (
	"github.com/DerZombiiie/fotos/fotos"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"

	"fmt"
	"log"
	"time"
)

const (
	onColor    = uint32(0x0000FF) // some blue
	offColor   = uint32(0x000000) // off
	dimColor   = uint32(0x00007F) // some lighter / less bright blue
	flashColor = uint32(0xFFFFFF) // some lighter / less bright blue
)

type Neopixel struct {
	*ws2811.WS2811
}

func NewPixel() (cw *Neopixel, err error) {
	cw = &Neopixel{}

	opt := ws2811.DefaultOptions
	opt.Channels[0].Brightness = 90
	opt.Channels[0].LedCount = 16
	opt.Channels[0].GpioPin = 12

	cw.WS2811, err = ws2811.MakeWS2811(&opt)

	return cw, cw.Init()
}

func (cw *Neopixel) SetAll(c uint32) {
	l := cw.Leds(0)

	for k, _ := range l {
		l[k] = c
	}
}

func (cw *Neopixel) ShowImage(*fotos.Image) {
	go func() {
		cw.SetAll(flashColor)
		time.Sleep(time.Millisecond * 250)
		cw.SetAll(0)
	}()
}

func (cw *Neopixel) SetProgress(c, ttl int) error {
	l := cw.Leds(0)

	leds := int(float32(c) / float32(ttl) * float32(ttl))

	for i := 0; i < ttl; i++ {
		if i+1 < leds {
			l[i] = onColor
		} else {
			l[i] = offColor
		}
	}

	if leds > 0 && leds < ttl {
		l[leds-1] = dimColor
	}

	if leds == 0 {
		l[0] = 0
	}

	if err := cw.Render(); err != nil {
		return err
	}

	return nil
}

func (cw *Neopixel) Countdown(duration time.Duration) {
	if err := cw.countdown(duration); err != nil {
		log.Printf("Error starting countdown: %s \n", err)
	}
}

func (cw *Neopixel) ShowMsg(msg string) {
	log.Printf("Msg: %s \n", msg)
}

func (cw *Neopixel) SetStatus(msg string) {
	log.Printf("Status: %s \n", msg)
}

func (cw *Neopixel) Run() {
	cw.Countdown(time.Second * 2)

	<-make(chan struct{})
}

func (cw *Neopixel) countdown(duration time.Duration) error {
	leds := len(cw.Leds(0))

	tick := duration / time.Duration(leds)

	fmt.Printf("Tick Duration is %s \n", tick)

	for d := 0; d < leds; d++ {
		if err := cw.SetProgress(leds-d, leds); err != nil {
			return err
		}

		time.Sleep(tick)
	}

	cw.SetProgress(0, 16)

	return nil
}

func init() {
	pix, err := NewPixel()
	if err != nil {
		log.Printf("Error creating Neopixel %s \n Maybe you're not on rpi? \n", err)
	}

	fotos.RegisterUI("neopixel", pix)
}
