package invader

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewAlien(t *testing.T) {
	city := NewCity("TestCity")
	alien := NewAlien(city)

	require.Greater(t, alien.ID, uint(0))
	require.Equal(t, city, alien.CurrentCity)
	require.Equal(t, city.Alien, alien)
}

func TestAlienMove(t *testing.T) {
	city := NewCity("TestCity")
	borderCity := NewCity("BorderCity")
	city.borderCities[North] = borderCity

	alien := NewAlien(city)
	_, ok := alien.Move(North)
	require.True(t, ok)

	require.Equal(t, borderCity, alien.CurrentCity)
}

func TestAlienRandomMove(t *testing.T) {
	city := NewCity("TestCity")
	borderCity := NewCity("BorderCity")
	city.borderCities[North] = borderCity

	alien := NewAlien(city)
	_, ok := alien.RandomMove()
	require.True(t, ok)

	require.Equal(t, borderCity, alien.CurrentCity)
}
