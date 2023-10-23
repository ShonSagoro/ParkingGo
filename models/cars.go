package models

import (
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Car struct {
	ID         int64
	rectangule *canvas.Rectangle
	time       int
	color      color.Color
}

const (
	minDuration = 1
	maxDuration = 5
)

var (
	exitCars []*Car
)

func NewSpaceCar() *Car {

	colorRectangle := color.RGBA{R: 30, G: 30, B: 30, A: 255}
	rectangule := canvas.NewRectangle(colorRectangle)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))
	car := &Car{
		ID:         -1,
		rectangule: rectangule,
		time:       0,
		color:      colorRectangle,
	}

	return car
}

func NewCar(id int64) *Car {
	rangR := rand.Intn(255-130) + 130
	rangG := rand.Intn(255-130) + 130
	rangB := rand.Intn(255-130) + 130
	colorRectangle := color.RGBA{R: uint8(rangR), G: uint8(rangG), B: uint8(rangB), A: 255}
	time := rand.Intn(maxDuration-minDuration) + minDuration

	rectangule := canvas.NewRectangle(colorRectangle)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	car := &Car{
		ID:         int64(id),
		rectangule: rectangule,
		time:       time,
		color:      colorRectangle,
	}

	return car
}

func (c *Car) StartCount() {
	for {
		c.time--
		time.Sleep(1 * time.Second)
		if c.time == 0 {
			exitCars = append(exitCars, c)
			return
		}
	}
}

func (c *Car) GetRectangle() *canvas.Rectangle {
	return c.rectangule
}

func (c *Car) GetTime() int {
	return c.time
}

func (c *Car) GetID() int64 {
	return c.ID
}

func GetWaitCars() []*Car {
	return exitCars
}

func WaitCarsIsEmpty() bool {
	return len(exitCars) == 0
}
