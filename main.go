package main

import (
	"parking/views"

	"fyne.io/fyne/v2/app"
)

func main() {
	app := app.New()
	window := app.NewWindow("Parking")
	window.CenterOnScreen()
	views.NewMainView(window)
	window.ShowAndRun()
}
