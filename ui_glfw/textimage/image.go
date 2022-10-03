package textimage

import (
	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"io/ioutil"
	"sync"
)

var f *truetype.Font
var fMu sync.RWMutex

// most systems have this one:
const fontPath = "/usr/share/fonts/liberation/LiberationSerif-Regular.ttf"

// UpdateFont sets the font used
func UpdateFont(path string) error {
	fMu.Lock()
	defer fMu.Unlock()

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	f, err = truetype.Parse(b)
	if err != nil {
		return err
	}

	return nil
}

// Generate Image creats an image thats centerd (no newline support)
func GenerateImage(width, height int, text string) *image.RGBA {
	if f == nil {
		UpdateFont(fontPath)
	}

	fMu.RLock()
	defer fMu.RUnlock()

	rgba := image.NewRGBA(image.Rect(0, 0, width, height))
	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFont(f)
	c.SetFontSize(125)
	c.SetClip(rgba.Bounds())
	c.SetDst(rgba)
	c.SetSrc(image.Black)
	c.SetHinting(font.HintingNone)

	// Truetype stuff
	opts := truetype.Options{
		Size: 125.0,
	}
	face := truetype.NewFace(f, &opts)

	advance, _ := StrAdvance(face, text)

	pt := freetype.Pt((width/2)-advance.Round()/2, height/2+int(125)/4)
	c.DrawString(text, pt)

	return rgba
}

// like GlyphAdvance but for entire strings
// if ok != true one glyph isn't in font
func StrAdvance(face font.Face, str string) (i fixed.Int26_6, ok bool) {
	for _, x := range str {
		awidth, k := face.GlyphAdvance(x)
		if !k {
			ok = k
		}

		i += awidth - fixed.I(9)
	}

	return i, true
}
