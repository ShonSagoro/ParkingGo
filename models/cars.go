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
	ID         int
	rectangule *canvas.Rectangle
	text       *canvas.Text
	time       int
	semQuit    chan bool
}

const (
	minDuration = 5
	maxDuration = 7
)

var (
	exitCars []*Car
)

func NewSpaceCar() *Car {

	rectangule := canvas.NewRectangle(Gray)

	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	text := canvas.NewText(fmt.Sprintf("%d", 0), White)

	car := &Car{
		ID:         -1,
		rectangule: rectangule,
		time:       0,
		text:       text,
	}

	return car
}

func NewCar(id int, sQ chan bool) *Car {
	rangR := rand.Intn(255-130) + 130
	rangG := rand.Intn(255-130) + 130
	rangB := rand.Intn(255-130) + 130
	colorRectangle := color.RGBA{R: uint8(rangR), G: uint8(rangG), B: uint8(rangB), A: 255}
	time := rand.Intn(maxDuration-minDuration) + minDuration

	rectangule := canvas.NewRectangle(colorRectangle)
	rectangule.SetMinSize(fyne.NewSquareSize(float32(30)))

	text := canvas.NewText(fmt.Sprintf("%d", time), White)

	car := &Car{
		ID:         id,
		rectangule: rectangule,
		time:       time,
		text:       text,
		semQuit:    sQ,
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
			c.text.Text = fmt.Sprintf("%d", c.time)
			c.text.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}

			if c.time <= 0 {
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

func (c *Car) GetText() *canvas.Text {
	return c.text
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

func PopExitWaitCars() *Car {
	car := exitCars[0]
	if !WaitExitCarsIsEmpty() {
		exitCars = exitCars[1:]
	}
	return car
}

func WaitExitCarsIsEmpty() bool {
	return len(exitCars) == 0
}
