package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/johngillott/gomarsrover/pkg/domain"
	"github.com/johngillott/gomarsrover/pkg/domain/control"
	"github.com/johngillott/gomarsrover/pkg/domain/direction"
	"github.com/johngillott/gomarsrover/pkg/service"
)

var (
	plateauCoordinatesCommandRE = regexp.MustCompile(`^\d+ \d+$`)
	deployRoverCommandRE        = regexp.MustCompile(`^\d+ \d+ [NSWE]$`)
	moveRoverCommandRE          = regexp.MustCompile(`^[LMR]+$`)
)

var (
	input string
)

func check(err error, format string, v ...interface{}) {
	if err != nil {
		log.Fatalf(format+"\n", v...)
	}
}

func main() {

	flag.StringVar(&input, "input", "testdata/data.txt", "Provide a text file with rover commands")
	flag.Parse()

	f, err := os.Open(input)
	check(err, "Unable to open input command file - %s", err)
	defer f.Close()

	comm := NewCommandCenter()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		command := scanner.Text()
		log.Println(command)
		message, err := comm.Send(command)
		check(err, "bad command: %q - %s", command, err)

		if len(message) > 0 {
			log.Println(message)
		}
	}

	err = scanner.Err()
	check(err, "unable scan input file %s contents - %s", input, err)
}

type CommandCenter struct {
	m *service.MarsRoverService
}

func NewCommandCenter() *CommandCenter {
	return &CommandCenter{}
}

func (comm *CommandCenter) Send(command string) (string, error) {

	switch {
	case moveRoverCommandRE.MatchString(command):
		if comm.m == nil {
			return "", errors.New("unable to move rover - marsroverservice not initialized")
		}

		v := comm.m.ActiveVehicle()

		if v == nil {
			return "", errors.New("unable to move rover - no vehicle is active")
		}

		for _, c := range command {
			if err := comm.m.Move(control.Control(c), v); err != nil {
				return "", err
			}
		}

		return fmt.Sprintf("Rover Position: %d %d %c", v.X(), v.Y(), v.Heading()), nil

	case deployRoverCommandRE.MatchString(command):

		if comm.m == nil {
			return "", errors.New("unable to deploy rover - marsroverservice not initialized")
		}

		args := strings.Fields(command)

		var (
			x, _ = strconv.ParseUint(args[0], 10, 32)
			y, _ = strconv.ParseUint(args[1], 10, 32)
			d    = direction.Direction(args[2][0])
		)

		if err := comm.m.Deploy(x, y, d); err != nil {
			return "", err
		}

	case plateauCoordinatesCommandRE.MatchString(command):

		if comm.m != nil {
			return "", errors.New("marsroverservice already initialized")
		}

		args := strings.Fields(command)

		var (
			x, _ = strconv.ParseUint(args[0], 10, 32)
			y, _ = strconv.ParseUint(args[1], 10, 32)
		)

		p := domain.NewPlateau(x, y)

		comm.m = service.NewMarsRoverService(p)

	default:
		return "", fmt.Errorf("unrecognized command:%q", command)
	}

	return "", nil
}
