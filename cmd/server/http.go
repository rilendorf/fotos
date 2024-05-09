package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	name := query.Get("name")
	t := query.Get("token")
	account := getNameByToken(t)

	if account == "" {
		fmt.Fprintf(w, "Token invalid \n")
		fmt.Printf("Token invalid %s\n", t)
		return
	}

	if name == "" {
		fmt.Fprintf(w, "Invalid file name!")
		return
	}

	fmt.Printf("%s just uploaded %s\n", account, name)

	buf := &bytes.Buffer{}
	n, err := io.Copy(buf, r.Body)
	if err != nil {
		fmt.Fprintf(w, "Can't copy file to Buffer! %s\n", err)
	}

	addImage(name, account, buf.Bytes())

	fmt.Fprintf(w, "%d bytes are recieved.\n", n)
}
func handlePrinted(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	t := query.Get("token")
	account := getNameByToken(t)

	if account == "" {
		fmt.Fprintf(w, "Token invalid \n")
		fmt.Printf("Token invalid %s\n", t)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	err := MarkImagePrinted(query.Get("name"), account)
	if err != nil {
		fmt.Printf("Faield to makr printed %s for %s : %s", query.Get("name"), account, err)
		fmt.Fprintf(w, "faield to makr printe!d!!")

		return
	}

	log.Printf("Marked printed %s %s", query.Get("name"), account)
	fmt.Fprintf(w, "marked image as printed")
}

func handlePrintingRequests(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	t := query.Get("token")
	account := getNameByToken(t)

	if account == "" {
		fmt.Fprintf(w, "Token invalid \n")
		fmt.Printf("Token invalid %s\n", t)
		return
	}

	w.Header().Set("Content-Type", "application/jsonjson")

	e := json.NewEncoder(w)
	img, err := GetImages2Print(account)
	if err != nil {
		fmt.Fprintf(w, "failed to get imgages to print")
		fmt.Printf("failed to get imgages to print for %s:%s\n", account, err)

		return
	}

	e.Encode(img)
}

func listHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	t := query.Get("token")
	account := getNameByToken(t)

	if account == "" {
		fmt.Fprintf(w, "Token invalid \n")
		fmt.Printf("Token invalid %s\n", t)
		return
	}

	w.Header().Set("Content-Type", "plain/text")

	s := readImages(account)
	for _, v := range s {
		fmt.Fprintf(w, "%s\n", v)
	}
}

func viewtokenHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	t := query.Get("token")
	account := getNameByToken(t)

	if account == "" {
		fmt.Fprintf(w, "Token invalid \n")
		fmt.Printf("Token invalid %s\n", t)
		return
	}

	_, viewtkn := getTokenByName(account)

	fmt.Fprint(w, viewtkn)
}
