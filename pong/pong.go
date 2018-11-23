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

type position struct {
	x, y float32
}

type ball struct {
	position
	radius int
	xv     float32
	yv     float32
	color  color.RGBA
}

func (b *ball) draw(pixels []byte) {
	for y := -b.radius; y < b.radius; y++ {
		for x := -b.radius; x < b.radius; x++ {
			if x*x+y*y < b.radius*b.radius {
				setPixel(int(b.x)+x, int(b.y)+y, b.color, pixels)
			}
		}
	}
}

func (b *ball) update(left *paddle, right *paddle) {
	b.x += b.xv
	b.y += b.yv

	// handle collisions
	if int(b.y)-b.radius < 0 || int(b.y)+b.radius > height {
		b.yv = -b.yv
	}

	if b.x < 0 || int(b.x) > width {
		b.x = 300
		b.y = 300
	}

	if int(b.x)-b.radius < int(left.x)+left.w/2 &&
		(int(b.y) > int(left.y)-left.h/2 && int(b.y) < int(left.y)+left.h/2) {
		b.xv = -b.xv
	}

	if int(b.x)+b.radius > int(right.x)-right.w/2 &&
		(int(b.y) > int(right.y)-right.h/2 && int(b.y) < int(right.y)+right.h/2) {
		b.xv = -b.xv
	}
}

type paddle struct {
	position
	w     int
	h     int
	color color.RGBA
}

func (p *paddle) draw(pixels []byte) {
	startX := int(p.x) - p.w/2
	startY := int(p.y) - p.h/2

	for y := 0; y < p.h; y++ {
		for x := 0; x < p.w; x++ {
			setPixel(startX+x, startY+y, p.color, pixels)
		}
	}
}

func (p *paddle) update(keyState []uint8) {
	if keyState[sdl.SCANCODE_UP] != 0 {
		p.y -= 10
	}

	if keyState[sdl.SCANCODE_DOWN] != 0 {
		p.y += 10
	}
}

func (p *paddle) aiUpdate(b *ball) {
	p.y = b.y
}

func clear(pixels []byte) {
	for i := range pixels {
		pixels[i] = 0
	}
}

func setPixel(x, y int, c color.RGBA, pixels []byte) {
	index := (y*width + x) * bytesPerPixel

	if index < len(pixels)-bytesPerPixel && index >= 0 {
		pixels[index] = c.R
		pixels[index+1] = c.G
		pixels[index+2] = c.B
	}
}

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer sdl.Quit()

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
			setPixel(x, y, color.RGBA{0, 0, 0, 0}, pixels)
		}
	}

	player1 := paddle{position{50, 100}, 20, 100, color.RGBA{255, 255, 255, 0}}
	player2 := paddle{position{float32(width) - 50, 100}, 20, 100, color.RGBA{255, 255, 255, 0}}
	ball := ball{position{300, 300}, 20, 10, 10, color.RGBA{255, 255, 255, 0}}

	keyState := sdl.GetKeyboardState()

	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}

		clear(pixels)

		player1.update(keyState)
		ball.update(&player1, &player2)
		player2.aiUpdate(&ball)

		player1.draw(pixels)
		ball.draw(pixels)
		player2.draw(pixels)

		texture.Update(nil, pixels, width*bytesPerPixel)
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}
