package ui_webview

import (
	"github.com/DerZombiiie/fotos/fotos"
	"github.com/webview/webview"

	"bytes"
	_ "embed"
	"encoding/base64"
	"fmt"
	"log"
	"strings"
	"sync"
	"text/template"
	"time"
)

type Webview struct {
	webview.WebView

	htmlChan chan string

	templateMu sync.RWMutex
	template   *template.Template
	ImageUrl   string
	Status     string
	Popup      string
	Msg        string
}

func (w *Webview) ShowImage(img *fotos.Image) {
	if img == nil {
		img = fotos.ImageFromBytes(testimg)
	}

	println("setimg")

	buf := &bytes.Buffer{}

	e := base64.NewEncoder(base64.StdEncoding, buf)
	e.Write(img.Bytes())
	e.Close()

	w.templateMu.Lock()
	defer w.templateMu.Unlock()
	w.ImageUrl = "data:image/jpeg;base64," + buf.String()

	w.rebuild()

}

func (w *Webview) Countdown(i int) {
	w.templateMu.Lock()
	defer w.templateMu.Unlock()

	if i < 0 {
		return
	}

	for i > 0 {
		w.Popup = fmt.Sprintf("%d", i)
		w.rebuild()
		time.Sleep(time.Second)

		i--
	}

	w.Popup = "cheese"
	w.rebuild()

	go func() {
		w.templateMu.Lock()
		defer w.templateMu.Unlock()

		time.Sleep(time.Second)

		w.Popup = ""
		w.rebuild()

	}()

	time.Sleep(time.Second / 2)
}

func (w *Webview) ShowMsg(str string) {
	w.templateMu.Lock()
	defer w.templateMu.Unlock()

	w.Msg = strings.ReplaceAll(str, "\n", "<br>")
	w.rebuild()
}

func (w *Webview) SetStatus(str string) {
	w.templateMu.Lock()
	defer w.templateMu.Unlock()

	w.Status = str
	w.rebuild()
}

func (w *Webview) rebuild() {
	buf := &bytes.Buffer{}

	err := w.template.Execute(buf, w)
	if err != nil {
		log.Fatal("[ui_webview.rebuild] template execution failed " + err.Error())
	}

	w.SetHtml(buf.String())
}

//go:embed image.jpg
var testimg []byte

//go:embed interface.html
var html string

func NewWebview(debug bool) *Webview {
	var err error

	w := &Webview{WebView: webview.New(debug), htmlChan: make(chan string)}

	w.template, err = template.New("site").Parse(html)
	if err != nil {
		log.Fatal("[ui_webview.NewWebview] template creation failed " + err.Error())
	}

	return w
}

func init() {
	w := NewWebview(false)

	w.SetTitle("Fotos")
	w.SetSize(480, 320, webview.HintMax)

	w.ShowImage(fotos.ImageFromBytes(testimg))

	fotos.RegisterUI("webview", w)
}
