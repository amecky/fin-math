package math

import (
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestHighest(t *testing.T) {
	m := NewMatrix(1)
	m.AddRow("1").Set(0, 10.0)
	m.AddRow("2").Set(0, 11.0)
	m.AddRow("3").Set(0, 12.0)
	m.AddRow("4").Set(0, 9.0)
	m.AddRow("5").Set(0, 8.0)
	m.AddRow("6").Set(0, 7.0)
	m.AddRow("7").Set(0, 13.0)
	m.AddRow("8").Set(0, 11.0)
	m.AddRow("9").Set(0, 9.0)
	hi := Highest(m, 10, 0)
	assert.Equal(t, 13, m.DataRows[8].Get(hi))
}

func TestLowest(t *testing.T) {
	m := NewMatrix(1)
	m.AddRow("1").Set(0, 10.0)
	m.AddRow("2").Set(0, 11.0)
	m.AddRow("3").Set(0, 12.0)
	m.AddRow("4").Set(0, 9.0)
	m.AddRow("5").Set(0, 8.0)
	m.AddRow("6").Set(0, 7.0)
	m.AddRow("7").Set(0, 13.0)
	m.AddRow("8").Set(0, 11.0)
	m.AddRow("9").Set(0, 9.0)
	hi := Lowest(m, 10, 0)
	assert.Equal(t, 7, m.DataRows[8].Get(hi))
}

func TestWeightedMA(t *testing.T) {
	m := NewMatrix(1)
	m.AddRow("1").Set(0, 4.0)
	m.AddRow("2").Set(0, 5.0)
	m.AddRow("3").Set(0, 6.0)
	m.AddRow("4").Set(0, 3.0)
	m.AddRow("5").Set(0, 2.0)
	hi := WeightedMA(m, 5, 0)
	assert.Equal(t, 3.6, m.DataRows[4].Get(hi))
}
