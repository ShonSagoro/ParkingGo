package models

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"time"
)

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

const (
	lambda = 2.0
)

const maxWait int = 10

type Parking struct {
	waitCars               []*Car
	parking                [20]*Car
	entrace                *Car
	exit                   *Car
	semRenderNewCarWait    chan bool
	semRenderNewCarParking chan bool
	semRenderNewCarEnter   chan bool
	semRenderNewCarExit    chan bool
	semQuit                chan bool
}

// var semFullWaitStation chan bool
var semHaveSpace chan int
var semWait chan bool

func NewParking(sRNCW chan bool, sRNCP chan bool, sRNCEn chan bool, sRNCEx chan bool, sQ chan bool) *Parking {
	// semFullWaitStation = make(chan bool)
	semHaveSpace = make(chan int)
	semWait = make(chan bool)
	parking := &Parking{
		entrace:                NewSpaceCar(),
		exit:                   NewSpaceCar(),
		semRenderNewCarWait:    sRNCW,
		semRenderNewCarParking: sRNCP,
		semRenderNewCarEnter:   sRNCEn,
		semRenderNewCarExit:    sRNCEx,
		semQuit:                sQ,
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
		select {
		case <-p.semQuit:
			fmt.Printf("CheckParking Close")
			return
		default:
			index := p.SearchSpace()
			if index != -1 {
				semHaveSpace <- index
				<-semWait
				p.semRenderNewCarParking <- true
			}
			time.Sleep(time.Duration(1 * time.Second))
		}
	}
}

func (p *Parking) CarEntrace() {
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("CarEntrace Close")
			return
		default:
			index := <-semHaveSpace
			p.PassToEntraceState()
			if index != -1 {
				p.parking[index] = p.entrace
				go p.parking[index].StartCount()
				semWait <- true
			}
		}
	}
}

func (p *Parking) GenerateCars() {
	i := 0
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("GenerateCars Close")
			return
		default:
			fmt.Printf("Nuevo")
			interarrivalTime := -math.Log(1-rand.Float64()) / lambda
			time.Sleep(time.Duration(interarrivalTime * float64(time.Second)))
			if len(p.waitCars) < maxWait {
				car := NewCar(int64(i))
				i++
				p.waitCars = append(p.waitCars, car)
				fmt.Printf("Renderizalo")
				p.semRenderNewCarWait <- true
			} else {
				i = 0
			}
		}
	}
}

func (p *Parking) PassToEntraceState() {
	if len(p.waitCars) != 0 {
		p.entrace = p.waitCars[0]
		p.waitCars = p.waitCars[1:]
		p.semRenderNewCarEnter <- true
		time.Sleep(3 * time.Second)
	}
}

func (p *Parking) SearchSpace() int {
	for s := range p.parking {
		if p.parking[s].GetID() == -1 {
			fmt.Printf("\n %d\n", s)
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
