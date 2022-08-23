# gomarsrover

Mars Rover is a simple command line program which is a solution to the Mars Rover technical challenge.

## MARS ROVERS

A squad of robotic rovers are to be landed by NASA on a plateau on Mars.This plateau, which is curiously rectangular, must be navigated by the rovers so that their on board cameras can get a complete view of the surrounding terrain to send back to Earth. A rover's position is represented by a combination of an x and y co-ordinates and a letter representing one of the four cardinal compass points.The plateau is divided up into a grid to simplify navigation. An example position might be 0, 0, N, which means the rover is in the bottom left corner and facing North. In order to control a rover, NASA sends a simple string of letters. The possible letters are 'L', 'R' and 'M'. 'L' and 'R' makes the rover spin 90 degrees left or right respectively, without moving from its current spot.'M' means move forward one grid point, and maintain the same heading.
Assume that the square directly North from (x, y) is (x, y+1).

### Input

The first line of input is the upper-right coordinates of the plateau, the lower-left coordinates are assumed to be 0,0. The rest of the input is information pertaining to the rovers that have been deployed. Each rover has two lines of input. The first line gives the rover's position, and the second line is a series of instructions telling the rover how to explore the plateau.The position is made up of two integers and a letter separated by spaces, corresponding to the x and y co-ordinates and the rover's orientation. Each rover will be finished sequentially, which means that the second rover won't start to move until the first one has finished moving.

### Output

The output for each rover should be its final co-ordinates and heading.

#### Test

Input:

```bash
5 5
1 2 N
LMLMLMLMM
3 3 E
MMRMMRMRRM
```

Expected Output:

```bash
1 3 N
5 1 E
```

## Usage

To use the CLI program, execute the following commands in a new terminal window:

```bash
go run cmd/gomarsrover/main.go -input testdata/data.txt
```

Build with version: `go1.16.4 darwin/amd64`

## Testing

To run the project tests with coverage, execute the following:

```bash
go test ./... -cover
```

## API Documentation

`go doc -all domain`

`go doc -all service`

## Next Stage for production

In order to take this service to production, I would consider the following steps:

- Add CI integration with GitHub Actions
- Add better support for logging e.g logging verbosity (error, critical, debug, info, etc...)
- Add service metric functionality, for reporting (ELK for logs and prometheus & grafana for metrics)
- Cucumber Tests for `CommandCenter` in `cmd/gomarsrover`
