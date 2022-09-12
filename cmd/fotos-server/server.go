package main

import (
	"flag"
	"fmt"
)

var flagUsername = flag.String("add-user", "", "Specify to generate a new user account")
var flagGetToken = flag.String("get-token", "", "Specify to print token of account")

func main() {
	openDB("users.sqlite")

	flag.Parse()

	if *flagUsername != "" {
		fmt.Printf("adding user with name %s!\n", *flagUsername)
		tkn := genTkn()
		addUser(*flagUsername, tkn)
		fmt.Println("done, token: " + tkn)

		return
	}

	if *flagGetToken != "" {
		fmt.Printf("querying token of user %s!\n", *flagGetToken)
		tkn := getToken(*flagGetToken)
		fmt.Println("done, token: " + tkn)

		return
	}
}
