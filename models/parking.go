package models

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"
)

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

type Parking struct {
	waitCars           []*Car
	parking            [20]*Car
	entrace            *Car
	exit               *Car
	semNewCar          chan bool
	semFullWaitStation chan bool
	semHaveSpace       chan bool
	semWait            chan bool
}

func NewParking(sN chan bool, sF chan bool, sH chan bool, sW chan bool) *Parking {
	parking := &Parking{
		entrace:            NewSpaceCar(),
		exit:               NewSpaceCar(),
		semNewCar:          sN,
		semFullWaitStation: sF,
		semHaveSpace:       sH,
		semWait:            sW,
	}
	return parking
}

func (p *Parking) MakeParking() {
	for i := range p.parking {
		carSpace := NewSpaceCar()
		p.parking[i] = carSpace
	}
}

func (p *Parking) CheckParking() {
	for {
		if p.SearchSpace() != -1 {
			p.semHaveSpace <- true
			<-p.semWait
		}
		time.Sleep(time.Duration(1 * time.Second))
	}
}

func (p *Parking) CarEntrace() {
	for {
		<-p.semHaveSpace

		p.PassToEntraceState()

		index := p.SearchSpace()

		if index != -1 {
			p.parking[index] = p.entrace
			if len(p.waitCars) != 20 {
				p.semFullWaitStation <- true
			}
			p.semWait <- true
		}
	}
}

func (p *Parking) GenerateCars() {
	i := 0
	for {
		fmt.Printf("Nuevo")
		interarrivalTime := -math.Log(1-rand.Float64()) / lambda
		time.Sleep(time.Duration(interarrivalTime * float64(time.Second)))

		car := NewCar(int64(i))
		i++
		if len(p.waitCars) < 20 {
			p.waitCars = append(p.waitCars, car)
			p.semNewCar <- true
			fmt.Printf("Renderizalo")
		} else {
			<-p.semFullWaitStation
		}
	}
}

func (p *Parking) PassToEntraceState() {
	if len(p.waitCars) != 0 {
		p.entrace = p.waitCars[0]
		p.waitCars = p.waitCars[1:]
	}
}

func (p *Parking) SearchSpace() int {
	for s := range p.parking {
		if p.parking[s].GetID() == -1 {
			return s
		}
	}
	return -1
}

func (p *Parking) GetWaitCars() []*Car {
	return p.waitCars
}
func (p *Parking) GetEntraceCar() *Car {
	return p.entrace
}

func (p *Parking) GetExitCar() *Car {
	return p.exit
}

func (p *Parking) GetParking() [20]*Car {
	return p.parking
}

// Funcion de entrada ->
//	La encargada de sacar un carro de la cola <- otra funcion
//  dejar que busque su espacio vacio <- otra funcion
//	colocarse en ese espacio vacio
//  Si hay mas autos esperando para entrar, esperar a que estos entren uno por uno.
//  Avisar que ya no van a entrar porque el estacionamiento se lleno carros para que salgan los carros
//	faltantes.

// Funcion de salida ->
// 	Sacar el carro
//  Esperar a que salga (1 seg)
//  Si hay mas autos esperando para salir, sacar uno por uno los carros.
//  Avisar que ya no van a salir carros para que entren
