package domain

import (
	"errors"
	"fmt"

	"github.com/johngillott/gomarsrover/pkg/domain/control"
	"github.com/johngillott/gomarsrover/pkg/domain/direction"
)

var (
	ErrUnableToMove                 = errors.New("unable to move")
	ErrUnableToMoveUndefinedHeading = errors.New("undefined directional move behaviour")
	ErrUndefinedControl             = errors.New("undefined control behaviour")
	ErrUndefinedRotationBehaviour   = errors.New("undefined rotation behaviour")
)

var _ Vehicle = (*Rover)(nil)

type Rover struct {
	x, y        uint64
	heading     direction.Direction
	landingZone LandingZone
}

// NewRover returns a Rover with x, y as start coordinates of the robot and d as the heading of the robot.
func NewRover(x, y uint64, h direction.Direction, l LandingZone) *Rover {
	return &Rover{
		x:           x,
		y:           y,
		heading:     h,
		landingZone: l,
	}
}

func (r Rover) X() uint64 {
	return r.x
}

func (r Rover) Y() uint64 {
	return r.y
}

func (r Rover) Heading() direction.Direction {
	return r.heading
}

const (
	MOVE_SIZE = uint64(1)
	MIN_X     = uint64(0)
	MIN_Y     = uint64(0)
)

func (r *Rover) NextMove() (uint64, uint64, error) {

	switch r.heading {
	case direction.North:
		i := r.Y() + MOVE_SIZE
		if i > r.landingZone.Height() {
			return 0, 0, ErrUnableToMove
		}

		return r.X(), i, nil
	case direction.East:
		i := r.X() + MOVE_SIZE
		if i > r.landingZone.Width() {
			return 0, 0, ErrUnableToMove
		}
		return i, r.Y(), nil
	case direction.South:
		if r.Y() == 0 {
			return 0, 0, ErrUnableToMove
		}
		return r.X(), r.Y() - MOVE_SIZE, nil
	case direction.West:
		if r.X() == 0 {
			return 0, 0, ErrUnableToMove
		}
		return r.X() - MOVE_SIZE, r.Y(), nil
	default:
		return 0, 0, ErrUnableToMoveUndefinedHeading
	}
}

func (r *Rover) Move() error {
	x, y, err := r.NextMove()
	if err != nil {
		return fmt.Errorf("unable to move rover - %w", err)
	}

	r.x, r.y = x, y
	return nil
}

func (r *Rover) Rotate(c control.Control) error {
	switch c {
	case control.Left:
		if err := r.rotateLeft(); err != nil {
			return fmt.Errorf("unable to rotate rover left with heading=%b - %w", r.Heading(), err)
		}
	case control.Right:
		if err := r.rotateRight(); err != nil {
			return fmt.Errorf("unable to rotate rover right with heading %b - %w", r.Heading(), err)
		}
	case control.Move:
	default:
		return fmt.Errorf("unable to rotate with undefined control %c - %w", c, ErrUndefinedControl)
	}

	return nil
}

func (r *Rover) rotateLeft() error {
	switch r.heading {
	case direction.North:
		r.heading = direction.West
	case direction.East:
		r.heading = direction.North
	case direction.South:
		r.heading = direction.East
	case direction.West:
		r.heading = direction.South
	default:
		return ErrUndefinedRotationBehaviour
	}

	return nil
}

func (r *Rover) rotateRight() error {
	switch r.heading {
	case direction.North:
		r.heading = direction.East
	case direction.East:
		r.heading = direction.South
	case direction.South:
		r.heading = direction.West
	case direction.West:
		r.heading = direction.North
	default:
		return ErrUndefinedRotationBehaviour
	}

	return nil
}
