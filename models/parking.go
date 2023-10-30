package models

import (
	"fmt"
	"image/color"
	"math"
	"math/rand"
	"sync"
	"time"
)

var Gray = color.RGBA{R: 30, G: 30, B: 30, A: 255}

const (
	lambda = 2.0
)

const MaxWait int = 10
const MaxParking int = 20

type Parking struct {
	waitCars               []*Car
	parking                [MaxParking]*Car
	entrace                *Car
	exit                   *Car
	semRenderNewCarWait    chan bool
	semRenderNewCarParking chan bool
	semRenderNewCarEnter   chan bool
	semRenderNewCarExit    chan bool
	semRenderIDExit        chan int
	senRenderTime          chan int
	semQuit                chan bool
}

var mutexExitEnter sync.Mutex

func NewParking(
	sRNCW chan bool,
	sRNCP chan bool,
	sRNCEn chan bool,
	sRNCEx chan bool,
	sICEx chan int,
	sRT chan int,
	sQ chan bool,
) *Parking {
	parking := &Parking{
		entrace:                NewSpaceCar(),
		exit:                   NewSpaceCar(),
		semRenderNewCarWait:    sRNCW,
		semRenderNewCarParking: sRNCP,
		semRenderNewCarEnter:   sRNCEn,
		semRenderNewCarExit:    sRNCEx,
		semRenderIDExit:        sICEx,
		senRenderTime:          sRT,
		semQuit:                sQ,
	}
	parking.ClearParking()
	return parking
}

func (p *Parking) MakeParking() {
	for i := range p.parking {
		p.parking[i] = NewSpaceCar()
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
				fmt.Printf("ENTRE")
				mutexExitEnter.Lock()
				p.PassToEntraceState()
				p.ParkingCar(index)
				p.semRenderNewCarEnter <- true
				p.semRenderNewCarParking <- true
				mutexExitEnter.Unlock()
			}
			time.Sleep(1 * time.Second)

		}
	}
}

func (p *Parking) PassToEntraceState() {
	if len(p.waitCars) != 0 {
		p.entrace = p.waitCars[0]
		p.waitCars = p.waitCars[1:]
		p.semRenderNewCarEnter <- true
		time.Sleep(1 * time.Second)
	}
}

func (p *Parking) ParkingCar(index int) {
	p.parking[index] = p.entrace
	p.entrace = NewSpaceCar()
	go p.parking[index].StartCount(index)
}

func (p *Parking) OutCarToExit() {
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("CarExit Close")
			return
		default:
			if !WaitExitCarsIsEmpty() {
				mutexExitEnter.Lock()
				car := PopWaitCars()
				p.semRenderIDExit <- car.GetID()
				p.parking[car.ID] = NewSpaceCar()
				<-p.semRenderNewCarExit
				mutexExitEnter.Unlock()
			}
		}
	}
}

func (p *Parking) GenerateCars() {
	i := 20
	for {
		select {
		case <-p.semQuit:
			fmt.Printf("GenerateCars Close")
			return
		default:
			interarrivalTime := -math.Log(1-rand.Float64()) / lambda
			time.Sleep(time.Duration(interarrivalTime * float64(time.Second)))
			if len(p.waitCars) < MaxWait {
				car := NewCar(i, p.senRenderTime, p.semQuit)
				i++
				p.waitCars = append(p.waitCars, car)
				p.semRenderNewCarWait <- true
			}
		}
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

func (p *Parking) GetParking() [MaxParking]*Car {
	return p.parking
}

func (p *Parking) ClearParking() {
	for i := range p.parking {
		p.parking[i] = nil
	}
}
