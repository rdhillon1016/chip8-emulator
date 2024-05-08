package io

import (
	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

const (
	windowWidth = 1024
	windowHeight = 512
)

type Display struct {
	Window *pixelgl.Window
	imd *imdraw.IMDraw
}

func NewDisplay() *Display {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync: true,
	}

	win, err := pixelgl.NewWindow(cfg)
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	win.Clear(colornames.Skyblue)
	win.Update()
	if err != nil {
		panic(err)
	}
	return &Display{
		Window: win,
		imd: imd,
	}
}

func (display *Display) UpdateScreen(pixels [][]bool) {
	scaledUpPixelWidth := windowWidth / len(pixels)
	scaledUpPixelHeight := windowHeight / len(pixels[0])
	
	display.Window.Clear(colornames.Skyblue)
	display.imd.Clear()

	for x, column := range pixels {
		for j, val := range column {
			if val {
				// Flip y-index since the drawing starts from lower left
				y := 32 - j - 1
				scaledXCoordinateBottomLeft := x*scaledUpPixelWidth
				scaledYCoordinateBottomLeft := y*scaledUpPixelHeight
				scaledXCoordinateTopRight := scaledXCoordinateBottomLeft + scaledUpPixelWidth
				scaledYCoordinateTopRight := scaledYCoordinateBottomLeft + scaledUpPixelHeight
				display.imd.Push(pixel.V(float64(scaledXCoordinateBottomLeft), float64(scaledYCoordinateBottomLeft)))
				display.imd.Push(pixel.V(float64(scaledXCoordinateTopRight), float64(scaledYCoordinateTopRight)))
				display.imd.Rectangle(0)
			}
		}
	}

	display.imd.Draw(display.Window)
	display.Window.Update()
}
