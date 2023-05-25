package invader

import (
	"fmt"
	"io"
	"strings"
)

type Direction string

const (
	North Direction = "north"
	East  Direction = "east"
	West  Direction = "west"
	South Direction = "south"
)

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

var (
	ErrCityOccupy = fmt.Errorf("city already occupy")
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

func (c *City) MoveAlien(dir Direction) (target *City, occupy *Alien) {
	target, ok := c.borderCities[dir]
	if !ok {
		return
	}

	occupy = target.Alien
	target.Alien = c.Alien
	c.Alien = nil
	return
}

func (c *City) Destroy() {
	c.IterateBorder(func(dir Direction, neighbor *City) {
		delete(neighbor.borderCities, dir.Opposite())
		delete(c.borderCities, dir)
	})
	c.Alien = nil
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

func (c *City) IterateBorder(call func(dir Direction, c *City)) {
	for _, dir := range c.GetAvailableDirections() {
		city, _ := c.GetDirection(dir)
		call(dir, city)
	}
}

func (c *City) Print(w io.Writer) {
	dirs := c.GetAvailableDirections()
	citydirs := make([]string, len(dirs))
	for i, dir := range dirs {
		citydir, _ := c.GetDirection(dir)
		citydirs[i] = fmt.Sprintf("%s=%s", dir, citydir.Name)
	}

	fmt.Fprintf(w, "%s %s\n", c.Name, strings.Join(citydirs, " "))
}
