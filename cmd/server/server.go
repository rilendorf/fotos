package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/user"
	"path/filepath"

	"flag"
)

var Account = flag.String("account", "", "specify account")
var ListImages = flag.Bool("list-images", false, "if true, programm lists name of all images and quits")
var RmImage = flag.String("rm-image", "", "delete a image (blob = 0 & mark delete)")

var Hard = flag.Bool("hard", false, "use with rm-image to remove entry")

var ReadImage = flag.String("read-image", "", "read image to stdout")

var AddImage = flag.String("add-image", "", "specify a file path for a image to add (use with -account)")

var AddAccount = flag.String("add-account", "", "add a new account")
var RmAccount = flag.String("rm-account", "", "remove a account")
var GetToken = flag.String("get-token", "", "get a accounts token")

var ConfigPath = flag.String("config", "/opt/fotos/fotos.cfg", "set configuration path")

var conf Config

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	currentUser, err := user.Current()
	if err != nil {
		log.Printf("WARN: failed to determin user: %s\n", currentUser)
	} else {
		if currentUser.Username != "fotos" {
			log.Printf("WARN: the server is best executed as 'fotos' user")
		}
	}

	flag.Parse()

	conf = ReadConfig()
	log.Printf("Configuration: %#v", conf)

	openDB(conf.DBPath)

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
		fmt.Printf("View: %s/?name=%s&view=%s\n", conf.PublicAccess, *AddAccount, viewtkn)

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

		if *Hard {
			fmt.Println("Removing Image entry...")
			err = removeImage(*RmImage, *Account)
		} else {
			fmt.Println("Deleting Image...")
			err = deleteImage(*RmImage, *Account)
		}

		if err != nil {
			log.Fatalf("Failed: %s", err)
		}
		fmt.Printf("Done")
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

	if *AddImage != "" {
		if *Account == "" {
			log.Fatal("No account specified!")
		}

		// read image
		b, err := os.ReadFile(*AddImage)
		if err != nil {
			log.Fatalf("Failed to read image file from '%s': %s", *AddImage, err)
		}

		_, name := filepath.Split(*AddImage)

		err = addImage(name, *Account, b)
		if err != nil {
			log.Fatalf("Failed: %s", err)
		}

		log.Printf("Ok")
		return
	}

	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/viewtoken", viewtokenHandler)
	http.HandleFunc("/", albumHandler)
	http.HandleFunc("/images/", imageHandler)
	http.HandleFunc("/export/", exportHander)
	http.HandleFunc("/printingrequests", handlePrintingRequests)
	http.HandleFunc("/printed", handlePrinted)

	log.Printf("Listening on %s\n", conf.Listen)

	err = http.ListenAndServe(conf.Listen, nil)
	if err != nil {
		log.Fatalf("ListenAndServe: %s", err)
	}
}
