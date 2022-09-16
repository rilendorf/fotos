package main

import (
	"archive/zip"
	"io"
	"log"
)

func createZip(acc string, w io.Writer) error {
	zipW := zip.NewWriter(w)

	list := readImages(acc)

	for _, name := range list {
		f, err := zipW.Create(name)
		if err != nil {
			log.Printf("Error creating zip file, can't create file %s: %s\n", name, err)
			continue
		}

		_, err = f.Write(readImage(acc, name))
		if err != nil {
			log.Printf("Error crating zip file, can't add image %s: %s\n", name, err)
			continue
		}
	}

	err := zipW.Close()
	if err != nil {
		log.Printf("Error creating zip file, error closing: %s", err)
		return err
	}

	return nil
}
