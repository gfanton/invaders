package invader

import (
	"context"
	"io"
	"log"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	testInvaderDefaultWriter = io.Discard
	testInvaderLogger        = log.New(testInvaderDefaultWriter, "", log.LstdFlags)
)

func newTestingAlienInvaders(t *testing.T) *AlienInvaders {
	t.Helper()

	// Generate cities and aliens
	aliens := make(map[*Alien]struct{})
	cities := NewCities()
	{
		city1 := NewCity("City1")
		a1 := NewAlien(city1)
		aliens[a1] = struct{}{}

		city2 := NewCity("City2")
		a2 := NewAlien(city2)
		aliens[a2] = struct{}{}

		city3 := NewCity("City3")
		// no alien for city3

		city1.SetDirection(North, city2)
		city3.SetDirection(South, city2)
		cities.Set(city1, city2, city3)
	}

	return &AlienInvaders{
		writer: io.Discard,
		logger: testInvaderLogger,
		cities: cities,
		aliens: aliens,
	}
}

func TestParseMap(t *testing.T) {
	ai := newTestingAlienInvaders(t)

	reader := strings.NewReader("City1 north=City2")
	err := ai.ParseMap(reader)

	require.NoError(t, err)
}

func TestParseMapError(t *testing.T) {
	ai := newTestingAlienInvaders(t)

	reader := strings.NewReader("City1=City2")
	err := ai.ParseMap(reader)

	require.Error(t, err)
}

func TestGenerateAliens(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)

	reader := strings.NewReader("city1 south=city2 north=city3 west=city3")
	err := ai.ParseMap(reader)
	require.NoError(t, err)

	err = ai.GenerateAliens(2)
	require.NoError(t, err)

	require.Equal(t, 2, len(ai.aliens))
}

func TestGenerateAliensNoMapError(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)
	err := ai.GenerateAliens(10)
	require.Error(t, err)
}

func TestAlienExhausted(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)

	reader := strings.NewReader("a north=b south=c east=d")
	ai.ParseMap(reader)
	ai.PrintMap()
	ai.GenerateAliens(1)

	ctx := context.Background()

	err := ai.Run(ctx, 100)
	require.NoError(t, err)
	require.Len(t, ai.aliens, 1)
}

func TestAlienBattle(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)

	reader := strings.NewReader("a north=b")
	ai.ParseMap(reader)
	ai.PrintMap()
	ai.GenerateAliens(2)

	ctx := context.Background()

	err := ai.Run(ctx, 100)
	require.Error(t, err)
	require.Equal(t, ErrAllAliensAreKO, err)
	require.Len(t, ai.aliens, 0)
}
