package view

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
)

type ParkingView struct {
	window fyne.Window
}

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window}
	parkingView.MakeScene()
	return parkingView
}

func (p *ParkingView) MakeScene() {
	containerParkingView := container.New(layout.NewVBoxLayout())
	containerParkingView.Add(p.MakeEnterAndExitStation())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParking())
	p.window.SetContent(containerParkingView)
}

func (p *ParkingView) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	for i := 0; i < 20; i++ {
		if i == 10 {
			addSpace(parkingContainer)
		}
		parkingContainer.Add(makeSquare())
	}
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeWaitStation() *fyne.Container {
	return nil

}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	parkingContainer.Add(makeSquare())
	parkingContainer.Add(layout.NewSpacer())
	parkingContainer.Add(makeSquare())
	parkingContainer.Add(layout.NewSpacer())
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) AddParkingLotEntrance() *fyne.Container {
	return nil
}

func addSpace(parkingContainer *fyne.Container) {
	for j := 0; j < 5; j++ {
		parkingContainer.Add(layout.NewSpacer())
	}
}

func makeSquare() *canvas.Rectangle {
	square := canvas.NewRectangle(Gray)
	square.SetMinSize(fyne.NewSquareSize(float32(30)))
	return square
}
