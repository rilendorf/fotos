package fotos

import (
	"github.com/go-errors/errors"

	"log"
)

func Panic(msg string) {
	err := errors.Errorf(msg)
	if err != nil {
		log.Fatal(err.ErrorStack())
	}
}
