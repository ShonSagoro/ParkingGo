package view

import (
	"fmt"
	"image/color"
	"parking/models"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

var semRenderNewCarWait chan bool
var semRenderNewCarParking chan bool
var semRenderNewCarEnter chan bool
var semRenderNewCarExit chan bool
var semRenderIDExit chan int
var senRenderTime chan int
var semQuit chan bool

const maxWait int = 10

type (
	ParkingView struct {
		window               fyne.Window
		waitRectangleStation [maxWait]*canvas.Rectangle
		carsParking          [20]*CarPark
		entrace              *canvas.Rectangle
		exit                 *canvas.Rectangle
		out                  *canvas.Rectangle
	}

	CarPark struct {
		text      *canvas.Text
		rectangle *canvas.Rectangle
	}
)

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

var parking *models.Parking

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window}

	semRenderNewCarWait = make(chan bool)
	semRenderNewCarParking = make(chan bool)
	semRenderNewCarEnter = make(chan bool)
	semRenderNewCarExit = make(chan bool)
	senRenderTime = make(chan int)
	semRenderIDExit = make(chan int)

	parking = models.NewParking(semRenderNewCarWait, semRenderNewCarParking, semRenderNewCarEnter, semRenderNewCarExit, semRenderIDExit, senRenderTime, semQuit)

	parkingView.MakeScene()
	parkingView.StartSimulation()

	return parkingView
}

func NewParkCar(text *canvas.Text, rectangle *canvas.Rectangle) *CarPark {
	return &CarPark{
		text:      text,
		rectangle: rectangle,
	}
}

func (p *ParkingView) MakeScene() {

	containerParkingView := container.New(layout.NewVBoxLayout())
	containerParkingOut := container.New(layout.NewHBoxLayout())
	containerButtons := container.New(layout.NewHBoxLayout())

	restart := widget.NewButton("Restart Simulation", func() {
		dialog.ShowConfirm("Salir", "¿Desea Reiniciar la aplicación?", func(response bool) {
			if response {
				p.RestartSimulation()
			} else {
				fmt.Println("No")
			}
		}, p.window)
	})

	exit := widget.NewButton("Menu Exit",
		func() {
			dialog.ShowConfirm("Salir", "¿Desea salir de la aplicación?", func(response bool) {
				if response {
					p.BackToMenu()
				} else {
					fmt.Println("No")
				}
			}, p.window)
		},
	)

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
	p.window.Resize(fyne.NewSize(300, 800))
	p.window.CenterOnScreen()
}

func (p *ParkingView) MakeParking() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parking.MakeParking()

	parkingArray := parking.GetParking()
	for i := 0; i < len(parkingArray); i++ {
		if i == 10 {
			addSpace(parkingContainer)
		}
		time := canvas.NewText(fmt.Sprintf("%d", parkingArray[i].GetTime()), color.RGBA{R: 255, G: 255, B: 255, A: 255})
		car := NewParkCar(time, parkingArray[i].GetRectangle())
		p.carsParking[i] = car
		parkingContainer.Add(container.NewCenter(p.carsParking[i].rectangle, p.carsParking[i].text))
	}
	return container.NewCenter(parkingContainer)
}

func (p *ParkingView) MakeWaitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	for i := len(p.waitRectangleStation) - 1; i >= 0; i-- {
		car := models.NewSpaceCar().GetRectangle()
		p.waitRectangleStation[i] = car
		p.waitRectangleStation[i].Hide()
		parkingContainer.Add(p.waitRectangleStation[i])
	}
	return parkingContainer
}

func (p *ParkingView) MakeExitStation() *fyne.Container {
	p.out = makeSquare()
	return container.NewCenter(p.out)
}

func (p *ParkingView) MakeEnterAndExitStation() *fyne.Container {
	parkingContainer := container.New(layout.NewGridLayout(5))
	parkingContainer.Add(layout.NewSpacer())
	p.entrace = parking.GetEntraceCar().GetRectangle()
	parkingContainer.Add(p.entrace)
	parkingContainer.Add(layout.NewSpacer())
	p.exit = parking.GetExitCar().GetRectangle()
	parkingContainer.Add(p.exit)
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

func (p *ParkingView) RenderNewCarStation() {
	for {
		select {
		case <-semQuit:
			return
		default:
			<-semRenderNewCarWait
			waitCars := parking.GetWaitCars()
			for i := len(waitCars) - 1; i >= 0; i-- {
				if waitCars[i].ID != -1 {
					p.waitRectangleStation[i].Show()
					p.waitRectangleStation[i].FillColor = waitCars[i].GetRectangle().FillColor
				}
			}
			p.window.Content().Refresh()
		}
	}
}

func (p *ParkingView) RenderParkCar() {
	for {
		select {
		case <-semQuit:
			return
		default:
			<-semRenderNewCarParking
			parkingArray := parking.GetParking()
			for i := range parkingArray {
				p.carsParking[i].rectangle.FillColor = parkingArray[i].GetRectangle().FillColor
				p.carsParking[i].text.Text = fmt.Sprintf("%d", parkingArray[i].GetTime())
				p.carsParking[i].text.Color = Gray
				p.carsParking[i].text.Show()
			}
			p.window.Content().Refresh()
		}
	}
}
func (p *ParkingView) RenderTimeCarPark() {
	for {
		select {
		case <-semQuit:
			return
		default:
			index := <-senRenderTime
			parkingArray := parking.GetParking()
			p.carsParking[index].text.Text = fmt.Sprintf("%d", parkingArray[index].GetTime())
			p.window.Content().Refresh()
		}
	}
}

func (p *ParkingView) RenderEnterCar() {
	for {
		select {
		case <-semQuit:
			return
		default:
			<-semRenderNewCarEnter
			p.entrace.FillColor = parking.GetEntraceCar().GetRectangle().FillColor
			fmt.Printf("Renderize")
			p.window.Content().Refresh()
		}
	}
}

func (p *ParkingView) RenderExitCar() {
	for {
		select {
		case <-semQuit:
			return
		default:
			id := <-semRenderIDExit

			p.exit.FillColor = p.carsParking[id].rectangle.FillColor
			p.carsParking[id].rectangle.FillColor = Gray
			p.carsParking[id].text.Hide()
			p.carsParking[id].text.Color = color.RGBA{R: 255, G: 255, B: 255, A: 255}
			p.window.Content().Refresh()
			time.Sleep(1 * time.Second)

			p.out.FillColor = p.exit.FillColor
			p.exit.FillColor = Gray
			p.window.Content().Refresh()
			time.Sleep(1 * time.Second)

			p.out.FillColor = Gray
			p.window.Content().Refresh()
			time.Sleep(1 * time.Second)

			p.window.Content().Refresh()
			semRenderNewCarExit <- true
		}
	}
}

func (p *ParkingView) StartSimulation() {
	go parking.GenerateCars()
	go parking.CarEntrace()
	go parking.CheckParking()
	go parking.CarExit()
	go p.RenderNewCarStation()
	go p.RenderParkCar()
	go p.RenderEnterCar()
	go p.RenderExitCar()
	go p.RenderTimeCarPark()
}

func (p *ParkingView) BackToMenu() {
	close(semQuit)
	NewMainView(p.window)
}

func (p *ParkingView) RestartSimulation() {
	close(semQuit)
	NewParkingView(p.window)
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
