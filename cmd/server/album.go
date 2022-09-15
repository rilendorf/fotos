package main

import (
	"github.com/DerZombiiie/fotos/album"

	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func albumHandler(w http.ResponseWriter, r *http.Request) {
	switch path.Ext(r.URL.Path) {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".js":
		w.Header().Set("Content-Type", "application/javascript")
	case ".jpg":
		w.Header().Set("Content-Type", "image/jpeg")
	case ".woff":
		w.Header().Set("Content-Type", "font/woff")
	case ".txt":
		w.Header().Set("Content-Type", "text/plain")
	default:
		w.Header().Set("Content-Type", "text/html")
	}

	_, file := filepath.Split(r.URL.Path)
	if file == "" {
		params := r.URL.Query()

		Qacc := params.Get("name")
		Qvtkn := params.Get("view")
		token, viewtoken := getTokenByName(Qacc)

		if token == "" || Qvtkn != viewtoken {
			w.WriteHeader(404)
			fmt.Fprint(w, "Sorry, this location dosn't exist\n")
			return
		}

		images := readImages(Qacc)
		paths := make([]string, len(images))

		for k, v := range images {
			paths[k] = "/images/" + v + "?name=" + Qacc + "&view=" + viewtoken
		}

		album.List(w, paths)
	} else {
		f, err := album.DATA.Open("html/" + file)
		if err != nil {
			w.WriteHeader(404)
			fmt.Fprint(w, "Sorry, this location dosn't exist\n")
			return
		}

		io.Copy(w, f)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()

	Qacc := params.Get("name")
	Qvtkn := params.Get("name")
	token, viewtoken := getTokenByName(Qacc)

	if token == "" && Qvtkn == viewtoken {
		w.WriteHeader(404)
		fmt.Fprint(w, "Sorry, this location dosn't exist\n")
		return
	}

	_, file := filepath.Split(r.URL.Path)

	f := readImage(Qacc, file)

	if len(f) == 0 {
		w.WriteHeader(404)
		fmt.Fprint(w, "Sorry, this location dosn't exist\n")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(f)
}
