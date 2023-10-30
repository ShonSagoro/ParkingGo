package models

import (
	"fmt"
	"image/color"
	"math/rand"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
)

type Car struct {
	ID            int
	rectangule    *canvas.Rectangle
	time          int
	senRenderTime chan int
	color         color.Color
	semQuit       chan bool
}

const (
	minDuration = 5
	maxDuration = 7
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

func NewCar(id int, sRT chan int, sQ chan bool) *Car {
	rangR := rand.Intn(255-130) + 130
	rangG := rand.Intn(255-130) + 130
	rangB := rand.Intn(255-130) + 130
	colorRectangle := color.RGBA{R: uint8(rangR), G: uint8(rangG), B: uint8(rangB), A: 255}
	time := rand.Intn(maxDuration-minDuration) + minDuration

	rectangule := canvas.NewRectangle(colorRectangle)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	car := &Car{
		ID:            id,
		rectangule:    rectangule,
		time:          time,
		senRenderTime: sRT,
		semQuit:       sQ,
		color:         colorRectangle,
	}

	return car
}

func (c *Car) StartCount(id int) {
	for {
		select {
		case <-c.semQuit:
			fmt.Printf("StartCount Close")
			return
		default:
			c.time--
			c.senRenderTime <- id
			if c.time == 0 {
				c.ID = id
				exitCars = append(exitCars, c)
				return
			}
			time.Sleep(1 * time.Second)
		}
	}
}

func (c *Car) GetRectangle() *canvas.Rectangle {
	return c.rectangule
}

func (c *Car) GetTime() int {
	return c.time
}

func (c *Car) GetID() int {
	return c.ID
}

func GetWaitCars() []*Car {
	return exitCars
}
func PopWaitCars() *Car {
	car := exitCars[0]
	if !WaitExitCarsIsEmpty() {
		exitCars = exitCars[1:]
	}
	return car
}

func WaitExitCarsIsEmpty() bool {
	return len(exitCars) == 0
}
