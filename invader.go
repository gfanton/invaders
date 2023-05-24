package main

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"time"
)

type Options struct {
	// iteration tick time
	Tick time.Duration
}

type AlienInvaders struct {
	tick   time.Duration
	cities Cities

	aliens       map[*Alien]struct{}
	aliensCities map[string] /* city name */ *Alien
}

func NewAlienInvaders(*Options) *AlienInvaders {
	return &AlienInvaders{
		aliens: make(map[*Alien]struct{}),
		cities: NewCities(),
	}
}

func (ai *AlienInvaders) ParseMap(r io.Reader) error {
	if err := ai.cities.Parse(r); err != nil {
		return fmt.Errorf("unable to parse map: %w", err)
	}

	return nil
}

func (ai *AlienInvaders) GenerateAliens(x int) error {
	cities := ai.cities.GetAll()
	switch {
	case x > len(cities):
		return fmt.Errorf("cannot have more aliens than cities: %d > %d", x, len(cities))
	case x == len(cities): // do nothing
	default:
		rand.Shuffle(len(cities), func(i, j int) {
			cities[i], cities[j] = cities[j], cities[i]
		})
		cities = cities[:x]
	}

	for _, city := range cities {
		alien := NewAlien(city)
		ai.aliens[alien] = struct{}{}
		ai.aliensCities[city.Name] = alien
	}

	return nil
}

func (ai *AlienInvaders) Iteration(ctx context.Context) error {
	for alien := range ai.aliens {
		// move
		newcity := alien.CurrentCity
		if occupy, ok := ai.aliensCities[newcity.Name]; ok {
			// we got a fight !
			fmt.Println("%s has been destroyed by alien %s and alien %s!", newcity.Name, alien.Name(), occupy.Name())

			// cleanup
			alien.Alive = false
			occupy.Alive = false
			delete(ai.aliensCities, newcity.Name)
			delete(ai.aliens, occupy)
			delete(ai.aliens, alien)
		}

		delete(ai.aliensCities, newcity.Name)
		ai.aliensCities[newcity.Name] = alien
		// enemy, ok := ai.aliens[alien.CurrentCity.Name]
	}

	return nil
}

func (ai *AlienInvaders) Run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(ai.tick): // ok
		}

	}
}
