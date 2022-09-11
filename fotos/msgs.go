package fotos

import (
	"log"
)

func Msgs(name string) string {
	msg, ok := Conf()["msg."+name]

	if ok {
		return msg
	} else {
		log.Println("Msgs: msg '" + name + "' not defined")
		return ""
	}
}
