package fotos

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "embed"
)

func save(img *Image) error {
	file := time.Now().Format("2006-01-02.030405.jpg")
	p := SaveDir() + file

	log.Println(fmt.Sprintf("[fotos.save] saving image to %s (%d bytes)", p, len(img.Bytes())))

	err := os.WriteFile(p, img.Bytes(), 0755)
	if err != nil {
		return err
	}

	addImage(file, img)

	return err
}

func SaveDir() string {
	d, ok := Conf()["savedir"]
	if !ok || d == "" {
		log.Fatal("No directory for saving configured")
	}

	if !filepath.IsAbs(d) {
		d = Path(d)
	}

	err := os.Mkdir(d, 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		log.Fatal("[fotos.SaveDir] Directory creation failed. " + err.Error())
	}

	return d
}

//go:embed LICENSE
var license string

func License() string {
	return license
}
