package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"github.com/disintegration/imaging"
	"image"
	"image/jpeg"
	"log"
)

var sha1hasher = sha1.New()

func hash(b []byte) (r [20]byte) {
	h := sha1hasher.Sum(b)

	for k := range r {
		r[k] = h[k]
	}

	return
}

const previewWidth = 512

var previewCache = make(map[[20]byte][]byte)

func PreviewImage(h []byte) []byte {
	ha := hash(h)

	if d, k := previewCache[ha]; k {
		return d
	} else {
		log.Printf("[ALBUM] Image not in cache, generating '%s' \n", hex.EncodeToString(ha[:]))
		image, _, err := image.Decode(bytes.NewBuffer(h))
		if err != nil {
			return h
		}

		b := image.Bounds()
		width, height := b.Max.X-b.Min.X, b.Max.Y-b.Min.Y
		if width < previewWidth || height < previewWidth { // image is small enough already
			buf := &bytes.Buffer{}
			jpeg.Encode(buf, image, &jpeg.Options{Quality: 50})

			previewCache[ha] = buf.Bytes()
			return buf.Bytes()

		}

		image = imaging.Resize(image,
			previewWidth,
			int(float32(height)/float32(width)*previewWidth),
			imaging.Box,
		)

		buf := &bytes.Buffer{}
		jpeg.Encode(buf, image, &jpeg.Options{Quality: 50})
		previewCache[ha] = buf.Bytes()

		return buf.Bytes()
	}
}
