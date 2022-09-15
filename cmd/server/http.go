package main

import (
	"bytes"
	"fmt"
	"io"
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
