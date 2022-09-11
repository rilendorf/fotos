package fotos

import (
	"bytes"
	"image"
	"image/jpeg"
	"log"
)

// Image can either contain a []byte with JPEG data
// or a image.Image which can both be semlessly converted between
type Image struct {
	bytes []byte
	image image.Image
}

func ImageFromBytes(b []byte) *Image {
	return &Image{bytes: b}
}

func ImageFromImage(i image.Image) *Image {
	return &Image{image: i}
}

func (i *Image) Bytes() []byte {
	if len(i.bytes) == 0 {
		buf := &bytes.Buffer{}

		err := jpeg.Encode(buf, i.image, nil)
		if err != nil {
			log.Fatal("[fotos.Image] can't encode JPEG buffer: " + err.Error())
		}

		i.bytes = buf.Bytes()
	}

	return i.bytes
}

func (i *Image) Image() image.Image {
	if i.image == nil {
		buf := bytes.NewBuffer(i.bytes)
		img, err := jpeg.Decode(buf)
		if err != nil {
			log.Fatal("[fotos.Image] can't decode JPEG buffer: " + err.Error())
		}

		i.image = img
	}

	return i.image
}
