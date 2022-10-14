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

	clearTimeMu sync.RWMutex
	clearTime   time.Time
}

func getImageTimeout() time.Duration {
	ds, ok := fotos.Conf()["ui.resetimg"]
	if !ok {
		ds = "8s"
	}

	d, err := time.ParseDuration(ds)
	if err != nil {
		log.Printf("[ui.resetimg] Error parsing specified duration '%s': %s\n", ds, err)
		d = time.Second * 8
	}

	return d
}

func (w *Webview) showImage(img *fotos.Image, clear bool) {
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

	if clear {
		go w.clearImage()
	}
}

func (w *Webview) ShowImage(img *fotos.Image) {
	w.showImage(img, true)
}

func (w *Webview) clearImage() {
	time.Sleep(getImageTimeout())

	w.showImage(fotos.ImageFromBytes(testimg), false)
}

func (w *Webview) Countdown(d time.Duration) {
	w.templateMu.Lock()
	defer w.templateMu.Unlock()

	if d < 0 {
		return
	}

	for d.Seconds() > 0 {
		w.Popup = fmt.Sprintf("%d", int(d.Seconds()))
		w.rebuild()
		time.Sleep(time.Second)

		d -= time.Second
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
	w.SetSize(1920, 1080, webview.HintFixed)

	w.showImage(fotos.ImageFromBytes(testimg), false)

	fotos.RegisterUI("webview", w)
}
