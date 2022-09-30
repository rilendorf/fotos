package ui

import (
	"github.com/DerZombiiie/fotos/fotos"
	"github.com/rpi-ws281x/rpi-ws281x-go"

	"log"
	"time"
)

type Neopixel struct {
	Pin   int
	Count int

	ws *ws2811.WS2811
}

func (c *Neopixel) Countdown(i int) {
	if i < 0 {
		i = 3
	}

	for i > 0 {
		time.Sleep(time.Second)
	}

	i--
}

func (c *Neopixel) ShowImage(img *fotos.Image) {
}

func (c *Neopixel) ShowMsg(msg string) {
}

func (c *Neopixel) SetStatus(str string) {
}

func (c *Neopixel) SetColor(col uint32) error {
	for i := 0; i < len(c.ws.Leds(0)); i++ {
		c.ws.Leds(0)[i] = col
	}

	return c.ws.Render()
}

func (c *Neopixel) Run() {
	if err := c.SetColor(0x0000ff); err != nil {
		log.Println("error setting init color")
	}
}

func MakeNeopixel(pin, count int) (n *Neopixel, err error) {
	n = &Neopixel{
		Pin:   pin,
		Count: count,
	}

	n.ws, err = ws2811.MakeWS2811(&ws2811.Option{
		Channels: []ws2811.ChannelOption{
			ws2811.ChannelOption{
				GpioPin:  n.Pin,
				LedCount: n.Count,
			},
		},
	})
	if err != nil {
		return
	}

	if err = n.ws.Init(); err != nil {
		return
	}

	return
}

func init() {
	n, err := MakeNeopixel(18, 16)
	if err != nil {
		log.Printf("Cant make neopixel: %s \n", err)
	}

	fotos.RegisterUI("neopixel", n)
}
