package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestRandomDurationAround(t *testing.T) {
	occurrences := 100
	for i := 0; i < occurrences; i++ {
		res := RandomDurationAround(time.Second, 0.25)
		require.GreaterOrEqual(t, res, 750*time.Millisecond)
		require.LessOrEqual(t, res, 1250*time.Millisecond)
	}

	for i := 0; i < occurrences; i++ {
		res := RandomDurationAround(time.Second, 0.50)
		require.GreaterOrEqual(t, res, 500*time.Millisecond)
		require.LessOrEqual(t, res, 1500*time.Millisecond)
	}
}
