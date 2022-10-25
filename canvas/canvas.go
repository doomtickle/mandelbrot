package canvas

import (
	"image"
	"image/png"
	"os"
)

type Canvas struct {
	Width, Height int
	Img           *image.RGBA
}

func New(w, h int) *Canvas {
	var canvas Canvas

	canvas.Width = w
	canvas.Height = h
	canvas.Img = image.NewRGBA(image.Rect(0, 0, w, h))

	return &canvas
}

func (c *Canvas) Save() {
	f, _ := os.Create("image.png")
	png.Encode(f, c.Img)
}
