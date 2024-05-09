package input_ttyusb

import (
	"github.com/DerZombiiie/fotos/fotos"

	"errors"
	"fmt"
	"github.com/tarm/serial"
	"io"
	"log"
	"time"
)

func init() {
	// do things and not block
	go do()
}

func Now() int64 {
	return time.Now().Unix()
}

func do() {
	LastPress := Now()
	
	Name := "/dev/ttyUSB0"
	NameI := 0

outer:
	for {
		Name = fmt.Sprintf("/dev/ttyUSB%d", NameI)
		NameI++
		if NameI > 5 {
			NameI = 0
		}

		time.Sleep(time.Millisecond * 250)

		config := &serial.Config{
			Name:        Name,
			Baud:        9600,
			ReadTimeout: 1,
			Size:        8,
		}

		stream, err := serial.OpenPort(config)
		if err != nil {
			log.Printf("Failed to open '%s': %s", Name, err)
			continue outer
		}

		log.Printf("Opened port '%s'", Name)

		buf := make([]byte, 1024)

		var errsthissec = 0
		var lasttime = time.Now()

	inner:
		for true {
			n, err := stream.Read(buf)
			if err != nil {
				if errors.Is(err, io.EOF) {
					errsthissec++

					// second passed
					if lasttime.Add(time.Second).Unix() < time.Now().Unix() {
						//log.Printf("%d errors this sec", errsthissec)
						if errsthissec > 50 {
							log.Printf("probably disconnect (%d > 50 errs in a sec)", errsthissec)
							continue outer
						}

						errsthissec = 0
						lasttime = time.Now()
					}

					continue inner
				}

				log.Printf("Failed to read: %s", err)
				continue outer
			}

			if n < 2 {
				n = 2 // ensure no negative indexes
			}
			s := string(buf[:n-2])

			switch s {
			case "PRESS":
				log.Printf("PRESSED! :D (last %d; now %d)", LastPress, Now())
				if LastPress < Now() {
					log.Printf("scheduling pic")
					go fotos.TakePicture()
				}

				LastPress = Now() + 5

			default:
				log.Printf("Failed to parse input '%s'", s)

			}
		}
	}
}
