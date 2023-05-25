package invader

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
)

var (
	ErrAllAliensAreDead = fmt.Errorf("all aliens are dead")
)

type AlienInvaders struct {
	logger *log.Logger
	cities Cities
	aliens map[*Alien]struct{}
}

func NewAlienInvaders(logger *log.Logger) *AlienInvaders {
	return &AlienInvaders{
		logger: logger,
		aliens: make(map[*Alien]struct{}),
		cities: NewCities(),
	}
}

func (ai *AlienInvaders) ParseMap(r io.Reader) error {
	if err := ai.cities.Parse(r); err != nil {
		return fmt.Errorf("unable to parse map: %w", err)
	}

	ai.logger.Printf("successfully parsed %d cities", len(ai.cities))
	return nil
}

func (ai *AlienInvaders) PrintMap() {
	ai.cities.Print(os.Stdout)
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
	}

	return nil
}

func (ai *AlienInvaders) nextIteration() (deadAliens []*Alien) {
	deadAliens = []*Alien{}

	for alien := range ai.aliens {
		if alien.State == Killed {
			// alien has been previously killed in this iteration, skip
			continue
		}

		currentCity := alien.CurrentCity

		// make a random move
		occupyAlien, ok := alien.RandomMove()
		if !ok {
			// alien has been trapped!
			// collect it has a dead body, but we dont have to kill
			// it, he can't move anyway

			deadAliens = append(deadAliens, alien)

			fmt.Printf("%s has been trapped!\n", alien.Name())
			continue
		}

		ai.logger.Printf("alien `%s` moved from `%s` to `%s`", alien.Name(), currentCity.Name, alien.CurrentCity.Name)

		// move was succefull check if someone is here

		if occupyAlien != nil { // someone already here !
			ai.logger.Printf("city `%s` already occuped by alien `%s`", occupyAlien.CurrentCity.Name, occupyAlien.Name())

			targetcity := alien.CurrentCity

			// we got a fight !

			// mark both aliens has dead
			alien.Kill()
			occupyAlien.Kill()

			// destroy the city
			targetcity.Destroy()

			// gather dead bodies (for later cleanup)
			deadAliens = append(deadAliens, alien, occupyAlien)

			fmt.Printf("%s has been destroyed by alien %s and alien %s!\n", targetcity.Name, alien.Name(), occupyAlien.Name())
		}
	}

	return
}

func (ai *AlienInvaders) Run(ctx context.Context, limit int) error {
	totalAliens := len(ai.aliens)

	for steps := 0; steps < limit && ctx.Err() == nil; steps++ {
		ai.logger.Printf("iteration: %d\n", steps)

		// generate next iteration, collect dead bodies
		deadAliens := ai.nextIteration()

		// cleanup dead bodies
		for _, deadAlien := range deadAliens {
			delete(ai.aliens, deadAlien)
		}

		if len(deadAliens) > 0 {
			ai.logger.Printf("iteration[%d]: %d/%d aliens have been killed", steps, totalAliens-len(ai.aliens), totalAliens)
		}

		// are we done ?
		if len(ai.aliens) == 0 {
			return ErrAllAliensAreDead
		}
	}

	ai.logger.Printf("%d/%d aliens left", len(ai.aliens), totalAliens)

	return ctx.Err()
}
