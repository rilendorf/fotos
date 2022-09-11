package fotos

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

var execDir string
var execDirOnce sync.Once

func Path(path ...string) string {
	execDirOnce.Do(func() {
		e, err := os.Executable()
		if err != nil {
			log.Fatal("[fotos.Path] " + err.Error())
		}

		execDir = filepath.Dir(e)
	})

	return execDir + "/" + strings.Join(path, "")
}
