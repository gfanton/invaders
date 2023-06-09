package invader

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewCities(t *testing.T) {
	cities := NewCities()

	require.Len(t, cities, 0)
}

func TestGetOrCreate(t *testing.T) {
	cities := NewCities()
	cityName := "TestCity"
	city := cities.GetOrCreate(cityName)

	require.Equal(t, cityName, city.Name)
	require.Equal(t, city, cities[cityName], "City was not added to the cities map.")
}

func TestParse(t *testing.T) {
	cities := NewCities()
	reader := strings.NewReader("TestCity North=BorderCity")

	err := cities.Parse(reader)
	require.NoError(t, err)

	require.Len(t, cities, 2)

	testCity, exists := cities.Get("TestCity")
	require.True(t, exists)

	borderCity, exists := cities.Get("BorderCity")
	require.True(t, exists)

	northCity, exists := testCity.GetDirection(North)
	require.True(t, exists)
	require.NotNil(t, northCity)
	require.Equal(t, northCity, borderCity)

	southCity, exists := testCity.GetDirection(South)
	require.False(t, exists)
	require.Nil(t, southCity)
}

func TestParseInput(t *testing.T) {
	testCases := []struct {
		Name       string
		Input      string
		WantError  bool
		WantCities int
	}{
		{
			Name:       "correct input",
			Input:      "TestCity North=BorderCity",
			WantError:  false,
			WantCities: 2,
		},
		{
			Name:       "correct input",
			Input:      "TestCity North=BorderCity South=BorderCity2 West=BorderCity3",
			WantError:  false,
			WantCities: 4,
		},
		{
			Name: "same input multiple time",
			Input: `TestCity North=BorderCity South=BorderCity2 West=BorderCity3
 TestCity North=BorderCity South=BorderCity2 West=BorderCity3
 TestCity North=BorderCity South=BorderCity2 West=BorderCity3
 TestCity North=BorderCity South=BorderCity2 West=BorderCity3
			`,
			WantError:  false,
			WantCities: 4,
		},
		{
			Name:       "no input",
			Input:      "",
			WantCities: 0,
		},
		{
			Name:       "trapped city",
			Input:      "TestCity",
			WantCities: 1,
		},
		{
			Name:      "malformed input",
			Input:     "TestCity North BorderCity",
			WantError: true,
		},
		{
			Name:      "unknown direction",
			Input:     "TestCity Northeast=BorderCity",
			WantError: true,
		},
		{
			Name:      "invalid city name",
			Input:     "Northeast=BorderCity",
			WantError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			cities := NewCities()
			reader := strings.NewReader(tc.Input)

			err := cities.Parse(reader)
			if tc.WantError {
				require.Error(t, err)
				return
			}

			require.Len(t, cities, tc.WantCities)

		})
	}
}

func TestGenerateRandomCity(t *testing.T) {
	cities := NewCities()
	depth := 10
	cities.GenerateRandomCity(depth)

	// Check if the number of generated cities is less than or equal to depth * depth
	require.Less(t, len(cities), depth*depth)

	// Check if every city has at least one border city
	for _, city := range cities {
		hasBorder := false
		for _, borderCity := range city.borderCities {
			if borderCity != nil {
				hasBorder = true
				break
			}
		}

		require.True(t, hasBorder)
	}
}
