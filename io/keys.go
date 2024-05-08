package io

import "github.com/gopxl/pixel/pixelgl"

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

func (display *Display) GetKeyPresses() [16]bool {
	var keyPresses [16]bool
	for k, v := range keysToIndexMap {
		if display.Window.Pressed(k) {
			keyPresses[v] = true
		}
	}
	return keyPresses
}
