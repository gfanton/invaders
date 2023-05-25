package invader

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
)

var (
	ErrAllAliensAreKO = fmt.Errorf("all aliens are dead/trapped")
)

type AlienInvaders struct {
	writer io.Writer
	logger *log.Logger
	cities Cities
	aliens map[*Alien]struct{} // Keeps track of all active aliens
}

func NewAlienInvaders(logger *log.Logger, writter io.Writer) *AlienInvaders {
	return &AlienInvaders{
		writer: writter,
		logger: logger,
		aliens: make(map[*Alien]struct{}),
		cities: NewCities(),
	}
}

// ParseMap parses the map from the provided reader.
func (ai *AlienInvaders) ParseMap(r io.Reader) error {
	if err := ai.cities.Parse(r); err != nil {
		return err
	}

	ai.logger.Printf("successfully parsed %d cities", len(ai.cities))
	return nil
}

// PrintMap prints the current map to stdout.
func (ai *AlienInvaders) PrintMap() {
	ai.cities.Print(ai.writer)
}

// GenerateAliens generates the given number of aliens and places them in random cities.
func (ai *AlienInvaders) GenerateAliens(x int) error {
	cities := ai.cities.GetAll()
	switch {
	case x > len(cities):
		return fmt.Errorf("cannot have more aliens than cities: %d > %d", x, len(cities))
	case x == len(cities): // If the number of aliens equals the number of cities, there's no need to shuffle
	default:
		// Shuffle the cities and slice to the desired number of aliens
		rand.Shuffle(len(cities), func(i, j int) {
			cities[i], cities[j] = cities[j], cities[i]
		})
		cities = cities[:x]
	}

	// Create new aliens in the chosen cities
	for _, city := range cities {
		alien := NewAlien(city)
		ai.aliens[alien] = struct{}{}
	}

	return nil
}

// nextIteration simulates the next iteration in the alien invasion.
// It returns the aliens that died during the iteration and an error, if any occurred.
func (ai *AlienInvaders) nextIteration(ctx context.Context) (deadAliens []*Alien, err error) {
	deadAliens = []*Alien{}

	for alien := range ai.aliens {
		if ctx.Err() != nil {
			return nil, ctx.Err() // If context is cancelled, return immediately
		}

		// If the alien has been killed/Trapped, skip it
		switch alien.State {
		case Killed, Trapped:
			continue
		default:
		}

		currentCity := alien.CurrentCity

		// Make a random move
		occupyAlien, ok := alien.RandomMove()
		if !ok { // Alien is trapped and cannot move
			deadAliens = append(deadAliens, alien)
			fmt.Fprintf(ai.writer, "alien %s has been trapped in `%s`!\n", alien.Name(), alien.CurrentCity.Name)
			continue
		}

		ai.logger.Printf("alien `%s` moved from `%s` to `%s`", alien.Name(), currentCity.Name, alien.CurrentCity.Name)

		// Move was succefull, Check if the city is already occupied
		if occupyAlien != nil {
			ai.logger.Printf("city `%s` already occupied by alien `%s`", occupyAlien.CurrentCity.Name, occupyAlien.Name())

			targetCity := alien.CurrentCity

			// We got a fight !

			// Mark both aliens as dead
			alien.Kill()
			occupyAlien.Kill()

			// Destroy the city
			targetCity.Destroy()

			// Gather the dead aliens body for later cleanup
			deadAliens = append(deadAliens, alien, occupyAlien)

			fmt.Fprintf(ai.writer, "%s has been destroyed by alien %s and alien %s!\n", targetCity.Name, alien.Name(), occupyAlien.Name())
		}
	}

	return
}

// Run starts the simulation and continues it for the specified number of
// iterations or until context is cancelled.
func (ai *AlienInvaders) Run(ctx context.Context, limit int) error {
	totalAliens := len(ai.aliens)
	trappedAliens := []*Alien{}

	for steps := 0; steps < limit && ctx.Err() == nil; steps++ {
		ai.logger.Printf("iteration: %d\n", steps)

		// Generate next iteration, collect dead bodies
		deadAliens, err := ai.nextIteration(ctx)
		if err != nil {
			return fmt.Errorf("failed to generate next iteration: %w", err)
		}

		// Remove dead aliens from the map
		for _, deadAlien := range deadAliens {
			if deadAlien.State == Trapped {
				// keep track of trapped aliens for later logging
				trappedAliens = append(trappedAliens, deadAlien)
			}
			delete(ai.aliens, deadAlien)
		}

		if len(deadAliens) > 0 {
			ai.logger.Printf("iteration[%d]: %d/%d aliens have been killed/trapped", steps, totalAliens-len(ai.aliens), totalAliens)
		}

		// If all aliens are dead, stop the simulation
		if len(ai.aliens) == 0 {
			return ErrAllAliensAreKO
		}
	}

	// Log the remaining alien and their position
	ai.logger.Printf("%d/%d aliens left", len(ai.aliens)+len(trappedAliens), totalAliens)
	for alien := range ai.aliens {
		ai.logger.Printf("alien `%s` live in `%s`", alien.Name(), alien.CurrentCity.Name)
	}
	for _, alien := range trappedAliens {
		ai.logger.Printf("alien `%s` trapped in `%s`", alien.Name(), alien.CurrentCity.Name)
	}

	// exit run
	return ctx.Err()
}
