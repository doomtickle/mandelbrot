package main

import (
	"flag"
	"log"
	"math"

	"github.com/doomtickle/mandelbrot/canvas"
)

// let's set some flags so we can modify stuff from the cli
var (
	bg         = flag.String("bg", "#fff", "hex value for background color")
	res        = flag.Int("res", 4096, "The width and height of the canvas.")
	min        = flag.Float64("min", -1.2, "The minimum value on the x and y axis.")
	max        = flag.Float64("max", 1.2, "The maximum value on the x and y axis.")
	out        = flag.String("out", "image.png", "output path for generated file.")
	cReal      = flag.Float64("real", 0, "c's real component.")
	palette    = flag.String("palette", "blue", "color palette from your mandelbrot.json config")
	cImaginary = flag.Float64("i", 0, "c's imaginary component.")
	iterations = flag.Int("iterations", 100, "how many operations until considering a point bounded.")
)

func main() {
	flag.Parse()
  bgColor, err := canvas.ParseHexColorFast(*bg)
  if err != nil {
    log.Fatal(err)
  }
	// Generate a new canvas
	c := canvas.New(*res, *palette, bgColor)
	// Instantiate our range maps to normalize x,y values to be between the min and max values specified at runtime.
	rm := newRangeMap(rangeBounds{0, float64(c.Width)}, rangeBounds{*min, *max})

	// for every pixel
	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {
			// get the normalized x
			a, ok := rm(float64(x))
			if !ok {
				log.Fatal("Rangemap Error")
			}
			// get the normalized y
			b, ok := rm(float64(y))
			if !ok {
				log.Fatal("Rangemap Error")
			}
			// iteration counter
			n := 0
			for n < *iterations {
				// Math stuff that I had to watch videos about and look up....
				// This is where we apply the zeta function to each pixel recursively.
				// The function will either converge to a value or blow up to infinity.
				aSquared := a*a - b*b
				twoAB := 2 * a * b

				a = aSquared + *cReal
				b = twoAB + *cImaginary

				// Looks like we're heading to infinity
				if math.Abs(aSquared+twoAB) > 16 {
					break
				}
				n++
			}
			c.Img.Set(x, y, c.Palette[n%len(c.Palette)])

			// stays bounded
			if n == *iterations {
				c.Img.Set(x, y, c.Bg)
			}

			// this calms down some of the color schemes.
			// Can be removed or tweaked based on your preference.
			if n <= 16 {
				c.Img.Set(x, y, c.Bg)
			}
		}
	}

	c.Save(*out)
}

type rangeBounds struct {
	Min, Max float64
}

// newRangeMap maps one range onto another (thank you, smart people)
func newRangeMap(xr, yr rangeBounds) func(float64) (float64, bool) {
	// normalize direction of ranges so that out-of-range test works
	if xr.Min > xr.Max {
		xr.Min, xr.Max = xr.Max, xr.Min
		yr.Min, yr.Max = yr.Max, yr.Min
	}
	// compute slope, intercept
	m := (yr.Max - yr.Min) / (xr.Max - xr.Min)
	b := yr.Min - m*xr.Min
	// return function literal
	return func(x float64) (y float64, ok bool) {
		if x < xr.Min || x > xr.Max {
			return 0, false // out of range
		}
		return m*x + b, true
	}
}
