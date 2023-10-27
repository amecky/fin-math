package math

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
	for i := period; i < m.Rows; i++ {
		l := m.DataRows[i].Get(field)
		for j := 1; j < period; j++ {
			if m.DataRows[i-j].Get(field) < l {
				l = m.DataRows[i-j].Get(field)
			}
		}
		m.DataRows[i].Set(i, l)
	}
	return ret
}

func Highest(m *Matrix, period, field int) int {
	ret := m.AddNamedColumn("Highest")
	for i := period; i < m.Rows; i++ {
		h := m.DataRows[i].Get(field)
		for j := 1; j < period; j++ {
			if m.DataRows[i-j].Get(field) > h {
				h = m.DataRows[i-j].Get(field)
			}
		}
		m.DataRows[i].Set(i, h)
	}
	return ret
}
