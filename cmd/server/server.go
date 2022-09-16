package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"flag"
)

var Account = flag.String("account", "", "specify account")
var ListImages = flag.Bool("list-images", false, "if true, programm lists name of all images and quits")
var RmImage = flag.String("rm-image", "", "delete a image")
var ReadImage = flag.String("read-image", "", "read image to stdout")

var AddAccount = flag.String("add-account", "", "add a new account")
var RmAccount = flag.String("rm-account", "", "remove a account")
var GetToken = flag.String("get-token", "", "get a accounts token")

var DatabaseFile = flag.String("db-file", "fotos.sqlite", "specify sqlite file used")

var addr = ":5050"
var pubAccess = "http://127.0.0.1:5050"

func main() {
	openDB(*DatabaseFile)

	flag.Parse()

	if *GetToken != "" {
		tkn, viewtkn := getTokenByName(*GetToken)
		if tkn == "" {
			fmt.Printf("Account %s dosn't exist!\n", *GetToken)
			return
		}
		fmt.Printf("%s has token %s and viewtoken %s\n", *GetToken, tkn, viewtkn)
		return
	}

	if *RmAccount != "" {
		fmt.Printf("Removing account '%s'\n", *RmAccount)
		removeAccountImages(*RmAccount)
		removeAccount(*RmAccount)
		return
	}

	if *AddAccount != "" {
		tkn := genTkn()
		viewtkn := genTkn()
		fmt.Printf("Adding account '%s' with token %s and viewtoken %s\n", *AddAccount, tkn, viewtkn)
		addAccount(*AddAccount, tkn, viewtkn)
		return
	}

	if *ListImages {
		if *Account == "" {
			log.Fatal("No account specified!")
			return
		}

		fmt.Println("Images:")
		for _, name := range readImages(*Account) {
			fmt.Printf("\t%s\n", name)
		}

		return
	}

	if *RmImage != "" {
		if *Account == "" {
			log.Fatal("No account specified!")
			return
		}

		fmt.Println("Removing Image...")
		removeImage(*RmImage, *Account)
		fmt.Println("Done!")
		return
	}

	if *ReadImage != "" {
		if *Account == "" {
			log.Fatal("No account specified!")
			return
		}
		os.Stdout.Write(readImage(*Account, *ReadImage))
		return
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/viewtoken", viewtokenHandler)
	http.HandleFunc("/", albumHandler)
	http.HandleFunc("/images/", imageHandler)
	http.HandleFunc("/export/", exportHander)

	log.Printf("listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}
