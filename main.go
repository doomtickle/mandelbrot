package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/doomtickle/mandelbrot/canvas"
)

// let's set some flags so we can modify stuff from the cli
var (
	bg         = flag.String("bg", "#fff", "hex value for background color")
	jobs       = flag.Bool("jobs", false, "Specifies whether to run jobs from ./mandelbrot.json")
	res        = flag.Int("res", 4096, "The width and height of the canvas.")
	xmin       = flag.Float64("xmin", -1.2, "The minimum value on the x axis.")
	xmax       = flag.Float64("xmax", 1.2, "The maximum value on the x axis.")
	ymin       = flag.Float64("ymin", -1.2, "The minimum value on the y axis.")
	ymax       = flag.Float64("ymax", 1.2, "The maximum value on the y axis.")
	out        = flag.String("o", "image.png", "output path for generated file.")
	cReal      = flag.Float64("r", 0, "c's real component.")
	palette    = flag.String("p", "blue", "color palette from your mandelbrot.json config")
	cImaginary = flag.Float64("im", 0, "c's imaginary component.")
	iterations = flag.Int("iter", 100, "how many operations until considering a point bounded.")
)

type rangemap func(float64) (float64, bool)

func mandelbrot(c canvas.Canvas) {
	// Instantiate our range maps to normalize x,y values to be between the min and max values specified at runtime.
	rmX := newRangeMap(rangeBounds{0, float64(c.Width)}, rangeBounds{*xmin, *xmax})
	rmY := newRangeMap(rangeBounds{0, float64(c.Height)}, rangeBounds{*ymin, *ymax})
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
			ca := a
			cb := b
			// iteration counter
			n := 0
			for n < *iterations {
				// Math stuff that I had to watch videos about and look up....
				// This is where we apply the zeta function to each pixel recursively.
				// The function will either converge to a value or blow up to infinity.
				aSquared := a*a - b*b
				twoAB := 2 * a * b

				a = aSquared + ca
				b = twoAB + cb

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
			if n <= 5 {
				c.Img.Set(x, y, c.Bg)
			}
		}
	}
	c.Save(*out)
}

func julia(c canvas.Canvas) {
	// Instantiate our range maps to normalize x,y values to be between the min and max values specified at runtime.
	rmX := newRangeMap(rangeBounds{0, float64(c.Width)}, rangeBounds{*xmin, *xmax})
	rmY := newRangeMap(rangeBounds{0, float64(c.Height)}, rangeBounds{*ymin, *ymax})
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
				aSquared := a*a - b*b
				twoAB := 2 * a * b

				a = aSquared + *cReal
				b = twoAB + *cImaginary
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

func runJobs() {
	var jc canvas.JsonConfig
	f, err := os.ReadFile("mandelbrot.json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(f, &jc)
	if err != nil {
		log.Fatal(err)
	}

	if len(jc.Jobs) > 0 {
		for _, j := range jc.Jobs {
			// set runtime variables
			*out = j.Out
			*res = j.Res
			*xmin = j.XMin
			*xmax = j.XMax
			*ymin = j.YMin
			*ymax = j.YMax
			*cReal = j.Real
			*cImaginary = j.Imaginary
			bgColor, err := canvas.ParseHexColorFast(j.Bg)
			if err != nil {
				log.Fatal(err)
			}
			// Generate a new canvas
			c := canvas.New(j.Res, j.Palette, bgColor)

			fmt.Printf("%#v", j)

			if j.Real != 0 || j.Imaginary != 0 {
				julia(*c)
			} else {
				mandelbrot(*c)
			}

		}
	}
}

func main() {
	flag.Parse()

	if *jobs {
		runJobs()
		return
	}

	bgColor, err := canvas.ParseHexColorFast(*bg)
	if err != nil {
		log.Fatal(err)
	}
	// Generate a new canvas
	c := canvas.New(*res, *palette, bgColor)

	if *cReal != 0 || *cImaginary != 0 {
		julia(*c)
		return
	}

	mandelbrot(*c)
	return
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
