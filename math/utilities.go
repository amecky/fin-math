package math

import (
	m "math"
)

func CeilNearest(value, precision float64) float64 {
	return m.Ceil(value/precision) * precision
}

func RoundNearest(value, precision float64) float64 {
	return m.Round(value/precision) * precision
}

func FloorNearest(value, precision float64) float64 {
	return m.Floor(value/precision) * precision
}
