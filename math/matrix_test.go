package math

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestFractals(t *testing.T) {
	m := NewMatrix(4)
	m.AddRow("1").Set(0, 10.0).Set(1, 10.0).Set(2, 7.0).Set(3, 8.0)
	m.AddRow("2").Set(0, 11.0).Set(1, 11.0).Set(2, 6.0).Set(3, 8.0)
	m.AddRow("3").Set(0, 12.0).Set(1, 10.0).Set(2, 7.0).Set(3, 8.0)
	m.AddRow("4").Set(0, 9.0).Set(1, 12.0).Set(2, 6.0).Set(3, 8.0)
	m.AddRow("5").Set(0, 8.0).Set(1, 10.0).Set(2, 7.0).Set(3, 8.0)
	sp := m.Fractals(3)
	assert.Equal(t, 4, len(sp))
}

func TestStd(t *testing.T) {
	m := NewMatrix(1)
	// 10, 12, 23, 23, 16, 23, 21, 16
	m.AddRow("1").Set(0, 10.0)
	m.AddRow("2").Set(0, 12.0)
	m.AddRow("3").Set(0, 23.0)
	m.AddRow("4").Set(0, 23.0)
	m.AddRow("5").Set(0, 16.0)
	m.AddRow("6").Set(0, 23.0)
	m.AddRow("7").Set(0, 21.0)
	m.AddRow("8").Set(0, 16.0)
	sp := m.StdDev(0, 7)
	for _, r := range m.DataRows {
		fmt.Println(r)
	}
	assert.Equal(t, "5.17", fmt.Sprintf("%.2f", m.DataRows[m.Rows-1].Get(sp)))
}
