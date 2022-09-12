package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type ApiResponse struct {
	Err  string `json:"err"`
	Addr string `json:"addr"`
	Port int    `json:"port"`
	Pass string `json:"pass"`
	SSH  string `json:"ssh"`
}

type ApiRequest struct {
	User  string `json:"user"`
	Token string `json:"token"`
}

func encodeJson(v any) string {
	buf := &bytes.Buffer{}

	encodeJsonW(buf, v)

	return buf.String()
}

func encodeJsonW(w io.Writer, v any) {
	e := json.NewEncoder(w)
	err := e.Encode(v)
	if err != nil {
		log.Println("error occurred while encoding json: " + err.Error())
	}
}

type User struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

var userDB = make(map[string]*User)

func readUserDB() {
	f, err := os.OpenFile("users.json", os.O_CREATE|os.O_RDONLY, 0755)
	if err != nil {
		log.Fatal("error reading users.json", err)
	}

	defer f.Close()

	d := json.NewDecoder(f)

	u := make([]*User, 1)

	d.Decode(&u)

	for _, v := range u {
		userDB[v.Name] = v
	}
}

func startHttp() {
	http.HandleFunc("/forward", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		reterr := func(s string) {
			w.WriteHeader(404)
			encodeJsonW(w, &ApiResponse{Err: s})
		}

		if r.Method != "GET" {
			reterr("invalid method")
		}

		q := r.URL.Query()

		user, ok := q["user"]
		if !ok || len(user) == 0 || user[0] == "" {
			reterr("invalid key")
			return
		}

		token, ok := q["token"]
		if !ok || len(token) == 0 || token[0] == "" {
			reterr("invalid key")
			return
		}

		fmt.Printf("Get forward with user %s and token %s\n", user, token)

		u, ok := userDB[user[0]]
		if !ok {
			reterr("invalid key")
			return
		}

		if u.Token != token[0] {
			reterr("invalid key")
			return
		}

		pwd := genPass()
		port := getPort(u.Name, pwd)

		resp := &ApiResponse{
			Err:  "",
			Addr: fmt.Sprintf("%s:%d/", forward, port),
			Port: port,
			Pass: pwd,
			SSH:  ip,
		}

		encodeJsonW(w, resp)
	})

	http.ListenAndServe(":8081", nil)
}

func genPass() (s string) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 10; i++ {
		s += fmt.Sprintf("%d", rand.Intn(9))
	}

	return
}
