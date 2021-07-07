package service

import (
	"errors"
	"fmt"

	"github.com/johngillott/gomarsrover/pkg/domain"
	"github.com/johngillott/gomarsrover/pkg/domain/control"
	"github.com/johngillott/gomarsrover/pkg/domain/direction"
)

var (
	// ErrOutOfBounds is returned when a vehicle is deployed or moved outside of the landing zone.
	ErrOutOfBounds = errors.New("position out of bounds")
	// ErrPositionBusy is returned when a vehicle is deployed or moved to a position occupied by abother vehicle.
	ErrPositionBusy = errors.New("vehicle currently deployed to position")
)

type MarsRoverService struct {
	landingZone domain.LandingZone
	fleet       []domain.Vehicle
}

func NewMarsRoverService(l domain.LandingZone) *MarsRoverService {
	return &MarsRoverService{
		landingZone: l,
		fleet:       make([]domain.Vehicle, 0, 1),
	}
}

func (m *MarsRoverService) Deploy(x, y uint64, h direction.Direction) error {

	if m.landingZone == nil {
		return errors.New("marsroverservice created without landing zone")
	}

	if !m.isPositionInLandingZone(x, y) {
		return fmt.Errorf("unable to deploy rover to position (%d,%d). Outside of landing zone height=%d width=%d - %w",
			x, y,
			m.landingZone.Height(),
			m.landingZone.Width(),
			ErrOutOfBounds)
	}

	if m.willCollideWithVehicle(x, y) {
		return fmt.Errorf("unable to deploy vehicle to position (%d,%d) - %w", x, y, ErrPositionBusy)
	}

	m.fleet = append(m.fleet, domain.NewRover(x, y, h, m.landingZone))

	return nil
}

func (m *MarsRoverService) Move(c control.Control, v domain.Vehicle) error {

	if v == nil {
		return errors.New("provide a vehicle value")
	}

	if c != control.Move {
		if err := v.Rotate(c); err != nil {
			return fmt.Errorf("service cannot process control state control=%c heading=%c - %w", c, v.Heading(), err)
		}
		return nil
	}

	nextX, nextY, err := v.NextMove()
	if err != nil {
		return fmt.Errorf("vehicle position=%d %d %c - %w", v.X(), v.Y(), v.Heading(), ErrOutOfBounds)
	}

	if m.willCollideWithVehicle(nextX, nextY) {
		return fmt.Errorf("unable to move vehicle to position (%d,%d) - %w", nextX, nextY, ErrPositionBusy)
	}

	if err := v.Move(); err != nil {
		return fmt.Errorf("unable to move vehicle to position (%d,%d) - %w", nextX, nextY, err)
	}

	return nil
}

func (m *MarsRoverService) ActiveVehicle() domain.Vehicle {

	if len(m.fleet) == 0 {
		return nil
	}

	return m.fleet[len(m.fleet)-1]
}

func (m MarsRoverService) isPositionInLandingZone(x, y uint64) bool {
	return x <= m.landingZone.Height() && y <= m.landingZone.Width()
}

func (m MarsRoverService) willCollideWithVehicle(x, y uint64) bool {

	for _, u := range m.fleet {
		if u.X() == x && u.Y() == y {
			return true
		}
	}

	return false
}
