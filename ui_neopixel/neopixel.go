package neopixel

import (
	"github.com/DerZombiiie/fotos/fotos"
	ws2811 "github.com/rpi-ws281x/rpi-ws281x-go"

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
	opt.Channels[0].LedCount = 32
	opt.Channels[0].GpioPin = 18 // PWD1 / pin33

	cw.WS2811, err = ws2811.MakeWS2811(&opt)
	if err != nil {
		return nil, err
	}

	return cw, cw.Init()
}

func (cw *Neopixel) SetAll(c uint32) {
	l := cw.Leds(0)

	for k, _ := range l {
		l[k] = c
	}
}

func (cw *Neopixel) ShowImage(*fotos.Image) {
}

func (cw *Neopixel) Countdown(duration time.Duration) {
	steps := int(duration.Seconds()) * 100
	l := cw.Leds(0)

	for i := 0; i < steps; i++ {
		m := Smooth(len(l), 1-(float32(i)/float32(steps)))

		for li := range l {
			l[li] = (uint32(m[li])<<8 | uint32(m[li]))
		}

		cw.Render()

		time.Sleep(time.Second / 100)
	}

	time.Sleep(time.Second / 4)

	// flash to signal ok
	for i := 0; i < len(l); i++ {
		l[i] = 0x00FFFF
	}
	cw.Render()

	go func() {
		time.Sleep(time.Second / 4)
		for i := 0; i < len(l); i++ {
			l[i] = 0x000000
		}
		cw.Render()
	}()
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

func init() {
	pix, err := NewPixel()
	if err != nil {
		log.Printf("Error creating Neopixel %s \n Maybe you're not on rpi or you're not root? \n", err)
	}

	fotos.RegisterUI("neopixel", pix)
}
