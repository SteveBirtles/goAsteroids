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

const initialAsteroids = 20

type etype int

const (
	Ship       etype = 1
	Asteroid   etype = 2
	Projectile etype = 3
)

type entity struct {
	etype
	x      float64
	y      float64
	dx     float64
	dy     float64
	radius float64
	angle  float64
	scale  float64
	sprite *pixel.Sprite
}

func (e entity) separation(e2 entity) float64 {

	return math.Sqrt(math.Pow(e.x-e2.x, 2) + math.Pow(e.y-e2.y, 2))

}

func (e entity) collidesWith(e2 entity) bool {

	return e.separation(e2) < e.radius+e2.radius

}

func (e entity) velocity() float64 {

	return math.Sqrt(math.Pow(e.dx, 2) + math.Pow(e.dy, 2))

}

var (
	windowTitlePrefix = "Go Asteroids"
	frames            = 0
	second            = time.Tick(time.Second)
	window            *pixelgl.Window
	frameLength       float64
	es                []entity
	shipPic           pixel.Picture
	asteroidPic       pixel.Picture
	fireballPic       pixel.Picture
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
	shipPic = pixel.PictureDataFromImage(shipImage)

	asteroidImage, initError := loadImageFile("asteroid.png")
	if initError != nil {
		panic(initError)
	}
	asteroidPic = pixel.PictureDataFromImage(asteroidImage)

	fireballImage, initError := loadImageFile("fireball.png")
	if initError != nil {
		panic(initError)
	}

	fireballPic = pixel.PictureDataFromImage(fireballImage)

	es = make([]entity, initialAsteroids+1)

	es[0] = entity{
		etype:  Ship,
		x:      float64(screenWidth / 2),
		y:      float64(screenHeight / 2),
		dx:     0,
		dy:     0,
		angle:  0.0,
		radius: 30,
		sprite: pixel.NewSprite(shipPic, shipPic.Bounds()),
		scale:  0.2,
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 1; i <= initialAsteroids; i++ {

		var x, y float64

		okPosition := false
		for !okPosition {
			x = r.Float64() * screenWidth
			y = r.Float64() * screenHeight
			okPosition = true
			for j := 0; j < i; j++ {
				if es[i].collidesWith(es[j]) {
					okPosition = false
				}
			}
		}

		es[i] = entity{
			etype:  Asteroid,
			x:      x,
			y:      y,
			dx:     r.Float64()*100 - 50,
			dy:     r.Float64()*100 - 50,
			angle:  r.Float64() * 2 * math.Pi,
			sprite: pixel.NewSprite(asteroidPic, asteroidPic.Bounds()),
			scale:  0.1,
			radius: 45,
		}
	}

}

func game() {

	initiate()

	for !window.Closed() {

		frameStart := time.Now()

		if window.Pressed(pixelgl.KeyLeft) {
			es[0].angle += 2 * frameLength
		}
		if window.Pressed(pixelgl.KeyRight) {
			es[0].angle -= 2 * frameLength
		}
		if window.Pressed(pixelgl.KeyW) {
			es[0].dx -= 25 * math.Sin(es[0].angle)
			es[0].dy += 25 * math.Cos(es[0].angle)
		}
		if window.Pressed(pixelgl.KeyS) {
			es[0].dx += 25 * math.Sin(es[0].angle)
			es[0].dy -= 25 * math.Cos(es[0].angle)
		}
		if window.Pressed(pixelgl.KeyA) {
			es[0].dx -= 25 * math.Cos(es[0].angle)
			es[0].dy -= 25 * math.Sin(es[0].angle)
		}
		if window.Pressed(pixelgl.KeyD) {
			es[0].dx += 25 * math.Cos(es[0].angle)
			es[0].dy += 25 * math.Sin(es[0].angle)
		}

		if window.Pressed(pixelgl.KeySpace) {

			projDx := -math.Sin(es[0].angle)
			projDy := math.Cos(es[0].angle)

			es = append(es, entity{
				etype:  Projectile,
				x:      es[0].x + es[0].radius*projDx,
				y:      es[0].y + es[0].radius*projDy,
				dx:     500 * projDx,
				dy:     500 * projDy,
				angle:  es[0].angle,
				radius: 5,
				sprite: pixel.NewSprite(fireballPic, fireballPic.Bounds()),
				scale:  0.05,
			})

		}

		for i := 0; i < len(es); {

			remove := false

			for j := 0; j < i; j++ {

				if es[i].collidesWith(es[j]) {

					if es[i].etype == Projectile {

						if es[j].etype == Asteroid {
							remove = true
							break
						} else {
							continue
						}
					}

					d := es[i].separation(es[j])
					dx := es[i].x - es[j].x
					dy := es[i].y - es[j].y

					v1 := es[i].velocity()
					v2 := es[j].velocity()

					es[i].dx = v2 * dx / d
					es[i].dy = v2 * dy / d

					es[j].dx = -v1 * dx / d
					es[j].dy = -v1 * dy / d

					break

				}

			}

			if remove {
				es = append(es[:i], es[i+1:]...)
			} else {
				i++
			}
		}

		for i := range es {

			es[i].x += es[i].dx * frameLength
			es[i].y += es[i].dy * frameLength

			if es[i].x < -50 {
				es[i].x += screenWidth + 100
			}
			if es[i].y < -50 {
				es[i].y += screenHeight + 100
			}
			if es[i].x > screenWidth+50 {
				es[i].x -= screenWidth + 100
			}
			if es[i].y > screenHeight+50 {
				es[i].y -= screenHeight + 100
			}

			v := es[i].velocity()
			if es[i].etype == Ship {
				if v > 256 {
					es[i].dx *= 256 / v
					es[i].dy *= 256 / v
				} else {
					es[i].dx *= 1 - frameLength
					es[i].dy *= 1 - frameLength
				}
			} else if es[i].etype == Asteroid {
				if v > 128 {
					es[i].dx *= 128 / v
					es[i].dy *= 128 / v
				}
			}
		}

		window.Clear(colornames.Black)

		for i := range es {

			matrix := pixel.IM.
				Rotated(pixel.ZV, es[i].angle).
				Scaled(pixel.ZV, es[i].scale).
				Moved(pixel.Vec{X: es[i].x, Y: es[i].y})

			es[i].sprite.Draw(window, matrix)

		}

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
