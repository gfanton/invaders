package main

import (
	"fmt"
	"math/rand"
	"strings"
)

type Direction string

const (
	North Direction = "north"
	East  Direction = "east"
	West  Direction = "west"
	South Direction = "south"
)

var (
	ErrCityOccupy            = fmt.Errorf("city already occupy")
	ErrTrapped               = fmt.Errorf("alien is trapped")
	ErrDirectionNotAvailable = fmt.Errorf("direction is not available")
)

var AllDirections = []Direction{North, East, West, South}

func ParseDirection(dir string) (Direction, error) {
	switch d := Direction(strings.ToLower(dir)); d {
	case North, East, West, South:
		return d, nil
	default:
		return "", fmt.Errorf("invalid direction")
	}
}

type City struct {
	Name  string
	Alien *Alien

	borderCities map[Direction]*City
}

func NewCity(name string) *City {
	return &City{
		Name:         name,
		borderCities: make(map[Direction]*City),
	}
}

func (c *City) MoveAlien(dir Direction) (target *City, err error) {
	target, ok := c.borderCities[dir]
	switch {
	case !ok:
		return nil, ErrDirectionNotAvailable
	case target.Alien != nil:
		return target, ErrCityOccupy
	default: //ok
	}

	c.Alien.Steps++
	target.Alien = c.Alien
	c.Alien = nil
	return
}

func (c *City) MoveAlienRandomly() (target *City, err error) {
	dirs := c.GetAvailableDirections()
	if len(dirs) == 0 {
		return nil, ErrTrapped
	}

	ndir := rand.Intn(len(dirs))
	target, _ = c.MoveAlien(dirs[ndir])
	return
}

func (c *City) GetDirection(dir Direction) (city *City, ok bool) {
	city, ok = c.borderCities[dir]
	return
}

func (c *City) GetAvailableDirections() (dirs []Direction) {
	dirs = make([]Direction, 0, len(AllDirections))
	for dir, c := range c.borderCities {
		if c != nil {
			dirs = append(dirs, dir)
		}
	}

	return
}
