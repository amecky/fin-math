package math

import (
	"math"
	"sort"
)

// -----------------------------------------------------------------------
// HL2 - (high+low)/2
// -----------------------------------------------------------------------
func HL2(m *Matrix) int {
	ret := m.AddNamedColumn("HL2")
	for i := 0; i < m.Rows; i++ {
		d := (m.DataRows[i].Get(1) + m.DataRows[i].Get(2)) / 2.0
		m.DataRows[i].Set(ret, d)
	}
	return ret
}

// -----------------------------------------------------------------------
// HLC3 - (high+low+close)/3
// -----------------------------------------------------------------------
func HLC3(m *Matrix) int {
	ret := m.AddNamedColumn("HLC3")
	for i := 0; i < m.Rows; i++ {
		d := (m.DataRows[i].Get(1) + m.DataRows[i].Get(2) + m.DataRows[i].Get(4)) / 3.0
		m.DataRows[i].Set(ret, d)
	}
	return ret
}

func OHLC4(m *Matrix) int {
	ret := m.AddNamedColumn("OHLC4")
	for i := 0; i < m.Rows; i++ {
		d := (m.DataRows[i].Get(0) + m.DataRows[i].Get(1) + m.DataRows[i].Get(2) + m.DataRows[i].Get(4)) / 4.0
		m.DataRows[i].Set(ret, d)
	}
	return ret
}

func AVG(m *Matrix, fields ...int) int {
	ret := m.AddNamedColumn("AVG")
	for i := 0; i < m.Rows; i++ {
		sum := 0.0
		for _, f := range fields {
			sum += m.DataRows[i].Get(f)
		}
		m.DataRows[i].Set(i, sum)
	}
	return ret
}

func Lowest(m *Matrix, period, field int) int {
	ret := m.AddNamedColumn("Lowest")
	for i := 0; i < m.Rows; i++ {
		h := m.DataRows[i].Get(field)
		for j := 1; j < period; j++ {
			idx := i - j
			if idx >= 0 && m.DataRows[idx].Get(field) < h {
				h = m.DataRows[idx].Get(field)
			}
		}
		m.DataRows[i].Set(ret, h)
	}
	return ret
}

func Highest(m *Matrix, period, field int) int {
	ret := m.AddNamedColumn("Highest")
	for i := 0; i < m.Rows; i++ {
		h := m.DataRows[i].Get(field)
		for j := 1; j < period; j++ {
			idx := i - j
			if idx >= 0 && m.DataRows[idx].Get(field) > h {
				h = m.DataRows[idx].Get(field)
			}
		}
		m.DataRows[i].Set(ret, h)
	}
	return ret
}

func Quantile(m *Matrix, limit float64, field int) int {
	ret := m.AddNamedColumn("Quantile")
	for i := 0; i < m.Rows; i++ {
		h := m.DataRows[i].Get(field)
		if h >= limit {
			m.DataRows[i].Set(ret, 1.0)
		}
	}
	return ret
}

func Quantiles(m *Matrix, lower, upper float64, field int) int {
	ret := m.AddNamedColumn("Quantiles")
	for i := 0; i < m.Rows; i++ {
		h := m.DataRows[i].Get(field)
		if h <= lower {
			m.DataRows[i].Set(ret, -1.0)
		}
		if h >= upper {
			m.DataRows[i].Set(ret, 1.0)
		}
	}
	return ret
}

// Standard pandas-style quantile (from earlier)
func internalQuantile(data []float64, q float64) float64 {
	tmp := make([]float64, len(data))
	copy(tmp, data)
	sort.Float64s(tmp)

	pos := (float64(len(tmp) - 1)) * q
	i := int(math.Floor(pos))
	f := pos - float64(i)

	if f == 0 {
		return tmp[i]
	}
	return tmp[i] + f*(tmp[i+1]-tmp[i])
}

// RollingQuantile calculates rolling quantiles like pandas.
// window = lookback period (e.g. 14)
// q = quantile (0.75 = 75th percentile)
func RollingQuantile(m *Matrix, field, window int, q float64) int {
	ret := m.AddColumn()
	//result := make([]float64, len(data))
	var win []float64
	for i := 0; i < m.Rows; i++ {

		// grow window until full
		win = append(win, m.DataRows[i].Get(field))

		// once full, drop oldest
		if len(win) > window {
			win = win[1:]
		}

		if len(win) < window {
			//	result[i] = math.NaN() // not enough data yet
			continue
		}
		m.DataRows[i].Set(ret, internalQuantile(win, q))
	}

	return ret
}
