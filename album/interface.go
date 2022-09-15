package album

import (
	"embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"path"
	"strings"
	"text/template"
)

//go:embed html
var DATA embed.FS

//go:embed LICENSE
var MIT []byte

//go:embed html/list.html
var htmllist string

type handler struct {
	getImage   func(string) []byte
	listImages func() []string
	msgs       func(string) string
	template   *template.Template
}

func List(w io.Writer, paths []string) {
	templateList.Execute(w, paths)
}

func (h *handler) ServeList(w http.ResponseWriter, r *http.Request) {
	if h.template == nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, h.msgs("error.http.500"))
		return
	}

	p := h.listImages()
	if p == nil {
		w.WriteHeader(500)
		fmt.Fprintf(w, h.msgs("error.http.500"))
		return
	}

	var paths = make([]string, len(p))
	for k, v := range p {
		paths[k] = "images/" + v
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
		f, err := DATA.Open(path)
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
			w.Write(MIT)
			return
		}
		sendFile("html/" + p[0])
		return
	} else {

	}

	if p[0] == "images" {
		if p[1] == "" {
			fmt.Fprintf(w, "<body style=\"font-family:monospace\">Files:<br>\n")

			pics := h.listImages()
			if pics == nil {
				w.WriteHeader(500)
				fmt.Fprintf(w, h.msgs("error.http.500"))
				return
			}

			for _, file := range pics {
				fmt.Fprintf(w, "<a style=\"padding-left: 3em;\" href=\"%s\">%s</a><br>\n", file)
			}

			fmt.Fprintf(w, "</body>")
		} else {
			d := h.getImage(p[1])
			if d == nil || len(d) == 0 {
				w.WriteHeader(404)
				fmt.Fprintf(w, h.msgs("error.http.404"))
				return
			}

			w.Write(d)
		}

		return
	}

	w.WriteHeader(404)
	fmt.Fprintf(w, h.msgs("error.http.404"))
}

var templateList *template.Template

func init() {
	var err error

	templateList, err = template.New("list").Parse(htmllist)
	if err != nil {
		log.Fatal("template creation failed: " + err.Error())
	}
}

func Listen(listImages func() []string,
	getImage func(string) []byte,
	msgs func(string) string,
	addr string) {

	http.ListenAndServe(addr, &handler{
		template:   templateList,
		listImages: listImages,
		getImage:   getImage,
		msgs:       msgs,
	})
}
