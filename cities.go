package invader

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
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

// GenerateRandomCity populates the Cities map with a collection of cities.
// Each city is connected to one or more neighboring cities, forming a random graph.
// The graph is organized as a square grid of a specified depth.
func (cs Cities) GenerateRandomCity(depth int) {
	// Initialize a square grid of pointers to City
	table := make([][]*City, depth)
	for y := range table {
		table[y] = make([]*City, depth)
	}

	// Generate random starting position
	x, y := rand.Intn(depth), rand.Intn(depth)

	// Helper function to generate city name
	var counterID int64
	genid := func() string {
		counterID++
		return fmt.Sprintf("city_%d", counterID)
	}

	// Helper function to calculate new position based on direction
	calculateNewPosition := func(dir Direction) (newx int, newy int) {
		switch dir {
		case North:
			return x - 1, y
		case South:
			return x + 1, y
		case East:
			return x, y + 1
		case West:
			return x, y - 1
		default:
			return x, y
		}
	}

	// Create root city and add it to cities map and grid
	root := NewCity(genid())
	cs[root.Name] = root
	table[y][x] = root

	// Initialize distance from root
	distance := 0

	// Recursive function to generate cities and connections
	var citygen func(dir Direction, c *City)
	citygen = func(dir Direction, c *City) {
		// If we've reached the maximum depth, stop
		if distance > depth {
			return
		}

		distance++

		// Remember the current position
		oldx, oldy := x, y

		// Calculate new position based on direction
		x, y = calculateNewPosition(dir)

		// If the new position is within the grid
		if x >= 0 && x < depth && y >= 0 && y < depth {
			// If a city already exists at the new position, connect it
			if table[y][x] != nil {
				c.borderCities[dir] = table[y][x]
				table[y][x].borderCities[dir.Opposite()] = c
			} else {
				// Create a new city and add it to cities map and grid
				newcity := NewCity(genid())
				table[y][x] = newcity
				cs[newcity.Name] = newcity

				// Connect the new city
				c.borderCities[dir] = newcity
				newcity.borderCities[dir.Opposite()] = c

				// Generate a random list of directions, with a random length of 1 or more
				ad := make([]Direction, len(AllDirections))
				copy(ad, AllDirections)
				rand.Shuffle(len(ad), func(i, j int) { ad[i], ad[j] = ad[j], ad[i] })
				ad = ad[0 : rand.Intn(len(AllDirections)-1)+1]

				// Recursively generate cities in these directions
				for _, dir := range ad {
					citygen(dir, newcity)
				}
			}
		}

		distance--
		// Return to previous position before returning
		x, y = oldx, oldy
	}

	// Start generating cities in all directions from the root
	for _, dir := range AllDirections {
		citygen(dir, root)
	}
}

func (cs Cities) Print(w io.Writer) {
	for _, city := range cs {
		city.Print(w)
	}
}
