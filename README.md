# Mandelbrot Shenanigans

A cli to print beautiful images of the Mandelbrot and Julia sets.

## Installation

After cloning the repo, `cd` into the project directory and build the executable with

`go build .`

## Usage

Out of the box, not much will happen, but you can use the commands below to set up and run the program.

- i <float>

  c's imaginary component.

- iterations <int>

  how many operations until considering a point bounded. (default 100)

- real <float>

  c's real component.

- res <int>

  The width and height of the canvas. (default 4096)

- xmax <float>

  The maximum value on the x axis. (default 1.2)

- xmin <float>

  The minimum value on the x axis. (default -1.2)

- ymax <float>

  The maximum value on the y axis. (default 1.2)

- ymin <float>

  The minimum value on the y axis. (default -1.2)

### Example

`./mandelbrot -real=0.285 -i=0.01`
