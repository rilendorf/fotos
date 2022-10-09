package main

import (
	"github.com/DerZombiiie/fotos/ui_neopixel"
	"log"
	"time"
)

func main() {
	pix, err := neopixel.NewPixel()
	if err != nil {
		log.Printf("Error creating Neopixel %s \n Maybe you're not on rpi? \n", err)
	}

	for {
		time.Sleep(time.Second * 5)

		log.Println("Countdown 10")
		pix.Countdown(time.Second * 10)
		log.Println("Countdown 10 done")
	}
}
