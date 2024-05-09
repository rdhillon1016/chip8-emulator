package io

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/rdhillon1016/chip8-emulator/chip8"
)

const (
	windowWidth  = 1024
	windowHeight = 512
)

var keysToIndexMap map[ebiten.Key]uint = map[ebiten.Key]uint{
	ebiten.KeyX: 0x0,
	ebiten.Key1: 0x1,
	ebiten.Key2: 0x2,
	ebiten.Key3: 0x3,
	ebiten.KeyQ: 0x4,
	ebiten.KeyW: 0x5,
	ebiten.KeyE: 0x6,
	ebiten.KeyA: 0x7,
	ebiten.KeyS: 0x8,
	ebiten.KeyD: 0x9,
	ebiten.KeyZ: 0xA,
	ebiten.KeyC: 0xB,
	ebiten.Key4: 0xC,
	ebiten.KeyR: 0xD,
	ebiten.KeyF: 0xE,
	ebiten.KeyV: 0xF,
}

type Game struct {
	chip *chip8.Chip
}

func (g *Game) Update() error {
	g.chip.SetKeys(getKeyPresses())
	g.chip.ExecuteCycle()
	return nil
}

func getKeyPresses() [16]bool {
	var keyPresses [16]bool
	for k, v := range keysToIndexMap {
		keyPresses[v] = ebiten.IsKeyPressed(k)
	}
	return keyPresses
}

func (g *Game) Draw(screen *ebiten.Image) {
	pixels := g.chip.Pixels

	scaledUpPixelWidth := windowWidth / len(pixels)
	scaledUpPixelHeight := windowHeight / len(pixels[0])
	pixelImg := ebiten.NewImage(scaledUpPixelWidth, scaledUpPixelHeight)
	pixelImg.Fill(color.RGBA{0x0b, 0xd3, 0xd3, 0xff})

	imgOptions := &ebiten.DrawImageOptions{}

	for _, column := range pixels {
		for _, val := range column {
			if val {
				screen.DrawImage(pixelImg, imgOptions)
			}
			imgOptions.GeoM.Translate(0, float64(scaledUpPixelHeight))
		}
		imgOptions.GeoM.Translate(0, -float64(windowHeight))
		imgOptions.GeoM.Translate(float64(scaledUpPixelWidth), 0)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1024, 512
}

func Run(c *chip8.Chip, executionRateHz int) {
	ebiten.SetWindowSize(windowWidth, windowHeight)
	ebiten.SetWindowTitle("Chip8")
	ebiten.SetTPS(executionRateHz)
	if err := ebiten.RunGame(&Game{chip: c}); err != nil {
		log.Fatal(err)
	}
}
