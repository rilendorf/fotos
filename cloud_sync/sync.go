package cloud

import (
	"github.com/DerZombiiie/fotos/fotos"

	"bytes"
	"log"
	"os"
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

				log.Println("Done with periodic check to sync images")
				sync()
				log.Println("Done with periodic check to sync images")
			}
		}()
	})

	fotos.ImageUpdater(func(name string, img *fotos.Image) {
		Upload(name, bytes.NewBuffer(img.Bytes()))
	})
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
