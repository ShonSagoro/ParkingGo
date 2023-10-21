package view

import "fyne.io/fyne/v2"

type ParkingView struct {
	window fyne.Window
}

func NewParkingView(window fyne.Window) *ParkingView {
	parkingView := &ParkingView{window: window}
	return parkingView
}
