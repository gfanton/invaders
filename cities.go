package main

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type Cities map[string] /* city name */ *City

func NewCities() Cities {
	return make(map[string]*City)
}

func (cs Cities) Get(name string) (city *City, ok bool) {
	city, ok = cs[name]
	return
}

func (cs Cities) Destroy(name string) {
	delete(cs, name)
}

func (cs Cities) GetAll() []*City {
	all := make([]*City, len(cs))
	i := 0
	for _, city := range cs {
		all[i] = city
		i++
	}

	return all
}

func (cs Cities) GetOrCreate(name string) (city *City) {
	var ok bool
	if city, ok = cs[name]; ok {
		return
	}

	city = NewCity(name)
	cs[name] = city
	return
}

// Parse method reads from an io.Reader and populates the Cities map with City structs.
// Each line from the reader should represent a city and its bordering cities.
// The format of each line should be: "CityName Border1=CityName Border2=CityName ..."
// For example: "Paris North=Lille South=Lyon East=Strasbourg West=Rouen"
func (cs Cities) Parse(r io.Reader) error {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)

		// If the line doesn't have at least two parts (a city and at least one border),
		// then return an error
		if len(parts) < 2 {
			return fmt.Errorf("malformed line: %s", line)
		}

		cityName := parts[0]
		city := cs.GetOrCreate(cityName)
		for _, border := range parts[1:] {
			borderParts := strings.Split(border, "=")

			// If the border part doesn't split into two parts, then return an error
			if len(borderParts) != 2 {
				return fmt.Errorf("malformed border: %s, in line: %s", border, line)
			}

			borderDirection, borderCityName := strings.ToLower(borderParts[0]), borderParts[1]
			dir, err := ParseDirection(borderDirection)
			if err != nil {
				return fmt.Errorf("unable to parse direction `%s`: %w ", dir, err)
			}

			// Get or create the border city
			borderCity := cs.GetOrCreate(borderCityName)
			if borderCity == city {
				return fmt.Errorf("a city cannot be bordered by itself")
			}

			switch Direction(borderDirection) {
			case North, South, West, East:
				city.borderCities[dir] = borderCity
			default: // unknown direction
				return fmt.Errorf("unknown border direction: %s for city: %s", borderDirection, borderCityName)
			}
		}
	}

	return scanner.Err()
}
