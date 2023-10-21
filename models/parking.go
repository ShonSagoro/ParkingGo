package models

import "image/color"

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

type Parking struct {
	capacity [4][5]*Car
}
