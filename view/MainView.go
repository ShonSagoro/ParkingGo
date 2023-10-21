package view

import (
	"image/color"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type MainView struct {
	window fyne.Window
}

func NewMainView(window fyne.Window) *MainView {
	MainView := &MainView{
		window: window,
	}
	MainView.InitApp()
	return MainView
}

func (m *MainView) InitApp() {
	m.DrawSceneMenu()
}

func (m *MainView) DrawSceneMenu() {
	title := canvas.NewText("Parking simulator", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	title.Resize(fyne.NewSize(20, 20))
	titleContainer := container.NewCenter(title)

	start := widget.NewButton("Start Simulation", m.StartParkingSimulation)

	credits := widget.NewButton("Credits", m.DrawCredits)

	exit := widget.NewButton("Exit", m.ExitGame)

	container_center := container.NewVBox(
		titleContainer,
		layout.NewSpacer(),
		start,
		credits,
		exit,
		layout.NewSpacer(),
	)

	m.window.SetContent(container_center)
	m.window.Resize(fyne.NewSize(400, 500))
	m.window.SetFixedSize(true)
}

func (m *MainView) ExitGame() {
	m.window.Close()
}

func (m *MainView) DrawCredits() {
	shonPhoto := canvas.NewImageFromFile("resources/shon_credits.png")
	shonPhoto.Resize(fyne.NewSize(300, 250))
	shonPhoto.Move(fyne.NewPos(0, 10))
	shonPhotoContainer := container.NewWithoutLayout(shonPhoto)
	goBackMenu := widget.NewButton("Menu", m.DrawSceneMenu)

	shonCard := widget.NewCard("Jonathan G. Shon sagoro", "Programmer", shonPhotoContainer)

	githubURL := "https://github.com/ShonSagoro"

	parsedURL, err := url.Parse(githubURL)
	if err != nil {
		panic(err)
	}

	github_widget := widget.NewHyperlink("Checa mi Github", parsedURL)

	container_center := container.NewHBox(
		container.NewCenter(github_widget),
		container.NewVBox(
			shonCard,
			layout.NewSpacer(),
			goBackMenu,
		),
	)

	m.window.Resize(fyne.NewSize(300, 400))
	m.window.SetContent(container_center)
}

func (m *MainView) StartParkingSimulation() {
	NewParkingView(m.window)
}
