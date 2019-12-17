package config

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPlusMinusDelta(t *testing.T) {
	occurrences := 100
	for i := 0; i < occurrences; i++ {
		res := PlusMinusDelta(1, DefaultDelta)
		require.Truef(t, res >= 0.75, "valus is %v, should be >= 0.75", res)
		require.Truef(t, res <= 1.25, "value is %v, should be <= 1.25", res)
	}

	for i := 0; i < occurrences; i++ {
		res := PlusMinusDelta(1, Delta(0.5))
		require.Truef(t, res >= 0.5, "valus is %v, should be >= 0.50", res)
		require.Truef(t, res <= 1.5, "value is %v, should be <= 1.50", res)
	}
}
