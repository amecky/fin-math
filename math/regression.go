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

func LinearRegressionLSE(m *Matrix, start, end, field int) (float64, float64) {
	q := end - start
	if q == 0 {
		return 0.0, 0.0
	}
	p := float64(q)
	sum_x, sum_y, sum_xx, sum_xy := 0.0, 0.0, 0.0, 0.0
	pX := 1.0
	for i := start; i < end; i++ {
		p := m.DataRows[i]
		sum_x += pX
		sum_y += p.Get(field)
		sum_xx += pX * pX
		sum_xy += pX * p.Get(field)
		pX += 1.0
	}

	lm := (p*sum_xy - sum_x*sum_y) / (p*sum_xx - sum_x*sum_x)
	lb := (sum_y / p) - (lm * sum_x / p)
	/*
	   for i, p := range series {
	           r[i] = Point{p.X, (p.X*m + b)}
	       }
	*/
	return lm, lb
}
