package signal

import (
	"github.com/DerZombiiie/fotos/fotos"

	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func init() {
	var lastimage = time.Now()

	fotos.Runner(func() {
		log.Println("Initializing input_signal")

		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGUSR1, syscall.SIGUSR2)
		for {
			if time.Now().Sub(lastimage) > time.Second*15 {
				s := <-sigChannel
				log.Printf("received signal: %s \n", s)
				go fotos.TakePicture()

				lastimage = time.Now()
			}
		}
	})
}
