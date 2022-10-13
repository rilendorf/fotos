package signal

import (
	"github.com/DerZombiiie/fotos/fotos"

	"log"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	fotos.Runner(func() {
		log.Printf("initializing input_signal")

		sigChannel := make(chan os.Signal, 1)
		signal.Notify(sigChannel, syscall.SIGUSR1, syscall.SIGUSR2)
		select {
		case s := <-sigChannel:
			log.Printf("received signal: %s \n", s)
			fotos.TakePicture()
		}
	})
}
