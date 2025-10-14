// LLD Walkthrough [Golang]: Parking Lot

// ----------------------------------------------------------------------
// Design a Parking Complex which has multiple levels, and on each level there are parking spots for different vehicles - for CAR and BIKE

// REQUIREMENTS
// 1. Create a parking lot, add levels with different spots types to it.
// 2. given a vehicle, park it, and unpark it.
// 3. display the number of available spots on each level.

// Q: Can you find a valid parking spot in O(1) time? Instead of looping through all the spots on all the levels?

// Note: If the requirements are not explicitly specified, please clarify them before proceeding
// ----------------------------------------------------------------------

// ----------------------------------------------------------------------
// IMPLEMENTATION STEPS

// STEP 1 - DISCUSS ENTITIES AND CORE METHODS TO SATISFY THE REQUIREMENTS

// ParkingLot → has multiple Levels.
// Level → has multiple Spots.
// Spot → knows its type (Car, Bike) and occupancy.
// Vehicle → type, registration.

// STEP 2 - FACTORY & CORE METHODS / APIs

// Park(vehicle) → returns spot info / ticket
// Unpark(vehicle or ticket) → frees spot
// GetAvailableSpots(level) → returns count

// STEP 3 - MAIN FUNC DEMO CODE
// ----------------------------------------------------------------------

// ----------------------------------------------------------------------
package main

import (
	"errors"
	"fmt"
)

// ----------------------------
// ENTITIES
// ----------------------------
type VehicleType int

const (
	CAR VehicleType = iota
	BIKE
)

// Vehicle represents a vehicle
type Vehicle struct {
	Type VehicleType
	Reg  string
}

// Spot represents a parking spot
type Spot struct {
	ID          int
	VehicleType VehicleType
	IsOccupied  bool
	Vehicle     *Vehicle
}

// Level represents one level in the parking lot
type Level struct {
	ID        int
	Spots     map[int]*Spot                    // spotID -> Spot
	FreeSpots map[VehicleType]map[int]struct{} // vehicleType -> set of free spot IDs
	FreeCount map[VehicleType]int              // vehicleType -> free spot count
}

// ParkingLot represents the whole parking complex
type ParkingLot struct {
	Levels []*Level
}

// ----------------------------
// FACTORY FUNCTIONS
// ✅ Factory Method Pattern: encapsulates creation logic for complex structs.
// ----------------------------
func NewSpot(id int, vType VehicleType) *Spot {
	return &Spot{ID: id, VehicleType: vType}
}

func NewLevel(id int, carSpots, bikeSpots int) *Level {
	l := &Level{
		ID:        id,
		Spots:     make(map[int]*Spot),
		FreeSpots: make(map[VehicleType]map[int]struct{}),
		FreeCount: make(map[VehicleType]int),
	}

	l.FreeSpots[CAR] = make(map[int]struct{})
	l.FreeSpots[BIKE] = make(map[int]struct{})

	// Initialize car spots
	for i := 1; i <= carSpots; i++ {
		spot := NewSpot(i, CAR)
		l.Spots[i] = spot
		l.FreeSpots[CAR][i] = struct{}{}
		l.FreeCount[CAR]++
	}

	// Initialize bike spots
	for i := carSpots + 1; i <= carSpots+bikeSpots; i++ {
		spot := NewSpot(i, BIKE)
		l.Spots[i] = spot
		l.FreeSpots[BIKE][i] = struct{}{}
		l.FreeCount[BIKE]++
	}

	return l
}

func NewParkingLot() *ParkingLot {
	return &ParkingLot{}
}

// ----------------------------
// CORE METHODS
// ----------------------------
// ✅ Open/Closed principle — new levels can be added without modifying core logic.
func (p *ParkingLot) AddLevel(level *Level) {
	p.Levels = append(p.Levels, level)
}

// Park parks a vehicle in O(1) time
func (p *ParkingLot) Park(v *Vehicle) (int, int, error) {
	for _, level := range p.Levels {
		freeSet := level.FreeSpots[v.Type]
		if len(freeSet) == 0 {
			continue // no free spot of this type on this level
		}

		// Pick any one free spot (O(1))
		for spotID := range freeSet {
			spot := level.Spots[spotID]
			spot.IsOccupied = true
			spot.Vehicle = v

			// Update free sets and counts
			delete(level.FreeSpots[v.Type], spotID)
			level.FreeCount[v.Type]--

			return level.ID, spot.ID, nil
		}
	}
	return -1, -1, errors.New("no available spot for vehicle type")
}

// Unpark frees a spot in O(1) time
func (p *ParkingLot) Unpark(levelID, spotID int) error {
	for _, level := range p.Levels {
		if level.ID == levelID {
			spot, exists := level.Spots[spotID]
			if !exists {
				return errors.New("invalid spot ID")
			}
			if !spot.IsOccupied {
				return errors.New("spot already empty")
			}

			vType := spot.VehicleType
			spot.IsOccupied = false
			spot.Vehicle = nil

			// Add back to free set and increment counter
			level.FreeSpots[vType][spotID] = struct{}{}
			level.FreeCount[vType]++

			return nil
		}
	}
	return errors.New("invalid level ID")
}

// DisplayAvailability prints available spots per level
func (p *ParkingLot) DisplayAvailability() {
	fmt.Println("----- Parking Availability -----")
	for _, level := range p.Levels {
		fmt.Printf("Level %d: ", level.ID)
		fmt.Printf("CAR=%d, BIKE=%d\n",
			level.FreeCount[CAR],
			level.FreeCount[BIKE])
	}
	fmt.Println("--------------------------------")
}

// ----------------------------
// DEMO
// ----------------------------
func main() {
	// Create parking lot
	lot := NewParkingLot()

	// Add levels
	lot.AddLevel(NewLevel(1, 3, 3)) // Level 1: 3 CAR, 3 BIKE spots
	lot.AddLevel(NewLevel(2, 1, 3)) // Level 2: 1 CAR, 3 BIKE spots

	lot.DisplayAvailability()

	// Park vehicles
	car := &Vehicle{Type: CAR, Reg: "CAR123"}
	bike := &Vehicle{Type: BIKE, Reg: "BIKE456"}

	vehicles := []*Vehicle{car, car, car, car, bike} // Attempt multiple parks

	for _, v := range vehicles {
		levelID, spotID, err := lot.Park(v)
		if err != nil {
			fmt.Printf("Failed to park %s: %v\n", v.Reg, err)
		} else {
			fmt.Printf("Parked %s at Level %d, Spot %d\n", v.Reg, levelID, spotID)
		}
	}

	lot.DisplayAvailability()

	// Unpark a vehicle
	err := lot.Unpark(1, 1)
	if err != nil {
		fmt.Printf("Failed to unpark vehicle from Level 1, Spot 1: %v\n", err)
	} else {
		fmt.Println("Unparked vehicle from Level 1, Spot 1")
	}

	lot.DisplayAvailability()

	// Attempt to unpark from an invalid spot
	err = lot.Unpark(1, 10)
	if err != nil {
		fmt.Printf("Failed to unpark vehicle from Level 1, Spot 10: %v\n", err)
	}
}

// ----- POSSIBLE FOLLOW UP QUESTIONS -----

/*
Question / Follow-up	Concept / Pattern	High-Level Answer

How would you handle concurrency (multiple vehicles parking/unparking)?
Concurrency control / sync primitives
Use sync.Mutex or sync.RWMutex to lock each level or spot during updates. Alternatively, use message queues or channels to serialize operations.


Can we notify users when a spot frees up?
Observer / Pub-Sub Pattern
Implement observer pattern: levels publish “spot available” events, subscribers (apps/users) listen.


How to charge customers dynamically?
Strategy Pattern (Payment Strategy)
Introduce PaymentStrategy interface with implementations for hourly, subscription, or flat-rate payments.


How to integrate payments (Stripe, PayTM)?
Adapter / Strategy
Add PaymentProcessor interface; inject concrete implementations (StripeAdapter, RazorpayAdapter).
*/

// *************************************************
// REMEMBER : CODE THAT GOOD ENOUGH AND ALL REQUIREMENTS SATISFIED >>>>>> CODE THAT IS PERFECT BUT INCOMPLETE ( for the interview )
// ALL THE BEST
// *************************************************
