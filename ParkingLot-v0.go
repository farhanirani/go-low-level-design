package main

import (
	"errors"
	"fmt"
	"log"
)

type VehicleType int

const (
	CAR VehicleType = iota
	BIKE
)

type Vehicle struct {
	Number string
	Type   VehicleType
}

type ParkingSpot struct {
	ID          int
	IsFree      bool
	VehicleType VehicleType
	Vehicle     *Vehicle
}

type ParkingLot struct {
	Spots     map[int]*ParkingSpot
	FreeCount map[VehicleType]int
	FreeSpots map[VehicleType]map[int]struct{}
}

func NewSpot(id int, vt VehicleType) *ParkingSpot {
	return &ParkingSpot{ID: id, IsFree: true, Vehicle: nil, VehicleType: vt}
}

func NewParkingLot(carSpots, bikeSpots int) *ParkingLot {
	lot := &ParkingLot{
		Spots:     make(map[int]*ParkingSpot),
		FreeCount: make(map[VehicleType]int),
		FreeSpots: make(map[VehicleType]map[int]struct{}),
	}

	lot.FreeSpots[CAR] = make(map[int]struct{})
	lot.FreeSpots[BIKE] = make(map[int]struct{})

	for i := 1; i <= carSpots; i++ {
		cs := NewSpot(i, CAR)
		lot.Spots[i] = cs
		lot.FreeSpots[CAR][i] = struct{}{}
		lot.FreeCount[CAR]++
	}

	for i := carSpots + 1; i <= carSpots+bikeSpots; i++ {
		cs := NewSpot(i, BIKE)
		lot.Spots[i] = cs
		lot.FreeSpots[BIKE][i] = struct{}{}
		lot.FreeCount[BIKE]++
	}
	return lot
}

func (p *ParkingLot) ParkVehicle(v *Vehicle) (int, error) {
	for spotId := range p.FreeSpots[v.Type] {
		spot := p.Spots[spotId]
		spot.IsFree = false
		spot.Vehicle = v

		p.FreeCount[v.Type]--
		delete(p.FreeSpots[v.Type], spotId)

		return spotId, nil
	}

	return -1, errors.New("no available spots")
}

func (p *ParkingLot) Status() {
	fmt.Println("Parking Lot Status:")
	fmt.Printf("CAR SPOTS %d: BIKE: %d\n", p.FreeCount[CAR], p.FreeCount[BIKE])
}

// --- Demo ---
func main() {
	lot := NewParkingLot(2, 2)

	car := &Vehicle{Number: "MH12AB1234", Type: CAR}
	bike := &Vehicle{Number: "MH14XY5678", Type: BIKE}

	lot.Status()
	t1, _ := lot.ParkVehicle(car)
	log.Printf("Car Parked in %d", t1)
	t2, _ := lot.ParkVehicle(bike)
	log.Printf("Bike Parked in %d", t2)
	t2, _ = lot.ParkVehicle(bike)
	log.Printf("Bike Parked in %d", t2)
	lot.Status()
	//
	//lot.UnparkVehicle(t1)
	//fmt.Println("\nAfter unparking car:")
	//lot.Status()
	//
	//_ = lot.UnparkVehicle(t2)
	//fmt.Println("\nAfter unparking bike:")
	//lot.Status()
}
