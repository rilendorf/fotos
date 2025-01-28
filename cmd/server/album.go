package main

import (
	"github.com/rilendorf/fotos/album"

	"fmt"
	"io"
	"net/http"
	"path"
	"path/filepath"
)

func verifyViewTkn(w http.ResponseWriter, r *http.Request) (string, string, bool) {
	params := r.URL.Query()

	Qacc := params.Get("name")
	Qvtkn := params.Get("view")
	token, viewtoken := getTokenByName(Qacc)

	if token == "" || Qvtkn != viewtoken {
		w.WriteHeader(404)
		fmt.Fprint(w, "Sorry, this location dosn't exist\n")
		return "", "", false
	}

	return Qacc, viewtoken, true
}

func contentType(w http.ResponseWriter, r *http.Request) {
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
}

func albumHandler(w http.ResponseWriter, r *http.Request) {
	contentType(w, r)

	_, file := filepath.Split(r.URL.Path)
	if file == "" {
		acc, viewTkn, ok := verifyViewTkn(w, r)
		if !ok {
			return
		}

		images := readImages(acc)
		paths := make([]string, len(images))
		l := len(paths)

		for k, v := range images {
			paths[l-1-k] = conf.PublicAccess + "/images/" + v + "?name=" + acc + "&view=" + viewTkn
		}

		album.List(w, album.ListParamsFromStrs(paths, "/export?name="+acc+"&view="+viewTkn))
	} else if file == "LICENSE.txt" {
		w.WriteHeader(200)
		w.Write(album.MIT)
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
	acc, _, ok := verifyViewTkn(w, r)
	if !ok {
		return
	}

	_, file := filepath.Split(r.URL.Path)

	f := readImage(acc, file)

	if len(f) == 0 {
		w.WriteHeader(404)
		fmt.Fprint(w, "Sorry, this location dosn't exist\n")
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")

	if r.URL.Query().Has("print") {
		w.Header().Set("Refresh: 5", "url=/?name="+acc+"&view="+r.URL.Query().Get("view"))
		fmt.Printf("printing image %s\n", file)
		err := MarkImagePrint(file, acc)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "Failed to print image!")
			fmt.Printf("Failed to mark image print: %s", err)
			return
		}

		w.Write(PreviewImage(f))
	} else if r.URL.Query().Has("preview") {
		w.Write(PreviewImage(f))
	} else {
		w.Write(f)
	}
}

func exportHander(w http.ResponseWriter, r *http.Request) {
	acc, _, ok := verifyViewTkn(w, r)
	if !ok {
		return
	}

	h := w.Header()

	h.Set("Content-Disposition", "attachment; filename=\""+acc+"-export.zip\"")
	h.Set("Content-Type", "application/zip")

	createZip(acc, w)
}
