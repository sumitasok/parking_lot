package models

import (
	"fmt"
	"strconv"
	"strings"
)

// StoreyResponse response object that encapsulates the data representation.
type StoreyResponse struct {
	slots   []Slot
	command string
}

// String returns string representation of the data for logging.
func (s StoreyResponse) String() string {
	switch s.command {
	case CmdPark:
		return fmt.Sprintf("Allocated slot number: %d", s.slots[0].Position())
	case CmdCreateParkingLot:
	case CmdStatus:
	case CmdLeave:
		return fmt.Sprintf("Slot number %d is free", s.slots[0].Position())
	case CmdRegistrationNumberByColor:
		regNumbers := []string{}
		for _, s := range s.slots {
			regNumbers = append(regNumbers, s.car.numberPlate)
		}

		return strings.Join(regNumbers, ", ")
	case CmdSlotnoByCarColor:
		positions := []string{}
		for _, s := range s.slots {
			positions = append(positions, strconv.Itoa(s.Position()))
		}

		return strings.Join(positions, ", ")
	case CmdSlotnoByRegNumber:
		return strconv.Itoa(s.slots[0].Position())
	}
	return ""
}
