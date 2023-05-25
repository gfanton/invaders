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

func TestParseMap(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)

	reader := strings.NewReader("City1 north=City2")
	err := ai.ParseMap(reader)

	require.NoError(t, err)
}

func TestParseMapError(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)

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

	ctx := context.Background()

	reader := strings.NewReader("a north=b south=c east=d")
	err := ai.ParseMap(reader)
	require.NoError(t, err)

	ai.PrintMap()

	err = ai.GenerateAliens(1)
	require.NoError(t, err)

	err = ai.Run(ctx, 100)
	require.NoError(t, err)
	require.Len(t, ai.aliens, 1)
}

func TestAlienBattle(t *testing.T) {
	ai := NewAlienInvaders(testInvaderLogger, testInvaderDefaultWriter)
	ctx := context.Background()

	reader := strings.NewReader("a north=b")
	err := ai.ParseMap(reader)
	require.NoError(t, err)

	err = ai.GenerateAliens(2)
	require.NoError(t, err)

	err = ai.Run(ctx, 100)
	require.Error(t, err)
	require.Equal(t, ErrAllAliensAreKO, err)
	require.Len(t, ai.aliens, 0)
}
