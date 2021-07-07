package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/johngillott/gomarsrover/pkg/domain"
	"github.com/johngillott/gomarsrover/pkg/domain/control"
	"github.com/johngillott/gomarsrover/pkg/domain/direction"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"
)

func TestMarsRoverService_Deploy(t *testing.T) {

	l := domain.NewPlateau(5, 5)

	want := domain.NewRover

	type args struct {
		x       uint64
		y       uint64
		heading direction.Direction

		m *MarsRoverService
	}

	tests := []struct {
		name    string
		args    args
		want    domain.Vehicle
		wantErr error
	}{
		{
			"Able to deploy inside Landing Zone",
			args{
				x:       0,
				y:       0,
				heading: direction.North,
				m:       NewMarsRoverService(l),
			},
			want(0, 0, direction.North, l),
			nil,
		},
		{
			"Unable deploy outside of Landing Zone",
			args{
				x:       6,
				y:       6,
				heading: direction.North,
				m:       NewMarsRoverService(l),
			},
			nil,
			ErrOutOfBounds,
		},
		{
			"Collision detection on deployment",
			args{
				x:       3,
				y:       3,
				heading: direction.North,
				m: func(t *testing.T) *MarsRoverService {

					m := NewMarsRoverService(l)

					const x, y = 3, 3

					if err := m.Deploy(x, y, direction.South); err != nil {
						t.Fatalf("Unable to initialize MarsRoverService with vehicle deployed to %d %d - %s", x, y, err)
					}

					return m

				}(t),
			},
			want(3, 3, direction.South, l),
			ErrPositionBusy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := qt.New(t)
			err := tt.args.m.Deploy(tt.args.x, tt.args.y, tt.args.heading)

			c.Assert(errors.Is(err, tt.wantErr), qt.IsTrue)
			c.Assert(tt.want, qt.CmpEquals(cmp.AllowUnexported(domain.Rover{}, domain.Plateau{})), tt.args.m.ActiveVehicle())
		})
	}
}

func TestMarsRoverService_Move(t *testing.T) {

	const MAX_X, MAX_Y = 5, 5

	l := domain.NewPlateau(MAX_X, MAX_Y)

	want := domain.NewRover

	type args struct {
		c control.Control
		m *MarsRoverService
	}

	tests := []struct {
		name    string
		args    args
		want    domain.Vehicle
		wantErr error
	}{
		{
			"Rover can move",
			args{
				c: control.Move,
				m: func(t *testing.T) *MarsRoverService {

					m := NewMarsRoverService(l)

					const x, y = 3, 3

					if err := m.Deploy(x, y, direction.North); err != nil {
						t.Fatalf("Unable to initialize MarsRoverService with vehicle deployed to %d %d - %s", x, y, err)
					}
					return m

				}(t),
			},
			want(3, 4, direction.North, l),
			nil,
		},
		{
			"Rover can turn",
			args{
				c: control.Left,
				m: func(t *testing.T) *MarsRoverService {

					m := NewMarsRoverService(l)

					const x, y = 3, 3

					if err := m.Deploy(x, y, direction.North); err != nil {
						t.Fatalf("Unable to initialize MarsRoverService with vehicle deployed to %d %d - %s", x, y, err)
					}
					return m

				}(t),
			},
			want(3, 3, direction.West, l),
			nil,
		},
		{
			"Rover cannot move outside landing zone",
			args{
				c: control.Move,
				m: func(t *testing.T) *MarsRoverService {

					m := NewMarsRoverService(l)

					const x, y = MAX_X, MAX_Y

					if err := m.Deploy(x, y, direction.North); err != nil {
						t.Fatalf("Unable to initialize MarsRoverService with vehicle deployed to %d %d - %s", x, y, err)
					}
					return m

				}(t),
			},
			want(5, 5, direction.North, l),
			ErrOutOfBounds,
		},
		{
			"Collision detection on move",
			args{
				c: control.Move,
				m: func(t *testing.T) *MarsRoverService {

					m := NewMarsRoverService(l)

					const x, y = 3, 3

					if err := m.Deploy(x, y, direction.North); err != nil {
						t.Fatalf("Unable to initialize MarsRoverService with vehicle deployed to %d %d - %s", x, y, err)
					}

					const u, v = x, y + 1

					if err := m.Deploy(u, v, direction.South); err != nil {
						t.Fatalf("Unable to initialize MarsRoverService with vehicle deployed to %d %d - %s", u, v, err)
					}

					return m

				}(t),
			},
			want(3, 4, direction.South, l),
			ErrPositionBusy,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			c := qt.New(t)
			v := tt.args.m.ActiveVehicle()
			err := tt.args.m.Move(tt.args.c, v)
			fmt.Printf("err = %s\n", err)

			c.Assert(errors.Is(err, tt.wantErr), qt.IsTrue)
			c.Assert(tt.want, qt.CmpEquals(cmp.AllowUnexported(domain.Rover{}, domain.Plateau{})), tt.args.m.ActiveVehicle())
		})
	}
}
