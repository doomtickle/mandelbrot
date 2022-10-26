# Mandelbrot Renders

A cli to print beautiful images of the Mandelbrot and Julia sets.

## Installation

After cloning the repo, `cd` into the project directory and build the executable with

`go build .`

## Usage

Out of the box, not much will happen, but you can use the commands below to set up and run the program.

- `-i` <float>

  c's imaginary component.

- `-real` <float>

  c's real component.

**Note: if you do not specify a real or imaginary component, then a standard mandelbrot
set will be rendered instead of a julia set.**

- `-iterations` <int>

  how many operations until considering a point bounded. (default 100)

- `-res` <int>

  The width and height of the canvas. (default 4096)

- `-xmax` <float>

  The maximum value on the x axis. (default 1.2)

- `-xmin` <float>

  The minimum value on the x axis. (default -1.2)

- `-ymax` <float>

  The maximum value on the y axis. (default 1.2)

- `-ymin` <float>

  The minimum value on the y axis. (default -1.2)

- `-out` <string>

  The output path for generated images. (default `image.png`)

- `-palette` <string>

  The name of a color palette specified in `mandelbrot.json`. (default `blue`)

- `-bg` <string>

  Hexadecimal string used for the canvas background color. (default `#fff`)

### Example

`./mandelbrot -real=0.285 -i=0.01 --out=blue_julia.png`
`/mandelbrot -out=mandelbrot.png -bg=#333 -xmin=-2 -xmax=0.8 -ymin=-1.4 -ymax=1.4`
