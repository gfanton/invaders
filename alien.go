package main

import (
	"math/rand"
	"strconv"
)

var counter uint64

type Alien struct {
	ID          uint64
	Steps       uint
	CurrentCity *City
	Alive       bool
}

func NewAlien(c *City) *Alien {
	counter++
	return &Alien{
		ID:          counter,
		CurrentCity: c,
	}
}

func (a *Alien) Name() string {
	return strconv.FormatUint(a.ID, 10)
}

func (a *Alien) Move(dir Direction) (old *City, ok bool) {
	var new *City
	new, ok = a.CurrentCity.GetDirection(dir)
	if !ok {
		return
	}

	old = a.CurrentCity
	a.CurrentCity = new
	return
}

func (a *Alien) RandomMove() (old *City) {
	if dirs := a.CurrentCity.GetAvailableDirections(); len(dirs) > 0 {
		ndir := rand.Intn(len(dirs))
		old, _ = a.Move(dirs[ndir])
	}
	return
}
