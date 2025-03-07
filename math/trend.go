package math

import (
	m "math"
	ma "math"
)

func HLCMAD(candles *Matrix, period int) int {
	// 0 = DL 1 = DC 2 = DH
	ret := candles.AddNamedColumn("DL")
	rc := candles.AddNamedColumn("DC")
	rh := candles.AddNamedColumn("DH")
	ei := EMA(candles, period, 4)
	candles.ApplyRow(ret, func(mr MatrixRow) float64 { return mr.Get(2) - mr.Get(ei) })
	candles.ApplyRow(rc, func(mr MatrixRow) float64 { return mr.Get(4) - mr.Get(ei) })
	candles.ApplyRow(rh, func(mr MatrixRow) float64 { return mr.Get(1) - mr.Get(ei) })
	candles.RemoveColumns(1)
	return ret
}

// https://www.investopedia.com/terms/p/pvi.asp
func PVI(candles *Matrix, period int) int {
	// 0 = PVI 1 = Signal
	ret := candles.AddNamedColumn("PVI")
	candles.DataRows[0].Set(ret, (candles.DataRows[1].Get(4)-candles.DataRows[0].Get(4))/candles.DataRows[0].Get(4))
	for i := 1; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		p := candles.DataRows[i-1]
		if c.Get(5) > p.Get(5) {
			c.Set(ret, (p.Get(ret) + ((c.Get(4)-p.Get(4))/p.Get(4))*p.Get(ret)))
		} else {
			c.Set(ret, p.Get(ret))
		}
	}
	EMA(candles, ret, period)
	return ret
}

// -----------------------------------------------------------------------
// Elder Ray Index
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/the-elder-ray-index-for-trading-b54c9b1741aa
func ElderRayIndex(candles *Matrix, period int) int {
	// 0 = Bull Power 1 = Bear Power
	bup := candles.AddNamedColumn("BUP")
	bep := candles.AddNamedColumn("BEP")
	ei := EMA(candles, period, ADJ_CLOSE)
	for i := 0; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		c.Set(bup, c.Get(1)-c.Get(ei))
		c.Set(bep, c.Get(2)-c.Get(ei))
	}
	candles.RemoveColumns(1)
	return bup
}

func f(x, a float64) float64 {

	return m.Exp(x) - a
}

func ln(n float64) float64 {

	var lo, hi, ma float64

	if n <= 0.0 {

		return -1
	}

	if n == 1 {

		return 0
	}

	EPS := 0.00001

	lo = 0.0

	hi = n

	for m.Abs(lo-hi) >= EPS {

		ma = float64((lo + hi) / 2.0)

		if f(ma, n) < 0 {

			lo = ma

		} else {

			hi = ma
		}
	}

	return float64((lo + hi) / 2.0)
}

func ParkinsonEstimator(candles *Matrix, period int) int {
	ret := candles.AddNamedColumn("PE")
	li := candles.Apply(func(mr MatrixRow) float64 {
		d := ln(mr.Get(1) / mr.Get(2))
		return d * d
	})
	e := 1.0 / (4.0 * float64(period) * m.Ln2)
	si := candles.Sum(li, period)
	for i := 0; i < candles.Rows; i++ {
		s := e * candles.DataRows[i].Get(si)
		candles.DataRows[i].Set(ret, m.Sqrt(s))
	}
	candles.RemoveColumns(1)
	return ret
}

func EMATrend(mat *Matrix, emas ...int) int {
	ret := mat.AddColumn()
	total := len(emas)
	eis := make([]int, 0)
	for _, e := range emas {
		eis = append(eis, EMA(mat, e, 4))
	}
	recent := emas[len(emas)-1]
	for i := recent + 1; i < mat.Rows; i++ {
		sum := 0.0
		div := 0.0
		// close above / below
		for k := 0; k < total; k++ {
			if mat.DataRows[i].Get(4) > mat.DataRows[i].Get(eis[k]) {
				sum += 1.0
			}
			div += 1.0
		}
		// rising
		for k := 0; k < total; k++ {
			if mat.DataRows[i-1].Get(eis[k]) < mat.DataRows[i].Get(eis[k]) {
				sum += 1.0
			}
			div += 1.0
		}
		// above each other
		for k := 1; k < total; k++ {
			if mat.DataRows[i].Get(eis[k-1]) > mat.DataRows[i].Get(eis[k]) {
				sum += 1.0
			}
			div += 1.0
		}
		mat.DataRows[i].Set(ret, (sum/div*2.0 - 1.0))
	}
	mat.RemoveColumns(total - 1)
	return ret
}

// should only be used on daily data
func Pivots(candles *Matrix) int {
	// 0 = PP 1 = R1 2 = R2 3 = S1 4 = S2
	ppi := candles.AddNamedColumn("PP")
	r1i := candles.AddNamedColumn("R1")
	r2i := candles.AddNamedColumn("R2")
	s1i := candles.AddNamedColumn("S1")
	s2i := candles.AddNamedColumn("S2")

	for i := 0; i < candles.Rows-1; i++ {
		cur := &candles.DataRows[i]
		n := &candles.DataRows[i+1]
		pp := (cur.Get(1) + cur.Get(2) + cur.Get(4)) / 3.0
		n.Set(ppi, pp)
		n.Set(s1i, pp*2.0-cur.Get(1))
		n.Set(r1i, pp*2.0-cur.Get(2))
		n.Set(r2i, pp+(cur.Get(1)-cur.Get(2)))
		n.Set(s2i, pp-(cur.Get(1)-cur.Get(2)))
	}
	return ppi
}

func HHLL(candles *Matrix) int {
	// 0 = Higher Highs 1 = Lower Lows
	ret := candles.AddNamedColumn("HH")
	lr := candles.AddNamedColumn("LL")
	for i := 1; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		p := candles.DataRows[i-1]
		if c.Get(1) > p.Get(1) {
			c.Set(ret, p.Get(ret)+1)
		} else {
			c.Set(ret, 0)
		}
		if c.Get(2) < p.Get(2) {
			c.Set(lr, p.Get(lr)+1)
		} else {
			c.Set(lr, 0)
		}
	}
	return ret
}

func UpDown(m *Matrix, days int) int {
	// 0 = Up 1 = Down 2 = Sum
	upCol := m.AddColumn()
	downCol := m.AddColumn()
	sumCol := m.AddColumn()
	delta := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		p := m.DataRows[i-1].Close()
		c := m.DataRows[i].Close()
		m.DataRows[i].Set(delta, c-p)
	}
	for i := 0; i < m.Rows; i++ {
		up := 0
		down := 0
		sum := 0.0
		end := i - days + 1
		if end < 0 {
			end = 0
		}
		for j := i; j >= end; j-- {
			c := m.DataRows[j]
			sum += m.DataRows[j].Get(delta)
			if c.IsGreen() {
				up++
			} else {
				down++
			}
		}
		m.DataRows[i].Set(upCol, float64(up))
		m.DataRows[i].Set(downCol, float64(down))
		m.DataRows[i].Set(sumCol, float64(sum))
	}
	m.RemoveColumns(1)
	return upCol
}

func SavGolFilter(m *Matrix, windowSize, polyOrder int) int {
	//windowSize := 5 // Must be odd
	//polyOrder := 2  // Polynomial degree
	coeffs := ComputeSavGolCoefficients(windowSize, polyOrder)

	// Apply the filter
	return applySavGolFilter(m, coeffs)
}

// ComputeSavGolCoefficients calculates the Savitzky-Golay coefficients
func ComputeSavGolCoefficients(windowSize, polyOrder int) []float64 {
	if windowSize%2 == 0 {
		panic("Window size must be odd.")
	}
	if polyOrder >= windowSize {
		panic("Polynomial order must be less than window size.")
	}

	// Define matrix A (Vandermonde matrix)
	halfWindow := windowSize / 2
	A := make([][]float64, windowSize)
	for i := -halfWindow; i <= halfWindow; i++ {
		row := make([]float64, polyOrder+1)
		for j := 0; j <= polyOrder; j++ {
			row[j] = ma.Pow(float64(i), float64(j))
		}
		A[i+halfWindow] = row
	}

	// Compute A^T * A
	ATA := make([][]float64, polyOrder+1)
	for i := range ATA {
		ATA[i] = make([]float64, polyOrder+1)
		for j := range ATA[i] {
			for k := 0; k < windowSize; k++ {
				ATA[i][j] += A[k][i] * A[k][j]
			}
		}
	}

	// Invert ATA (for small polyOrder this is feasible)
	ATAInv := invertMatrix(ATA)

	// Compute the filter coefficients
	AT := transposeMatrix(A)
	coeffs := make([]float64, windowSize)
	for i := 0; i < windowSize; i++ {
		for j := 0; j <= polyOrder; j++ {
			coeffs[i] += AT[0][j] * ATAInv[j][0] // Only first row corresponds to smoothing
		}
	}

	return coeffs
}

// ApplySavGolFilter applies the Savitzky-Golay filter to the data using the given coefficients
func applySavGolFilter(m *Matrix, coeffs []float64) int {
	windowSize := len(coeffs)
	halfWindow := windowSize / 2
	ret := m.AddColumn()
	//smoothed := make([]float64, len(data))

	for i := range m.DataRows {
		var sum float64
		for j := -halfWindow; j <= halfWindow; j++ {
			idx := i + j
			if idx >= 0 && idx < m.Rows {
				sum += m.DataRows[idx].Get(4) * coeffs[j+halfWindow]
			}
		}
		m.DataRows[i].Set(ret, sum)
	}
	return ret
}

// Utility: Transpose a matrix
func transposeMatrix(matrix [][]float64) [][]float64 {
	m := len(matrix)
	n := len(matrix[0])
	transposed := make([][]float64, n)
	for i := range transposed {
		transposed[i] = make([]float64, m)
		for j := range transposed[i] {
			transposed[i][j] = matrix[j][i]
		}
	}
	return transposed
}

// Utility: Invert a small matrix using Gaussian elimination
func invertMatrix(matrix [][]float64) [][]float64 {
	n := len(matrix)
	inv := make([][]float64, n)
	for i := range inv {
		inv[i] = make([]float64, n)
		inv[i][i] = 1
	}
	for i := 0; i < n; i++ {
		// Normalize the row
		diag := matrix[i][i]
		for j := 0; j < n; j++ {
			matrix[i][j] /= diag
			inv[i][j] /= diag
		}
		// Eliminate other rows
		for k := 0; k < n; k++ {
			if k != i {
				factor := matrix[k][i]
				for j := 0; j < n; j++ {
					matrix[k][j] -= factor * matrix[i][j]
					inv[k][j] -= factor * inv[i][j]
				}
			}
		}
	}
	return inv
}
