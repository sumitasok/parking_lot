package models

import (
	"errors"
)

var (
	// ErrMaxSlotReached - error when maximum number of slots allowed in a storey is reached.
	ErrMaxSlotReached = errors.New("Sorry, parking lot is full")
	// ErrNoCarsParked - error no cars parked.
	ErrNoCarsParked = errors.New("No cars parked")
	// ErrCarNotFound - error no cars parked with this registration number.
	ErrCarNotFound = errors.New("Not found")
	// ErrCarWithColorNotFound - error no cars parked with this registration number.
	ErrCarWithColorNotFound = errors.New("Car with specified color not found")
)

// Storey - holds the slots and the details about the storey. (Multi storey
// can have many Storey objects)
type Storey struct {
	maxSlots int
	// Linked list should keep it memory efficient.
	slotList *Slot
	db       *storeyDB
}

// Park - check if the Slot is available
// if available Create Slot in the vacancy and associate with adjacent slots
// return Slot
func (s *Storey) Park(numberPlate, color string) (*Slot, error) {
	slot := &Slot{}
	if s.OccupancyCount() >= s.maxSlots {
		return slot, ErrMaxSlotReached
	}

	car := NewCar(numberPlate, color)

	if s.OccupancyCount() == 0 {
		slot := NewSlot(car, 1)
		s.slotList = slot
		return slot, nil
	}

	// Added comments in func (s *Slot) AddNext(sc *Slot) error {
	// about how to improve the code.
	if s.slotList.Position() > 1 {
		currSlot := s.slotList
		s.slotList = NewSlot(car, 1)
		s.slotList.AddNext(currSlot)
		currSlot.prevSlot = s.slotList
	}

	slot = NewSlot(car, 0)
	s.slotList.AddNext(slot)

	return slot, nil
}

// Leave - find the vehicle by numberplate and free the slot
// and associate plevious slot with adjacent slots
// return Slot
func (s *Storey) Leave(numberPlate string) (*Slot, error) {
	if s.slotList == nil {
		return &Slot{}, ErrNoCarsParked
	}

	slotFound, err := s.slotList.FindCar(numberPlate)
	if err != nil {
		return &Slot{}, ErrCarNotFound
	}

	slotFound.Leave()
	if slotFound.prevSlot == nil {
		s.slotList = slotFound.nextSlot
	}

	return slotFound, nil
}

// LeaveByPosition - find the slot by position and free the slot
// and associate plevious slot with adjacent slots
// return Slot
func (s *Storey) LeaveByPosition(position int) (*Slot, error) {
	if s.slotList == nil {
		return &Slot{}, ErrNoCarsParked
	}

	slotFound, err := s.slotList.FindPosition(position)
	if err != nil {
		return &Slot{}, ErrCarNotFound
	}

	slotFound.Leave()
	if slotFound.prevSlot == nil {
		s.slotList = slotFound.nextSlot
	}

	return slotFound, nil
}

// FindByRegistrationNumber - find the slot which has car with the provided
// registration number in the storey.
func (s Storey) FindByRegistrationNumber(numberPlate string) (*Slot, error) {
	if s.slotList == nil {
		return &Slot{}, ErrNoCarsParked
	}

	return s.slotList.FindCar(numberPlate)
}

// FindAllByColor find all cars parked with the color
func (s Storey) FindAllByColor(color string) ([]Slot, error) {
	if s.slotList == nil {
		return []Slot{}, ErrNoCarsParked
	}

	slots, err := s.slotList.FindColor(color)
	if err != nil {
		return []Slot{}, err
	}

	if len(slots) == 0 {
		return []Slot{}, ErrCarWithColorNotFound
	}

	slotsList := []Slot{}
	// mismatch between []Slot and []*Slot as Resposne should not have []*Slot
	for i := range slots {
		slotsList = append(slotsList, *slots[i])
	}
	return slotsList, nil
}

// AllSlots returns all the slots
func (s *Storey) AllSlots() ([]Slot, error) {
	if s.slotList == nil {
		return []Slot{}, ErrNoCarsParked
	}

	return s.slotList.ListSelf(), nil
}

// OccupancyCount returns the number of slots occupied in this storey.
func (s *Storey) OccupancyCount() int {
	if s.slotList == nil {
		return 0
	}

	return s.slotList.CountSelf()
}

// NewStorey returns a Storey object
func NewStorey(maxSlots int) *Storey {
	return &Storey{
		maxSlots: maxSlots,
	}
}

// Slot - each Storey has slots <= maxSlots
type Slot struct {
	prevSlot *Slot
	car      *Car
	position int
	nextSlot *Slot
}

// ListSelf list self and help others list themselves
func (s *Slot) ListSelf() []Slot {
	if s.nextSlot == nil {
		return []Slot{*s}
	}

	return append([]Slot{*s}, s.nextSlot.ListSelf()...)
}

// Leave - leave the Car, and connect the prev slot with next
func (s *Slot) Leave() error {
	if s.prevSlot != nil {
		s.prevSlot.nextSlot = s.nextSlot
	}
	return nil
}

// FindCar - finds if the slot has the car or else check in the next slot
func (s *Slot) FindCar(numberPlate string) (*Slot, error) {
	if s.car.numberPlate == numberPlate {
		return s, nil
	}

	if s.nextSlot == nil {
		return &Slot{}, ErrCarNotFound
	}

	return s.nextSlot.FindCar(numberPlate)
}

// FindPosition - finds if the slot has the car or else check in the next slot
func (s *Slot) FindPosition(position int) (*Slot, error) {
	if s.position == position {
		return s, nil
	}

	if s.nextSlot == nil {
		return &Slot{}, ErrCarNotFound
	}

	return s.nextSlot.FindPosition(position)
}

// FindColor find the cars parked with the color specified, and pass the query to next slot.
func (s *Slot) FindColor(color string) ([]*Slot, error) {
	if s.car.color == color {
		if s.nextSlot == nil {
			return []*Slot{
				s,
			}, nil
		}

		slots, err := s.nextSlot.FindColor(color)
		if err == nil {
			slots = append([]*Slot{s}, slots...)
		}
		return slots, err
	}

	if s.nextSlot == nil {
		return []*Slot{}, nil
	}

	return s.nextSlot.FindColor(color)
}

// AddNext - add a new Slot after the current and associate the current next to the new.
func (s *Slot) AddNext(sc *Slot) error {
	// This requires, Storey to be informed to start looking at the new Slot
	// Or, let each slot point to Storey, so we can traverse and edit easily.
	// if s.prevSlot == nil && s.position > 1 {
	// 	sc.UpdatePosition(1).nextSlot = s
	// }

	if s.nextSlot == nil {
		s.nextSlot = sc.UpdatePosition(s.position + 1)
		sc.prevSlot = s
		return nil
	}

	if s.nextSlot.position > (s.position + 1) {
		currentNext := s.nextSlot
		s.nextSlot = sc.UpdatePosition(s.position + 1)
		sc.prevSlot = s
		sc.nextSlot = currentNext
		currentNext.prevSlot = sc
		return nil
	}

	s.nextSlot.AddNext(sc)

	return nil
}

// CountSelf counts 1 for self and relayes the count Self to next Slot.
func (s Slot) CountSelf() int {
	if s.nextSlot == nil {
		return 1
	}

	return 1 + s.nextSlot.CountSelf()
}

// UpdatePosition updates the position ofthe slot to the specified position value.
func (s *Slot) UpdatePosition(position int) *Slot {
	s.position = position
	return s
}

// Position return the position of the Slot
func (s Slot) Position() int {
	return s.position
}

// RegistrationNumber returns the registration number of the car parked in the slot
func (s Slot) RegistrationNumber() string {
	if s.car == nil {
		return ""
	}

	return s.car.numberPlate
}

// Color returns the color of the car parked in the slot
func (s Slot) Color() string {
	if s.car == nil {
		return ""
	}

	return s.car.color
}

// NewSlot returns a slot object
func NewSlot(car *Car, position int) *Slot {
	return &Slot{car: car, position: position}
}

// Car - define the car preoperties
type Car struct {
	numberPlate string
	color       string
}

// NewCar returns a new car object
func NewCar(numberPlate, color string) *Car {
	return &Car{
		numberPlate: numberPlate,
		color:       color,
	}
}
