package web

import (
	"github.com/DerZombiiie/fotos/album"
	"github.com/DerZombiiie/fotos/fotos"

	"log"
	"os"
	"path/filepath"
)

func listImages() (s []string) {
	filepath.Walk(fotos.SaveDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			s = append(s, info.Name())
		}

		return nil
	})

	return
}

func getImage(name string) (b []byte) {
	b, err := os.ReadFile(fotos.SaveDir() + name)
	if err != nil {
		log.Printf("error readimg image with name %s\n", name)
	}

	return
}

func init() {
	fotos.Runner(func() {
		c := fotos.Conf()
		addr, ok := c["http.album.addr"]
		if !ok {
			addr = "localhost:8080"
		}

		album.Listen(listImages, getImage, fotos.Msgs, addr)
	})
}
