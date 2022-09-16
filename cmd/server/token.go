package main

import (
	"math/rand"
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
