package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCity(t *testing.T) {
	cityName := "TestCity"
	city := NewCity(cityName)

	require.Equal(t, city.Name, cityName)
	require.Len(t, city.borderCities, 0)
}
