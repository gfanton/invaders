package invader

import (
	"fmt"
	"io"
	"strings"
)

type City struct {
	Name  string
	Alien *Alien

	// borderCities maps each direction to a neighbouring city.
	borderCities map[Direction]*City
}

// NewCity creates a new city with its name and an empty map for borderCities.
func NewCity(name string) *City {
	return &City{
		Name:         name,
		borderCities: make(map[Direction]*City),
	}
}

// MoveAlien moves the alien in the specified direction and returns the target city and the alien occupying it (if any).
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

// Destroy removes all the neighbouring cities from the borderCities map and sets the Alien property to nil.
func (c *City) Destroy() {
	c.IterateBorder(func(dir Direction, neighbor *City) {
		delete(neighbor.borderCities, dir.Opposite())
		delete(c.borderCities, dir)
	})
	c.Alien = nil
}

// GetDirection returns the city in the given direction if it exists.
func (c *City) GetDirection(dir Direction) (city *City, ok bool) {
	city, ok = c.borderCities[dir]
	return
}

// SetDirection sets the city in the given direction and also sets the calling city as the opposite direction in the city map.
func (c *City) SetDirection(dir Direction, city *City) {
	city.borderCities[dir.Opposite()] = c
	c.borderCities[dir] = city
}

// GetAvailableDirections returns a slice of directions that have cities.
func (c *City) GetAvailableDirections() (dirs []Direction) {
	dirs = make([]Direction, 0, len(AllDirections))
	for dir, c := range c.borderCities {
		if c != nil {
			dirs = append(dirs, dir)
		}
	}

	return
}

// IterateBorder iterates over all available directions and performs the given function.
func (c *City) IterateBorder(call func(dir Direction, c *City)) {
	for _, dir := range c.GetAvailableDirections() {
		city, _ := c.GetDirection(dir)
		call(dir, city)
	}
}

// Print prints the city name and all its neighbouring cities to the writer.
func (c *City) Print(w io.Writer) {
	dirs := c.GetAvailableDirections()
	citydirs := make([]string, len(dirs))
	for i, dir := range dirs {
		citydir, _ := c.GetDirection(dir)
		citydirs[i] = fmt.Sprintf("%s=%s", dir, citydir.Name)
	}

	fmt.Fprintf(w, "%s %s\n", c.Name, strings.Join(citydirs, " "))
}
