package main

import (
	"flag"
	"github.com/DerZombiiie/fotos/ui_neopixel"
	"log"
	"time"
)

var (
	countdownTime = flag.String("time", "5s", "Set time used for countdown")
)

func main() {
	flag.Parse()

	d, err := time.ParseDuration(*countdownTime)
	if err != nil {
		log.Fatalf("Error parsing duration '%s': %s \n", *countdownTime, err)
	}

	pix, err := neopixel.NewPixel()
	if err != nil {
		log.Fatalf("Error creating Neopixel %s \n Maybe you're not on rpi? \n", err)
	}

	for {
		time.Sleep(time.Second * 2)

		log.Printf("Countdown %s \n", d)
		pix.Countdown(d)
		log.Printf("Countdown %s done \n", d)
	}
}
