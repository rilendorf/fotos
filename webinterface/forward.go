package web

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"

	"github.com/DerZombiiie/easyssh"
	"github.com/DerZombiiie/fotos/fotos"
)

func Addr() string {
	cloud, ok := fotos.Conf()["cloud.addr"]
	if !ok {
		log.Println("no cloud address specified!, not using cloud!")
	}

	return cloud
}

type ApiRequest struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

func getForward() *ApiResponse {
	c := fotos.Conf()

	user, ok := c["cloud.user"]
	if !ok {
		log.Println("no cloud user specified!, not using cloud!")
	}

	token, ok := c["cloud.token"]
	if !ok {
		log.Println("no cloud token specified!, not using cloud!")
	}

	i := 0

	var resp *http.Response
	var err error

	params := url.Values{}
	params.Add("token", token)
	params.Add("user", user)

	for i < 10 {
		// ask server for port
		resp, err = http.Get(fmt.Sprintf("%s/forward?user=%s&token=%s", Addr(), user, token))
		if err != nil {
			log.Println("err getting reverse proxy token, retrying", i, err)
			time.Sleep(time.Second * 5)
		} else {
			break
		}

		i++
	}

	if err != nil {
		log.Fatal("exceeded 10 tries, quitting, " + err.Error())
	}

	defer resp.Body.Close()

	d := json.NewDecoder(resp.Body)

	r := &ApiResponse{}
	d.Decode(r)

	return r
}

var cert []byte

func Cert() []byte {
	if len(cert) == 0 {
		p, ok := fotos.Conf()["cloud.pubcert"]
		if !ok {
			log.Fatal("No cloud public certificate!")
		}

		if !filepath.IsAbs(p) {
			p = fotos.Path(p)
		}

		f, err := os.Open(p)
		if err != nil {
			log.Fatal("Can't read public certificate!")
		}
		defer f.Close()

		r := bufio.NewReader(f)
		r.ReadBytes(' ')

		encodedCert, _ := r.ReadBytes(' ')
		certBuf := bytes.NewBuffer(encodedCert)

		buf := &bytes.Buffer{}
		e := base64.NewDecoder(base64.StdEncoding, certBuf)

		io.Copy(buf, e)

		cert = buf.Bytes()
	}

	return cert
}

func User() string {
	user, ok := fotos.Conf()["cloud.user"]
	if !ok {
		log.Println("no cloud user specified!, not using cloud!")
	}

	return user
}

func accessUrl() string {
	a, ok := fotos.Conf()["http.album.addr"]
	if !ok {
		log.Fatal("No Album addr specified")
	}

	_, portString, err := net.SplitHostPort(a)
	if err != nil {
		return ""
	}

	return "127.0.0.1:" + portString
}

type ApiResponse struct {
	Port int    `json:"port"`
	Addr string `json:"addr"`
	Pass string `json:"pass"`
	SSH  string `json:"ssh"`
	Err  string `json:"err"`
}

func init() {
	fotos.Runner(func() {
		resp := getForward()

		if resp.Err != "" {
			log.Println("[CLOUD] API error " + resp.Err)
			return
		}

		log.Println("[CLOUD] got addr " + resp.Addr)

		config := &ssh.ClientConfig{
			User: User(),
			Auth: []ssh.AuthMethod{
				ssh.Password(resp.Pass),
			},

			HostKeyCallback: func(_ string, _ net.Addr, key ssh.PublicKey) error {
				if bytes.Compare(key.Marshal(), Cert()) != 0 {
					log.Fatal("ssh key dosn't match")
					return fmt.Errorf("Dosn't match")
				}

				return nil
			},
		}

		conn, err := easyssh.Dial("tcp", resp.SSH, config)
		if err != nil {
			log.Fatalf("unable to connect: %s", err)
		}
		defer conn.Close()

		err = conn.RemoteForward(fmt.Sprintf("localhost:%d", resp.Port), accessUrl())
		if err != nil {
			log.Fatalf("unable to forward local port: %s", err)
		}
	})

}
