# Chip8 Emulator

![Pong](<pong.png>)

![Tetris](<tetris.png>)

![IBM Logo](<ibm_logo.png>)

## Setup

According to [https://ebitengine.org/en/documents/install.html?os=linux](https://ebitengine.org/en/documents/install.html?os=linux), you need the following packages for Ubuntu:
`sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config`

The program takes three flags:

- -filePath path/to/rom/file
  - default: ./roms/Pong.ch8
- -executionRate 1234
  - default: 700 (Hz)

In order to run it for the ebitengine on WSL, you need to set the `GOOS` variable when building. E.g:

`GOOS=windows go run .`

## Sources

Game roms are available from [https://github.com/kripod/chip8-roms/tree/master/games](https://github.com/kripod/chip8-roms/tree/master/games).
A useful testing suite of ROMs is [https://github.com/Timendus/chip8-test-suite](https://github.com/Timendus/chip8-test-suite)