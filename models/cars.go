package models

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

const (
	minDuration = 1
	maxDuration = 5
	lamba       = 2.0
)

type Car struct {
	ID         int64
	state      string
	rectangule *canvas.Rectangle
	time       int
}

func NewSpaceCar() *Car {

	rectangule := canvas.NewRectangle(color.RGBA{R: 30, G: 30, B: 30, A: 255})
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	car := &Car{
		ID:         -1,
		state:      "space",
		rectangule: rectangule,
		time:       0,
	}

	return car
}

const (
	lambda = 2.0
)

func NewCar(id int64) *Car {
	rangR := rand.Intn(255-130) + 130
	rangG := rand.Intn(255-130) + 130
	rangB := rand.Intn(255-130) + 130
	time := rand.Intn(5-1) + 1

	rectangule := canvas.NewRectangle(color.RGBA{R: uint8(rangR), G: uint8(rangG), B: uint8(rangB), A: 255})
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	car := &Car{
		ID:         int64(id),
		state:      "waitStation",
		rectangule: rectangule,
		time:       time,
	}

	return car
}

func (c *Car) StartCount() {
	for {
		c.time--
		time.Sleep(1 * time.Second)
		if c.time == 0 {
			return
		}
	}
}

func (c *Car) GetRectangle() *canvas.Rectangle {
	return c.rectangule
}

func (c *Car) GetState() string {
	return c.state
}

func (c *Car) GetTime() int {
	return c.time
}

func (c *Car) GetID() int64 {
	return c.ID
}

func (c *Car) SetState(state string) {
	c.state = state
}
