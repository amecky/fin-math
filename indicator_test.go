package main

import (
	"testing"

	"github.com/amecky/fin-math/assert"
	"github.com/amecky/fin-math/math"
)

func TestSMA(t *testing.T) {
	m := math.NewMatrix(1)
	m.AddRow("1").Set(0, 1.0)
	m.AddRow("2").Set(0, 2.0)
	m.AddRow("3").Set(0, 3.0)
	m.AddRow("4").Set(0, 4.0)
	m.AddRow("5").Set(0, 5.0)
	m.AddRow("6").Set(0, 6.0)
	sma := math.SMA(m, 5, 0)
	assert.Equals(t, "Wrong number of entries", 1.0, float64(sma))
	assert.Equals(t, "Wrong value", 3.0, m.FindRow("5").Get(sma))
	assert.Equals(t, "Wrong value", 4.0, m.FindRow("6").Get(sma))
}
