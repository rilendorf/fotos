package cloud

import (
	"github.com/DerZombiiie/fotos/fotos"

	"io"
	"net/http"
	"net/url"

	"fmt"
	"log"
	"strings"
)

func baseUrl() string {
	s, ok := fotos.Conf()["cloud.addr"]
	if !ok {
		log.Println("no cloud addr configured!")
	}

	return s
}

func getToken() string {
	s, ok := fotos.Conf()["cloud.token"]
	if !ok {
		log.Println("no cloud token configured!")
	}

	return s
}

func FilterNeeded(in []string) (s []string) {
	// download filelist:
	path, err := url.JoinPath(baseUrl(), "/list")
	path += "?token=" + getToken()
	if err != nil {
		log.Printf("Filelist download failed: %s\n", err)
	}

	res, err := http.Get(path)
	if err != nil {
		log.Printf("Filelist download failed: %s\n", err)
		return
	}

	defer res.Body.Close()
	rawList, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Filelist download failed: %s\n", err)
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

func Upload(name string, file io.Reader) {
	path, err := url.JoinPath(baseUrl(), "/upload")
	path += "?name=" + name + "&token=" + getToken()
	if err != nil {
		log.Printf("Upload failed: %s\n", err)
	}

	res, err := http.Post(path, "binary/octet-stream", file)
	if err != nil {
		log.Printf("Upload failed: %s\n", err)
		return
	}
	defer res.Body.Close()

	message, _ := io.ReadAll(res.Body)
	fmt.Printf("Remote msg: " + string(message))
}
