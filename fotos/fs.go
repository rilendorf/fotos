package fotos

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"time"
)

func save(img *Image) error {
	err := os.Mkdir(SaveDir(), 0755)
	if err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}

	p := SaveDir() + time.Now().Format("2006-01-02.030405.jpg")

	log.Println(fmt.Sprintf("[fotos.save] saving image to %s (%d bytes)", p, len(img.Bytes())))

	err = os.WriteFile(p, img.Bytes(), 0755)
	if err != nil {
		return err
	}

	return err
}
