package main

import (
	"parking/view"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Parking")
	window.CenterOnScreen()
	view.NewMainView(window)
	window.ShowAndRun()
}
