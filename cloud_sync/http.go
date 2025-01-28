package cloud

import (
	"github.com/rilendorf/fotos/fotos"

	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"log"
	"strings"
)

func baseUrl() string {
	s, ok := fotos.Conf()["cloud.addr"]
	if !ok {
		log.Println("[sync] no cloud addr configured!")
	}

	return s
}

func getToken() string {
	s, ok := fotos.Conf()["cloud.token"]
	if !ok {
		log.Println("[sync] no cloud token configured!")
	}

	return s
}

func FilterNeeded(in []string) (s []string) {
	// download filelist:
	path, err := url.JoinPath(baseUrl(), "/list")
	path += "?token=" + getToken()
	if err != nil {
		log.Printf("[sync] Filelist download failed: %s\n", err)
	}

	res, err := http.Get(path)
	if err != nil {
		log.Printf("[sync] Filelist download failed: %s\n", err)
		return
	}

	defer res.Body.Close()
	rawList, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[sync] Filelist download failed: %s\n", err)
	}

	list := strings.Split(string(rawList), "\n")
	m := make(map[string]struct{})

	for _, v := range list {
		m[v] = struct{}{}
	}

	for _, v := range in {
		if _, ok := m[v]; !ok {
			s = append(s, v)
		}
	}

	return
}

func SetPrinted(name string) {
	path, err := url.JoinPath(baseUrl(), "/printed")
	path += "?token=" + getToken() + "&name=" + name
	if err != nil {
		log.Printf("[printing] set printed failed: %s\n", err)

		return
	}

	res, err := http.Get(path)
	if err != nil {
		log.Printf("[printing] set printed failed: %s\n", err)

		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("[printing] markprinted failed connerr: %s", err)
	}

	log.Printf("[printing] markprinted, server says: %s", body)

	return
}

func GetPrintingRequests() []string {
	path, err := url.JoinPath(baseUrl(), "/printingrequests")
	path += "?token=" + getToken()
	if err != nil {
		log.Printf("[printing] get printing requests failed: %s\n", err)
	}

	res, err := http.Get(path)
	if err != nil {
		log.Printf("[printing] get printing requests failed : %s\n", err)

		return nil
	}
	defer res.Body.Close()

	d := json.NewDecoder(res.Body)
	var s = make([]string, 0)
	err = d.Decode(&s)
	if err != nil {
		log.Printf("[printing] Failed to decode printing requests, %s", err)
		return nil
	}

	log.Printf("[printing] requests: '%s'", s)
	return s
}

func Upload(name string, file io.Reader) {
	path, err := url.JoinPath(baseUrl(), "/upload")
	path += "?name=" + name + "&token=" + getToken()
	if err != nil {
		log.Printf("[sync] Upload failed: %s\n", err)
	}

	res, err := http.Post(path, "binary/octet-stream", file)
	if err != nil {
		log.Printf("[sync] Upload failed: %s\n", err)
		return
	}
	defer res.Body.Close()

	message, _ := io.ReadAll(res.Body)
	log.Printf("[sync] Remote msg: '%s'", string(message))
	fotos.ShowMsg(strings.ReplaceAll(string(message), "\n", ""))
}
