package main

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	_ "image/png"
	"math"
	"math/rand"
	"os"
	"time"
)

const screenWidth = 1024
const screenHeight = 768

const initialAsteroids = 75

type entity struct {
	x      float64
	y      float64
	dx     float64
	dy     float64
	angle  float64
	scale  float64
	sprite *pixel.Sprite
}

func (e entity) collidesWith(x, y, r float64) bool {

	return math.Sqrt(math.Pow(x-e.x, 2)+math.Pow(y-e.y, 2)) < r

}

var (
	windowTitlePrefix = "Go Asteroids"
	frames            = 0
	second            = time.Tick(time.Second)
	window            *pixelgl.Window
	ship              entity
	asteroids         []entity
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

	shipImage, initError := loadImageFile("ship.png")
	if initError != nil {
		panic(initError)
	}

	asteroidImage, initError := loadImageFile("asteroid.png")
	if initError != nil {
		panic(initError)
	}

	shipPic := pixel.PictureDataFromImage(shipImage)

	ship = entity{
		x:      float64(screenWidth / 2),
		y:      float64(screenHeight / 2),
		dx:     0,
		dy:     0,
		angle:  0.0,
		sprite: pixel.NewSprite(shipPic, shipPic.Bounds()),
		scale:  0.2,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	asteroidPic := pixel.PictureDataFromImage(asteroidImage)

	asteroids = make([]entity, initialAsteroids)

	for i := 0; i < initialAsteroids; i++ {

		var x, y float64

		okPosition := false
		for !okPosition {
			x = r.Float64() * screenWidth
			y = r.Float64() * screenHeight
			okPosition = true
			for j := 0; j < i; j++ {
				if asteroids[j].collidesWith(x, y, 80) {
					okPosition = false
				}
			}
		}

		asteroids[i] = entity{
			x:      x,
			y:      y,
			dx:     r.Float64()*100 - 50,
			dy:     r.Float64()*100 - 50,
			angle:  r.Float64() * 2 * math.Pi,
			sprite: pixel.NewSprite(asteroidPic, asteroidPic.Bounds()),
			scale:  0.1,
		}
	}

}

func game() {

	initiate()

	for !window.Closed() {

		frameStart := time.Now()

		ship.dx, ship.dy = 0.0, 0.0

		if window.Pressed(pixelgl.KeyLeft) {
			ship.angle += 2 * frameLength
		}
		if window.Pressed(pixelgl.KeyRight) {
			ship.angle -= 2 * frameLength
		}
		if window.Pressed(pixelgl.KeyW) {
			ship.dx -= 512 * math.Sin(ship.angle)
			ship.dy += 512 * math.Cos(ship.angle)
		}
		if window.Pressed(pixelgl.KeyS) {
			ship.dx += 512 * math.Sin(ship.angle)
			ship.dy -= 512 * math.Cos(ship.angle)
		}
		if window.Pressed(pixelgl.KeyA) {
			ship.dx -= 512 * math.Cos(ship.angle)
			ship.dy -= 512 * math.Sin(ship.angle)
		}
		if window.Pressed(pixelgl.KeyD) {
			ship.dx += 512 * math.Cos(ship.angle)
			ship.dy += 512 * math.Sin(ship.angle)
		}

		window.Clear(colornames.Black)

		for i, a := range asteroids {

			asteroids[i].x += a.dx * frameLength
			asteroids[i].y += a.dy * frameLength

			if asteroids[i].x < -50 {
				asteroids[i].x += screenWidth + 100
			}
			if asteroids[i].y < -50 {
				asteroids[i].y += screenHeight + 100
			}
			if asteroids[i].x > screenWidth+50 {
				asteroids[i].x -= screenWidth + 100
			}
			if asteroids[i].y > screenHeight+50 {
				asteroids[i].y -= screenHeight + 100
			}

			asteroidMatrix := pixel.IM.Rotated(pixel.ZV, a.angle).Scaled(pixel.ZV, a.scale).Moved(pixel.Vec{X: a.x, Y: a.y})
			a.sprite.Draw(window, asteroidMatrix)
		}

		ship.x += ship.dx * frameLength
		ship.y += ship.dy * frameLength

		shipMatrix := pixel.IM.Rotated(pixel.ZV, ship.angle).Scaled(pixel.ZV, ship.scale).Moved(pixel.Vec{X: ship.x, Y: ship.y})
		ship.sprite.Draw(window, shipMatrix)
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
