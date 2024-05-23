package math

import (
	"fmt"
	"strings"
)

func (m *MatrixRow) Set(index int, value float64) *MatrixRow {
	if m != nil {
		if index < m.Num {
			m.Values[index] = value
		}
		return m
	}
	return nil
}

func (m *MatrixRow) SetComment(cmt string) *MatrixRow {
	if m != nil {
		m.Comment = cmt
		return m
	}
	return nil
}

func (m *MatrixRow) Get(index int) float64 {
	if index >= 0 && index < m.Num {
		return m.Values[index]
	}
	return 0.0
}

func (m *MatrixRow) High() float64 {
	return m.Get(1)
}

func (m *MatrixRow) Low() float64 {
	return m.Get(2)
}

func (m *MatrixRow) Open() float64 {
	return m.Get(0)
}

func (m *MatrixRow) Close() float64 {
	return m.Get(4)
}

func (m *MatrixRow) IsGreen() bool {
	return m.Get(4) > m.Get(0)
}

func (m *MatrixRow) IsRed() bool {
	return m.Get(4) < m.Get(0)
}

func (m *MatrixRow) IsInRange(other MatrixRow) bool {
	if m.Close() > other.Low() && m.Close() < other.High() {
		return true
	}
	return false
}
func (m *MatrixRow) ShortKey() string {
	idx := strings.Index(m.Key, " ")
	if idx != -1 {
		return m.Key[0:idx]
	}
	return m.Key
}

func (m *MatrixRow) Sum() float64 {
	sum := 0.0
	for i := 0; i < m.Num; i++ {
		sum += m.Values[i]
	}
	return sum
}

func (m *MatrixRow) PartialSum(start, count int) float64 {
	sum := 0.0
	if start >= m.Num {
		return 0.0
	}
	end := start + count
	if end > m.Num {
		end = m.Num
	}
	for i := start; i < end; i++ {
		sum += m.Values[i]
	}
	return sum
}

func (mr MatrixRow) String() string {
	var builder strings.Builder
	builder.WriteString("Key: ")
	builder.WriteString(mr.Key)
	for _, n := range mr.Values {
		builder.WriteString(fmt.Sprintf(" %.2f", n))
	}
	builder.WriteString(" C: ")
	builder.WriteString(mr.Comment)
	return builder.String()
}
