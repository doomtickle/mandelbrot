package canvas

import (
	"encoding/json"
	"errors"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

type Canvas struct {
	Width, Height int
	Img           *image.RGBA
	Palette       color.Palette
	Bg            color.Color
}

type Json struct {
	Palettes map[string][][]int `json:"palettes"`
}

func New(res int, plt string, bg color.Color) *Canvas {
	var canvas Canvas


	canvas.Width = res
	canvas.Height = res
	canvas.Img = image.NewRGBA(image.Rect(0, 0, res, res))
	canvas.Palette = getPalette(plt)
  canvas.Bg = bg

	return &canvas
}

func (c *Canvas) Save(s string) {
	f, _ := os.Create(s)
	png.Encode(f, c.Img)
}

func getPalette(s string) []color.Color {
	var (
		j Json
		p []color.Color
	)
	f, err := os.ReadFile("mandelbrot.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &j)
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range j.Palettes[s] {
		col := color.RGBA{uint8(c[0]), uint8(c[1]), uint8(c[2]), uint8(c[3])}
		p = append(p, col)
	}

	return p
}

var errInvalidFormat = errors.New("invalid format")

//https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func ParseHexColorFast(s string) (c color.RGBA, err error) {
    c.A = 0xff

    if s[0] != '#' {
        return c, errInvalidFormat
    }

    hexToByte := func(b byte) byte {
        switch {
        case b >= '0' && b <= '9':
            return b - '0'
        case b >= 'a' && b <= 'f':
            return b - 'a' + 10
        case b >= 'A' && b <= 'F':
            return b - 'A' + 10
        }
        err = errInvalidFormat
        return 0
    }

    switch len(s) {
    case 7:
        c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
        c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
        c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
    case 4:
        c.R = hexToByte(s[1]) * 17
        c.G = hexToByte(s[2]) * 17
        c.B = hexToByte(s[3]) * 17
    default:
        err = errInvalidFormat
    }
    return
}
