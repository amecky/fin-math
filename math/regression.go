package math

import ma "math"

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

func CorrelationCoefficient(candles *Matrix, first, second, period int) int {
	n := float64(period)
	ret := candles.AddColumn()
	for j := period; j < candles.Rows; j++ {
		sum_X := 0.0
		sum_Y := 0.0
		sum_XY := 0.0
		squareSum_X := 0.0
		squareSum_Y := 0.0
		for i := 0; i < period; i++ {
			// sum of elements of array X.
			sum_X = sum_X + candles.DataRows[i+j-period].Get(first)

			// sum of elements of array Y.
			sum_Y = sum_Y + candles.DataRows[i+j-period].Get(second)

			// sum of X[i] * Y[i].
			sum_XY = sum_XY + candles.DataRows[i+j-period].Get(first)*candles.DataRows[i+j-period].Get(second)

			// sum of square of array elements.
			squareSum_X = squareSum_X + candles.DataRows[i+j-period].Get(first)*candles.DataRows[i+j-period].Get(first)
			squareSum_Y = squareSum_Y + candles.DataRows[i+j-period].Get(second)*candles.DataRows[i+j-period].Get(second)
		}

		// use formula for calculating correlation coefficient.
		corr := float64((n*sum_XY - sum_X*sum_Y)) /
			(ma.Sqrt(float64((n*squareSum_X - sum_X*sum_X) * (n*squareSum_Y - sum_Y*sum_Y))))
		candles.DataRows[j].Set(ret, corr)
	}
	return ret

}

/*
func Correlation(candles *Matrix, first, second, period int) []float64 {

	outReal := make([]float64, len(inReal0))

	inTimePeriodF := float64(period)
	lookbackTotal := period - 1
	startIdx := lookbackTotal
	trailingIdx := startIdx - lookbackTotal
	sumXY, sumX, sumY, sumX2, sumY2 := 0.0, 0.0, 0.0, 0.0, 0.0
	today := trailingIdx
	for today = trailingIdx; today <= startIdx; today++ {
		x := candles.DataRows[today].Get(first)
		sumX += x
		sumX2 += x * x
		y := candles.DataRows[today].Get(second)
		sumXY += x * y
		sumY += y
		sumY2 += y * y
	}
	trailingX := candles.DataRows[trailingIdx].Get(first)
	trailingY := candles.DataRows[trailingIdx].Get(second)
	trailingIdx++
	tempReal := (sumX2 - ((sumX * sumX) / inTimePeriodF)) * (sumY2 - ((sumY * sumY) / inTimePeriodF))
	if !(tempReal < 0.00000000000001) {
		outReal[period-1] = (sumXY - ((sumX * sumY) / inTimePeriodF)) / math.Sqrt(tempReal)
	} else {
		outReal[period-1] = 0.0
	}
	outIdx := period
	for today < len(inReal0) {
		sumX -= trailingX
		sumX2 -= trailingX * trailingX
		sumXY -= trailingX * trailingY
		sumY -= trailingY
		sumY2 -= trailingY * trailingY
		x := inReal0[today]
		sumX += x
		sumX2 += x * x
		y := inReal1[today]
		today++
		sumXY += x * y
		sumY += y
		sumY2 += y * y
		trailingX = inReal0[trailingIdx]
		trailingY = inReal1[trailingIdx]
		trailingIdx++
		tempReal = (sumX2 - ((sumX * sumX) / inTimePeriodF)) * (sumY2 - ((sumY * sumY) / inTimePeriodF))
		if !(tempReal < (0.00000000000001)) {
			outReal[outIdx] = (sumXY - ((sumX * sumY) / inTimePeriodF)) / math.Sqrt(tempReal)
		} else {
			outReal[outIdx] = 0.0
		}
		outIdx++
	}
	return outReal
}
*/
