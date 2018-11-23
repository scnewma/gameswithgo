package main

import (
	"fmt"
	"image/color"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height int = 800, 600
	bytesPerPixel     = 4
)

func setPixel(x, y int, c color.RGBA, pixels []byte) {
	index := (y*width + x) * bytesPerPixel

	if index < len(pixels)-bytesPerPixel && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

func main() {
	window, err := sdl.CreateWindow("Testing SDL2",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(width), int32(height), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_ABGR8888, sdl.TEXTUREACCESS_STREAMING, int32(width), int32(height))
	if err != nil {
		fmt.Println(err)
		return
	}
	defer texture.Destroy()

	pixels := make([]byte, width*height*bytesPerPixel)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			setPixel(x, y, color.RGBA{byte(x % 255), byte(y % 255), 0, 0}, pixels)
		}
	}
	texture.Update(nil, pixels, width*bytesPerPixel)
	renderer.Copy(texture, nil, nil)
	renderer.Present()

	sdl.Delay(5000)
}
