package invader

import (
	"fmt"
	"strings"
)

type Direction string

// Define the four possible directions.
const (
	North Direction = "north"
	East  Direction = "east"
	West  Direction = "west"
	South Direction = "south"
)

// ParseDirection converts a string to Direction type.
func ParseDirection(dir string) (Direction, error) {
	switch d := Direction(strings.ToLower(dir)); d {
	case North, East, West, South:
		return d, nil
	default:
		return "", fmt.Errorf("invalid direction")
	}
}

// AllDirections is a slice containing all possible directions.
var AllDirections = []Direction{North, East, West, South}

// Opposite method returns the opposite direction.
func (d Direction) Opposite() Direction {
	switch d {
	case North:
		return South
	case South:
		return North
	case East:
		return West
	case West:
		return East
	}

	return Direction("")
}
