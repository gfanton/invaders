package invader

import (
	"math/rand"
	"strconv"
)

type AlienState int

const (
	Alive AlienState = iota
	Killed
	Trapped
)

var counterID uint

type Alien struct {
	ID          uint
	Steps       uint
	CurrentCity *City
	State       AlienState
}

func NewAlien(c *City) *Alien {
	counterID++
	a := &Alien{
		ID:          counterID,
		CurrentCity: c,
	}

	// Put the alien inside the city
	c.Alien = a
	return a
}

func (a *Alien) Kill() {
	a.State = Killed
}

func (a *Alien) Name() string {
	return strconv.FormatUint(uint64(a.ID), 10)
}

// Move moves the Alien in a given direction. It returns a reference to an Alien
// occupying the city in the direction of the move (if any) and a boolean indicating
// whether the move was successful.
func (a *Alien) Move(dir Direction) (occupy *Alien, ok bool) {
	if a.State == Killed {
		return
	}

	var newcity *City
	newcity, occupy = a.CurrentCity.MoveAlien(dir)
	if newcity == nil { // not target available, should be trapped
		a.State = Trapped
		return
	}

	ok = true

	// Replace current city
	a.CurrentCity = newcity
	return
}

// RandomMove makes the Alien move in a random available direction.
// It returns a reference to an Alien occupying the city in the direction of the move (if any)
// and a boolean indicating whether the move was successful.
func (a *Alien) RandomMove() (occupy *Alien, ok bool) {
	if dirs := a.CurrentCity.GetAvailableDirections(); len(dirs) > 0 {
		ndir := rand.Intn(len(dirs))
		return a.Move(dirs[ndir])
	}

	// No target available, should be trapped
	a.State = Trapped
	return
}
