package fotos

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

type Config map[string]string

var conf Config
var confMu sync.RWMutex

func readConfig() error {
	confMu.Lock()
	defer confMu.Unlock()

	f, err := os.OpenFile(Path("config.json"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}

	defer f.Close()

	s, err := f.Stat()
	if err != nil {
		return err
	}

	if s.Size() == 0 {
		f.Write([]byte("{\n\t\n}\n"))
		f.Seek(0, 0)
	}

	d := json.NewDecoder(f)
	err = d.Decode(&conf)
	if err != nil {
		return err
	}

	return nil
}

func Conf() Config {
	confMu.RLock()
	defer confMu.RLock()

	c := make(Config)
	for k, v := range conf {
		c[k] = v
	}

	return c
}

func Device() string {
	d, ok := Conf()["device"]
	if !ok || d == "" {
		log.Fatal("No device configured")
	}

	return d
}

func CountdownTime() time.Duration {
	raw, ok := Conf()["countdown"]
	if !ok {
		log.Println("Error, no countdown in config! using 5s")
		raw = "5s"
	}

	d, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("Error parsing duration '%s': %s\n", raw, err)
		d = time.Second * 5
	}

	return d
}
