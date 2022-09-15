package main

import (
	"bytes"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

func genTkn() (s string) {
	rand.Seed(time.Now().UnixNano())

	for i := 0; i < 25; i++ {
		s += string(alphabet[rand.Intn(len(alphabet))])
	}

	return
}

func readToken() {
	f, err := os.OpenFile("token.txt", os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		log.Fatal("Error Opening token file!")
	}

	info, err := f.Stat()
	if err != nil {
		log.Fatal("Error getting info for token.txt")
	}
	if info.Size() == 0 {
		log.Println("Token size is 0, generating new one")

		f.Write([]byte(genTkn()))
		f.Seek(0, 0)
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, f)

	token = buf.String()
}
