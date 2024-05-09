# Chip8 Emulator

## Games

The emulator was tested using the IBM Logo, Pong, and Tetris ROMs.

-- insert pictures here --

## Setup

Since this emulator supports running with either the [https://github.com/gopxl/pixel](https://github.com/gopxl/pixel) game engine or the [https://ebitengine.org](https://ebitengine.org) game engine, you need to follow their requirements for compilation depending on your platform.

For Ubuntu, you need the `libgl1-mesa-dev` and `xorg-dev` packages according to [https://github.com/gopxl/pixel#requirements](https://github.com/gopxl/pixel#requirements).

According to [https://ebitengine.org/en/documents/install.html?os=linux](https://ebitengine.org/en/documents/install.html?os=linux), you also need the following packages for Ubuntu:
`sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config`

The program takes three flags:

- -filePath path/to/rom/file
  - default: ./roms/Pong.ch8
- -usePixelEngine
  - default: false
- -executionRate 1234
  - default: 700 (Hz)

In order to run it for the ebitengine on WSL, you need to set the `GOOS` variable when starting the process. E.g:

`GOOS=windows go run .`