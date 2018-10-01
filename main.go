package main

import (
	_ "image/png"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel"
	"golang.org/x/image/colornames"
	"fmt"
	"time"
	"image"
	"os"
	"math"
)

const screenWidth = 1024
const screenHeight = 768

var (
	windowTitlePrefix   = "Go Asteroids"
	frames                                    = 0
	second                                    = time.Tick(time.Second)
	window            *pixelgl.Window
	shipSprite        *pixel.Sprite
	frameLength		float64
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

	shipImage, initError := loadImageFile("ship.png")
	if initError != nil {
		panic(initError)
	}

	shipPic := pixel.PictureDataFromImage(shipImage)

	shipSprite = pixel.NewSprite(shipPic, shipPic.Bounds())

}

func game() {

	initiate()

	x, y, angle := float64(screenWidth/2), float64(screenHeight/2), 0.0

	for !window.Closed() {

		frameStart := time.Now()

		if window.Pressed(pixelgl.KeyLeft) {
			angle += 2 * frameLength
		}
		if window.Pressed(pixelgl.KeyRight) {
			angle -= 2 * frameLength
		}
		if window.Pressed(pixelgl.KeyW) {
			x -= 512 * frameLength * math.Sin(angle)
			y += 512 * frameLength * math.Cos(angle)
		}
		if window.Pressed(pixelgl.KeyS) {
			x += 512 * frameLength * math.Sin(angle)
			y -= 512 * frameLength * math.Cos(angle)
		}
		if window.Pressed(pixelgl.KeyA) {
			x -= 512 * frameLength * math.Cos(angle)
			y -= 512 * frameLength * math.Sin(angle)
		}
		if window.Pressed(pixelgl.KeyD) {
			x += 512 * frameLength * math.Cos(angle)
			y += 512 * frameLength * math.Sin(angle)
		}

		matrix := pixel.IM.Rotated(pixel.ZV, angle).Scaled(pixel.ZV, 0.2).Moved(pixel.Vec{X: x, Y: y})

		window.Clear(colornames.Black)
		shipSprite.Draw(window, matrix)
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

