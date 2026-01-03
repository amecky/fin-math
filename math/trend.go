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

// https://wire.insiderfinance.io/trading-tip-practical-way-to-know-when-the-markets-out-of-fuel-0be834979f25
func StableATR(m *Matrix, period int) float64 {
	// Step 1: Compute True Range (TR)
	tr := TrueRange(m)
	// Step 2: Take last N TR values
	mn := m.DataRows[m.Rows-period].Get(tr)
	for i := m.Rows - period + 1; i < m.Rows; i++ {
		if m.DataRows[i].Get(tr) < mn {
			mn = m.DataRows[i].Get(tr)
		}
	}
	// Step 3: Outlier check
	cnt := 0
	for i := m.Rows - period; i < m.Rows; i++ {
		if m.DataRows[i].Get(tr) <= 2.0*mn {
			cnt++
		}
	}
	if cnt < 3 {
		return 0.0
	}

	// Step 4: Stability filter
	//SET median_tr = MEDIAN(last_tr)
	md := SMA(m, period, tr)
	median := m.Last().Get(md)
	filtered := make([]float64, 0)
	for i := m.Rows - period; i < m.Rows; i++ {
		if m.DataRows[i].Get(tr) >= 0.5*median && m.DataRows[i].Get(tr) <= 1.5*median {
			filtered = append(filtered, m.DataRows[i].Get(tr))
		}
	}
	//SET filtered = values in last_tr WHERE value BETWEEN 0.5 * median_tr AND 1.5 * median_tr
	//IF LENGTH(filtered) < 3:
	//    RETURN None
	if len(filtered) < 3 {
		return 0.0
	}
	// Step 5: Compute average of filtered
	val := 0.0
	for _, v := range filtered {
		val += v
	}
	return val / float64(len(filtered))

}

type TrendAggregate struct {
	Start     int
	End       int
	Direction int
	Count     int
}

func myTrend(m *Matrix, field int) int {
	// 0 = Trend
	ret := m.AddNamedColumn("Trend")
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		p := m.DataRows[i-1]
		val := p.Get(ret)
		if c.IsRed() {
			if p.Get(ret) < 0.0 {
				val -= 1.0
			} else {
				val = -1.0
			}
		}
		if c.IsGreen() {
			if p.Get(ret) > 0.0 {
				val += 1.0
			} else {
				val = 1.0
			}
		}
		m.DataRows[i].Set(ret, val)
	}
	return ret
}

func AggregateTrend(m *Matrix) []TrendAggregate {
	trend := myTrend(m, 4)
	start := 0
	ret := make([]TrendAggregate, 1)
	for i := 2; i < m.Rows; i++ {
		c := m.DataRows[i].Get(trend)
		p := m.DataRows[i-1].Get(trend)

		if c > 0.0 && p < 0.0 {
			ret = append(ret, TrendAggregate{
				Start:     start,
				End:       i - 1,
				Direction: -1,
				Count:     int(ma.Abs(p)),
			})
			start = i
		}
		if c < 0.0 && p > 0.0 {
			ret = append(ret, TrendAggregate{
				Start:     start,
				End:       i - 1,
				Direction: 1,
				Count:     int(ma.Abs(p)),
			})
			start = i
		}
	}
	if start != m.Rows-1 {
		dir := 1
		if m.Last().Get(trend) < 0.0 {
			dir = -1
		}
		ret = append(ret, TrendAggregate{
			Start:     start,
			End:       m.Rows - 1,
			Direction: dir,
			Count:     int(ma.Abs(m.Last().Get(trend))),
		})
	}

	return ret
}

func PercentageChange(m *Matrix, field int) int {
	// 0 = change percentage
	ret := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i].Get(field)
		p := m.DataRows[i-1].Get(field)
		m.DataRows[i].Set(ret, (c-p)/p*100.0)
	}
	return ret
}

// calculateReturns computes percentage change between consecutive closes.
func calculateReturns(m *Matrix) int {
	ret := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i-1].Set(ret, (m.DataRows[i].Close()-m.DataRows[i-1].Close())/m.DataRows[i-1].Close())
	}
	return ret
}

// internalShannonEntropy calculates Shannon entropy (in bits) from float64 values.
// It bins the values into a histogram and then computes Σ -p * log2(p)
func internalShannonEntropy(values []float64, bins int) float64 {
	if len(values) == 0 {
		return 0
	}

	// Find min and max
	minVal, maxVal := values[0], values[0]
	for _, v := range values {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	if minVal == maxVal {
		return 0 // no variation
	}

	// Build histogram
	hist := make([]int, bins)
	for _, v := range values {
		idx := int((v - minVal) / (maxVal - minVal) * float64(bins-1))
		hist[idx]++
	}

	// Convert to probabilities
	total := float64(len(values))
	entropy := 0.0
	for _, count := range hist {
		if count == 0 {
			continue
		}
		p := float64(count) / total
		entropy -= p * ma.Log2(p)
	}
	return entropy
}

func ShannonEntropy(m *Matrix, windowSize int, bins int) int {
	returns := calculateReturns(m)
	entropies := m.AddNamedColumn("Shannon Entropy")
	for i := 0; i < m.Rows-windowSize; i++ {
		window := m.GetPartialColumn(returns, i, i+windowSize)
		m.DataRows[i+windowSize].Set(entropies, internalShannonEntropy(window, bins))
	}
	return entropies
}

func ConvertRBD(v float64) string {
	types := []string{"RBD", "RBR", "DBD", "DBR"}
	it := max(min(int(v)-1, 3), 0)
	return types[it]
}

const (
	RDB_R   = 1.0
	RDB_D   = 2.0
	RDB_B   = 3.0
	RDB_RBD = 1.0
	RDB_RBR = 2.0
	RDB_DBD = 3.0
	RDB_DBR = 4.0
)

// https://www.youtube.com/watch?v=4Jwq2SALZKA
func RBD(m *Matrix) int {
	// 0 = RBD Type
	tmp := m.AddNamedColumn("RBDType")
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		p := m.DataRows[i-1]
		gc := c.Close() > c.Open()
		rc := c.Close() < c.Open()
		pgc := p.Close() > p.Open()
		prc := p.Close() < p.Open()

		rally := gc && pgc
		drop := rc && prc
		base := !rally && !drop

		tp := 0.0
		if rally {
			tp = RDB_R
		}
		if drop {
			tp = RDB_D
		}
		if base {
			tp = RDB_B
		}
		m.DataRows[i].Set(tmp, tp)
	}
	ret := m.AddNamedColumn("RBD")
	for i := 2; i < m.Rows; i++ {
		c := m.DataRows[i].Get(tmp)
		p1 := m.DataRows[i-1].Get(tmp)
		p2 := m.DataRows[i-2].Get(tmp)
		if p2 == RDB_R && p1 == RDB_B && c == RDB_D {
			m.DataRows[i].Set(ret, RDB_RBD)
		}
		if p2 == RDB_R && p1 == RDB_B && c == RDB_R {
			m.DataRows[i].Set(ret, RDB_RBR)
		}
		if p2 == RDB_D && p1 == RDB_B && c == RDB_B {
			m.DataRows[i].Set(ret, RDB_DBD)
		}
		if p2 == RDB_D && p1 == RDB_B && c == RDB_R {
			m.DataRows[i].Set(ret, RDB_DBR)
		}

	}
	return ret
}

func NRX(m *Matrix, period int) int {
	// 0 = 1 if NR
	ret := m.AddColumn()
	rng := Range(m)
	for i := period; i < m.Rows; i++ {
		cur := m.DataRows[i].Get(rng)
		cnt := 0
		for j := range period - 1 {
			cmp := m.DataRows[i-j-1].Get(rng)
			if cmp > cur {
				cnt++
			}
		}
		if cnt == period-1 {
			m.DataRows[i].Set(ret, 1.0)
		}
	}
	return ret
}

func HLBand(m *Matrix, period int) int {
	// 0 = Upper 1 = Lower
	upper := m.AddColumn()
	lower := m.AddColumn()
	high := SMA(m, period, 4)
	low := SMA(m, period, 4)
	hlAvg := m.Apply(func(mr MatrixRow) float64 {
		return mr.Get(high) - mr.Get(low)
	})
	for i := range m.Rows {
		c := m.DataRows[i]
		m.DataRows[i].Set(upper, c.Get(high))
		m.DataRows[i].Set(lower, c.Get(high)-c.Get(hlAvg))
	}
	//data['HL_avg'] = data['High'].rolling(window=25).mean() - data['Low'].rolling(window=25).mean()
	//data['Band'] = data['High'].rolling(window=25).mean() - (data['HL_avg'] * 2.25)
	m.RemoveColumns(2)
	return upper
}

/*
Calculate HLC3 (average of high, low, close prices)
Apply 9-period EMA to HLC3
Calculate deviation from the EMA
Create Fast Wave using 12-period EMA
Create Slow Wave as 3-period moving average of Fast Wave
Calculate Delta as the difference between Fast and Slow waves
*/
func Wave(m *Matrix, period, fast, slow int) int {
	// 0 = Fast Wave 1 = Slow Wave 2 = Delta
	fc := m.AddColumn()
	sc := m.AddColumn()
	delta := m.AddColumn()
	hlc := HLC3(m)
	eh := EMA(m, period, hlc)
	std := m.StdDev(eh, period)
	fe := EMA(m, fast, std)
	se := EMA(m, slow, fe)
	deltaWaves := m.Apply(func(mr MatrixRow) float64 {
		return mr.Get(fe) - mr.Get(se)
	})
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(fc, m.DataRows[i].Get(fe))
		m.DataRows[i].Set(sc, m.DataRows[i].Get(se))
		m.DataRows[i].Set(delta, m.DataRows[i].Get(deltaWaves))
	}
	m.RemoveColumns(6)
	return fc
}

func Correlation(m *Matrix, period, first, second int) int {
	ret := m.AddColumn()
	n := float64(period)
	var sumX, sumY, sumXY, sumX2, sumY2 float64
	for i := period; i < m.Rows; i++ {
		sumX, sumY, sumXY, sumX2, sumY2 = 0.0, 0.0, 0.0, 0.0, 0.0
		for j := range period {
			idx := i - period + j
			sumX += m.DataRows[idx].Get(first)
			sumY += m.DataRows[idx].Get(second)
			sumXY += m.DataRows[idx].Get(first) * m.DataRows[idx].Get(second)
			sumX2 += m.DataRows[idx].Get(first) * m.DataRows[idx].Get(first)
			sumY2 += m.DataRows[idx].Get(second) * m.DataRows[idx].Get(second)
		}
		num := n*sumXY - sumX*sumY
		den := ma.Sqrt((n*sumX2 - sumX*sumX) * (n*sumY2 - sumY*sumY))
		if den == 0 {
			return 0
		}
		m.DataRows[i].Set(ret, num/den)
	}

	return ret
}

func StretchMove(m *Matrix, atrPeriod, lookback int) int {
	// 0 = Low StretchMove 1 = High StretchMove
	lo := m.AddColumn()
	hi := m.AddColumn()
	atr := ATR(m, atrPeriod)
	for i := lookback; i < m.Rows; i++ {
		// low
		idx := m.FindLowestIndex(i-lookback, lookback)
		cur := m.DataRows[i].Close()
		delta := ma.Abs(m.DataRows[idx].Low() - cur)
		m.DataRows[i].Set(lo, delta/m.DataRows[i].Get(atr))
		// high
		idx = m.FindHighestIndex(i-lookback, lookback)
		delta = ma.Abs(m.DataRows[idx].High() - cur)
		m.DataRows[i].Set(hi, delta/m.DataRows[i].Get(atr))
	}
	return lo
}

func ATRRegime(m *Matrix, atrPeriod, lookback int) int {
	// 0 = ATR Regime
	ret := m.AddNamedColumn("ATRRegime")
	atr := ATR(m, atrPeriod)
	// 	vol_threshold = df['ATR'].rolling(100).quantile(0.7)
	q := RollingQuantile(m, atr, lookback, 0.7)
	//df['regime'] = (df['ATR'] < vol_threshold).astype(int)
	for i := range m.Rows {
		c := m.DataRows[i]
		if c.Get(atr) >= c.Get(q) {
			m.DataRows[i].Set(ret, 1.0)
		}
	}
	return ret
}
