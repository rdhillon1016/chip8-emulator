package pixelengine

import (
	"time"

	"github.com/gopxl/pixel"
	"github.com/gopxl/pixel/imdraw"
	"github.com/gopxl/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"github.com/rdhillon1016/chip8-emulator/chip8"
)

const (
	windowWidth  = 1024
	windowHeight = 512
)

type display struct {
	window *pixelgl.Window
	imd    *imdraw.IMDraw
}

func newDisplay() *display {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, windowWidth, windowHeight),
		VSync:  true,
	}

	win, err := pixelgl.NewWindow(cfg)
	imd := imdraw.New(nil)
	imd.Color = colornames.Black

	win.Clear(colornames.Skyblue)
	win.Update()
	if err != nil {
		panic(err)
	}
	return &display{
		window: win,
		imd:    imd,
	}
}

func (d *display) updateScreen(pixels [][]bool) {
	scaledUpPixelWidth := windowWidth / len(pixels)
	scaledUpPixelHeight := windowHeight / len(pixels[0])

	d.window.Clear(colornames.Skyblue)
	d.imd.Clear()

	for x, column := range pixels {
		for j, val := range column {
			if val {
				// Flip y-index since the drawing starts from lower left
				y := 32 - j - 1
				scaledXCoordinateBottomLeft := x * scaledUpPixelWidth
				scaledYCoordinateBottomLeft := y * scaledUpPixelHeight
				scaledXCoordinateTopRight := scaledXCoordinateBottomLeft + scaledUpPixelWidth
				scaledYCoordinateTopRight := scaledYCoordinateBottomLeft + scaledUpPixelHeight
				d.imd.Push(pixel.V(float64(scaledXCoordinateBottomLeft), float64(scaledYCoordinateBottomLeft)))
				d.imd.Push(pixel.V(float64(scaledXCoordinateTopRight), float64(scaledYCoordinateTopRight)))
				d.imd.Rectangle(0)
			}
		}
	}

	d.imd.Draw(d.window)
	d.window.Update()
}

func Run(c *chip8.Chip, executionRateHz int) {
	d := newDisplay()

	for !d.window.Closed() {
		
		if screenUpdated {
			d.updateScreen(c.Pixels)
		} else {
			d.window.UpdateInput()
		}
		time.Sleep(time.Second / time.Duration(executionRateHz))
	}
}

var keysToIndexMap map[pixelgl.Button]uint = map[pixelgl.Button]uint{
	pixelgl.KeyX: 0x0,
	pixelgl.Key1: 0x1,
	pixelgl.Key2: 0x2,
	pixelgl.Key3: 0x3,
	pixelgl.KeyQ: 0x4,
	pixelgl.KeyW: 0x5,
	pixelgl.KeyE: 0x6,
	pixelgl.KeyA: 0x7,
	pixelgl.KeyS: 0x8,
	pixelgl.KeyD: 0x9,
	pixelgl.KeyZ: 0xA,
	pixelgl.KeyC: 0xB,
	pixelgl.Key4: 0xC,
	pixelgl.KeyR: 0xD,
	pixelgl.KeyF: 0xE,
	pixelgl.KeyV: 0xF,
}

func (d *display) getKeyPresses() [16]bool {
	var keyPresses [16]bool
	for k, v := range keysToIndexMap {
		if d.window.Pressed(k) {
			keyPresses[v] = true
		}
	}
	return keyPresses
}