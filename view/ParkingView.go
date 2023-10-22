package view

import (
	"fmt"
	"image/color"
	"parking/models"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var semNewCar chan bool
var semFullWaitStation chan bool
var semHaveSpace chan bool
var semWait chan bool

type ParkingView struct {
	window           fyne.Window
	waitParkingArray *fyne.Container
	parkingArray     [20]*models.Car
}

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

var parking *models.Parking

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window, waitParkingArray: container.New(nil)}

	semNewCar = make(chan bool)
	semFullWaitStation = make(chan bool)
	semHaveSpace = make(chan bool)
	semWait = make(chan bool)

	parking = models.NewParking(semNewCar, semFullWaitStation, semHaveSpace, semWait)

	parkingView.RenderScene()
	parkingView.StartSimulation()

	return parkingView
}

func (p *ParkingView) RenderScene() {

	containerParkingView := container.New(layout.NewVBoxLayout())
	containerParkingOut := container.New(layout.NewHBoxLayout())
	containerButtons := container.New(layout.NewHBoxLayout())

	restart := widget.NewButton("Restart Simulation", nil)
	exit := widget.NewButton("Menu Exit",
		func() {
			dialog.ShowConfirm("Salir", "¿Desea salir de la aplicación?", func(response bool) {
				if response {
					p.BackToMenu()
				} else {
					close(semNewCar)
					close(semFullWaitStation)
					close(semHaveSpace)
					close(semWait)
					fmt.Println("No")
				}
			}, p.window)
		},
	)

	containerButtons.Add(restart)
	containerButtons.Add(exit)

	containerParkingOut.Add(p.waitParkingArray)
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
	p.window.CenterOnScreen()
}

func (p *ParkingView) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parking.MakeParking()
	fmt.Println("Mira hice algo xDDD")

	p.parkingArray = parking.GetParking()
	for i := 0; i < len(p.parkingArray); i++ {
		if i == 10 {
			addSpace(parkingContainer)
		}

		time := canvas.NewText(fmt.Sprintf("%d", p.parkingArray[i].GetTime()), color.RGBA{R: 255, G: 255, B: 255, A: 255})

		parkingContainer.Add(container.NewCenter(p.parkingArray[i].GetRectangle(), time))
	}
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) RenderWaitStation() {
	for {
		fmt.Printf("Renderizando")
		<-semNewCar
		waitContainer := container.New(layout.NewGridLayout(5))
		for i := len(parking.GetWaitCars()) - 1; i >= 0; i-- {
			car := parking.GetWaitCars()[i]
			waitContainer.Add(car.GetRectangle())
		}
		p.waitParkingArray = waitContainer
		p.RenderScene()
		fmt.Printf("Ya")
	}
}

func (p *ParkingView) MakeExitStation() *fyne.Container {
	return container.NewCenter(makeSquare())

}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	enterSquare := parking.GetEntraceCar().GetRectangle()
	// enterSquare.Hide()
	parkingContainer.Add(enterSquare)
	parkingContainer.Add(layout.NewSpacer())
	exitSquare := parking.GetExitCar().GetRectangle()
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

func (p *ParkingView) BackToMenu() {
	NewMainView(p.window)
}

func (p *ParkingView) StartSimulation() {
	go parking.GenerateCars()
	go p.RenderWaitStation()
	go parking.CheckParking()
	go parking.CarEntrace()
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
