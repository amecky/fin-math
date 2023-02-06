package math

// https://www.socscistatistics.com/tests/regression/default.aspx
func SimpleLinearRegression(m *Matrix, start, end, field int) (float64, float64) {
	xv := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(xv, float64(i)+1.0)
	}
	period := end - start
	xm := SMA(m, period, xv)
	ym := SMA(m, period, field)
	dx := m.AddColumn()

	for j := start; j < end; j++ {
		m.DataRows[j].Set(dx, m.DataRows[j].Get(xv)-m.DataRows[end-1].Get(xm))
	}

	dy := m.AddColumn()
	for j := start; j < end; j++ {
		m.DataRows[j].Set(dy, m.DataRows[j].Get(field)-m.DataRows[end-1].Get(ym))
	}
	sxy := 0.0
	for i := start; i < end; i++ {
		sxy += m.DataRows[i].Get(dx) * m.DataRows[i].Get(dy)
	}
	sxx := 0.0
	for i := start; i < end; i++ {
		sxx += m.DataRows[i].Get(dx) * m.DataRows[i].Get(dx)
	}
	rm := 0.0
	rc := 0.0
	if sxx != 0.0 {
		rm = sxy / sxx
	}
	rc = m.DataRows[end-1].Get(ym) - m.DataRows[end-1].Get(xm)*rm
	m.RemoveColumns(5)
	return rm, rc
}
