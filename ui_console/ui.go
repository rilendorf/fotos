package ui

import (
	"github.com/rilendorf/fotos/fotos"

	"fmt"
	"os"
	"strings"
	"time"
)

type Console struct {
	*os.File
}

func (c *Console) Countdown(i time.Duration) {
	fmt.Fprintf(c, "Countdown: %s \n", i)
	time.Sleep(i)
	fmt.Fprintf(c, "Countdown: Done \n")
}

func (c *Console) ShowImage(img *fotos.Image) {
	b := img.Image().Bounds()

	fmt.Fprintf(c, "Image!\n\tBounds: %d %d\n", b.Min, b.Max)
}

func (c *Console) ShowMsg(msg string) {
	fmt.Fprintf(c, "Message:\n%s\n", strings.ReplaceAll(msg, "\n", "\n\t"))
}

func (c *Console) SetStatus(str string) {
	fmt.Fprintf(c, "Status update: %s\n", str)
}

func (c *Console) Run() {
	for {
		fmt.Println("Press enter to take picture!")
		fmt.Scanf("\n")

		fotos.TakePicture()
	}
}

func init() {
	fotos.RegisterUI("console", &Console{os.Stdout})
}
