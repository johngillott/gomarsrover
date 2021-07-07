package domain_test

import (
	"errors"
	"testing"

	"github.com/johngillott/gomarsrover/pkg/domain"
	"github.com/johngillott/gomarsrover/pkg/domain/control"
	"github.com/johngillott/gomarsrover/pkg/domain/direction"

	qt "github.com/frankban/quicktest"
	"github.com/google/go-cmp/cmp"
)

const (
	UNDEFINED_CONTROL   = control.Control('X')
	UNDEFINED_DIRECTION = direction.Direction('X')
)

func TestRover_Move(t *testing.T) {

	const MAX_X, MAX_Y = 5, 5

	l := domain.NewPlateau(MAX_X, MAX_Y)

	want := domain.NewRover

	tests := []struct {
		name    string
		args    *domain.Rover
		want    *domain.Rover
		wantErr error
	}{
		{
			"North",
			domain.NewRover(2, 2, direction.North, l),
			want(2, 3, direction.North, l),
			nil,
		},
		{
			"East",
			domain.NewRover(2, 2, direction.East, l),
			want(3, 2, direction.East, l),
			nil,
		},
		{
			"South",
			domain.NewRover(2, 2, direction.South, l),
			want(2, 1, direction.South, l),
			nil,
		},
		{
			"West",
			domain.NewRover(2, 2, direction.West, l),
			want(1, 2, direction.West, l),
			nil,
		},
		{
			"Cannot move North",
			domain.NewRover(4, MAX_Y, direction.North, l),
			want(4, MAX_Y, direction.North, l),
			domain.ErrUnableToMove,
		},
		{
			"Cannot move East",
			domain.NewRover(MAX_Y, 1, direction.East, l),
			want(MAX_Y, 1, direction.East, l),
			domain.ErrUnableToMove,
		},
		{
			"Cannot move South",
			domain.NewRover(1, 0, direction.South, l),
			want(1, 0, direction.South, l),
			domain.ErrUnableToMove,
		},
		{
			"Cannot move West",
			domain.NewRover(0, 1, direction.West, l),
			want(0, 1, direction.West, l),
			domain.ErrUnableToMove,
		},

		{
			"Moving rover with undefined heading has error",
			domain.NewRover(2, 2, UNDEFINED_DIRECTION, l),
			want(2, 2, UNDEFINED_DIRECTION, l),
			domain.ErrUnableToMoveUndefinedHeading,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			err := tt.args.Move()

			c.Assert(errors.Is(err, tt.wantErr), qt.IsTrue)
			c.Assert(tt.args, qt.CmpEquals(cmp.AllowUnexported(domain.Rover{}, domain.Plateau{})), tt.want)
		})
	}
}

func TestRover_Rotate(t *testing.T) {

	type args struct {
		x           uint64
		y           uint64
		heading     direction.Direction
		landingZone domain.LandingZone

		turn control.Control
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
		want    direction.Direction
	}{
		{
			"Heading North when turning Right should face East",
			args{
				heading:     direction.North,
				turn:        control.Right,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.East,
		},
		{
			"Heading East when turning Right should face South",
			args{
				heading:     direction.East,
				turn:        control.Right,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.South,
		},
		{
			"Heading South when turning Right should face West",
			args{
				heading:     direction.South,
				turn:        control.Right,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.West,
		},
		{
			"Heading West when turning Right should face North",
			args{
				heading:     direction.West,
				turn:        control.Right,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.North,
		},
		{
			"Heading North when turning Left should face West",
			args{
				heading:     direction.North,
				turn:        control.Left,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.West,
		},
		{
			"Heading West when turning Left should face South",
			args{
				heading:     direction.West,
				turn:        control.Left,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.South,
		},
		{
			"Heading South when turning Left should face East",
			args{
				heading:     direction.South,
				turn:        control.Left,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.East,
		},
		{
			"Heading East when turning Left should face North",
			args{
				heading:     direction.East,
				turn:        control.Left,
				landingZone: domain.NewPlateau(5, 5),
			},
			nil,
			direction.North,
		},
		{
			"Rotation with undefined control has error",
			args{
				heading:     direction.East,
				turn:        UNDEFINED_CONTROL,
				landingZone: domain.NewPlateau(5, 5),
			},
			domain.ErrUndefinedControl,
			direction.East,
		},
		{
			"Turning left with undefined heading has error",
			args{
				heading:     UNDEFINED_DIRECTION,
				turn:        control.Left,
				landingZone: domain.NewPlateau(5, 5),
			},
			domain.ErrUndefinedRotationBehaviour,
			UNDEFINED_DIRECTION,
		},
		{
			"Turning right with undefined heading has error",
			args{
				heading:     UNDEFINED_DIRECTION,
				turn:        control.Right,
				landingZone: domain.NewPlateau(5, 5),
			},
			domain.ErrUndefinedRotationBehaviour,
			UNDEFINED_DIRECTION,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := qt.New(t)
			r := domain.NewRover(tt.args.x, tt.args.y, tt.args.heading, tt.args.landingZone)
			err := r.Rotate(tt.args.turn)
			c.Assert(errors.Is(err, tt.wantErr), qt.IsTrue)
			c.Assert(r.Heading(), qt.Equals, tt.want)
		})
	}
}
