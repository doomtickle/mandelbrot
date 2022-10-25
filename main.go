package main

import (
	"flag"
	"image"
	"image/color"
	"log"
	"math"

	"github.com/doomtickle/mandelbrot/canvas"
)

// Different palettes I tried out.
//
// blue_plt is the best so far.
//
// You can add as many palettes as you want and each palette can contain any amount of colors.
// It's totally up to your preferences. Go wild.
var (
	psychedelic_plt = []color.Color{
		color.RGBA{245, 76, 0, 255},
		color.RGBA{255, 150, 0, 255},
		color.RGBA{244, 78, 163, 255},
		color.RGBA{199, 26, 53, 255},
		color.RGBA{73, 65, 248, 255},
		color.RGBA{0, 154, 228, 255},
	}
	soft_plt = []color.Color{
		color.RGBA{34, 22, 43, 255},
		color.RGBA{69, 31, 85, 255},
		color.RGBA{114, 78, 145, 255},
		color.RGBA{229, 79, 109, 255},
		color.RGBA{73, 65, 248, 255},
		color.RGBA{248, 198, 48, 255},
	}
	blue_plt = []color.Color{
		color.RGBA{202, 229, 255, 255},
		color.RGBA{172, 237, 255, 255},
		color.RGBA{137, 187, 254, 255},
		color.RGBA{111, 138, 183, 255},
		color.RGBA{97, 93, 108, 255},
	}
	spring_plt = []color.Color{
		color.RGBA{245, 255, 198, 255},
		color.RGBA{255, 235, 250, 255},
		color.RGBA{171, 135, 255, 255},
		color.RGBA{255, 172, 228, 255},
		color.RGBA{193, 255, 155, 255},
	}

	pink_plt = []color.Color{
		color.RGBA{236, 145, 216, 255},
		color.RGBA{255, 170, 234, 255},
		color.RGBA{255, 190, 239, 255},
		color.RGBA{255, 211, 218, 255},
		color.RGBA{255, 255, 255, 255},
	}
)

func main() {
	// let's set some flags so we can modify stuff from the cli
	res := flag.Int("res", 4096, "The width and height of the canvas.")
	xMin := flag.Float64("xmin", -1.2, "The minimum value on the x axis.")
	yMin := flag.Float64("ymin", -1.2, "The minimum value on the y axis.")
	xMax := flag.Float64("xmax", 1.2, "The maximum value on the x axis.")
	yMax := flag.Float64("ymax", 1.2, "The maximum value on the y axis.")
	cReal := flag.Float64("real", 0, "c's real component.")
	cImaginary := flag.Float64("i", 0, "c's imaginary component.")
	iterations := flag.Int("iterations", 100, "how many operations until considering a point bounded.")

	flag.Parse()

	// Generate a new canvas
	c := canvas.New(*res, *res)
	// Instantiate our range maps to normalize x,y values to be between the min and max values specified at runtime.
	rmX := newRangeMap(rangeBounds{0, float64(c.Width)}, rangeBounds{*xMin, *xMax})
	rmY := newRangeMap(rangeBounds{0, float64(c.Height)}, rangeBounds{*yMin, *yMax})

	// for every pixel
	for x := 0; x < c.Width; x++ {
		for y := 0; y < c.Height; y++ {

			// get the normalized x
			a, ok := rmX(float64(x))
			if !ok {
				log.Fatal("Rangemap Error")
			}

			// get the normalized y
			b, ok := rmY(float64(y))
			if !ok {
				log.Fatal("Rangemap Error")
			}

			// iteration counter
			n := 0

			for n < *iterations {
				// Math stuff that I had to watch videos about and look up....
				// This is where we apply the zeta function to each pixel recursively.
				// The function will either converse to a value or blow up to infinity.
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

			c.Img.Set(x, y, blue_plt[n%len(blue_plt)])

			// stayed bounded
			if n == *iterations {
				c.Img.Set(x, y, image.White)
			}

			// this calms down some of the color schemes.
			// Can be removed or tweaked based on your preference.
			if n <= 16 {
				c.Img.Set(x, y, image.White)
			}
		}
	}

	c.Save()
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
