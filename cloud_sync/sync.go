package cloud

import (
	"github.com/DerZombiiie/fotos/fotos"

	"bytes"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func init() {
	fotos.Runner(func() {
		sync()

		go func() {
			t := time.NewTicker(time.Minute * 10)

			for {
				_, ok := <-t.C
				if !ok {
					return
				}

				log.Println("Start periodic check to sync images")
				sync()
				log.Println("Done with periodic check to sync images")
			}
		}()

		_, ok := fotos.Conf()["printing.printer"]
		if ok {
			log.Printf("Printing enabled")

			go func() {
				t := time.NewTicker(time.Second * 10)

				for {
					_, ok := <-t.C
					if !ok {
						return
					}

					log.Println("Check for printing requests")
					checkprintingrequests()
					log.Println("Done with printing requests")
				}
			}()
		} else {
			log.Printf("Printing disabled")

		}
	})

	fotos.ImageUpdater(func(name string, img *fotos.Image) {
		Upload(name, bytes.NewBuffer(img.Bytes()))
	})
}

func checkprintingrequests() {
	reqests := GetPrintingRequests()
	if reqests == nil {
		return
	}

	for _, r := range reqests {
		start := time.Now()
		fotos.ShowMsg("printing " + r)

		log.Printf("Printing %s", r)
		args := []string{"-d", fotos.Conf()["printing.printer"],
			"-o", "media=" + fotos.Conf()["printing.media"],
			fotos.SaveDir() + r,
		}
		log.Printf("Printing lp %v", args)

		cmd := exec.Command("lp", args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			log.Printf("Printing Run failed with %s", err)
		}

		fotos.ShowMsg("Printed " + r + " in " + time.Now().Sub(start).String())
		SetPrinted(r)
	}
}

func sync() {
	list := listPictures()
	names := make([]string, len(list))
	for k, v := range list {
		names[k] = v.Name()
	}

	for _, pic := range FilterNeeded(names) {
		f, err := os.Open(fotos.SaveDir() + pic)
		if err != nil {
			continue
		}

		Upload(pic, f)

		f.Close()
	}
}

func listPictures() (s []os.FileInfo) {
	filepath.Walk(fotos.SaveDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			s = append(s, info)
		}

		return nil
	})

	return
}
