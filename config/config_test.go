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
		require.Truef(t, res >= 750*time.Millisecond, "valus is %v, should be >= 750ms", res)
		require.Truef(t, res <= 1250*time.Millisecond, "value is %v, should be <= 1.25s", res)
	}

	for i := 0; i < occurrences; i++ {
		res := RandomDurationAround(time.Second, 0.50)
		require.Truef(t, res >= 500*time.Millisecond, "valus is %v, should be >= 500ms", res)
		require.Truef(t, res <= 1500*time.Millisecond, "value is %v, should be <= 1.50s", res)
	}
}
