package view

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
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
	containerParkingOut := container.New(layout.NewHBoxLayout())
	containerButtons := container.New(layout.NewHBoxLayout())

	restart := widget.NewButton("Restart Simulation", nil)
	exit := widget.NewButton("Menu Exit", nil)

	containerButtons.Add(restart)
	containerButtons.Add(exit)

	containerParkingOut.Add(p.MakeWaitStation())
	containerParkingOut.Add(layout.NewSpacer())
	containerParkingOut.Add(p.MakeExitStation())
	containerParkingOut.Add(layout.NewSpacer())

	containerParkingView.Add(containerParkingOut)
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParkingLotEntrance())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeEnterAndExitStation())
	containerParkingView.Add(layout.NewSpacer())
	containerParkingView.Add(p.MakeParking())
	containerParkingView.Add(layout.NewSpacer())

	containerParkingView.Add(container.NewCenter(containerButtons))
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
	carsWaitContainer := container.New(layout.NewGridLayout(5))
	for i := 0; i < 20; i++ {
		carsWaitContainer.Add(makeSquare())
	}
	return carsWaitContainer

}
func (p *ParkingView) MakeExitStation() *fyne.Container {
	return container.NewCenter(makeSquare())

}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	enterSquare := makeSquare()
	// enterSquare.Hide()
	parkingContainer.Add(enterSquare)
	parkingContainer.Add(layout.NewSpacer())
	exitSquare := makeSquare()
	// exitSquare.Hide()
	parkingContainer.Add(exitSquare)
	parkingContainer.Add(layout.NewSpacer())
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeParkingLotEntrance() *fyne.Container {
	EntraceContainer := container.New(layout.NewGridLayout(3))
	EntraceContainer.Add(makeBorder())
	EntraceContainer.Add(layout.NewSpacer())
	EntraceContainer.Add(makeBorder())
	return EntraceContainer
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

func makeBorder() *canvas.Rectangle {
	square := canvas.NewRectangle(color.RGBA{R: 255, G: 255, B: 255, A: 0})
	square.SetMinSize(fyne.NewSquareSize(float32(30)))
	square.StrokeColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
	square.StrokeWidth = float32(1)
	return square
}
