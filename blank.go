package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"os"
	"time"
)

const screenWidth = 1024
const screenHeight = 768

var (
	windowTitlePrefix = "Go Asteroids"
	frames            = 0
	second            = time.Tick(time.Second)
	window            *pixelgl.Window
	frameLength       float64
)

func loadImageFile(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func initiate() {

	var initError error

	cfg := pixelgl.WindowConfig{
		Bounds: pixel.R(0, 0, screenWidth, screenHeight),
		VSync:  true,
	}

	window, initError = pixelgl.NewWindow(cfg)
	if initError != nil {
		panic(initError)
	}

}

func game() {

	initiate()

	for !window.Closed() {

		frameStart := time.Now()

		window.Clear(colornames.Black)

		window.Update()

		frames++
		select {
		case <-second:
			window.SetTitle(fmt.Sprintf("%s | FPS: %d", windowTitlePrefix, frames))
			frames = 0
		default:
		}

		frameLength = time.Since(frameStart).Seconds()

	}
}

func main() {

	pixelgl.Run(game)

}
