package math

import (
	"fmt"
	m "math"
	"strconv"
	"strings"
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

func Identical(a, b, epsilon float64) bool {
	if a == b {
		return true
	}
	d := m.Abs(a - b)
	if b == 0 {
		return d < epsilon
	}
	return (d / m.Abs(b)) < epsilon
}

func Copy(src *Matrix) *Matrix {
	m := Matrix{
		Cols:    1,
		Headers: []string{"Key"},
	}
	for i, h := range src.Headers {
		if i > 0 {
			m.Headers = append(m.Headers, h)
			m.Cols++
		}
	}
	return &m
}

func Aggregate(m *Matrix, intervall string) *Matrix {
	steps := 5
	offset := 0
	key := m.DataRows[0].Key
	fmt.Println("first", key)
	if intervall == "5m" {
		steps = 5
		dt := key[strings.Index(key, ":")+1:]
		iv, _ := strconv.Atoi(dt)
		offset = steps - iv%steps
		if offset >= steps {
			offset = 0
		}
		fmt.Println("=> dt", dt, "iv", iv, "offset", offset)
	}
	if intervall == "15m" {
		steps = 15
		dt := key[strings.Index(key, ":")+1:]
		iv, _ := strconv.Atoi(dt)
		offset = steps - iv%steps
	}
	ret := Copy(m)
	s := m.Rows / steps
	fmt.Println("rows", m.Rows, "steps", steps, "s", s)
	for i := 0; i < s; i++ {
		idx := i*steps + offset
		fmt.Println("idx", idx, "key:", m.DataRows[idx])
		if idx+steps < m.Rows {
			row := ret.AddRow(m.DataRows[idx].Key)
			row.Set(0, m.DataRows[idx].Open())
			row.Set(4, m.DataRows[idx+steps-1].Close())
			row.Set(3, m.DataRows[idx+steps-1].Close())
			h := m.DataRows[idx].High()
			l := m.DataRows[idx].Low()
			v := 0.0
			for j := 0; j < steps; j++ {
				v += m.DataRows[idx+j].Get(5)
				if i > 0 {
					if m.DataRows[idx+j].High() > h {
						h = m.DataRows[idx+j].High()
					}
					if m.DataRows[idx+j].Low() < l {
						l = m.DataRows[idx+j].Low()
					}
				}
			}
			row.Set(1, h)
			row.Set(2, l)
			row.Set(5, v)
			cols := len(m.DataRows[idx].Values)
			if cols >= 6 {
				for j := 6; j < cols; j++ {
					row.Set(j, m.DataRows[idx].Get(j))
				}
			}
		}
	}
	return ret
}

func Subset(start, end int, m *Matrix) *Matrix {
	ret := Matrix{
		Cols:    1,
		Headers: []string{"Key"},
	}
	for i, h := range m.Headers {
		if i > 0 {
			m.Headers = append(m.Headers, h)
			m.Cols++
		}
	}
	if start < 0 {
		start = 0
	}
	if end > m.Rows {
		end = m.Rows
	}
	for i := start; i < end; i++ {
		ret.DataRows = append(ret.DataRows, m.DataRows[i])
		ret.Rows++
	}
	return &ret
}

func CalculateRow(m *Matrix, fn func(col, index int, m *Matrix)) int {
	ret := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		fn(ret, i, m)
	}
	return ret
}
