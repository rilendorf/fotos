package web

import (
	"github.com/DerZombiiie/fotos/fotos"

	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed html
var data embed.FS

//go:embed html/list.html
var htmllist string

type handler struct {
	template *template.Template
}

func listPictures() (s []os.FileInfo) {
	filepath.Walk(fotos.SaveDir(), func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			s = append(s, info)
		}

		return nil
	})

	return
}

func (h *handler) ServeList(w http.ResponseWriter, r *http.Request) {
	if h.template == nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, fotos.Msgs("error.http.500"))
		return
	}

	p := listPictures()
	if p == nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, fotos.Msgs("error.http.500"))
		return
	}

	var paths = make([]string, len(p))
	for k, v := range p {
		paths[k] = "images/" + v.Name()
	}

	h.template.Execute(w, paths)
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("[HTTP] access url " + r.URL.Path)

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

	sendFile := func(path string) error {
		f, err := data.Open(path)
		if err != nil {
			w.WriteHeader(404)
			return err
		}

		w.WriteHeader(200)
		io.Copy(w, f)
		return nil
	}

	if r.URL.Path == "/" || r.URL.Path == "" {
		h.ServeList(w, r)

		return
	}

	p := strings.Split(r.URL.Path, "/")
	if p[0] == "" {
		p = p[1:]
	}

	if len(p) == 1 {
		if p[0] == "LICENSE.txt" {
			w.WriteHeader(200)
			w.Write([]byte(fotos.License()))
			return
		}
		sendFile("html/" + p[0])
		return
	} else {

	}

	if p[0] == "images" {
		if p[1] == "" {
			fmt.Fprintf(w, "<body style=\"font-family:monospace\">Files:<br>\n")

			pics := listPictures()
			if pics == nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, fotos.Msgs("error.http.500"))
				return
			}

			for _, info := range pics {
				fmt.Fprintf(w, "<a style=\"padding-left: 3em;\" href=\"%s\">%s</a><br>\n", info.Name(), info.Name())
			}

			fmt.Fprintf(w, "</body>")
		} else {
			f, err := os.OpenFile(fotos.SaveDir()+p[1], os.O_RDONLY, 0)
			if err != nil {
				w.WriteHeader(404)
				fmt.Fprintf(w, fotos.Msgs("error.http.404"))
			}

			io.Copy(w, f)
		}

		return
	}

	w.WriteHeader(404)
	fmt.Fprintf(w, fotos.Msgs("error.http.404"))
}

func init() {
	fotos.Runner(func() {
		c := fotos.Conf()
		addr, ok := c["http.album.addr"]
		if !ok {
			addr = "localhost:8080"
		}

		t, err := template.New("list").Parse(htmllist)
		if err != nil {
			log.Fatal("[http.Runner] template creation failed: " + err.Error())
		}

		http.ListenAndServe(addr, &handler{
			template: t,
		})
	})
}
