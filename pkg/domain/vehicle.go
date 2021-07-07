package domain

import (
	"github.com/johngillott/gomarsrover/pkg/domain/control"
	"github.com/johngillott/gomarsrover/pkg/domain/direction"
)

type Vehicle interface {
	X() uint64
	Y() uint64
	Heading() direction.Direction
	Move() error
	NextMove() (uint64, uint64, error)
	Rotate(c control.Control) error
}
