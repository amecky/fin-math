package math

import (
	"fmt"
	"math"
	m "math"
)

// Normalizes Value into -1 / 1 range
func NormalizeRange(v, min, max float64) float64 {
	d := 0.0
	if v > max {
		v = max
	}
	if v < min {
		v = min
	}
	if (max - min) != 0.0 {
		d = (v-min)*2.0/(max-min) - 1
	}
	return d
}

func NormalizeRange2(v, min, max float64) float64 {
	d := 0.0
	if v > max {
		v = max
	}
	if v < min {
		v = min
	}
	if (max - min) != 0.0 {
		d = (v - min) / (max - min)
	}
	return d
}

func CategorizeRelation(u, v, min, max float64) (float64, int) {
	de := 0.0
	if v != 0.0 {
		de = (u/v - 1.0) * 100.0
	}
	return CategorizeNormalizeRange(de, min, max)
}

func CategorizeNormalizeRange(v, min, max float64) (float64, int) {
	d := NormalizeRange(v, min, max)
	ret := 0
	if d == -1 {
		ret = 1
	} else if d <= -0.75 {
		ret = 2
	} else if d <= -0.5 {
		ret = 3
	} else if d <= -0.25 {
		ret = 4
	} else if d == 0.0 {
		ret = 5
	} else if d <= 0.25 {
		ret = 6
	} else if d <= 0.5 {
		ret = 7
	} else if d <= 0.75 {
		ret = 8
	} else {
		ret = 9
	}
	return d, ret
}

// Categorizes value based on 4 values and returns the pocket
func CategorizeFiveSteps(v, s1, s2, s3, s4 float64) int {
	if v <= s1 {
		return 0
	}
	if v <= s2 {
		return 1
	}
	if v <= s3 {
		return 2
	}
	if v <= s4 {
		return 3
	}
	return 4
}

func ChangePercentage(first, second float64) float64 {
	if second != 0.0 {
		return (first/second - 1.0) * 100.0
	}
	return 0.0
}

type MAFunc func(m *Matrix, days, field int) int

func MAFuncByName(name string) MAFunc {
	mf := SMA
	if name == "EMA" {
		mf = EMA
	}
	if name == "HMA" {
		mf = HMA
	}
	if name == "WMA" {
		mf = WMA
	}
	if name == "RMA" {
		mf = RMA
	}
	if name == "DEMA" {
		mf = DEMA
	}
	if name == "TEMA" {
		mf = TEMA
	}
	return mf
}

// -----------------------------------------------------------------------
// SMA - Simple moving average
// -----------------------------------------------------------------------
func SMA(m *Matrix, days, field int) int {
	ret := m.AddNamedColumn(fmt.Sprintf("SMA%d", days))
	n := float64(days)
	for i := days - 1; i < m.Rows; i++ {
		sum := 0.0
		for j := 0; j < days; j++ {
			idx := i - days + 1 + j
			sum += m.DataRows[idx].Get(field)
		}
		avg := sum / n
		m.DataRows[i].Set(ret, avg)
	}
	return ret
}

// -----------------------------------------------------------------------
// SWMA - Simple weighted moving average (P1 + 2*P2 + 2*P3 + P4) / 6
// -----------------------------------------------------------------------
func SWMA(m *Matrix, field int) int {
	ret := m.AddColumn()
	for i := 3; i < m.Rows; i++ {
		sum := (m.DataRows[i-3].Get(field) + 2.0*m.DataRows[i-2].Get(field) + 2.0*m.DataRows[i-1].Get(field) + m.DataRows[i].Get(field)) / 6.0
		m.DataRows[i].Set(ret, sum)
	}
	return ret
}

// -----------------------------------------------------------------------
// EMA
// -----------------------------------------------------------------------
func EMA(m *Matrix, days, field int) int {
	ret := m.AddNamedColumn(fmt.Sprintf("EMA%d", days))
	if m.Rows > days {
		n := float64(days)
		sma := SMA(m, days, field)
		multiplier := 2.0 / (n + 1)
		m.DataRows[days].Set(ret, m.DataRows[days-1].Get(sma))
		for i := days + 1; i < m.Rows; i++ {
			v := m.DataRows[i-1].Get(ret)*(1.0-multiplier) + m.DataRows[i].Get(field)*multiplier
			m.DataRows[i].Set(ret, v)
		}
		m.RemoveColumn()
	}
	return ret
}

// -----------------------------------------------------------------------
// RMA
// -----------------------------------------------------------------------
func RMA(m *Matrix, days, field int) int {
	ret := m.AddNamedColumn(fmt.Sprintf("RMA%d", days))
	total := m.Rows
	n := float64(days)
	if total >= days {
		sum := 0.0
		for i := 0; i < days; i++ {
			c := m.DataRows[i].Get(field)
			sum += c
		}
		avg := sum / float64(days)
		m.DataRows[days-1].Set(ret, avg)
		prev := avg
		for i := days; i < total; i++ {
			c := m.DataRows[i].Get(field)
			v := (prev*(n-1.0) + c) / n
			m.DataRows[i].Set(ret, v)
			prev = v
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// WMA
// -----------------------------------------------------------------------
func WMA(m *Matrix, days, field int) int {
	ret := m.AddNamedColumn(fmt.Sprintf("WMA%d", days))
	n := float64(days)
	for i := days - 1; i < m.Rows; i++ {
		sum := 0.0
		for j := 0; j < days; j++ {
			idx := i - days + 1 + j
			row := m.DataRows[idx]
			sum += row.Get(field) * float64(j+1)
		}
		avg := sum / (n * (n + 1.0) / 2.0)
		m.DataRows[i].Set(ret, avg)
	}
	return ret
}

// -----------------------------------------------------------------------
// TEMA is calculated as 3*MA - (3*MA(MA)) + (MA(MA(MA)))
// -----------------------------------------------------------------------
func TEMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
	ema1 := EMA(m, days, field)
	ema2 := EMA(m, days, ema1)
	ema3 := EMA(m, days, ema2)
	for i := 0; i < m.Rows; i++ {
		e1 := m.DataRows[i].Get(ema1)
		e2 := m.DataRows[i].Get(ema2)
		e3 := m.DataRows[i].Get(ema3)
		m.DataRows[i].Set(ret, 3.0*e1-3.0*e2+e3)
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// DEMA is calculated as 2*EMA - (EMA(EMA))
// -----------------------------------------------------------------------
// https://www.youtube.com/watch?v=HE6XDux4Ig4
func DEMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
	ema1 := EMA(m, days, field)
	ema2 := EMA(m, days, ema1)
	ema3 := EMA(m, days, ema2)
	for i := 0; i < m.Rows; i++ {
		e1 := m.DataRows[i].Get(ema1)
		e3 := m.DataRows[i].Get(ema3)
		m.DataRows[i].Set(ret, 2.0*e1-e3)
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

func ZLEMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
	lag := (days - 1) / 2
	d := m.AddColumn()
	for i := lag; i < m.Rows; i++ {
		m.DataRows[i].Set(d, 2.0*m.DataRows[i].Get(field)-m.DataRows[i-lag].Get(field))
	}
	ei := EMA(m, days, d)
	m.CopyColumn(ei, ret)
	m.RemoveColumns(2)
	return ret
}

func ZLSMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
	lag := (days - 1) / 2
	d := m.AddColumn()
	for i := lag; i < m.Rows; i++ {
		m.DataRows[i].Set(d, 2.0*m.DataRows[i].Get(field)-m.DataRows[i-lag].Get(field))
	}
	ei := SMA(m, days, d)
	m.CopyColumn(ei, ret)
	m.RemoveColumns(2)
	return ret
}

func MASlope(prices *Matrix, operator MAFunc, days, lookback int) int {
	// 0 = MA Slope
	ret := prices.AddNamedColumn(fmt.Sprintf("MAS%d", days))
	steps := float64(lookback)
	si := operator(prices, days, ADJ_CLOSE)
	for i := lookback; i < prices.Rows; i++ {
		cur := prices.DataRows[i].Get(si)
		prev := prices.DataRows[i-lookback].Get(si)
		prices.DataRows[i].Set(ret, (cur-prev)/steps)
	}
	prices.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Disparity
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/d/disparityindex.asp
//
// A value greater than zero—a positive percentage—shows that the price is rising, suggesting that the asset is gaining upward momentum.
// Conversely, a value less than zero—a negative percentage—can be interpreted as a sign that selling pressure is increasing, forcing the price to drop.
// A value of zero means that the asset’s current price is exactly consistent with its moving average.
func Disparity(prices *Matrix, days int) int {
	// 0 = Disparity
	ret := prices.AddColumn()
	// DI 14 = [ C.PRICE - MOVING  AVG 14 ] / [ MOVING AVG 14 ] * 100
	sma := EMA(prices, days, 4)
	for _, p := range prices.DataRows {
		p.Set(ret, (p.Get(ADJ_CLOSE)-p.Get(sma))/p.Get(sma)*100.0)
	}
	prices.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Awesome Oscilator
// -----------------------------------------------------------------------
// https://de.tradingview.com/scripts/awesomeoscillator/
func AO(m *Matrix, short, long int) int {
	// 0 = AO 1 = Color
	ao := m.AddColumn()
	clr := m.AddColumn()
	mid := HL2(m)
	sma5 := SMA(m, short, mid)
	sma34 := SMA(m, long, mid)
	// AO = sma((high+low)/2, LängeAO1) - sma((high+low)/2, LängeAO2)
	oldHist := 0.0
	for i := 0; i < m.Rows; i++ {
		d := m.DataRows[i].Get(sma5) - m.DataRows[i].Get(sma34)
		if d > oldHist {
			m.DataRows[i].Set(clr, 1.0)
		} else {
			m.DataRows[i].Set(clr, -1.0)
		}
		m.DataRows[i].Set(ao, d)
		oldHist = d
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ao
}

// -----------------------------------------------------------------------
// Linear Regression Line
// -----------------------------------------------------------------------
func LinearRegression(m *Matrix, period int) int {
	// 0 = M 1 = C
	lr := m.AddColumn()
	ci := m.AddColumn()
	for i := period; i < m.Rows; i++ {
		s, c := SimpleLinearRegression(m, i-period, i, 4)
		m.DataRows[i].Set(lr, s)
		m.DataRows[i].Set(ci, c)
	}
	return lr
}

// -----------------------------------------------------------------------
// ACC
// -----------------------------------------------------------------------
// https://forextester.com/blog/accelerator-oscillator
// https://admiralmarkets.com/education/articles/forex-indicators/accelerator-oscillator
func ACC(m *Matrix, short, long, s int) int {
	// 0 = ACC
	ret := m.AddColumn()
	ao := AO(m, short, long)
	sma := SMA(m, s, ao)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(ao)-m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// MACD
// -----------------------------------------------------------------------
func MACD(m *Matrix, short, long, signal int) int {
	// 0 = Line 1 = Signal 2 = Diff
	ret := m.AddNamedColumn("MACD-Line")
	sig := m.AddNamedColumn("MACD-Signal")
	diff := m.AddNamedColumn("MACD-Diff")

	f := EMA(m, short, 4)
	s := EMA(m, long, 4)

	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(f)-m.DataRows[i].Get(s))
	}
	signalPairs := EMA(m, signal, ret)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(sig, m.DataRows[i].Get(signalPairs))
		m.DataRows[i].Set(diff, m.DataRows[i].Get(ret)-m.DataRows[i].Get(signalPairs))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// MACD
// -----------------------------------------------------------------------
// https://cmtassociation.org/wp-content/uploads/2022/05/MACD-V-Volatility-Normalised-Momentum-by-Alex-Spiroglou-DipTA-ATAA-CFTe.pdf
func MACDV(m *Matrix, short, long, signal int) int {
	// 0 = Line 1 = Signal 2 = Diff
	ret := m.AddNamedColumn("MACD-Line")
	sig := m.AddNamedColumn("MACD-Signal")
	diff := m.AddNamedColumn("MACD-Diff")

	f := EMA(m, short, 4)
	s := EMA(m, long, 4)
	a := ATR(m, long)
	m.ApplyRow(ret, func(mr MatrixRow) float64 {
		if mr.Get(a) != 0.0 {
			return (mr.Get(f) - mr.Get(s)) / mr.Get(a) * 100.0
		}
		return 0.0
	})

	signalPairs := EMA(m, signal, ret)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(sig, m.DataRows[i].Get(signalPairs))
		m.DataRows[i].Set(diff, m.DataRows[i].Get(ret)-m.DataRows[i].Get(signalPairs))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// MACD
// -----------------------------------------------------------------------
func MACDZL(m *Matrix, short, long, signal int) int {
	// 0 = Line 1 = Signal 2 = Diff
	ret := m.AddColumn()
	sig := m.AddColumn()
	diff := m.AddColumn()

	f := ZLEMA(m, short, 4)
	s := ZLEMA(m, long, 4)

	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(f)-m.DataRows[i].Get(s))
	}
	signalPairs := EMA(m, signal, ret)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(sig, m.DataRows[i].Get(signalPairs))
		m.DataRows[i].Set(diff, m.DataRows[i].Get(ret)-m.DataRows[i].Get(signalPairs))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

func MACDExt(m *Matrix, field, short, long, signal int) int {
	// 0 = Line 1 = Signal 2 = Diff
	ret := m.AddColumn()
	sig := m.AddColumn()
	diff := m.AddColumn()

	f := EMA(m, short, field)
	s := EMA(m, long, field)

	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(f)-m.DataRows[i].Get(s))
	}
	signalPairs := EMA(m, signal, ret)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(sig, m.DataRows[i].Get(signalPairs))
		m.DataRows[i].Set(diff, m.DataRows[i].Get(ret)-m.DataRows[i].Get(signalPairs))
	}
	m.RemoveColumns(3)
	return ret
}

// -----------------------------------------------------------------------
//
//	Momentum
//
// -----------------------------------------------------------------------
func Momentum(m *Matrix, days, smoothed int) int {
	// 0 = Momentum 1 = Momentum Percentage 2 = EMA Momentum
	ret := m.AddNamedColumn("Momentum")
	per := m.AddNamedColumn("Momentum-Pct")
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(ADJ_CLOSE) - m.DataRows[i-days].Get(ADJ_CLOSE)))
		m.DataRows[i].Set(per, (m.DataRows[i].Get(ADJ_CLOSE)-m.DataRows[i-days].Get(ADJ_CLOSE))/m.DataRows[i-days].Get(ADJ_CLOSE)*100.0)
	}
	EMA(m, smoothed, ret)
	return ret
}

// -----------------------------------------------------------------------
//
//	MomentumExt
//
// -----------------------------------------------------------------------
func MomentumExt(m *Matrix, days, smoothed, field int) int {
	// 0 = Momentum 1 = Momentum Percentage 2 = EMA Momentum
	ret := m.AddNamedColumn("Momentum")
	per := m.AddNamedColumn("Momentum-Pct")
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(field) - m.DataRows[i-days].Get(field)))
		m.DataRows[i].Set(per, (m.DataRows[i].Get(field)-m.DataRows[i-days].Get(field))/m.DataRows[i-days].Get(field)*100.0)
	}
	EMA(m, smoothed, ret)
	return ret
}

// -----------------------------------------------------------------------
//
//	Daily Percentage Change
//
// -----------------------------------------------------------------------
func DPC(m *Matrix) int {
	// 0 = DPC
	ret := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(ADJ_CLOSE)-m.DataRows[i-1].Get(ADJ_CLOSE))/m.DataRows[i-1].Get(ADJ_CLOSE)*100.0)
	}
	return ret
}

// -----------------------------------------------------------------------
// MeanBreakout
// -----------------------------------------------------------------------
// https://medium.com/superalgos/all-in-one-indicator-for-exponential-moving-average-crossover-strategy-b6b4b0da957e
func MeanBreakout(m *Matrix, period int) int {
	// 0 = MBO
	ret := m.AddColumn()
	sma := EMA(m, period, 4)
	for i := period; i < m.Rows; i++ {
		cp := m.DataRows[i]
		cs := m.DataRows[i].Get(sma)
		min, max := m.FindMinMaxBetween(4, i-period, period)
		if max != min {
			d := (cp.Get(ADJ_CLOSE) - cs) / (max - min)
			m.DataRows[i].Set(ret, d)
		}
	}
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Consolidated Price Difference
// -----------------------------------------------------------------------
func ConsolidatedPriceDifference(m *Matrix, min int) int {
	// 0 = CPD
	sub := m.AddColumn()
	ret := m.AddColumn()
	md := 0.0
	for i := 1; i < m.Rows; i++ {
		cp := m.DataRows[i]
		pp := m.DataRows[i-1]
		cnt := 0
		counting := true
		price := cp.Get(ADJ_CLOSE)
		for j := i + 1; j < m.Rows; j++ {
			if counting == true {
				p := m.DataRows[j]
				if p.Get(ADJ_CLOSE) > pp.Get(2) && p.Get(ADJ_CLOSE) < cp.Get(HIGH) {
					cnt++
					price += p.Get(ADJ_CLOSE)
				} else {
					counting = false
				}
			}
		}
		if cnt >= min {
			md = price / float64(cnt+1)
		}
		m.DataRows[i].Set(sub, md)
	}
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(sub) != 0.0 {
			m.DataRows[i].Set(ret, (m.DataRows[i].Get(ADJ_CLOSE)-m.DataRows[i].Get(sub))/m.DataRows[i].Get(sub)*100.0)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
//
//	RSI
//
// -----------------------------------------------------------------------
func RSI(m *Matrix, days, field int) int {
	// 0 = RSI
	ret := m.AddNamedColumn("RSI")
	diff := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(diff, m.DataRows[i].Get(field)-m.DataRows[i-1].Get(field))
	}
	//up = ta.rma(math.max(ta.change(rsiSourceInput), 0), rsiLengthInput)
	//down = ta.rma(-math.min(ta.change(rsiSourceInput), 0), rsiLengthInput)
	//rsi = down == 0 ? 100 : up == 0 ? 0 : 100 - (100 / (1 + up / down))
	cgi := m.Apply(func(mr MatrixRow) float64 {
		return math.Max(mr.Get(diff), 0.0)
	})
	cli := m.Apply(func(mr MatrixRow) float64 {
		return -1.0 * math.Min(mr.Get(diff), 0.0)
	})
	up := RMA(m, days, cgi)
	down := RMA(m, days, cli)

	for i := 0; i < m.Rows; i++ {
		rsi := 0.0
		if m.DataRows[i].Get(down) == 0.0 {
			rsi = 100.0
		} else if m.DataRows[i].Get(up) == 0.0 {
			rsi = 0.0
		} else {
			rsi = 100.0 - (100.0 / (1.0 + m.DataRows[i].Get(up)/m.DataRows[i].Get(down)))
		}
		m.DataRows[i].Set(ret, rsi)
	}
	m.RemoveColumns(5)
	return ret
}

func ModifiedRSI(m *Matrix, days, field int) int {
	// 0 = RSI
	ret := m.AddNamedColumn("MRSI")
	diff := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(diff, m.DataRows[i].Get(field)-m.DataRows[i-1].Get(field))
	}
	//up = ta.rma(math.max(ta.change(rsiSourceInput), 0), rsiLengthInput)
	//down = ta.rma(-math.min(ta.change(rsiSourceInput), 0), rsiLengthInput)
	//rsi = down == 0 ? 100 : up == 0 ? 0 : 100 - (100 / (1 + up / down))
	cgi := m.Apply(func(mr MatrixRow) float64 {
		return math.Max(mr.Get(diff), 0.0)
	})
	cli := m.Apply(func(mr MatrixRow) float64 {
		return -1.0 * math.Min(mr.Get(diff), 0.0)
	})
	up := SMA(m, days, cgi)
	down := SMA(m, days, cli)

	for i := 0; i < m.Rows; i++ {
		rsi := 0.0
		if m.DataRows[i].Get(down) == 0.0 {
			rsi = 100.0
		} else if m.DataRows[i].Get(up) == 0.0 {
			rsi = 0.0
		} else {
			rsi = 100.0 - (100.0 / (1.0 + m.DataRows[i].Get(up)/m.DataRows[i].Get(down)))
		}
		m.DataRows[i].Set(ret, rsi)
	}
	m.RemoveColumns(5)
	return ret
}

func RSISMA(m *Matrix, days, smoothing, field int) int {
	ret := RSI(m, days, field)
	SMA(m, smoothing, ret)
	return ret
}

// -----------------------------------------------------------------------
//
//	RSI-Trend
//
// -----------------------------------------------------------------------
func RSITrend(m *Matrix, days, sma, field int) int {
	// 0 = RSI
	ret := m.AddNamedColumn("RSITrend")
	ri := RSI(m, days, field)
	si := EMA(m, sma, ri)
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		p := m.DataRows[i-1]
		cnt := 0.0
		if c.Get(ri) > p.Get(ri) {
			cnt += 1.0
		}
		if c.Get(si) > p.Get(si) {
			cnt += 1.0
		}
		if c.Get(ri) > c.Get(si) {
			cnt += 1.0
		}
		if c.Get(ri) >= 50.0 {
			cnt += 1.0
		}
		d1 := c.Get(ri) - c.Get(si)
		d2 := p.Get(ri) - p.Get(si)
		if d1 > d2 {
			cnt += 1.0
		}
		m.DataRows[i].Set(ret, cnt/5.0)
	}
	m.RemoveColumns(2)
	return ret
}

// -----------------------------------------------------------------------
//
//	RSI
//
// -----------------------------------------------------------------------
func RSI_BB(m *Matrix, days, field int) int {
	// 0 = RSI 1 = Upper 2 = Lower 3 = Mid
	ri := RSI(m, days, 4)
	BollingerBandExt(m, ri, days, 2.0, 2.0)
	return ri
}

// -----------------------------------------------------------------------
//
//	RSI Momentum
//
// -----------------------------------------------------------------------
func RSIMomentum(m *Matrix, short, long, field int) int {
	// 0 = RSI Momentum
	ret := m.AddColumn()
	shortRSI := RSI(m, short, field)
	longRSI := RSI(m, long, field)
	for i := 0; i < m.Rows; i++ {
		l := m.DataRows[i].Get(longRSI)
		if l != 0.0 {
			m.DataRows[i].Set(ret, m.DataRows[i].Get(shortRSI)/l)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// ATR
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/a/atr.asp
func ATR(m *Matrix, days int) int {
	// 0 = ATR 1 = smoothed
	ret := m.AddNamedColumn("ATR")
	smoothed := m.AddNamedColumn("Smoothed")
	trIdx := m.AddColumn()
	if m.Rows < 1 {
		return -1
	}
	for i := 1; i < m.Rows; i++ {
		curH := m.DataRows[i].Get(HIGH)
		curL := m.DataRows[i].Get(LOW)
		prevC := m.DataRows[i-1].Get(ADJ_CLOSE)
		tr := myMax(curH-curL, math.Abs(curH-prevC), math.Abs(curL-prevC))
		m.DataRows[i].Set(trIdx, tr)
	}
	rma := RMA(m, days, trIdx)
	ema := EMA(m, days, rma)
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(rma))
		m.DataRows[i].Set(smoothed, m.DataRows[i].Get(ema))
	}
	m.RemoveColumns(3)
	return ret
}

func ATRExt(m *Matrix, days int, ops MAFunc) int {
	// 0 = ATR 1 = smoothed
	ret := m.AddColumn()
	smoothed := m.AddColumn()
	trIdx := m.AddColumn()
	if m.Rows < 1 {
		return -1
	}
	for i := 1; i < m.Rows; i++ {
		curH := m.DataRows[i].Get(HIGH)
		curL := m.DataRows[i].Get(LOW)
		prevC := m.DataRows[i-1].Get(ADJ_CLOSE)
		d1 := curH - prevC
		d2 := math.Abs(curL - prevC)
		d3 := math.Abs(curH - curL)
		tr := d1
		if d2 > tr {
			tr = d2
		}
		if d3 > tr {
			tr = d3
		}
		m.DataRows[i].Set(trIdx, tr)
	}
	rma := ops(m, days, trIdx)
	ema := EMA(m, days, rma)
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(rma))
		m.DataRows[i].Set(smoothed, m.DataRows[i].Get(ema))
	}
	m.RemoveColumns(3)
	return ret
}

// -----------------------------------------------------------------------
// ADR
// -----------------------------------------------------------------------
// Average Daily Range
// Links
// https://de.tradingview.com/script/afwR0BdW-Average-Daily-Range/
func ADR(m *Matrix, days int) int {
	// 0 = ADR
	ret := m.AddColumn()
	tmp := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(tmp, (m.DataRows[i].Get(HIGH) / m.DataRows[i].Get(LOW)))
	}
	si := SMA(m, days, tmp)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, 100.0*(m.DataRows[i].Get(si)-1.0))
	}
	m.RemoveColumns(2)
	return ret
}

// -----------------------------------------------------------------------
// DailyRange
// -----------------------------------------------------------------------
// Daily Range
// https://www.stockfetcher.com/forums/Indicators/Day-Range-Average-Day-Range-Day-Point-Range/26144
// https://de.tradingview.com/script/6KVjtmOY-ADR-Average-Daily-Range-by-MikeC-AKA-TheScrutiniser/
func DailyRange(m *Matrix, days int) int {
	// 0 = DailyRange
	ret := m.AddColumn()
	ti := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(2) != 0.0 {
			m.DataRows[i].Set(ti, m.DataRows[i].Get(HIGH)/m.DataRows[i].Get(2))
		}
	}
	si := SMA(m, days, ti)
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, 100.0*m.DataRows[i].Get(si))
	}
	return ret
}

func RVA(m *Matrix, days int) int {
	// 0 = RVA
	ret := m.AddColumn()
	ti := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ti, m.DataRows[i].Get(HIGH)-m.DataRows[i].Get(LOW))
	}
	si := SMA(m, days, ti)
	for i := days; i < m.Rows; i++ {
		rng := m.DataRows[i].Get(HIGH) - m.DataRows[i].Get(LOW)
		if m.DataRows[i].Get(si) != 0.0 {
			m.DataRows[i].Set(ret, rng/m.DataRows[i].Get(si))
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// ROC
// -----------------------------------------------------------------------
func ROC(m *Matrix, days, field int) int {
	// 0 = ROC 1 = Diff
	ret := m.AddNamedColumn(fmt.Sprintf("ROC%d", days))
	di := m.AddNamedColumn("ROC-D")
	for i := days; i < m.Rows; i++ {
		current := m.DataRows[i].Get(field)
		prev := m.DataRows[i-days].Get(field)
		v := 0.0
		if prev != 0.0 {
			//v = (current - prev) / prev * 100.0
			v = (current/prev - 1.0) * 100.0
		}
		m.DataRows[i].Set(ret, v)
		m.DataRows[i].Set(di, current-prev)
	}
	return ret
}

func Trend(m *Matrix) int {
	// 0 = Trend
	ret := m.AddNamedColumn("Trend")
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(OPEN) > m.DataRows[i].Get(ADJ_CLOSE) {
			m.DataRows[i].Set(ret, -1.0)
		} else {
			m.DataRows[i].Set(ret, 1.0)
		}

	}
	return ret
}

func DiffTrendCounter(m *Matrix, field int) int {
	// 0 = Trend
	ret := m.AddNamedColumn("Trend")
	tc := 0.0
	cnt := 1.0
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i].Get(field)
		dir := -1.0
		if c >= 0.0 {
			dir = 1.0
		}
		if dir != tc {
			tc = dir
			cnt = dir
		} else {
			cnt += dir
		}
		m.DataRows[i].Set(ret, cnt)
	}
	return ret
}

func TrendCounter(m *Matrix, field int) int {
	// 0 = Trend
	ret := m.AddNamedColumn("Trend")
	tc := 0.0
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i].Get(field)
		p := m.DataRows[i-1].Get(field)
		if c >= p {
			if tc > 0.0 {
				tc += 1.0
			} else {
				tc = 1.0
			}
		} else {
			if tc < 0.0 {
				tc -= 1.0
			} else {
				tc = -1.0
			}
		}
		m.DataRows[i].Set(ret, tc)
	}
	return ret
}

func CompareCounter(m *Matrix, first, second int) int {
	// 0 = Trend
	ret := m.AddNamedColumn("Compare")
	tc := 0.0
	for i := 0; i < m.Rows; i++ {
		c := m.DataRows[i].Get(first)
		p := m.DataRows[i].Get(second)
		if c >= p {
			if tc > 0.0 {
				tc += 1.0
			} else {
				tc = 1.0
			}
		} else {
			if tc < 0.0 {
				tc -= 1.0
			} else {
				tc = -1.0
			}
		}
		m.DataRows[i].Set(ret, tc)
	}
	return ret
}

// -----------------------------------------------------------------------
// Tom Denmark ROC
// -----------------------------------------------------------------------
func TDROC(m *Matrix, days, field int) int {
	// 0 = ROC
	ret := m.AddNamedColumn(fmt.Sprintf("TDROC%d", days))
	start := days
	for i := days; i < m.Rows; i++ {
		if i >= start {
			current := m.DataRows[i].Get(field)
			prev := m.DataRows[i-days].Get(field)
			v := 0.0
			if prev != 0.0 {
				v = current / prev * 100.0
			}
			m.DataRows[i].Set(ret, v)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// Stochastic
// -----------------------------------------------------------------------
func Stochastic(m *Matrix, days, ema int) int {
	// 0 = K 1 = D
	k := m.AddNamedColumn("Stoch")
	d := m.AddNamedColumn("Stoch-K")
	v := m.AddColumn()
	total := m.Rows
	if total < days {
		return -1
	}
	for i := days; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-days+1, days)
		high := m.FindMaxBetween(1, i-days+1, days)
		m.DataRows[i].Set(v, (m.DataRows[i].Get(ADJ_CLOSE)-low)/(high-low)*100.0)
	}
	slowData := SMA(m, ema, v)
	dData := SMA(m, 3, slowData)
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(k, m.DataRows[i].Get(slowData))
		m.DataRows[i].Set(d, m.DataRows[i].Get(dData))
	}
	m.RemoveColumns(3)
	return k
}

// -----------------------------------------------------------------------
// Stochastic
// -----------------------------------------------------------------------
func StochasticExt(m *Matrix, days, ema, highField, lowField, priceField int) int {
	// 0 = K 1 = D
	k := m.AddColumn()
	d := m.AddColumn()
	v := m.AddColumn()
	total := m.Rows
	if total < days {
		return -1
	}
	for i := days; i < m.Rows; i++ {
		low := m.FindMinBetween(lowField, i-days+1, days)
		high := m.FindMaxBetween(highField, i-days+1, days)
		m.DataRows[i].Set(v, (m.DataRows[i].Get(priceField)-low)/(high-low)*100.0)
	}
	slowData := SMA(m, ema, v)
	dData := SMA(m, 3, slowData)
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(k, m.DataRows[i].Get(slowData))
		m.DataRows[i].Set(d, m.DataRows[i].Get(dData))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return k
}

// -----------------------------------------------------------------------
// Stochastic
// -----------------------------------------------------------------------
func Stoch(m *Matrix, days, field int) int {
	// 0 = K
	k := m.AddColumn()
	total := m.Rows
	if total < days {
		return -1
	}
	// 100 * (close - lowest(low, length)) / (highest(high, length) - lowest(low, length)).
	for i := days; i < m.Rows; i++ {
		low := m.FindMinBetween(field, i-days+1, days)
		high := m.FindMaxBetween(field, i-days+1, days)
		d := high - low
		if d != 0.0 {
			m.DataRows[i].Set(k, (m.DataRows[i].Get(field)-low)/d*100.0)
		}
	}
	return k
}

func StochasticSMA(m *Matrix, sma, days, ema int) int {
	smaIdx := SMA(m, sma, 4)
	return StochasticExt(m, days, ema, smaIdx, smaIdx, smaIdx)
}

// -----------------------------------------------------------------------
// Stochastic RSI
// -----------------------------------------------------------------------
/*
smoothK = input(3, "K", minval=1)
smoothD = input(3, "D", minval=1)
lengthRSI = input(14, "RSI Length", minval=1)
lengthStoch = input(14, "Stochastic Length", minval=1)
src = input(close, title="RSI Source")
rsi1 = rsi(src, lengthRSI)
k = sma(stoch(rsi1, rsi1, rsi1, lengthStoch), smoothK)
d = sma(k, smoothD)
*/
func StochasticRSI(m *Matrix, rsi, stoch, smoothK, smoothD int) int {
	// 0 = K 1 = D
	ki := m.AddColumn()
	di := m.AddColumn()
	ri := RSI(m, rsi, 4)
	sr := StochasticExt(m, stoch, smoothK, ri, ri, ri)
	k := SMA(m, smoothK, sr)
	d := SMA(m, smoothD, k)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ki, m.DataRows[i].Get(k))
		m.DataRows[i].Set(di, m.DataRows[i].Get(d))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ki
}

// RSS measures the relative strength index (RSI) of the spread between two SMA indicators.
// "Readings above 70 or below 30 merely identify the potential for price to reverse and should not be taken as a trade signal.
// When an extreme is made, you should study the lower time frame to look for a trade signal.
// The trade signal could be a break in the trendline or confirmation of a reversal pattern."
// -----------------------------------------------------------------------
// RSS
// -----------------------------------------------------------------------
func RSS(m *Matrix, slow, fast, rsi, smoothing int) int {
	// 0 = RSS
	ret := m.AddColumn()
	spread := m.AddColumn()
	emaFast := EMA(m, fast, 4)
	emaSlow := EMA(m, slow, 4)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(spread, m.DataRows[i].Get(emaSlow)-m.DataRows[i].Get(emaFast))
	}
	rd := RSI(m, rsi, spread)
	rss := SMA(m, smoothing, rd)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(rss))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// PPO
// -----------------------------------------------------------------------
func PPO(m *Matrix, short, long, signal int) int {
	// 0 = Line 1 = Signal 2 = Diff
	line := m.AddColumn()
	signalIdx := m.AddColumn()
	diff := m.AddColumn()

	emaShort := EMA(m, short, 4)
	emaLong := EMA(m, long, 4)

	for i := 0; i < m.Rows; i++ {
		vl := m.DataRows[i].Get(emaLong)
		if vl != 0.0 {
			vs := m.DataRows[i].Get(emaShort)
			vd := (vs - vl) / vl * 100.0
			m.DataRows[i].Set(line, vd)
		}
	}
	signalPairs := EMA(m, signal, line)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(signalIdx, m.DataRows[i].Get(signalPairs))
		m.DataRows[i].Set(diff, m.DataRows[i].Get(line)-m.DataRows[i].Get(signalIdx))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return line

}

func standardAbbreviation(m *Matrix, v float64, offset, cnt int) float64 {
	sum := 0.0
	end := offset + cnt
	if end > m.Rows {
		end = m.Rows
	}
	for i := offset; i < end; i++ {
		d := m.DataRows[i].Get(ADJ_CLOSE) - v
		sum += d * d
	}
	return math.Sqrt(sum / float64(cnt))
}

// -----------------------------------------------------------------------
// BollingerBand
// -----------------------------------------------------------------------
func BollingerBand(m *Matrix, ema int, upper, lower float64) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddNamedColumn("BB-UP")
	lowIdx := m.AddNamedColumn("BB-LOW")
	midIdx := m.AddNamedColumn("BB-MID")
	sma := SMA(m, ema, 4)
	std := m.StdDev(4, ema)
	for i := 0; i < m.Rows; i++ {
		sa := m.DataRows[i].Get(std)
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(sma)+sa*upper)
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(sma)-sa*lower)
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return upIdx
}

// -----------------------------------------------------------------------
// BollingerBand Price Position
// -----------------------------------------------------------------------
func BollingerBand_Price_Relation(m *Matrix, ema int, upper, lower float64) int {
	// 0 = BBPR
	ret := m.AddColumn()
	bb := BollingerBand(m, ema, upper, lower)
	for i := ema; i < m.Rows; i++ {
		cb := m.DataRows[i]
		d := cb.Get(bb) - cb.Get(bb+1)
		if d != 0.0 {
			per := (1.0 - (cb.Get(bb)-cb.Get(ADJ_CLOSE))/d) * 100.0
			m.DataRows[i].Set(ret, per)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Any channel and Price Position
// -----------------------------------------------------------------------
func ChannelPriceRelation(m *Matrix, upperIndex, lowerIndex int) int {
	// 0 = BBPR
	ret := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		cb := m.DataRows[i]
		d := cb.Get(upperIndex) - cb.Get(lowerIndex)
		if d != 0.0 {
			per := (1.0 - (cb.Get(upperIndex)-cb.Get(ADJ_CLOSE))/d) * 100.0
			m.DataRows[i].Set(ret, per)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// ESDBand - like Bollinger but uses EMA
// -----------------------------------------------------------------------
func ESDBand(m *Matrix, ema int, upper, lower float64) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	midIdx := m.AddColumn()
	sma := EMA(m, ema, 4)
	std := m.StdDev(4, ema)
	for i := 0; i < m.Rows; i++ {
		sa := m.DataRows[i].Get(std)
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(sma)+sa*upper)
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(sma)-sa*lower)
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return upIdx
}

func EMAChannelPriceRelation(m *Matrix, ema int, upper, lower float64) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	midIdx := m.AddColumn()
	sma := EMA(m, ema, 4)
	std := m.StdDev(4, ema)
	for i := 0; i < m.Rows; i++ {
		sa := m.DataRows[i].Get(std)
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(sma)+sa*upper)
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(sma)-sa*lower)
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ChannelPriceRelation(m, upIdx, lowIdx)
}

// -----------------------------------------------------------------------
// BollingerBand
// -----------------------------------------------------------------------
func BollingerBandExt(m *Matrix, field, ema int, upper, lower float64) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	midIdx := m.AddColumn()
	sma := SMA(m, ema, field)
	std := m.StdDev(field, ema)
	for i := 0; i < m.Rows; i++ {
		sa := m.DataRows[i].Get(std)
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(sma)+sa*upper)
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(sma)-sa*lower)
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return upIdx
}

// -----------------------------------------------------------------------
// BollingerBand
// -----------------------------------------------------------------------
func BollingerBandSqueeze(m *Matrix, ema int, upper, lower float64, period int) int {
	ret := m.AddColumn()
	bb := BollingerBandWidth(m, ema, upper, lower)
	si := m.Stochastic(period, bb)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(si))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// VolatilityIndex
// -----------------------------------------------------------------------
func VolatilityIndex(m *Matrix, ema int, upper, lower float64) int {
	ret := m.AddColumn()
	bb := BollingerBand(m, ema, upper, lower)
	for i := 0; i < m.Rows; i++ {
		c := m.DataRows[i]
		m.DataRows[i].Set(ret, (c.Get(bb)-c.Get(bb+1))/c.Get(bb))
	}
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Elder Bars
// -----------------------------------------------------------------------
// https://school.stockcharts.com/doku.php?id=chart_analysis:elder_impulse_system
func ElderBars(m *Matrix) int {
	ret := m.AddColumn()
	ei := EMA(m, 13, ADJ_CLOSE)
	mi := MACD(m, 12, 26, 9)
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		p := m.DataRows[i-1]
		if c.Get(ei) > p.Get(ei) && c.Get(mi+2) > p.Get(mi+2) {
			m.DataRows[i].Set(ret, 1.0)
		}
		if c.Get(ei) < p.Get(ei) && c.Get(mi+2) < p.Get(mi+2) {
			m.DataRows[i].Set(ret, -1.0)
		}
	}
	m.RemoveColumns(4)
	return ret
}

// -----------------------------------------------------------------------
// BollingerBand Width
// -----------------------------------------------------------------------
func BollingerBandWidth(m *Matrix, ema int, upper, lower float64) int {
	ret := m.AddColumn()
	bb := BollingerBand(m, ema, upper, lower)
	for i := 0; i < m.Rows; i++ {
		cb := m.DataRows[i]
		if cb.Get(bb+2) != 0.0 {
			m.DataRows[i].Set(ret, (cb.Get(bb)-cb.Get(bb+1))/cb.Get(bb+2))
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// BollingerBand Width
// -----------------------------------------------------------------------
func BollingerBandWidthRatio(m *Matrix, ema int, upper, lower float64, avg int) int {
	ret := m.AddColumn()
	bb := BollingerBand(m, ema, upper, lower)
	bw := m.Apply(func(mr MatrixRow) float64 {
		if mr.Get(bb+2) != 0.0 {
			return (mr.Get(bb) - mr.Get(bb+1)) / mr.Get(bb+2)
		}
		return 0.0
	})
	si := EMA(m, avg, bw)
	for i := 0; i < m.Rows; i++ {
		c := &m.DataRows[i]
		if c.Get(si) != 0.0 {
			c.Set(ret, c.Get(bw)/c.Get(si))
		}
	}
	m.RemoveColumns(4)
	return ret
}

// -----------------------------------------------------------------------
// BollingerBand Width
// -----------------------------------------------------------------------
func BollingerBandPercentage(m *Matrix, ema int, upper, lower float64) int {
	ret := m.AddColumn()
	bb := BollingerBand(m, ema, upper, lower)
	for i := 0; i < m.Rows; i++ {
		cb := m.DataRows[i]
		if cb.Get(bb+2) != 0.0 {
			m.DataRows[i].Set(ret, (cb.Get(4)-cb.Get(bb+1))/(cb.Get(bb)-cb.Get(bb+1)))
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// K-Envelope
// -----------------------------------------------------------------------
// https://medium.com/swlh/detecting-support-resistance-levels-with-ks-envelopes-8c391ef4a471
//
/*
What it is simply saying is that we can:
Go long (Buy) whenever the market price enters the K’s Envelopes with the previous value above the K’s Envelopes so that it knows markets are seeing the Envelopes as a support.
Go short (Sell) whenever the market price enters the K’s Envelopes with the previous value below the K’s Envelopes so that it knows markets are seeing the Envelopes as a resistance.
*/
func KEnvelope(m *Matrix, days int) int {
	// 0 = Upper 1 = Lower 2 = Mid
	ret := SMA(m, days, 1)
	SMA(m, days, 2)
	midIdx := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(ADJ_CLOSE))
	}
	m.RemoveColumn()
	return ret
}

func TrueRange(cn *Matrix) int {
	ret := cn.AddColumn()
	for i := 1; i < cn.Rows; i++ {
		cur := &cn.DataRows[i]
		p := cn.DataRows[i-1]
		cur.Set(ret, myMax(cur.Get(1)-cur.Get(2), m.Abs(cur.Get(1)-p.Get(4)), m.Abs(cur.Get(2)-p.Get(4))))
	}
	return ret
}

// -----------------------------------------------------------------------
// KeltnerChannel
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/k/keltnerchannel.asp
func Keltner(m *Matrix, ema, atrLength int, multiplier float64) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddNamedColumn("Upper")
	lowIdx := m.AddNamedColumn("Lower")
	midIdx := m.AddNamedColumn("Mid")
	tp := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		cur := m.DataRows[i]
		m.DataRows[i].Set(tp, (cur.Get(HIGH)+cur.Get(2)+cur.Get(ADJ_CLOSE))/3.0)
	}
	ed := EMA(m, ema, tp)
	ad := ATR(m, atrLength)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(ed)+multiplier*m.DataRows[i].Get(ad))
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(ed)-multiplier*m.DataRows[i].Get(ad))
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(ed))
	}
	m.RemoveColumns(4)
	return upIdx
}

// -----------------------------------------------------------------------
//
//	DonchianChannel
//
// -----------------------------------------------------------------------
func DonchianChannel(m *Matrix, days int) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddNamedColumn("Upper")
	lowIdx := m.AddNamedColumn("Lower")
	midIdx := m.AddNamedColumn("Mid")
	for i := days; i < m.Rows; i++ {
		h := m.FindMaxBetween(1, i-days, days)
		l := m.FindMinBetween(2, i-days, days)
		m.DataRows[i].Set(upIdx, h)
		m.DataRows[i].Set(lowIdx, l)
		m.DataRows[i].Set(midIdx, (h+l)/2.0)
	}
	return upIdx
}

// -----------------------------------------------------------------------
//
//	RSI - ATR - RSI
//
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/coding-the-volatility-adjusted-rsi-in-tradingview-d109ef151724
func RAR(prices *Matrix, days int) int {
	// 0 = RAR
	ret := prices.AddColumn()
	tmp := prices.AddColumn()
	ri := RSI(prices, days, ADJ_CLOSE)
	ai := ATR(prices, days)
	for _, p := range prices.DataRows {
		if p.Get(ai) != 0.0 {
			p.Set(tmp, p.Get(ri)/p.Get(ai))
		}
	}
	rn := RSI(prices, days, tmp)
	prices.CopyColumn(rn, ret)
	prices.RemoveColumns(4)
	return ret
}

// -----------------------------------------------------------------------
//
//	WilliamsRange
//
// -----------------------------------------------------------------------
// W%R 14 = [ H.HIGH - C.PRICE ] / [ L.LOW - C.PRICE ] * ( - 100 )
// where,
// W%R 14 = 14-day Williams %R of the stock
// H.HIGH = 14-day Highest High of the stock
// L.LOW = 14-day Lowest Low of the stock
// C.PRICE = Closing price of the stock
func WilliamsRange(m *Matrix, days int) int {
	// 0 = %R
	ret := m.AddColumn()
	for i := days; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-days+1, days)
		high := m.FindMaxBetween(1, i-days+1, days)
		dl := high - low
		value := 1.0
		if dl != 0.0 {
			value = (high - m.DataRows[i].Get(ADJ_CLOSE)) / dl * -100.0
		}
		m.DataRows[i].Set(ret, value)
	}
	return ret

}

// -----------------------------------------------------------------------
// Mean distance
// Calculating the difference between price and SMA and then using RSI
// -----------------------------------------------------------------------
func MeanDistance(m *Matrix, lookback int) int {
	// 0 = Mean Distance
	ret := m.AddColumn()
	d := m.AddColumn()
	sd := SMA(m, lookback, 4)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(d, m.DataRows[i].Get(ADJ_CLOSE)-m.DataRows[i].Get(sd))
	}
	n := RSI(m, lookback, d)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(n))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// PER
// -----------------------------------------------------------------------
func PER(m *Matrix, ema, smoothing int) int {
	// 0 = Spread 1 = Spread Percentage 2 = Smoothed
	ret := m.AddColumn()
	sp := m.AddColumn()
	sm := m.AddColumn()
	ed := EMA(m, ema, 4)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(ADJ_CLOSE)-m.DataRows[i].Get(ed))
		if m.DataRows[i].Get(ed) != 0.0 {
			m.DataRows[i].Set(sp, m.DataRows[i].Get(ret)/m.DataRows[i].Get(ed)*100.0)
		}
	}

	sma := SMA(m, smoothing, ret)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(sm, m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// ATR
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/a/atr.asp
func StochasticATR(m *Matrix, days int) int {
	// 0 = Stoch ATR
	ret := m.AddNamedColumn("StochATR")
	ai := ATR(m, days)
	stoch := StochasticExt(m, days, ai, ai, ai, ai)
	m.CopyColumn(stoch, ret)
	m.RemoveColumns(4)
	return ret
}

// -----------------------------------------------------------------------
// Candles
// -----------------------------------------------------------------------
func Candles(m *Matrix, sizeSmoothening int) int {
	// 0 = BodySize 1 = BodyPos 2 = Mid 3 = RelBodySize 4 = RelAvg 5 = Upper 6 = Lower 7 = Trend 8 = Spread 9 = RelSpread
	bodySize := m.AddColumn()
	bodyPos := m.AddColumn()
	mid := m.AddColumn()
	relBS := m.AddColumn()
	relAvg := m.AddColumn()
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	trendIdx := m.AddColumn()
	spreadIdx := m.AddColumn()
	relSpreadIdx := m.AddColumn()

	for i := 0; i < m.Rows; i++ {
		p := m.DataRows[i]
		top := math.Max(p.Get(OPEN), p.Get(ADJ_CLOSE))
		bottom := math.Min(p.Get(OPEN), p.Get(ADJ_CLOSE))

		d := p.Get(0) - p.Get(ADJ_CLOSE)
		if d < 0.0 {
			d *= -1.0
		}
		td := p.Get(HIGH) - p.Get(LOW)
		uwAbs := p.Get(HIGH) - top
		lwAbs := bottom - p.Get(2)
		uwRel := uwAbs / td * 100.0
		lwRel := lwAbs / td * 100.0
		rb := d / td * 100.0

		//r.Set("Top", top)
		//r.Set("Bottom", bottom)
		m.DataRows[i].Set(spreadIdx, td)
		m.DataRows[i].Set(bodySize, d)
		m.DataRows[i].Set(upIdx, uwRel)
		m.DataRows[i].Set(lowIdx, lwRel)
		m.DataRows[i].Set(mid, (top+bottom)/2.0)
		m.DataRows[i].Set(relBS, rb)
		if p.Get(0) > p.Get(ADJ_CLOSE) {
			m.DataRows[i].Set(trendIdx, -1.0)
		} else {
			m.DataRows[i].Set(trendIdx, 1.0)
		}
		m.DataRows[i].Set(bodyPos, (1.0-((p.Get(HIGH)-top)/td))*100.0)
	}
	sma := SMA(m, sizeSmoothening, bodyPos)
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(sma) != 0.0 {
			m.DataRows[i].Set(relAvg, m.DataRows[i].Get(bodySize)/m.DataRows[i].Get(sma)*100.0)
		}
	}
	sma = SMA(m, sizeSmoothening, spreadIdx)
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(sma) != 0.0 {
			m.DataRows[i].Set(relSpreadIdx, m.DataRows[i].Get(spreadIdx)/m.DataRows[i].Get(sma)*100.0)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return bodySize

}

// -----------------------------------------------------------------------
// Candles
// -----------------------------------------------------------------------
func CandleWicks(m *Matrix) int {
	// 0 = BodySize 1 = Upper Wick 2 = Lower Wick 3 = BodyPos 4 = Range 5 = BodySize
	relBS := m.AddNamedColumn("RelBodySize")
	upIdx := m.AddNamedColumn("Upper")
	lowIdx := m.AddNamedColumn("Lower")
	trendIdx := m.AddNamedColumn("Trend")
	rngIdx := m.AddNamedColumn("Range")
	bs := m.AddNamedColumn("BodySize")
	for i := 0; i < m.Rows; i++ {
		p := m.DataRows[i]
		top := math.Max(p.Get(OPEN), p.Get(ADJ_CLOSE))
		bottom := math.Min(p.Get(OPEN), p.Get(ADJ_CLOSE))

		d := p.Get(OPEN) - p.Get(ADJ_CLOSE)
		if d < 0.0 {
			d *= -1.0
		}
		m.DataRows[i].Set(bs, d)
		td := p.Get(HIGH) - p.Get(LOW)
		uwAbs := p.Get(HIGH) - top
		lwAbs := bottom - p.Get(2)
		uwRel := uwAbs / td * 100.0
		lwRel := lwAbs / td * 100.0
		rb := d / td * 100.0

		m.DataRows[i].Set(upIdx, uwRel)
		m.DataRows[i].Set(lowIdx, lwRel)
		m.DataRows[i].Set(relBS, rb)
		if p.Get(OPEN) > p.Get(ADJ_CLOSE) {
			m.DataRows[i].Set(trendIdx, -1.0)
		} else {
			m.DataRows[i].Set(trendIdx, 1.0)
		}
		m.DataRows[i].Set(rngIdx, td)
	}
	return relBS
}

func StochasticBodySize(m *Matrix, days int) int {
	ret := m.AddColumn()
	cndIdx := Candles(m, days)
	for i := days; i < m.Rows; i++ {
		low := m.FindMinBetween(cndIdx, i-days, days+1)
		high := m.FindMaxBetween(cndIdx, i-days, days+1)
		if high != low {
			m.DataRows[i].Set(ret, (m.DataRows[i].Get(cndIdx)-low)/(high-low)*100.0)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	Relative Volume
//
// -----------------------------------------------------------------------
func RelativeVolume(m *Matrix, period int) int {
	// 0 = normalized stochastic volume
	ret := m.AddColumn()
	si := StochasticExt(m, period, 3, VOLUME, VOLUME, VOLUME)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(si)/100.0)
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

func AveragePrice(m *Matrix) int {
	ret := m.Apply(func(mr MatrixRow) float64 {
		return (mr.Get(HIGH) + mr.Get(LOW) + mr.Get(ADJ_CLOSE)) / 3.0
	})
	return ret
}

// -----------------------------------------------------------------------
//
//	VO
//
// -----------------------------------------------------------------------
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/volume-oscillator
// https://commodity.com/technical-analysis/volume-oscillator/
func VO(m *Matrix, fast, slow int) int {
	ret := m.AddColumn()
	emaSlow := EMA(m, slow, VOLUME)
	emaFast := EMA(m, fast, VOLUME)
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(emaSlow) != 0.0 {
			vo := (m.DataRows[i].Get(emaFast) - m.DataRows[i].Get(emaSlow)) / m.DataRows[i].Get(emaSlow) * 100.0
			m.DataRows[i].Set(ret, vo)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	AverageVolume
//
// -----------------------------------------------------------------------
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/volume-oscillator
// https://commodity.com/technical-analysis/volume-oscillator/
func AverageVolume(m *Matrix, lookback int) int {
	// 0 = SMA of volume
	sma := SMA(m, lookback, 5)
	return sma
}

// -----------------------------------------------------------------------
//
//	Ichimoku
//
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/i/ichimoku-cloud.asp
func Ichimoku(m *Matrix, short, mid, long int) int {
	// 0 = Conversion line (Tenkan) 1 = Base Line (Kijun) 2 = Leading Span A 3 = Leading Span B 4 = Lagging Span (Chikou)
	ret := m.AddNamedColumn("Tenkan")
	midIdx := m.AddNamedColumn("Kijun")
	lsaIdx := m.AddColumn()
	lsbIdx := m.AddColumn()
	lsIdx := m.AddNamedColumn("Chikou")
	// Tenkan
	for i := short; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-short, short)
		high := m.FindMaxBetween(1, i-short, short)
		m.DataRows[i].Set(ret, (high+low)/2.0)
	}
	// Kijun
	for i := mid; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-mid, mid)
		high := m.FindMaxBetween(1, i-mid, mid)
		m.DataRows[i].Set(midIdx, (high+low)/2.0)
	}

	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(lsaIdx, (m.DataRows[i].Get(ret)+m.DataRows[i].Get(midIdx))/2.0)
	}
	m.Shift(lsaIdx, 26)

	for i := long; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-long, long)
		high := m.FindMaxBetween(1, i-long, long)
		m.DataRows[i].Set(lsbIdx, (high+low)/2.0)
	}
	m.Shift(lsbIdx, 26)
	// Chikou - shift price 26 periods back
	m.CopyColumn(4, lsIdx)
	m.Shift(lsIdx, -26)
	//for i := 0; i < m.Rows-mid; i++ {
	//	m.DataRows[i].Set(lsIdx, m.DataRows[i+mid].Get(ADJ_CLOSE))
	//}
	return ret
}

// -----------------------------------------------------------------------
//
//	Weighted Trend Intensity
//
// -----------------------------------------------------------------------
// something like this: https://medium.com/geekculture/the-psychological-line-indicator-coding-back-testing-in-python-cf5210d9e079
func WeightedTrendIntensity(m *Matrix, period int) int {
	ret := m.AddColumn()
	n := float64(period)
	for i := period; i < m.Rows; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			if m.DataRows[i-j].Get(ADJ_CLOSE) > m.DataRows[i-j].Get(0) {
				sum += float64(period - j + 1)
			}
		}
		avg := sum / (n * (n + 1.0) / 2.0) * 100.0
		m.DataRows[i].Set(ret, avg)
	}
	return ret
}

// -----------------------------------------------------------------------
//  Supertrend
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/one-of-the-best-trend-following-strategies-893f903230e7

/*
//@version=5
indicator("Pine Supertrend")

[supertrend, direction] = ta.supertrend(3, 10)
plot(direction < 0 ? supertrend : na, "Up direction", color = color.green, style=plot.style_linebr)
plot(direction < 0? na : supertrend, "Down direction", color = color.red, style=plot.style_linebr)

// The same on Pine
pine_supertrend(factor, atrPeriod) =>

	src = hl2
	atr = ta.atr(atrPeriod)
	upperBand = src + factor * atr
	lowerBand = src - factor * atr
	prevLowerBand = nz(lowerBand[1])
	prevUpperBand = nz(upperBand[1])




	[superTrend, direction]
*/
func Supertrend(m *Matrix, period int, multiplier float64) int {
	ret := m.AddColumn()
	atrIdx := ATR(m, period)
	buIdx := m.AddColumn()
	blIdx := m.AddColumn()
	bfuIdx := m.AddColumn()
	bflIdx := m.AddColumn()
	/*
		src = hl2
		atr = ta.atr(atrPeriod)
		upperBand = src + factor * atr
		lowerBand = src - factor * atr
		prevLowerBand = nz(lowerBand[1])
		prevUpperBand = nz(upperBand[1])
	*/
	for i := 0; i < m.Rows; i++ {
		cp := m.DataRows[i]
		m.DataRows[i].Set(buIdx, (cp.Get(HIGH)+cp.Get(2))/2.0+multiplier*cp.Get(atrIdx))
		m.DataRows[i].Set(blIdx, (cp.Get(HIGH)+cp.Get(2))/2.0-multiplier*cp.Get(atrIdx))
	}

	//lowerBand := lowerBand > prevLowerBand or close[1] < prevLowerBand ? lowerBand : prevLowerBand
	//upperBand := upperBand < prevUpperBand or close[1] > prevUpperBand ? upperBand : prevUpperBand
	for i := 1; i < m.Rows; i++ {
		cp := m.DataRows[i]
		pp := m.DataRows[i-1]

		if cp.Get(buIdx) < pp.Get(buIdx) || pp.Get(ADJ_CLOSE) > pp.Get(buIdx) {
			m.DataRows[i].Set(bfuIdx, cp.Get(buIdx))
		} else {
			m.DataRows[i].Set(bfuIdx, pp.Get(bfuIdx))
		}

		if cp.Get(blIdx) > pp.Get(blIdx) || pp.Get(ADJ_CLOSE) < pp.Get(blIdx) {
			m.DataRows[i].Set(bflIdx, cp.Get(blIdx))
		} else {
			m.DataRows[i].Set(bflIdx, pp.Get(bflIdx))
		}

	}
	/*
	   int direction = na
	   	float superTrend = na
	   	prevSuperTrend = superTrend[1]
	   	if na(atr[1])
	   		direction := 1
	   	else if prevSuperTrend == prevUpperBand
	   		direction := close > upperBand ? -1 : 1
	   	else
	   		direction := close < lowerBand ? 1 : -1
	   	superTrend := direction == -1 ? lowerBand : upperBand
	*/
	for i := 1; i < m.Rows; i++ {
		cp := m.DataRows[i]
		pp := m.DataRows[i-1]

		dir := 0.0

		if pp.Get(ret) == pp.Get(bfuIdx) {
			if cp.Get(ADJ_CLOSE) > cp.Get(bfuIdx) {
				dir = 1.0
			} else {
				dir = -1.0
			}
		} else {
			if cp.Get(ADJ_CLOSE) < cp.Get(bflIdx) {
				dir = 1.0
			} else {
				dir = -1.0
			}
		}

		if dir == 1.0 {
			m.DataRows[i].Set(ret, cp.Get(bfuIdx))
		} else {
			m.DataRows[i].Set(ret, cp.Get(bflIdx))
		}
	}
	return ret
}

func GAP_ATR(m *Matrix) int {
	ret := m.AddColumn()
	atrIdx := ATR(m, 14)
	for i := 15; i < m.Rows; i++ {
		cp := m.DataRows[i]
		prev := m.DataRows[i-1]
		value := (cp.Get(0) - prev.Get(ADJ_CLOSE)) / cp.Get(atrIdx) * 100.0
		m.DataRows[i].Set(ret, value)
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

func GAP(m *Matrix) int {
	// 0 = GAP % 1 = GAP / ATR
	ret := m.AddNamedColumn("GAP")
	gi := m.AddNamedColumn("GAP/ATR")
	ai := ATR(m, 14)
	for i := 1; i < m.Rows; i++ {
		cp := m.DataRows[i]
		prev := m.DataRows[i-1]
		value := (cp.Get(0)/prev.Get(ADJ_CLOSE) - 1.0) * 100.0
		m.DataRows[i].Set(ret, value)
		if cp.Get(ai) != 0.0 {
			m.DataRows[i].Set(gi, (cp.Get(0)-prev.Get(ADJ_CLOSE))/cp.Get(ai))
		}
	}
	m.RemoveColumns(2)
	return ret
}

func PriceATR(prices *Matrix, period int) int {
	ret := prices.AddColumn()
	atrIdx := ATR(prices, period)
	for i := 1; i < prices.Rows; i++ {
		cp := prices.DataRows[i]
		prev := prices.DataRows[i-1]
		sl := m.Abs(cp.Get(ADJ_CLOSE) - prev.Get(ADJ_CLOSE))
		value := sl / cp.Get(atrIdx)
		prices.DataRows[i].Set(ret, value)
	}
	prices.RemoveColumns(2)
	return ret
}

func RangeATR(prices *Matrix, period int) int {
	ret := prices.AddColumn()
	atrIdx := ATR(prices, period)
	for i := 0; i < prices.Rows; i++ {
		cp := prices.DataRows[i]
		value := (1.0 - (cp.Get(HIGH)-cp.Get(LOW))/cp.Get(atrIdx))
		prices.DataRows[i].Set(ret, value)
	}
	prices.RemoveColumns(2)
	return ret
}

// -----------------------------------------------------------------------
// Kairi Relative Index
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/quantifying-the-deviation-of-price-from-its-equilibrium-98dae4fe9818
func KRI(m *Matrix, period int) int {
	// 0 = KRI
	ret := m.AddColumn()
	si := SMA(m, period, ADJ_CLOSE)
	for i := period; i < m.Rows; i++ {
		sma := m.DataRows[i].Get(si)
		if sma != 0.0 {
			d := (m.DataRows[i].Get(ADJ_CLOSE) - sma) / sma * 100.0
			m.DataRows[i].Set(ret, d)
		}
	}
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	STD
//
// -----------------------------------------------------------------------
func STD(prices *Matrix, days int) int {
	// 0 = STD
	s := prices.StdDev(4, days)
	return s
}

// -----------------------------------------------------------------------
//
//	STDChannel
//
// -----------------------------------------------------------------------
func STDChannel(prices *Matrix, days int, std float64) int {
	// 0 = Upper 1 = Lower
	ui := prices.AddColumn()
	li := prices.AddColumn()
	s := prices.StdDev(4, days)
	for i := 0; i < prices.Rows; i++ {
		prices.DataRows[i].Set(ui, prices.DataRows[i].Get(ADJ_CLOSE)+std*prices.DataRows[i].Get(s))
		prices.DataRows[i].Set(li, prices.DataRows[i].Get(ADJ_CLOSE)-std*prices.DataRows[i].Get(s))
	}
	return ui
}

// -----------------------------------------------------------------------
//
//	SupportResistanceChannel
//
// -----------------------------------------------------------------------
func SupportResistanceChannel(prices *Matrix, days int, std float64) int {
	// 0 = Upper 1 = Lower
	ma := prices.AddColumn()
	mi := prices.AddColumn()
	for i := days; i < prices.Rows; i++ {
		lm := prices.FindMinBetween(LOW, i-days, days)
		r := prices.FindMaxBetween(HIGH, i-days, days) - lm
		upper := lm + r*(1.0-std)
		prices.DataRows[i].Set(ma, upper)
		lower := lm + r*std
		prices.DataRows[i].Set(mi, lower)
	}
	return ma
}

// -----------------------------------------------------------------------
//
//	STD
//
// -----------------------------------------------------------------------
func STDStochastic(prices *Matrix, days int) int {
	// 0 = K 1 = D
	s := prices.StdDev(4, days)
	return StochasticExt(prices, days, 3, s, s, s)
}

// -----------------------------------------------------------------------
// DeMark
// -----------------------------------------------------------------------
func DeMark(candles *Matrix) int {
	// 0 = Trend 1 = Count
	trend := candles.AddColumn()
	count := candles.AddColumn()
	ct := 0
	cur := 1
	for i := 4; i < candles.Rows; i++ {
		c := candles.DataRows[i].Get(ADJ_CLOSE)
		p := candles.DataRows[i-4].Get(ADJ_CLOSE)
		t := 0
		if c > p {
			t = 1
			cur++
		} else {
			t = -1
			cur++
		}
		if t != ct {
			cur = 1
			ct = t
		}
		candles.DataRows[i].Set(trend, float64(t))
		candles.DataRows[i].Set(count, float64(cur))
	}
	return trend
}

// https://kaabar-sofien.medium.com/the-demarker-contrarian-indicator-a-study-in-python-2caa066a30e1
func DeMarker(m *Matrix, days int) int {
	// 0 = Trend
	ret := m.AddNamedColumn("Demark")
	ma := m.AddColumn()
	mi := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		p := m.DataRows[i-1]
		dh := c.Get(HIGH) - p.Get(HIGH)
		if dh > 0.0 {
			m.DataRows[i].Set(ma, dh)
		}
		dl := p.Get(LOW) - c.Get(LOW)
		if dl > 0.0 {
			m.DataRows[i].Set(mi, dl)
		}

	}
	s1 := EMA(m, days, ma)
	s2 := EMA(m, days, mi)
	for i := 0; i < m.Rows; i++ {
		s := m.DataRows[i].Get(s1) + m.DataRows[i].Get(s2)
		if s != 0.0 {
			m.DataRows[i].Set(ret, m.DataRows[i].Get(s1)/s)
		}
	}
	m.RemoveColumns(4)
	return ret
}

// -----------------------------------------------------------------------
// Bullish Bearish Power
// -----------------------------------------------------------------------
func BullishBearish(m *Matrix, period int) int {
	// 0 = Green / Red 1 = Trend Count 2 = High 3 = Low
	ret := m.AddColumn()
	ti := m.AddColumn()
	hi := m.AddColumn()
	li := m.AddColumn()
	for i := period; i < m.Rows; i++ {
		bu := 0
		tr := 0
		h := 0
		l := 0
		for j := 0; j < period; j++ {
			idx := i - j
			if m.DataRows[idx].Get(0) <= m.DataRows[idx].Get(ADJ_CLOSE) {
				bu++
			}
			if j > 0 {
				if m.DataRows[idx].Get(1) >= m.DataRows[idx-1].Get(1) {
					h++
				}
				if m.DataRows[idx].Get(2) <= m.DataRows[idx-1].Get(2) {
					l++
				}
				if m.DataRows[idx].Get(4) >= m.DataRows[idx-1].Get(4) {
					tr++
				}
			}
		}
		d := float64(bu) / float64(period) * 100.0
		m.DataRows[i].Set(ret, d)
		m.DataRows[i].Set(ti, float64(tr)/float64(period-1)*100.0)
		m.DataRows[i].Set(hi, float64(h)/float64(period-1)*100.0)
		m.DataRows[i].Set(li, float64(l)/float64(period-1)*100.0)
	}
	return ret
}

// -----------------------------------------------------------------------
//
// OBV
// https://www.investopedia.com/terms/o/onbalancevolume.asp
// -----------------------------------------------------------------------
func OBV(m *Matrix, scale float64) int {
	// 0 = OBV
	ret := m.AddColumn()
	if m.Rows > 2 {
		prev := m.DataRows[0]
		obv := prev.Get(5)
		sign := 1.0
		for i, p := range m.DataRows {
			if i > 0 {
				if p.Get(ADJ_CLOSE) > prev.Get(ADJ_CLOSE) {
					sign = 1.0
				}
				if p.Get(ADJ_CLOSE) < prev.Get(ADJ_CLOSE) {
					sign = -1.0
				}
				obv += p.Get(5) * sign
				m.DataRows[i].Set(ret, obv/scale)
			}
			prev = p
		}
	}
	return ret
}

// -----------------------------------------------------------------------
//
//	Aroon
//
// -----------------------------------------------------------------------
func Aroon(m *Matrix, days int) int {
	// 0 = AroonUp 1 = AroonDown
	ui := m.AddColumn()
	di := m.AddColumn()
	idx := 0
	n := float64(days)
	for i := days - 1; i < m.Rows; i++ {
		lIdx, hIdx := m.FindHighLowIndex(idx, days)
		lIdx = days - (lIdx - idx)
		hIdx = days - (hIdx - idx)
		m.DataRows[i].Set(ui, (n-float64(hIdx))/n*100.0)
		m.DataRows[i].Set(di, (n-float64(lIdx))/n*100.0)
		idx++
	}
	return ui
}

// -----------------------------------------------------------------------
//
//	TrendIntensity
//
// -----------------------------------------------------------------------
func TrendIntensity(m *Matrix, days int) int {
	// 0 = TS
	ret := m.AddColumn()
	sma := SMA(m, days, 4)
	for i := days; i < m.Rows; i++ {
		tu := 0.0
		tl := 0.0
		for j := 0; j < days; j++ {
			cp := m.DataRows[j-days+i]
			s := cp.Get(sma)
			if cp.Get(ADJ_CLOSE) > s {
				tu += 1.0
			} else {
				tl += 1.0
			}
		}
		v := 0.0
		if tu != 0.0 && tl != 0.0 {
			v = tu / (tu + tl) * 100.0
		}
		m.DataRows[i].Set(ret, v)
	}
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	Chaikin A/D Line
//
// -----------------------------------------------------------------------
// https://www.boerse.de/technische-indikatoren/A-D-Linie-1
func AD(m *Matrix) int {
	// 0 = A/D Line
	// ad = ta.cum(close==high and close==low or high==low ? 0 : ((2*close-low-high)/(high-low))*volume)
	ret := m.AddColumn()
	ad := 0.0
	for i := 1; i < m.Rows; i++ {
		// Acc./Distr. Line =[((C-L) - (H-C))/(H-L) * V] + I
		cur := m.DataRows[i]
		high := cur.Get(HIGH)
		close := cur.Get(ADJ_CLOSE)
		low := cur.Get(2)
		hl := high - low
		if hl != 0.0 {
			ad += ((close - low) - (high - close)) / hl * cur.Get(5)
		}
		m.DataRows[i].Set(ret, ad)
	}
	return ret
}

// -----------------------------------------------------------------------
//
//	TSI
//
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/t/tsi.asp
func TSI(m *Matrix, short, long, signal int) int {
	// 0 = TSI 1 = Signal 2 = Diff
	ret := m.AddNamedColumn("TSI")
	si := m.AddNamedColumn("Signal")
	di := m.AddNamedColumn("Diff")
	mi := m.AddColumn()
	// PC = CCP − PCP
	for i := 1; i < m.Rows; i++ {
		cur := m.DataRows[i]
		prev := m.DataRows[i-1]
		m.DataRows[i].Set(mi, cur.Get(ADJ_CLOSE)-prev.Get(ADJ_CLOSE))
	}
	pcs := EMA(m, long, mi)
	pcds := EMA(m, short, pcs)

	ami := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ami, math.Abs(m.DataRows[i].Get(mi)))
	}

	apcs := EMA(m, long, ami)
	apcds := EMA(m, short, apcs)

	for i := 0; i < m.Rows; i++ {
		tsi := 0.0
		if m.DataRows[i].Get(apcds) != 0.0 {
			tsi = 100.0 * (m.DataRows[i].Get(pcds) / m.DataRows[i].Get(apcds))
		}
		m.DataRows[i].Set(ret, tsi)
	}

	tsiEMA := EMA(m, signal, ret)

	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(si, m.DataRows[i].Get(tsiEMA))
		m.DataRows[i].Set(di, m.DataRows[i].Get(ret)-m.DataRows[i].Get(tsiEMA))
	}
	m.RemoveColumns(7)
	return ret
}

// -----------------------------------------------------------------------
//
//	Divergence
//
// -----------------------------------------------------------------------
func Divergence(m *Matrix, first, second, period int) int {
	// 0 = Divergence (1=bullish -1=bearish)
	ret := m.AddColumn()
	for i := period; i < m.Rows-1; i++ {
		cur := m.DataRows[i]
		next := m.DataRows[i+1]
		mi := m.FindMinBetween(first, i-period, period)
		if mi > cur.Get(first) && next.Get(first) > cur.Get(first) {
			mi = m.FindMinBetween(second, i-period, period)
			if cur.Get(second) > mi && next.Get(second) > cur.Get(second) {
				m.DataRows[i].Set(ret, 1.0)
			}
		}
		ma := m.FindMaxBetween(first, i-period, period)
		if cur.Get(first) > ma && cur.Get(first) > next.Get(first) {
			ma = m.FindMaxBetween(second, i-period, period)
			if cur.Get(second) < ma && cur.Get(second) > next.Get(second) {
				m.DataRows[i].Set(ret, -1.0)
			}
		}
	}
	return ret
}

// -----------------------------------------------------------------------
//
//	Divergence
//
// -----------------------------------------------------------------------
func HighLowChannel(m *Matrix, highPeriod, lowPeriod int) int {
	// 0 = High 1 = Low
	hs := SMA(m, highPeriod, 1)
	SMA(m, lowPeriod, 2)
	return hs
}

func HighLowEMAChannel(m *Matrix, highPeriod, lowPeriod int) int {
	// 0 = High 1 = Low
	hs := EMA(m, highPeriod, 1)
	EMA(m, lowPeriod, 2)
	return hs
}

func HighestLowestChannel(m *Matrix, period int) int {
	// 0 = High 1 = Low
	hi := m.AddNamedColumn("High")
	li := m.AddNamedColumn("Low")
	for i := period; i < m.Rows; i++ {
		m.DataRows[i].Set(hi, m.FindMaxBetween(HIGH, i-period, period))
		m.DataRows[i].Set(li, m.FindMinBetween(LOW, i-period, period))
	}
	return hi
}

// -----------------------------------------------------------------------
//
//	ADX
//
// -----------------------------------------------------------------------
func ADX(m *Matrix, lookback int) int {
	// 0 = ADX 1 = PDI 2 = MDI 3 = Diff
	lbf := float64(lookback)
	adx := m.AddNamedColumn("ADX")
	pdi := m.AddNamedColumn("PDI")
	mdi := m.AddNamedColumn("MDI")
	di := m.AddNamedColumn("Diff")
	if m.Rows > 2*lookback {
		plusDM := m.AddColumn()
		minusDM := m.AddColumn()
		tr := m.AddColumn()
		trur := m.AddColumn()
		pdm14 := m.AddColumn()
		mdm14 := m.AddColumn()
		pd14 := m.AddColumn()
		md14 := m.AddColumn()
		dx := m.AddColumn()
		// Calculate +DM, -DM, and the true range (TR) for each period. Fourteen periods are typically used.
		var fa []float64
		fa = append(fa, 0.0)
		fa = append(fa, 0.0)
		fa = append(fa, 0.0)
		for i := 1; i < m.Rows; i++ {
			c := m.DataRows[i]
			p := m.DataRows[i-1]
			// WENN(CH-PH>PL-CL;MAX(CH-PH;0);0)
			if (c.Get(HIGH) - p.Get(HIGH)) > (p.Get(2) - c.Get(2)) {
				m.DataRows[i].Set(plusDM, math.Max(c.Get(HIGH)-p.Get(HIGH), 0.0))
			}
			// WENN(D3-D4>C4-C3;MAX(D3-D4;0);0)
			if (p.Get(2) - c.Get(2)) > (c.Get(HIGH) - p.Get(HIGH)) {
				m.DataRows[i].Set(minusDM, math.Max(p.Get(2)-c.Get(2), 0.0))
			}
			// MAX(High-Low;ABS(High-PP);ABS(Low-PP))
			fa[0] = c.Get(HIGH) - c.Get(2)
			fa[1] = math.Abs(c.Get(HIGH) - m.DataRows[i-1].Get(ADJ_CLOSE))
			fa[2] = math.Abs(c.Get(2) - m.DataRows[i-1].Get(ADJ_CLOSE))
			trv := fa[0]
			if fa[1] > trv {
				trv = fa[1]
			}
			if fa[2] > trv {
				trv = fa[2]
			}
			m.DataRows[i].Set(tr, trv)

		}

		m.DataRows[lookback-1].Set(trur, m.PartialSum(tr, 0, lookback))
		for i, t := range m.DataRows {
			if i >= lookback {
				p := m.DataRows[i-1].Get(trur)
				c := t.Get(tr)
				v := p - (p / lbf) + c
				m.DataRows[i].Set(trur, v)
			}
		}

		m.DataRows[lookback-1].Set(pdm14, m.PartialSum(plusDM, 0, lookback))
		m.DataRows[lookback-1].Set(mdm14, m.PartialSum(minusDM, 0, lookback))

		for i, d := range m.DataRows {
			if i >= lookback {
				pp := m.DataRows[i-1].Get(pdm14)
				pm := m.DataRows[i-1].Get(mdm14)
				m.DataRows[i].Set(pdm14, pp-(pp/lbf)+d.Get(plusDM))
				m.DataRows[i].Set(mdm14, pm-(pm/lbf)+d.Get(minusDM))
			}
		}
		for i, d := range m.DataRows {
			if i >= (lookback - 1) {
				p := 100.0 * (d.Get(pdm14) / d.Get(trur))
				m.DataRows[i].Set(pd14, p)
				mv := 100.0 * (d.Get(mdm14) / d.Get(trur))
				m.DataRows[i].Set(md14, mv)
				diff := math.Abs(p - mv)
				sum := p + mv
				if sum == 0.0 {
					sum = 1.0
				}
				m.DataRows[i].Set(dx, 100.0*diff/sum)
			}
		}

		avg := m.PartialSum(dx, 0, 2*lookback+1) / lbf
		//nr.Set("ADX", avg)
		pa := avg
		for i := 26; i < m.Rows; i++ {
			v := (pa*(lbf-1.0) + m.DataRows[i].Get(dx)) / lbf
			m.DataRows[i].Set(pdi, m.DataRows[i].Get(pd14))
			m.DataRows[i].Set(mdi, m.DataRows[i].Get(md14))
			m.DataRows[i].Set(adx, v)
			m.DataRows[i].Set(di, (m.DataRows[i].Get(pd14) - m.DataRows[i].Get(md14)))
			pa = v
		}
		for i := 0; i < 9; i++ {
			m.RemoveColumn()
		}
	}
	return adx
}

// -----------------------------------------------------------------------
// RVI Relative Vigor index
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/r/relative_vigor_index.asp
func RVI(m *Matrix, lookback, signal int) int {
	// 0 = RVI 1 = Signal
	li := m.AddColumn()
	si := m.AddColumn()
	co := m.Subtract(4, 0)
	num := SWMA(m, co)
	hl := m.Subtract(1, 2)
	dem := SWMA(m, hl)
	sn := m.Sum(num, lookback)
	dn := m.Sum(dem, lookback)
	for i := lookback; i < m.Rows; i++ {
		if m.DataRows[i].Get(dn) != 0.0 {
			m.DataRows[i].Set(li, m.DataRows[i].Get(sn)/m.DataRows[i].Get(dn))
		}
	}
	for i := 3; i < m.Rows; i++ {
		d := (m.DataRows[i].Get(li) + 2.0*m.DataRows[i-1].Get(li) + 2.0*m.DataRows[i-2].Get(li) + m.DataRows[i-3].Get(li)) / 6.0
		m.DataRows[i].Set(si, d)
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return li
}

func RVIStochastic(m *Matrix, lookback, signal int) int {
	// 0 = RVI 1 = Signal 2 = RVI K 3 = RVI D
	li := RVI(m, lookback, signal)
	StochasticExt(m, 14, 3, li, li, li)
	return li
}

// -----------------------------------------------------------------------
// KD
// -----------------------------------------------------------------------
// https://medium.com/tej-api-%E9%87%91%E8%9E%8D%E8%B3%87%E6%96%99%E5%88%86%E6%9E%90/quant-10-kd-indicator-85f8944be83f
func KD(m *Matrix, period int) int {
	// 0 = RSV 1 = K 2 = D
	ret := m.AddColumn()
	ki := m.AddColumn()
	di := m.AddColumn()
	for i := period; i < m.Rows; i++ {
		cur := m.DataRows[i]
		mi := m.FindMinBetween(4, i-period, period)
		ma := m.FindMaxBetween(4, i-period, period)
		if ma != mi {
			rsv := (cur.Get(ADJ_CLOSE) - mi) / (ma - mi) * 100.0
			k := m.DataRows[i-1].Get(ki)*2.0/3.0 + rsv/3.0
			d := m.DataRows[i-1].Get(di)*2.0/3.0 + k/3.0
			m.DataRows[i].Set(ret, rsv)
			m.DataRows[i].Set(ki, k)
			m.DataRows[i].Set(di, d)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// DPO
// -----------------------------------------------------------------------
// https://medium.datadriveninvestor.com/the-detrended-price-oscillator-creating-back-testing-in-python-9ada6461bda5
// https://www.investopedia.com/terms/d/detrended-price-oscillator-dpo.asp
// https://school.stockcharts.com/doku.php?id=technical_indicators:detrended_price_osci
func DPO(m *Matrix, period int) int {
	// 0 = DPO
	ret := m.AddColumn()
	si := SMA(m, period, 4)
	k := period/2 + 1
	for i := period; i < m.Rows; i++ {
		cur := m.DataRows[i]
		prev := m.DataRows[i-k]
		m.DataRows[i].Set(ret, prev.Get(ADJ_CLOSE)-cur.Get(si))
	}
	m.RemoveColumn()
	return ret
}

func MACD_BB(m *Matrix, short, long, period int, s float64) int {
	// 0 = Upper 1 = Lower 2 = MACD
	upper := m.AddColumn()
	lower := m.AddColumn()
	macd := m.AddColumn()
	e1 := EMA(m, short, 4)
	e2 := EMA(m, long, 4)
	d := m.Subtract(e1, e2)
	si := SMA(m, period, d)
	std := m.StdDev(si, period)
	for i := 0; i < m.Rows; i++ {
		cur := m.DataRows[i]
		m.DataRows[i].Set(upper, cur.Get(si)+s*cur.Get(std))
		m.DataRows[i].Set(lower, cur.Get(si)-s*cur.Get(std))
		m.DataRows[i].Set(macd, cur.Get(d))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return upper
}

// -----------------------------------------------------------------------
// CCI
// -----------------------------------------------------------------------
func CCI(m *Matrix, days, smoothed int) int {
	// 0 = CCI 1 = Smoothed
	ret := m.AddNamedColumn(fmt.Sprintf("CCI %d", days))
	hlc := m.AddColumn()
	md := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		p := m.DataRows[i]
		m.DataRows[i].Set(hlc, (p.Get(HIGH)+p.Get(2)+p.Get(ADJ_CLOSE))/3.0)
	}
	sma := SMA(m, days, hlc)

	for i := days; i < m.Rows; i++ {
		sum := 0.0
		for j := 0; j < days; j++ {
			p := m.DataRows[j+i-days].Get(hlc)
			s := m.DataRows[j+i-days].Get(sma)
			sum += math.Abs(p - s)
		}
		sum = sum / float64(days)
		m.DataRows[i].Set(md, sum)
	}

	for i := 0; i < m.Rows; i++ {
		s := m.DataRows[i].Get(sma)
		md := m.DataRows[i].Get(md)
		hlc := m.DataRows[i].Get(hlc)
		if md != 0.0 {
			//cci = (src - ma) / (0.015 * dev(src, length))
			current := (hlc - s) / (0.015 * md)
			m.DataRows[i].Set(ret, current)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	SMA(m, smoothed, ret)
	return ret
}

// -----------------------------------------------------------------------
// VWAP
// -----------------------------------------------------------------------
func VWAP(m *Matrix, days int, std float64) int {
	// 0 = VWAP 1 = Upper 2 = Lower
	ret := m.AddColumn()
	upper := m.AddColumn()
	lower := m.AddColumn()
	hlc := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		p := m.DataRows[i]
		m.DataRows[i].Set(hlc, (p.Get(HIGH)+p.Get(2)+p.Get(ADJ_CLOSE))/3.0*p.Get(5))
	}
	for i := days; i < m.Rows; i++ {
		sumTP := 0.0
		sumV := 0.0
		for j := 0; j < days; j++ {
			sumTP += m.DataRows[j+i-days].Get(hlc)
			sumV += m.DataRows[j+i-days].Get(5)
		}
		vw := sumTP / sumV
		m.DataRows[i].Set(ret, vw)
	}
	si := m.StdDev(ret, days)
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(upper, m.DataRows[i].Get(ret)+std*m.DataRows[i].Get(si))
		m.DataRows[i].Set(lower, m.DataRows[i].Get(ret)-std*m.DataRows[i].Get(si))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

/*
study("Derivative Oscillator", shorttitle="DOSC")

rsiLength = input(title="RSI Length", type=input.integer, defval=14)
ema1Length = input(title="1st EMA Smoothing Length", type=input.integer, defval=5)
ema2Length = input(title="2nd EMA Smoothing Length", type=input.integer, defval=3)
smaLength = input(title="3rd SMA Smoothing Length", type=input.integer, defval=9)
signalLength = input(title="Signal Length", type=input.integer, defval=9)
src = input(title="Source", type=input.source, defval=close)

smoothedRSI = ema(ema(rsi(src, rsiLength), ema1Length), ema2Length)
dosc = smoothedRSI - sma(smoothedRSI, smaLength)
signal = sma(dosc, signalLength)

doscColor = dosc >= 0 ? dosc[1] < dosc ? #26A69A : #B2DFDB : dosc[1] < dosc ? #FFCDD2 : #EF5350
plot(dosc, title="DOSC", style=plot.style_columns, linewidth=2, color=doscColor, transp=0)
plot(signal, title="Signal", linewidth=2, color=color.black, transp=0)
*/
// -----------------------------------------------------------------------
// DOSC
// -----------------------------------------------------------------------
// https://www.daytrading.com/derivative-oscillator
// Defaults: 14,5,3,9,9
func DOSC(m *Matrix, r, e1, e2, s, sl int) int {
	// 0 = DOSC 1 = Signal
	ret := m.AddColumn()
	si := m.AddColumn()
	ri := RSI(m, r, 4)
	ei1 := EMA(m, e1, ri)
	ei2 := EMA(m, e2, ei1)
	smi := SMA(m, s, ei2)
	start := e1 + e2 + s
	for i := start; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(ei2)-m.DataRows[i].Get(smi))
	}
	si2 := SMA(m, sl, ret)
	start += sl
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(si, m.DataRows[i].Get(si2))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	HeikinAshi
//
// -----------------------------------------------------------------------
func HeikinAshi(m *Matrix) int {
	// 0 = Open 1 = High 2 = Low 3 = Close 4 = AdjClose 5 = Volume
	oi := m.AddColumn()
	hi := m.AddColumn()
	li := m.AddColumn()
	ci := m.AddColumn()
	aci := m.AddColumn()
	vi := m.AddColumn()
	c := &m.DataRows[0]
	c.Set(oi, 0.5*(c.Get(0)+c.Get(4)))
	c.Set(ci, 0.25*(c.Get(0)+c.Get(HIGH)+c.Get(2)+c.Get(ADJ_CLOSE)))
	c.Set(hi, c.Get(HIGH))
	c.Set(li, c.Get(2))
	for i := 1; i < m.Rows; i++ {
		c = &m.DataRows[i]
		prev := m.DataRows[i-1]
		// Open = [Open (previous bar) + Close (previous bar)] /2
		c.Set(oi, 0.5*(prev.Get(oi)+prev.Get(ci)))
		// Close = (Open + High + Low + Close) / 4
		c.Set(ci, 0.25*(c.Get(0)+c.Get(HIGH)+c.Get(2)+c.Get(ADJ_CLOSE)))
		// High = Maximum of High, Open, or Close (whichever is highest)
		c.Set(hi, FindMax([]float64{c.Get(HIGH), c.Get(oi), c.Get(ci)}))
		// Low = Minimum of Low, Open, or Close (whichever is lowest)
		c.Set(li, FindMin([]float64{c.Get(2), c.Get(oi), c.Get(ci)}))

		c.Set(aci, c.Get(ci))
		c.Set(vi, float64(c.Get(5)))
	}
	return oi
}

// -----------------------------------------------------------------------
//
//	Hull Moving Average
//
// -----------------------------------------------------------------------
// https://school.stockcharts.com/doku.php?id=technical_indicators:hull_moving_average
// https://blog.earn2trade.com/hull-moving-average/
func HMA(prices *Matrix, period, field int) int {
	// 0 = HMA
	ret := prices.AddNamedColumn(fmt.Sprintf("HMA%d", period))
	wi1 := WMA(prices, period/2, field)
	wi2 := WMA(prices, period, field)
	ri := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		prices.DataRows[i].Set(ri, 2.0*prices.DataRows[i].Get(wi1)-prices.DataRows[i].Get(wi2))
	}
	ri2 := WMA(prices, int(m.Sqrt(float64(period))), ri)
	for i := 0; i < prices.Rows; i++ {
		prices.DataRows[i].Set(ret, prices.DataRows[i].Get(ri2))
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Centre of gravity
// -----------------------------------------------------------------------
// https://www.stockmaniacs.net/center-of-gravity-indicator/
func COG(m *Matrix, period int) int {
	// 0 = COG
	ret := m.AddColumn()
	for i := period; i < m.Rows; i++ {
		s1 := 0.0
		s2 := 0.0
		for j := 0; j < period; j++ {
			s1 += m.DataRows[i-j].Get(ADJ_CLOSE)
			s2 += m.DataRows[i-j].Get(ADJ_CLOSE) * (float64(j) + 1.0)
		}

		m.DataRows[i].Set(ret, -1.0*s2/s1)
	}
	return ret
}

// -----------------------------------------------------------------------
// GRI
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/the-gri-range-index-trading-strategy-f49873266e28
func GRI(prices *Matrix, period int) int {
	// 0 = GRI
	ret := prices.AddColumn()
	n := float64(period)
	for i := period; i < prices.Rows; i++ {
		h, l := prices.FindHighestHighLowestLow(i-period, period)
		v := m.Log(h-l) / m.Log(n)
		prices.DataRows[i].Set(ret, v)
	}
	return ret
}

// -----------------------------------------------------------------------
// CMF Chaikin Money Flow
// -----------------------------------------------------------------------
// https://corporatefinanceinstitute.com/resources/knowledge/trading-investing/chaikin-money-flow-cmf/
func CMF(prices *Matrix, period int) int {
	// 0 = CMF
	ret := prices.AddColumn()
	mf := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		cp := prices.DataRows[i]
		hl := cp.Get(HIGH) - cp.Get(2)
		if hl != 0.0 {
			// [(Close  -  Low) - (High - Close)] /(High - Low) = Money Flow Multiplier
			m := ((cp.Get(ADJ_CLOSE) - cp.Get(2)) - (cp.Get(HIGH) - cp.Get(ADJ_CLOSE))) / hl
			prices.DataRows[i].Set(mf, m*cp.Get(5))
		}
	}
	for i := period; i < prices.Rows; i++ {
		s1 := 0.0
		s2 := 0.0
		for j := i - period; j < i; j++ {
			//21 Period Sum of Money Flow Volume / 21 Period Sum of Volume = 21 Period CMF
			s1 += prices.DataRows[j].Get(mf)
			s2 += prices.DataRows[j].Get(5)
		}
		if s2 != 0.0 {
			prices.DataRows[i].Set(ret, s1/s2)
		}
	}
	prices.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// STC
// -----------------------------------------------------------------------
/* 23,50,10,3,3*/
func STC(prices *Matrix, short, long, cycle, firstLength, secondLength int) int {
	// 0 = STC
	ret := prices.AddNamedColumn("STC")
	// macd = ema(src, fastLength) - ema(src, slowLength)
	es := EMA(prices, short, 4)
	el := EMA(prices, long, 4)
	macd := prices.Apply(func(mr MatrixRow) float64 {
		return mr.Get(es) - mr.Get(el)
	})
	//k = nz(fixnan(stoch(macd, macd, macd, cycleLength)))
	k := Stoch(prices, cycle, macd)
	// d = ema(k, d1Length)
	d := EMA(prices, firstLength, k)
	// kd = nz(fixnan(stoch(d, d, d, cycleLength)))
	kd := Stoch(prices, cycle, d)
	//stc = ema(kd, d2Length)
	stc := EMA(prices, secondLength, kd)

	for i := 0; i < prices.Rows; i++ {
		c := &prices.DataRows[i]
		// stc := max(min(stc, 100), 0)
		cv := c.Get(stc)
		if cv > 100.0 {
			cv = 100.0
		}
		if cv < 0.0 {
			cv = 0.0
		}
		c.Set(ret, cv)
	}
	prices.RemoveColumns(7)
	return ret
}

func FindMajorLevels(prices *Matrix, threshold float64) *MajorLevels {
	levels := NewMajorLevels(threshold)
	sp := prices.FindSwingPoints()
	for i := len(sp) - 1; i >= 0; i-- {
		cp := sp[i]
		cnt := 0
		for j := 0; j < len(sp); j++ {
			if i != j {
				d := (cp.Value/sp[j].Value - 1.0) * 100.0
				if m.Abs(d) < 2.0 {
					cnt++
				}
			}
		}
		levels.Add(cp.Value, cnt, cp.Timestamp)
	}
	return levels
}

func FindMajorGaps(prices *Matrix, threshold float64) *MajorLevels {
	levels := NewMajorLevels(threshold)
	gi := GAP(prices)
	for i := 0; i < prices.Rows; i++ {
		cp := prices.DataRows[i]
		if cp.Get(gi) > threshold || cp.Get(gi) < -threshold {
			levels.Add(cp.Get(gi+1), 1, cp.Key)
			levels.Add(cp.Get(gi+2), 1, cp.Key)
		}
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	return levels
}

func FindInsideBarMajorLevels(prices *Matrix, threshold float64) *MajorLevels {
	levels := NewMajorLevels(1.0)
	for i := prices.Rows - 2; i >= 0; i-- {
		cur := prices.DataRows[i]
		cnt := 0
		for j := i + 1; j < prices.Rows; j++ {
			th := prices.DataRows[j]
			if cur.Get(HIGH) > th.Get(HIGH) && cur.Get(2) < th.Get(2) {
				cnt++
			} else {
				upper := th.Get(0)
				lower := th.Get(ADJ_CLOSE)
				if upper < lower {
					upper = th.Get(ADJ_CLOSE)
					lower = th.Get(0)
				}
				if cur.Get(HIGH) > upper && cur.Get(2) < lower {
					cnt++
				} else {
					break
				}
			}
		}
		if cnt > 2 {
			m := (cur.Get(HIGH) + cur.Get(2) + cur.Get(ADJ_CLOSE)) / 3.0
			levels.Add(m, cnt, cur.Key)
			i -= cnt + 1
		}
	}
	return levels
}

func FindFibonacciLevels(prices *Matrix, lookback int) *MajorLevels {
	sp := prices.FindSwingPoints()
	levels := NewMajorLevels(0.0)
	hsl := sp.FilterByType(High)
	hi := len(hsl) - 1
	for i, p := range hsl {
		if p.Type == HigherHigh {
			hi = i
		}
	}
	lsl := sp.FilterByType(Low)
	li := len(lsl) - 1
	for i, p := range hsl {
		if p.Type == LowerLow {
			li = i
		}
	}
	if hi >= 0 && li >= 0 {
		cur := prices.Last().Key
		hs := hsl[hi]
		ls := lsl[li]
		pp := (hs.Value + ls.Value) / 2.0
		levels.Add(pp, 1, cur)
		diff := hs.Value - ls.Value
		levels.Add(pp+diff*0.382, 1, cur)
		levels.Add(pp+diff*0.618, 1, cur)
		levels.Add(pp-diff*0.382, 1, cur)
		levels.Add(pp-diff*0.618, 1, cur)
	}
	return levels
}

// -----------------------------------------------------------------------
//
//	Choppiness
//
// -----------------------------------------------------------------------
// https://medium.com/codex/detecting-ranging-and-trending-markets-with-choppiness-index-in-python-1942e6450b58
// https://www.tradingview.com/support/solutions/43000501980-choppiness-index-chop/
func Choppiness(prices *Matrix, days int) int {
	// 0 = Choppiness
	ret := prices.AddColumn()
	atr := ATR(prices, 1)
	// CI14 = 100 * LOG10 [14D ATR1 SUM/(14D HIGHH - 14D LOWL)] / LOG10(14)
	for i := days; i < prices.Rows; i++ {
		sum := 0.0
		for j := 0; j < days; j++ {
			cp := prices.DataRows[i-j]
			ca := cp.Get(atr)
			sum += ca
		}
		h := prices.FindMaxBetween(1, i-days+1, days)
		l := prices.FindMinBetween(2, i-days+1, days)
		if h != l {
			ci := 100.0 * m.Log10(sum/(h-l)) / m.Log10(float64(days))
			prices.DataRows[i].Set(ret, ci)
		}
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	return ret
}

func Slope(prices *Matrix, days, lookback int, f MAFunc) int {
	// 0 = SMA Slope
	ret := prices.AddColumn()
	steps := float64(lookback)
	si := f(prices, days, ADJ_CLOSE)
	for i := lookback; i < prices.Rows; i++ {
		cur := prices.DataRows[i].Get(si)
		prev := prices.DataRows[i-lookback].Get(si)
		prices.DataRows[i].Set(ret, (cur-prev)/steps)
	}
	prices.RemoveColumn()
	return ret
}

func Spread(prices *Matrix, lookback int) int {
	// 0 = Spread 1 = Spread Stoch K 2 = Spread Stoch K
	ret := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		cur := prices.DataRows[i]
		prices.DataRows[i].Set(ret, cur.Get(HIGH)-cur.Get(LOW))
	}
	StochasticExt(prices, lookback, 3, ret, ret, ret)
	return ret
}

// -----------------------------------------------------------------------
//
//	Guppy GMMA
//
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/g/guppy-multiple-moving-average.asp
func GMMA(prices *Matrix) int {
	// 0 = GMMA 1 = Signal
	ret := prices.AddColumn()
	fe := make([]int, 0)
	fast := []int{3, 5, 8, 10, 12, 15}
	for _, f := range fast {
		fe = append(fe, EMA(prices, f, 4))
	}
	fa := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		sum := 0.0
		for j := 0; j < len(fe); j++ {
			sum += prices.DataRows[i].Get(fe[j])
		}
		prices.DataRows[i].Set(fa, sum)
	}

	se := make([]int, 0)
	slow := []int{30, 35, 40, 45, 50, 60}
	for _, f := range slow {
		se = append(se, EMA(prices, f, 4))
	}
	sa := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		sum := 0.0
		for j := 0; j < len(se); j++ {
			sum += prices.DataRows[i].Get(se[j])
		}
		prices.DataRows[i].Set(sa, sum)
	}

	for i := 0; i < prices.Rows; i++ {
		if prices.DataRows[i].Get(sa) != 0.0 {
			d := (prices.DataRows[i].Get(fa) - prices.DataRows[i].Get(sa)) / prices.DataRows[i].Get(sa) * 100.0
			prices.DataRows[i].Set(ret, d)
		}
	}
	for i := 0; i < 14; i++ {
		prices.RemoveColumn()
	}
	EMA(prices, 13, ret)
	return ret
}

// -----------------------------------------------------------------------
//
//	Volatility
//
// -----------------------------------------------------------------------
func Volatility(m *Matrix, lookback, field int) int {
	// 0 = Volatility
	ret := m.AddColumn()
	for i := lookback; i < m.Rows; i++ {
		s := standardAbbreviation(m, m.DataRows[i].Get(field), i-lookback, lookback)
		m.DataRows[i].Set(ret, s)
	}
	return ret
}

// -----------------------------------------------------------------------
//
//	z-Score
//
// -----------------------------------------------------------------------
func ZScore(m *Matrix, lookback, field int) int {
	return ZNormalization(m, lookback, field)
}

// https://kaabar-sofien.medium.com/using-z-score-in-trading-a-python-study-5f4b21b41aa0
// z = (close - ta.sma(close, len)) / ta.stdev(close, len)
func ZNormalization(m *Matrix, lookback, field int) int {
	// 0 = Z-Score
	ret := m.AddColumn()
	sma := SMA(m, lookback, field)
	std := m.StdDev(field, lookback)
	for i := lookback; i < m.Rows; i++ {
		if m.DataRows[i].Get(std) != 0.0 {
			s := (m.DataRows[i].Get(field) - m.DataRows[i].Get(sma)) / m.DataRows[i].Get(std)
			m.DataRows[i].Set(ret, s)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	z-Normalization Bollinger Bands
//
// -----------------------------------------------------------------------
// https://medium.com/superalgos/normalization-of-oscillating-indicators-to-create-dynamic-overbought-and-oversold-levels-338d4ef72914
func ZNormalizationBollinger(m *Matrix, period, field int) int {
	// 0 = zScore 1 = Upper 2 = Lower 3 = Mid
	zi := ZNormalization(m, period, field)
	BollingerBandExt(m, zi, period, 2.0, 2.0)
	return zi
}

// -----------------------------------------------------------------------
// Smoothed candles - using SMA on every entry
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/chart-pattern-recognition-in-python-fda5c7efe7bf
func SmoothedCandles(m *Matrix, lookback int) int {
	// 0 = Open 1 = High 2 = Low 3 = Close 4 = AdjClose 5 = Volume
	oi := SMA(m, lookback, 0)
	for i := 1; i < 6; i++ {
		SMA(m, lookback, i)
	}
	return oi
}

// -----------------------------------------------------------------------
// Triple EMA Categorization
// -----------------------------------------------------------------------
func TripleEMA(m *Matrix, l1, l2, l3 int) int {
	// 0 = Normalized Value
	ret := m.AddColumn()
	e1 := EMA(m, l1, 4)
	e2 := EMA(m, l2, 4)
	e3 := EMA(m, l3, 4)
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		p := m.DataRows[i-1]
		sum := 0.0
		if p.Get(e1) < c.Get(e1) {
			sum += 1.0
		}
		if p.Get(e2) < c.Get(e2) {
			sum += 1.0
		}
		if p.Get(e3) < c.Get(e3) {
			sum += 1.0
		}
		if c.Get(e1) > c.Get(e2) {
			sum += 1.0
		}
		if c.Get(e1) > c.Get(e3) {
			sum += 1.0
		}
		if c.Get(e2) > c.Get(e3) {
			sum += 1.0
		}
		if c.Get(ADJ_CLOSE) > c.Get(e1) {
			sum += 1.0
		}
		if c.Get(ADJ_CLOSE) > c.Get(22) {
			sum += 1.0
		}
		if c.Get(ADJ_CLOSE) > c.Get(e3) {
			sum += 1.0
		}
		m.DataRows[i].Set(ret, sum/9.0)
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// TRIX
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/t/trix.asp
func TRIX(prices *Matrix, lookback int) int {
	// 0 = TRIX 1 = Signal
	ret := prices.AddColumn()
	li := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		prices.DataRows[i].Set(li, m.Log(c.Get(ADJ_CLOSE)))
	}
	e1 := EMA(prices, lookback, li)
	e2 := EMA(prices, lookback, e1)
	e3 := EMA(prices, lookback, e2)
	for i := 1; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		p := prices.DataRows[i-1]
		if p.Get(e3) != 0.0 {
			d := (c.Get(e3) - p.Get(e3)) / p.Get(e3) * 10000.0
			prices.DataRows[i].Set(ret, d)
		}
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	EMA(prices, 9, ret)
	return ret
}

/*
//@version=4
//Basic Hull Ma Pack tinkered by InSilico
study("Hull Suite by InSilico", overlay=true)

//INPUT
src = input(close, title="Source")
modeSwitch = input("Hma", title="Hull Variation", options=["Hma", "Thma", "Ehma"])
length = input(55, title="Length(180-200 for floating S/R , 55 for swing entry)")
lengthMult = input(1.0, title="Length multiplier (Used to view higher timeframes with straight band)")

useHtf = input(false, title="Show Hull MA from X timeframe? (good for scalping)")
htf = input("240", title="Higher timeframe", type=input.resolution)

switchColor = input(true, "Color Hull according to trend?")
candleCol = input(false,title="Color candles based on Hull's Trend?")
visualSwitch  = input(true, title="Show as a Band?")
thicknesSwitch = input(1, title="Line Thickness")
transpSwitch = input(40, title="Band Transparency",step=5)

//FUNCTIONS
//HMA
HMA(_src, _length) =>  wma(2 * wma(_src, _length / 2) - wma(_src, _length), round(sqrt(_length)))
//EHMA
EHMA(_src, _length) =>  ema(2 * ema(_src, _length / 2) - ema(_src, _length), round(sqrt(_length)))
//THMA
THMA(_src, _length) =>  wma(wma(_src,_length / 3) * 3 - wma(_src, _length / 2) - wma(_src, _length), _length)

//SWITCH
Mode(modeSwitch, src, len) =>
      modeSwitch == "Hma"  ? HMA(src, len) :
      modeSwitch == "Ehma" ? EHMA(src, len) :
      modeSwitch == "Thma" ? THMA(src, len/2) : na

//OUT
_hull = Mode(modeSwitch, src, int(length * lengthMult))
HULL = useHtf ? security(syminfo.ticker, htf, _hull) : _hull
MHULL = HULL[0]
SHULL = HULL[2]

//COLOR
hullColor = switchColor ? (HULL > HULL[2] ? #00ff00 : #ff0000) : #ff9800

//PLOT
///< Frame
Fi1 = plot(MHULL, title="MHULL", color=hullColor, linewidth=thicknesSwitch, transp=50)
Fi2 = plot(visualSwitch ? SHULL : na, title="SHULL", color=hullColor, linewidth=thicknesSwitch, transp=50)
alertcondition(crossover(MHULL, SHULL), title="Hull trending up.", message="Hull trending up.")
alertcondition(crossover(SHULL, MHULL), title="Hull trending down.", message="Hull trending down.")
///< Ending Filler
fill(Fi1, Fi2, title="Band Filler", color=hullColor, transp=transpSwitch)
///BARCOLOR
barcolor(color = candleCol ? (switchColor ? hullColor : na) : na)
*/

/*
length = input(20, title="BB Length")
mult = input(2.0,title="BB MultFactor")
lengthKC=input(20, title="KC Length")
multKC = input(1.5, title="KC MultFactor")

useTrueRange = input(true, title="Use TrueRange (KC)", type=bool)

// Calculate BB
source = close
basis = sma(source, length)
dev = multKC * stdev(source, length)
upperBB = basis + dev
lowerBB = basis - dev

// Calculate KC
ma = sma(source, lengthKC)
range = useTrueRange ? tr : (high - low)
rangema = sma(range, lengthKC)
upperKC = ma + rangema * multKC
lowerKC = ma - rangema * multKC

sqzOn  = (lowerBB > lowerKC) and (upperBB < upperKC)
sqzOff = (lowerBB < lowerKC) and (upperBB > upperKC)
noSqz  = (sqzOn == false) and (sqzOff == false)

val = linreg(source  -  avg(avg(highest(high, lengthKC), lowest(low, lengthKC)),sma(close,lengthKC)),
            lengthKC,0)

bcolor = iff( val > 0,
            iff( val > nz(val[1]), lime, green),
            iff( val < nz(val[1]), red, maroon))
scolor = noSqz ? blue : sqzOn ? black : gray
plot(val, color=bcolor, style=histogram, linewidth=4)
plot(0, color=scolor, style=cross, linewidth=2)

*/
// -----------------------------------------------------------------------
// Squeeze Momentum
// -----------------------------------------------------------------------
// https://medium.com/geekculture/implementing-the-most-popular-indicator-on-tradingview-using-python-239d579412ab
// https://school.stockcharts.com/doku.php?id=technical_indicators:ttm_squeeze
func SqueezeMomentum(prices *Matrix, length int, std, k1, k2, k3 float64) int {
	// 0 = State 1 = Line
	ret := prices.AddNamedColumn("Squeeze")
	li := prices.AddNamedColumn("SQ-MOM")
	bb := BollingerBand(prices, length, std, std)
	ki1 := Keltner(prices, length, length, k1)
	ki2 := Keltner(prices, length, length, k2)
	ki3 := Keltner(prices, length, length, k3)
	//SQUEEZE CONDITIONS
	//NoSqz = BB_lower < KC_lower_low or BB_upper > KC_upper_low //NO SQUEEZE: GREEN
	//LowSqz = BB_lower >= KC_lower_low or BB_upper <= KC_upper_low //LOW COMPRESSION: BLACK
	//MidSqz = BB_lower >= KC_lower_mid or BB_upper <= KC_upper_mid //MID COMPRESSION: RED
	//HighSqz = BB_lower >= KC_lower_high or BB_upper <= KC_upper_high //HIGH COMPRESSION: ORANGE

	//kc := Keltner(prices, keltner, keltner, mulKC)
	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		//NoSqz = BB_lower < KC_lower_low or BB_upper > KC_upper_low //NO SQUEEZE
		if c.Get(bb+1) < c.Get(ki3+1) || c.Get(bb) > c.Get(ki3+1) {
			prices.DataRows[i].Set(ret, 0.0)
		}
		//LowSqz = BB_lower >= KC_lower_low or BB_upper <= KC_upper_low //LOW COMPRESSION
		if c.Get(bb+1) >= c.Get(ki3+1) || c.Get(bb) < c.Get(ki3+1) {
			prices.DataRows[i].Set(ret, 0.25)
		}
		//MidSqz = BB_lower >= KC_lower_mid or BB_upper <= KC_upper_mid //MID COMPRESSION
		if c.Get(bb+1) >= c.Get(ki2+1) || c.Get(bb) <= c.Get(ki2+1) {
			prices.DataRows[i].Set(ret, 0.75)
		}
		//HighSqz = BB_lower >= KC_lower_high or BB_upper <= KC_upper_high //HIGH COMPRESSION
		if c.Get(bb+1) >= c.Get(ki1+1) || c.Get(bb) <= c.Get(ki1) {
			prices.DataRows[i].Set(ret, 1.0)
		}
	}
	/*
		Step 3: Calculate the highest high in the last 20 periods.
		Step 4: Calculate the lowest low in the last 20 periods.
		Step 5: Find the mean between the two above results.
		Step 6: Calculate a 20-period simple moving average on the closing price
		Step 7: Calculate the delta between the closing price and the mean between the result from step 5 and 6.
	*/
	di := prices.AddColumn()
	for i := length; i < prices.Rows; i++ {
		h, l := prices.FindHighLowIndex(i-length, length)
		d := (prices.DataRows[h].Get(ADJ_CLOSE) + prices.DataRows[l].Get(ADJ_CLOSE)) / 2.0
		prices.DataRows[i].Set(di, d)

	}

	si := SMA(prices, length, 4)
	for i := length; i < prices.Rows; i++ {
		d := prices.DataRows[i].Get(ADJ_CLOSE) - (prices.DataRows[i].Get(di)+prices.DataRows[i].Get(si))/2.0
		prices.DataRows[i].Set(li, d)

	}
	lr := LinearRegression(prices, length)
	prices.CopyColumn(lr, li)
	prices.RemoveColumns(16)
	return ret
}

func TTMSqueeze(prices *Matrix, length int, std, kc float64) int {
	// 0 = State 1 = Line
	ret := prices.AddNamedColumn("TTM-Squeeze")
	//mom := prices.AddNamedColumn("TTM-Hist")
	bb := BollingerBand(prices, length, std, std)
	ki := Keltner(prices, length, length, kc)

	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		if c.Get(bb+1) < c.Get(ki+1) && c.Get(bb) > c.Get(ki) {
			prices.DataRows[i].Set(ret, 0.0)
		}
		if c.Get(bb+1) > c.Get(ki+1) && c.Get(bb) < c.Get(ki) {
			prices.DataRows[i].Set(ret, 1.0)
		}
	}
	// (Highest high in 20 periods + lowest low in 20 periods) / 2
	// Close - ( (Donchian midline + SMA) / 2 )
	// linear regression on this
	prices.RemoveColumns(6)
	return ret
}

/*
//@version=2
study("SMA slope", overlay=false)
rad2degree=180/3.14159265359  //pi
iSMA = input(defval=50,title="Period SMA",type=integer)
iBarsBack=input(defval=10,title="Bars Back",type=integer)
hline(0)
sma2sample=sma(close,iSMA)
slopeD=rad2degree*atan((sma2sample[0]-nz(sma2sample[iBarsBack]))/iBarsBack)
plot(slopeD,color=black)
*/
// -----------------------------------------------------------------------
// Spread Range Relation / Small Range versus Large Range relation
// -----------------------------------------------------------------------
// https://www.reddit.com/r/thinkorswim/comments/p8b8ti/how_to_scan_for_volatility_contraction_pattern/
func SpreadRangeRelation(prices *Matrix, lookback int) int {
	ret := prices.AddColumn()
	sri := prices.AddColumn()
	for i := 1; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		p := prices.DataRows[i-1]
		sr := m.Max(c.Get(HIGH), p.Get(ADJ_CLOSE)) - m.Min(c.Get(2), p.Get(ADJ_CLOSE))
		prices.DataRows[i].Set(sri, sr)
	}
	sui := prices.AddColumn()
	for i := lookback - 1; i < prices.Rows; i++ {
		sum := 0.0
		for j := 0; j < lookback; j++ {
			sum += prices.DataRows[i-j].Get(sri)
		}
		prices.DataRows[i].Set(sui, sum)
	}

	for i := lookback; i < prices.Rows; i++ {
		h, l := prices.FindHighestHighLowestLow(i-lookback, lookback)
		lr := h - l
		r := m.Log(prices.DataRows[i].Get(sui)/lr) / m.Log(float64(lookback))
		prices.DataRows[i].Set(ret, r)
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// Historical Volatility
// -----------------------------------------------------------------------
// https://corporatefinanceinstitute.com/resources/knowledge/trading-investing/historical-volatility-hv/
func HistoricalVolatility(prices *Matrix, lookback int) int {
	ret := prices.AddColumn()
	si := SMA(prices, lookback, 4)
	di := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		sr := c.Get(ADJ_CLOSE) - c.Get(si)
		prices.DataRows[i].Set(di, sr*sr)
	}
	for i := lookback; i < prices.Rows; i++ {
		sum := 0.0
		for j := 0; j < lookback; j++ {
			sum += prices.DataRows[i-j].Get(di)
		}
		v := sum / float64(lookback)
		prices.DataRows[i].Set(ret, m.Sqrt(v))
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	return ret
}

func MinerviniScore(prices *Matrix) int {
	ret := prices.AddColumn()
	sma50 := SMA(prices, 50, 4)
	sma150 := SMA(prices, 150, 4)
	sma200 := SMA(prices, 200, 4)
	rsi := RSI(prices, 14, 4)
	for i := 20; i < prices.Rows; i++ {
		cp := prices.DataRows[i]
		price := cp.Get(ADJ_CLOSE)
		low, high := prices.FindMinMaxBetween(4, prices.Rows-250, prices.Rows)
		cnt := 0.0
		if price > cp.Get(sma150) && price > cp.Get(sma200) {
			cnt += 1.0
		}
		if cp.Get(sma150) > cp.Get(sma200) {
			cnt += 1.0
		}
		if cp.Get(sma200) > prices.DataRows[i-20].Get(sma200) {
			cnt += 1.0
		}
		if cp.Get(sma50) > cp.Get(sma150) && cp.Get(sma50) > cp.Get(sma200) {
			cnt += 1.0
		}
		if price > cp.Get(sma50) {
			cnt += 1.0
		}
		if price > (low * 1.3) {
			cnt += 1.0
		}
		if price > (high * 0.75) {
			cnt += 1.0
		}
		if cp.Get(rsi) >= 70.0 {
			cnt += 1.0
		}
		prices.DataRows[i].Set(ret, cnt/8.0)
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	return ret
}

func PriceMovingAverageDistance(prices *Matrix, f MAFunc, length int) int {
	// 0 = Price MA Diff 1 = Diff Percentage
	ret := prices.AddNamedColumn("PMADiff")
	per := prices.AddNamedColumn("PMADiffPer")
	si := f(prices, length, ADJ_CLOSE)
	for i := length; i < prices.Rows; i++ {
		prices.DataRows[i].Set(ret, prices.DataRows[i].Get(ADJ_CLOSE)-prices.DataRows[i].Get(si))
		prices.DataRows[i].Set(per, ChangePercentage(prices.DataRows[i].Get(ADJ_CLOSE), prices.DataRows[i].Get(si)))
	}
	prices.RemoveColumn()
	return ret
}

func RS(prices *Matrix, index *Matrix) int {
	// 0 = Strength
	ret := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		row := index.FindRow(prices.DataRows[i].Key)
		if row != nil && row.Get(ADJ_CLOSE) != 0.0 {
			prices.DataRows[i].Set(ret, prices.DataRows[i].Get(ADJ_CLOSE)/row.Get(ADJ_CLOSE)*1000.0)
		}
	}
	return ret
}

func IntradayIntensityTrend(prices *Matrix) int {
	// 0 = IIT
	ret := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		u := 2.0*c.Get(ADJ_CLOSE) - c.Get(HIGH) - c.Get(LOW)
		l := (c.Get(HIGH) - c.Get(LOW)) * c.Get(VOLUME)
		if l != 0.0 {
			prices.DataRows[i].Set(ret, u/l*10000000.0)
		}

	}
	return ret
}

func PercentRank(prices *Matrix, period, field int) int {
	// 0 = PercentRank
	ret := prices.AddNamedColumn("PercentRank")
	for i := period; i < prices.Rows; i++ {
		ref := prices.DataRows[i].Get(field)
		cnt := 0.0
		for j := 1; j < period; j++ {
			if prices.DataRows[i-j].Get(field) > ref {
				cnt += 1.0
			}
		}
		prices.DataRows[i].Set(ret, 100.0*cnt/float64(period))
	}
	return ret
}

func TSV(prices *Matrix) int {
	period := 13
	ret := prices.AddNamedColumn("TSV")
	tmp := prices.AddColumn()
	for i := 1; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		p := prices.DataRows[i-1]
		d := c.Get(ADJ_CLOSE) - p.Get(ADJ_CLOSE)
		sig := 1.0
		if d < 0.0 {
			sig = -1.0
		}
		c.Set(tmp, sig*m.Abs(d)*c.Get(VOLUME))
	}
	for i := period; i < prices.Rows; i++ {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices.DataRows[i-j].Get(tmp)
		}
		prices.DataRows[i].Set(ret, sum)
	}

	/*
			l  = input(13, title="Length")
		l_ma = input(7, title="MA Length")

		//t = sum(close>close[1]?volume*close-close[1]:close<close[1]?(volume*-1)*close-close:0,l)
		// previous line is non sensical. The correct version follows
		t = sum(close>close[1]?volume*(close-close[1]):close<close[1]?volume*(close-close[1]):0,l)
		m = sma(t ,l_ma )

		plot(t, color=red, style=histogram)
		plot(m, color=green)
	*/
	//TSV=(Sum( IIf( C > Ref(C,-1), V * ( C-Ref(C,-1) ),
	// IIf( C < Ref(C,-1),-V * ( C-Ref(C,-1) ), 0 ) ) ,18));
	return ret
}

// round_(val) => val > .99 ? .999 : val < -.99 ? -.999 : val
func myRound(value float64) float64 {
	if value > 0.99 {
		return 0.999
	}
	if value < -0.99 {
		return -0.999
	}
	return value
}

func FisherTransform(prices *Matrix, period int) int {
	// 0 = FT
	ret := prices.AddNamedColumn("FT")
	trig := prices.AddNamedColumn("FT-Trigger")
	tmp := prices.AddColumn()
	hl2 := HL2(prices)
	for i := period; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		p := prices.DataRows[i-1]
		l, h := prices.FindMinMaxBetween(hl2, i-period, period)
		//high_ = ta.highest(hl2, len)
		//low_ = ta.lowest(hl2, len)
		//value := round_(.66 * ((hl2 - low_) / (high_ - low_) - .5) + .67 * nz(value[1]))
		value := myRound(0.66*((c.Get(hl2)-l)/(h-l)-0.5) + 0.67*p.Get(tmp))
		prices.DataRows[i].Set(tmp, value)
		//fish1 := .5 * math.log((1 + value) / (1 - value)) + .5 * nz(fish1[1])
		if value != 1.0 {
			prices.DataRows[i].Set(ret, 0.5*math.Log((1.0+value)/(1.0-value))+0.5*p.Get(ret))
		}
		prices.DataRows[i].Set(trig, prices.DataRows[i-1].Get(ret))
	}
	prices.RemoveColumns(1)
	return ret
}

func LaguerreRSI(prices *Matrix, alpha float64) int {
	ret := prices.AddNamedColumn("Laguerre")
	l0s := prices.AddColumn()
	l1s := prices.AddColumn()
	l2s := prices.AddColumn()
	l3s := prices.AddColumn()
	gamma := 1.0 - alpha
	for i := 1; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		p := prices.DataRows[i-1]
		// L0 := (1-gamma) * src + gamma * nz(L0[1])
		l0 := (1.0-gamma)*c.Get(ADJ_CLOSE) + gamma*p.Get(l0s)
		c.Set(l0s, l0)
		// L1 := -gamma * L0 + nz(L0[1]) + gamma * nz(L1[1])
		l1 := -gamma*l0 + p.Get(l0s) + gamma*p.Get(l1s)
		c.Set(l1s, l1)
		// L2 := -gamma * L1 + nz(L1[1]) + gamma * nz(L2[1])
		l2 := -gamma*c.Get(l1s) + l1 + gamma*p.Get(l2s)
		c.Set(l2s, l2)
		// L3 := -gamma * L2 + nz(L2[1]) + gamma * nz(L3[1])
		l3 := -gamma*l2 + p.Get(l2s) + gamma*p.Get(l3s)
		c.Set(l3s, l3)

		// (L0>L1 ? L0-L1 : 0) + (L1>L2 ? L1-L2 : 0) + (L2>L3 ? L2-L3 : 0)
		cu := 0.0
		if l0 > l1 {
			cu = l0 - l1
		}
		if l1 > l2 {
			cu += (l1 - l2)
		}
		if l2 > l3 {
			cu += (l2 - l3)
		}
		//cd= (L0<L1 ? L1-L0 : 0) + (L1<L2 ? L2-L1 : 0) + (L2<L3 ? L3-L2 : 0)
		cd := 0.0
		if l0 < l1 {
			cd = l1 - l0
		}
		if l1 < l2 {
			cd += (l2 - l1)
		}
		if l2 < l3 {
			cd += (l3 - l2)
		}

		// temp= cu+cd==0 ? -1 : cu+cd
		temp := -1.0
		if cu+cd != 0.0 {
			temp = cu + cd
		}
		// LaRSI=temp==-1 ? 0 : cu/temp
		if temp != -1 {
			c.Set(ret, cu/temp*100.0)
		}
	}
	prices.RemoveColumns(4)
	return ret
}

// https://github.com/MathisWellmann/go_ehlers_indicators
// LaguerreFilter from paper: http://mesasoftware.com/papers/TimeWarp.pdf
func LaguerreFilter(prices *Matrix, gamma float64) int {
	ret := prices.AddNamedColumn("Laguerre")
	l0s := prices.AddColumn()
	l1s := prices.AddColumn()
	l2s := prices.AddColumn()
	l3s := prices.AddColumn()
	for i := 1; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		p := prices.DataRows[i-1]
		c.Set(l0s, (1.0-gamma)*c.Get(ADJ_CLOSE)+gamma*p.Get(l0s))
		c.Set(l1s, gamma*c.Get(l0s)+p.Get(l0s)+gamma*p.Get(l1s))
		c.Set(l2s, -gamma*c.Get(l1s)+c.Get(l1s)+gamma*p.Get(l2s))
		c.Set(l3s, -gamma*c.Get(l2s)+p.Get(l2s)+gamma*p.Get(l3s))

		c.Set(ret, (c.Get(l0s)+2.0*c.Get(l1s)+2.0*c.Get(l2s)+c.Get(l3s))/6.0)
		// firs[i] = (vals[i] + 2*vals[i-1] + 2*vals[i-2] + vals[i-3]) / 6
	}
	prices.RemoveColumns(4)
	return ret
}

func LaguerreFilterDefault(prices *Matrix, gamma float64) int {
	return LaguerreFilter(prices, 0.8)
}

// -------------------------------------------------------------
// APZ
// -------------------------------------------------------------
// https://www.investopedia.com/articles/trading/10/adaptive-price-zone-indicator-explained.asp
func APZ(prices *Matrix, period int, dev float64) int {
	// 0 = Upper 1 = Lower
	upper := prices.AddNamedColumn("Upper")
	lower := prices.AddNamedColumn("Lower")
	//nP = ceil(sqrt(nPeriods))
	np := int(m.Ceil(m.Sqrt(float64(period))))
	// xVal1 = ema(ema(close,nP), nP)
	e1 := EMA(prices, np, ADJ_CLOSE)
	e2 := EMA(prices, np, e1)
	// xVal2 = ema(ema(xHL, nP), nP)
	hl := prices.Apply(func(mr MatrixRow) float64 {
		return mr.Get(1) - mr.Get(2)
	})
	e3 := EMA(prices, np, hl)
	e4 := EMA(prices, np, e3)
	// UpBand = nBandPct*xVal2 + xVal1
	prices.ApplyRow(upper, func(mr MatrixRow) float64 {
		return mr.Get(e2) + dev*mr.Get(e4)
	})
	// DnBand = xVal1 - nBandPct*xVal2
	prices.ApplyRow(lower, func(mr MatrixRow) float64 {
		return mr.Get(e2) - dev*mr.Get(e4)
	})
	prices.RemoveColumns(5)
	return upper
}

/*
pine_alma(series, windowsize, offset, sigma) =>

	m = offset * (windowsize - 1)
	//m = math.floor(offset * (windowsize - 1)) // Used as m when math.floor=true
	s = windowsize / sigma
	norm = 0.0
	sum = 0.0
	for i = 0 to windowsize - 1
	    weight = math.exp(-1 * math.pow(i - m, 2) / (2 * math.pow(s, 2)))
	    norm := norm + weight
	    sum := sum + series[windowsize - i - 1] * weight
	sum / norm
*/
func ALMA(prices *Matrix, windowSize int, offset, sigma float64) int {
	// 0 = ALMA
	ret := prices.AddNamedColumn("ALMA")
	mv := offset * float64((windowSize - 1))
	//m = math.floor(offset * (windowsize - 1)) // Used as m when math.floor=true
	s := float64(windowSize) / sigma
	for i := windowSize; i < prices.Rows; i++ {
		norm := 0.0
		sum := 0.0
		for j := 0; j < windowSize; j++ {
			weight := m.Exp(-1 * m.Pow(float64(j)-mv, 2.0) / (2.0 * m.Pow(s, 2.0)))
			norm += weight
			sum += prices.DataRows[i+windowSize-j].Get(ADJ_CLOSE) * weight
		}
		prices.DataRows[i].Set(ret, sum/norm)
	}
	return ret
}

// https://kaabar-sofien.medium.com/using-the-distance-to-trade-b23c6e37552e
func RAD(prices *Matrix, ma, period int) int {
	// 0 = RAD
	ret := prices.AddNamedColumn("RAD")
	si := EMA(prices, ma, ADJ_CLOSE)
	di := prices.Apply(func(r MatrixRow) float64 {
		return r.Get(ADJ_CLOSE) - r.Get(si)
	})
	//ri := StochasticExt(prices, period, period/2, di, di, di)
	ri := RSI(prices, period, di)
	prices.CopyColumn(ri, ret)
	prices.RemoveColumns(3)
	return ret
}

func AnalyzeTrend(prices *Matrix, field int) int {
	ret := prices.AddNamedColumn("Count")
	ci := prices.AddNamedColumn("Trend")
	dir := 0
	for i := 1; i < prices.Rows; i++ {
		cnt := 1
		if prices.DataRows[i].Get(field) >= 0.0 {
			dir = 1
		} else {
			dir = -1
		}
		for j := i - 1; j >= 0; j-- {
			cc := prices.DataRows[j].Get(field)
			cd := 0
			if cc >= 0.0 {
				cd = 1
			} else {
				cd = -1
			}
			if cd != dir {
				dir = cd
				break
			} else {
				cnt++
			}
		}
		prices.DataRows[i].Set(ret, float64(cnt))
		prices.DataRows[i].Set(ci, -1.0*float64(dir))
	}
	return ret
}

func AnalyzeTrendRange(prices *Matrix, field int, lower, upper float64) int {
	ret := prices.AddNamedColumn("Count")
	ci := prices.AddNamedColumn("Trend")
	dir := 0
	for i := 1; i < prices.Rows; i++ {
		cnt := 1
		if prices.DataRows[i].Get(field) >= upper {
			dir = 1
		} else if prices.DataRows[i].Get(field) <= lower {
			dir = -1
		} else {
			dir = 0
		}
		if dir != 0 {
			for j := i - 1; j >= 0; j-- {
				cc := prices.DataRows[j].Get(field)
				cd := 0
				if cc >= upper {
					cd = 1
				} else if cc <= lower {
					cd = -1
				}
				if cd != dir {
					break
				} else {
					cnt++
				}
			}
			prices.DataRows[i].Set(ret, float64(cnt))
			prices.DataRows[i].Set(ci, float64(dir))
		}
	}
	return ret
}

// -------------------------------------------------------
// Volume Price Trend
// -------------------------------------------------------
func VPT(prices *Matrix) int {
	// 0 = VPT
	ret := prices.AddColumn()
	for i := 1; i < prices.Rows; i++ {
		// VPT = Previous VPT + Volume x (Today’s Close – Previous Close) / Previous Close
		c := &prices.DataRows[i]
		p := prices.DataRows[i-1]
		c.Set(ret, p.Get(ret)+c.Get(VOLUME)*(c.Get(ADJ_CLOSE)-p.Get(ADJ_CLOSE))/p.Get(ADJ_CLOSE))
	}
	return ret
}

// -------------------------------------------------------
// Time Weighted Moving Average
// -------------------------------------------------------
// https://blog.quantinsti.com/twap/
func TWAP(prices *Matrix, period int) int {
	// 0 = TWAP 1 = Upper 2 = Lower
	ret := prices.AddColumn()
	ui := prices.AddColumn()
	li := prices.AddColumn()
	sum := prices.Apply(func(mr MatrixRow) float64 {
		return (mr.Get(OPEN) + mr.Get(HIGH) + mr.Get(LOW) + mr.Get(ADJ_CLOSE)) / 4.0
	})
	si := SMA(prices, period, sum)
	prices.CopyColumn(si, ret)

	sti := prices.StdDev(si, period)
	prices.ApplyRow(ui, func(mr MatrixRow) float64 {
		return mr.Get(si) + mr.Get(sti)
	})
	prices.ApplyRow(li, func(mr MatrixRow) float64 {
		return mr.Get(si) - mr.Get(sti)
	})

	prices.RemoveColumns(2)
	return ret
}

func PriceTWAP(prices *Matrix, period int) int {
	ret := prices.AddColumn()
	twi := TWAP(prices, period)
	prices.ApplyRow(ret, func(mr MatrixRow) float64 {
		return (mr.Get(ADJ_CLOSE)/mr.Get(twi) - 1.0)
	})
	prices.RemoveColumns(3)
	return ret
}

func PSARTrend(prices *Matrix) int {
	ret := prices.AddNamedColumn("PSAR-Trend")
	pi := ParabolicSAR(prices)
	prices.ApplyRow(ret, func(mr MatrixRow) float64 {
		return (mr.Get(pi) - mr.Get(ADJ_CLOSE)) / mr.Get(ADJ_CLOSE) * 100.0
	})
	prices.RemoveColumns(2)
	return ret
}

func ParabolicSAR(prices *Matrix) int {
	// 0 = PSAR 1 = Trend
	const (
		psarAfStep = 0.02
		psarAfMax  = 0.20
	)

	psar := prices.AddNamedColumn("PSAR")
	trend := prices.AddNamedColumn("Trend")

	var af, ep float64

	prices.DataRows[0].Set(trend, -1.0)
	prices.DataRows[0].Set(psar, prices.DataRows[0].Get(HIGH))
	af = psarAfStep
	ep = prices.DataRows[0].Get(LOW)

	for i := 1; i < prices.Rows; i++ {
		c := &prices.DataRows[i]
		p := prices.DataRows[i-1]
		c.Set(psar, p.Get(psar)-(p.Get(psar)-ep)*af)

		if p.Get(trend) == -1.0 {
			c.Set(psar, m.Max(c.Get(psar), p.Get(HIGH)))
			if i > 1 {
				c.Set(psar, m.Max(c.Get(psar), prices.DataRows[i-2].Get(HIGH)))
			}

			if c.Get(HIGH) >= c.Get(psar) {
				c.Set(psar, ep)
			}
		} else {
			c.Set(psar, m.Min(c.Get(psar), p.Get(LOW)))
			if i > 1 {
				c.Set(psar, m.Min(c.Get(psar), prices.DataRows[i-2].Get(LOW)))
			}

			if c.Get(LOW) <= c.Get(psar) {
				c.Set(psar, ep)
			}
		}

		prevEp := ep

		if c.Get(psar) > c.Get(ADJ_CLOSE) {
			c.Set(trend, -1.0)
			ep = m.Min(ep, c.Get(LOW))
		} else {
			c.Set(trend, 1.0)
			ep = m.Max(ep, c.Get(HIGH))
		}

		if c.Get(trend) != p.Get(trend) {
			af = psarAfStep
		} else if prevEp != ep && af < psarAfMax {
			af += psarAfStep
		}
	}

	return psar
}

// The ChandelierExit function sets a trailing stop-loss based on the Average True Value (ATR).
/*
length = input(title="ATR Period", type=input.integer, defval=22)
mult = input(title="ATR Multiplier", type=input.float, step=0.1, defval=3.0)
showLabels = input(title="Show Buy/Sell Labels ?", type=input.bool, defval=true)
useClose = input(title="Use Close Price for Extremums ?", type=input.bool, defval=true)
highlightState = input(title="Highlight State ?", type=input.bool, defval=true)

atr = mult * atr(length)

longStop = (useClose ? highest(close, length) : highest(length)) - atr
longStopPrev = nz(longStop[1], longStop)
longStop := close[1] > longStopPrev ? max(longStop, longStopPrev) : longStop

shortStop = (useClose ? lowest(close, length) : lowest(length)) + atr
shortStopPrev = nz(shortStop[1], shortStop)
shortStop := close[1] < shortStopPrev ? min(shortStop, shortStopPrev) : shortStop

var int dir = 1
dir := close > shortStopPrev ? 1 : close < longStopPrev ? -1 : dir
*/
func ChandelierExit(prices *Matrix, period int, multiplier float64) int {
	// 0 = Long 1 = Short
	l := prices.AddNamedColumn("Long")
	s := prices.AddNamedColumn("Short")
	ai := ATR(prices, period)
	sh := SMA(prices, period, HIGH)
	lh := SMA(prices, period, LOW)
	for i := 0; i < prices.Rows; i++ {
		c := &prices.DataRows[i]
		c.Set(l, c.Get(sh)-multiplier*c.Get(ai))
		c.Set(s, c.Get(lh)+multiplier*c.Get(ai))
	}
	prices.RemoveColumns(4)
	//Chandelier Exit Long = 22-Period SMA High - ATR(22) * 3
	//Chandelier Exit Short = 22-Period SMA Low + ATR(22) * 3
	return l
}

func StratClassification(prices *Matrix) int {
	// 1 = Inside 2 = 2 down 3 = 2 Up 4 = Outside
	r := prices.AddNamedColumn("Strat")
	for i := 1; i < prices.Rows; i++ {
		c := &prices.DataRows[i]
		p := &prices.DataRows[i-1]
		if c.Get(HIGH) > p.Get(HIGH) {
			c.Set(r, 3.0)
			c.SetComment("2u")
		}
		if c.Get(LOW) < p.Get(LOW) {
			c.Set(r, 2.0)
			c.SetComment("2d")
		}
		if c.Get(HIGH) <= p.Get(HIGH) && c.Get(LOW) >= p.Get(LOW) {
			c.Set(r, 1.0)
			c.SetComment("1")
		}
		if c.Get(HIGH) >= p.Get(HIGH) && c.Get(LOW) <= p.Get(LOW) {
			c.Set(r, 4.0)
			c.SetComment("3")
		}
	}
	return r
}

func StratPMG(prices *Matrix) int {
	r := prices.AddNamedColumn("PMG")
	cnt := 5
	for i := cnt + 1; i < prices.Rows; i++ {
		cn := 1
		for j := 0; j < cnt; j++ {
			idx := i - cnt + j
			if prices.DataRows[idx].Low() > prices.DataRows[idx-1].Low() {
				cn++
			} else {
				break
			}
		}
		if cn > 4 {
			c := &prices.DataRows[i]
			c.Set(r, float64(cn))
			c.SetComment(fmt.Sprintf("PMG UP %d", cn))
		}
	}
	for i := cnt + 1; i < prices.Rows; i++ {
		cn := 1
		for j := 0; j < cnt; j++ {
			idx := i - cnt + j
			if prices.DataRows[idx].High() < prices.DataRows[idx-1].High() {
				cn++
			} else {
				break
			}
		}
		if cn > 4 {
			c := &prices.DataRows[i]
			c.Set(r, float64(cn)*-1.0)
			c.SetComment(fmt.Sprintf("PMG DOWN %d", cn))
		}
	}
	return r
}

type Gap struct {
	Timestamp string
	Upper     float64
	Lower     float64
	Filled    int
	Index     int
}

type FairValueGap struct {
	Timestamp string
	Upper     float64
	Lower     float64
	Filled    int
	Index     int
	Type      int
}

func getGapState(lower, upper, low, high float64) int {
	if high > upper && low < lower {
		return 1
	}
	if high > upper && low > upper {
		return 0
	}
	if high < lower && low < lower {
		return 0
	}
	if high > upper && low > lower {
		return 2
	}
	if high < upper && low < lower {
		return 3
	}
	if high < upper && low > lower {
		return 4
	}
	return 0
}

func MergeGaps(gaps []FairValueGap) []FairValueGap {
	for j := 1; j < len(gaps); j++ {
		cg := &gaps[j]
		if cg.Filled == 0 {
			s := j - 10
			if s < 0 {
				s = 0
			}
			for i := s; i < j; i++ {
				ng := &gaps[i]
				if ng.Filled == 0 {
					t := 0
					if cg.Upper <= ng.Upper && cg.Upper >= ng.Lower && cg.Lower >= ng.Lower {
						t = 1
					}
					if cg.Upper >= ng.Upper && cg.Lower <= ng.Lower {
						t = 1
					}
					if cg.Upper >= ng.Upper && cg.Lower >= ng.Lower && cg.Lower <= ng.Upper {
						t = 1
					}
					if t == 1 {
						ng.Filled = 1
						cg.Lower = m.Min(cg.Lower, ng.Lower)
						cg.Upper = m.Max(cg.Upper, ng.Upper)
					}
				}
			}
		}
	}
	var ret = make([]FairValueGap, 0)
	for _, g := range gaps {
		if g.Filled == 0 {
			ret = append(ret, g)
		}
	}
	return ret
}

func FindFairValueGaps3(candles *Matrix) []FairValueGap {
	var gaps = make([]FairValueGap, 0)
	for i := 2; i < candles.Rows; i++ {
		c := candles.Row(i)
		p := candles.Row(i - 2)
		if c.Get(HIGH) < p.Get(LOW) {
			gaps = append(gaps, FairValueGap{
				Timestamp: p.Key,
				Filled:    0,
				Upper:     p.Get(LOW),
				Lower:     c.Get(HIGH),
				Index:     i - 2,
				Type:      -1,
			})
		}
		if c.Get(LOW) > p.Get(HIGH) {
			gaps = append(gaps, FairValueGap{
				Timestamp: p.Key,
				Filled:    0,
				Upper:     c.Get(LOW),
				Lower:     p.Get(HIGH),
				Index:     i - 2,
				Type:      1,
			})
		}
	}
	for j, g := range gaps {
		cg := &gaps[j]
		for i := g.Index + 3; i < candles.Rows; i++ {
			c := candles.Row(i)
			state := getGapState(cg.Lower, cg.Upper, c.Get(LOW), c.Get(HIGH))
			if state == 1 {
				cg.Filled = 1
				break
			} else if state == 2 {
				cg.Upper = c.Get(LOW)
			} else if state == 3 {
				cg.Lower = c.Get(HIGH)
			} else if state == 4 {
				fmt.Println("SPLIT IT:", c.Key)
			}
			if cg.Upper < cg.Lower {
				cg.Filled = 1
				break
			}
		}
	}
	return MergeGaps(gaps)
}

func FindFairValueGaps(candles *Matrix) []Gap {
	var gaps = make([]Gap, 0)
	for i := 2; i < candles.Rows; i++ {
		c := candles.Row(i)
		p := candles.Row(i - 2)
		if c.Get(HIGH) < p.Get(LOW) {
			gaps = append(gaps, Gap{
				Timestamp: c.Key,
				Filled:    0,
				Upper:     p.Get(LOW),
				Lower:     c.Get(HIGH),
				Index:     i - 2,
			})
		}
		if c.Get(LOW) > p.Get(HIGH) {
			gaps = append(gaps, Gap{
				Timestamp: c.Key,
				Filled:    0,
				Upper:     c.Get(LOW),
				Lower:     p.Get(HIGH),
				Index:     i - 2,
			})
		}
	}
	for j, g := range gaps {
		cg := &gaps[j]
		for i := g.Index + 3; i < candles.Rows; i++ {
			c := candles.Row(i)
			state := getGapState(cg.Lower, cg.Upper, c.Get(LOW), c.Get(HIGH))
			if state == 1 {
				cg.Filled = 1
				break
			} else if state == 2 {
				cg.Upper = c.Get(LOW)
			} else if state == 3 {
				cg.Lower = c.Get(HIGH)
			} else if state == 4 {
				fmt.Println("SPLIT IT:", c.Key)
			}
			if cg.Upper < cg.Lower {
				cg.Filled = 1
				break
			}
		}
	}
	for j := 1; j < len(gaps); j++ {
		cg := &gaps[j]
		if cg.Filled == 0 {
			s := j - 10
			if s < 0 {
				s = 0
			}
			for i := s; i < len(gaps); i++ {
				ng := &gaps[i]
				if ng.Filled == 0 {
					t := 0
					if cg.Upper < ng.Upper && cg.Upper > ng.Lower && cg.Lower < ng.Lower {
						t = 1
					}
					if cg.Upper > ng.Upper && cg.Lower > ng.Lower && cg.Lower < ng.Upper {
						t = 1
					}
					if t == 1 {
						ng.Filled = 1
						cg.Lower = m.Min(cg.Lower, ng.Lower)
						cg.Upper = m.Max(cg.Upper, ng.Upper)
					}
				}
			}
		}
	}
	return gaps
}

func FindFairValueGaps2(candles *Matrix) []Gap {
	var gaps = make([]Gap, 0)
	for i := 2; i < candles.Rows; i++ {
		c := candles.Row(i)
		p := candles.Row(i - 1)
		p2 := candles.Row(i - 2)
		if c.Get(HIGH) < p2.Get(LOW) && p.Get(ADJ_CLOSE) < p2.Get(LOW) {
			gaps = append(gaps, Gap{
				Timestamp: c.Key,
				Filled:    0,
				Upper:     p2.Get(LOW),
				Lower:     c.Get(HIGH),
				Index:     i - 2,
			})
		}
		if c.Get(LOW) > p2.Get(HIGH) && p.Get(ADJ_CLOSE) > p2.Get(HIGH) {
			gaps = append(gaps, Gap{
				Timestamp: c.Key,
				Filled:    0,
				Upper:     c.Get(LOW),
				Lower:     p2.Get(HIGH),
				Index:     i - 2,
			})
		}
	}
	for j, g := range gaps {
		cg := &gaps[j]
		for i := g.Index + 3; i < candles.Rows; i++ {
			c := candles.Row(i)
			state := getGapState(cg.Lower, cg.Upper, c.Get(LOW), c.Get(HIGH))
			if state == 1 {
				cg.Filled = 1
				break
			} else if state == 2 {
				cg.Upper = c.Get(LOW)
			} else if state == 3 {
				cg.Lower = c.Get(HIGH)
			} else if state == 4 {
				fmt.Println("SPLIT IT:", c.Key)
			}
			if cg.Upper < cg.Lower {
				cg.Filled = 1
				break
			}
		}
	}
	return gaps
}

func DMA(candles *Matrix, sma, ema int) int {
	// 0 = distance 1 = EMA of distance 2 = upper 3 = lower
	// total = 6
	si := SMA(candles, sma, ADJ_CLOSE)
	di := candles.Apply(func(mr MatrixRow) float64 {
		return mr.Get(ADJ_CLOSE) - mr.Get(si)
	})
	ei := EMA(candles, ema, di)
	st := candles.StdDev(di, 10)
	candles.Apply(func(mr MatrixRow) float64 {
		return mr.Get(ei) + mr.Get(st)
	})
	candles.Apply(func(mr MatrixRow) float64 {
		return mr.Get(ei) - mr.Get(st)
	})
	return di
}

func Overlap(candles *Matrix) int {
	ret := candles.AddNamedColumn("Overlap")
	for i := 1; i < candles.Rows; i++ {
		c := candles.DataRows[i]
		p := candles.DataRows[i-1]
		cl := c.Get(LOW)
		ch := c.Get(HIGH)
		pl := p.Get(LOW)
		ph := p.Get(HIGH)
		if ch > ph && cl < pl {
			candles.DataRows[i].Set(ret, 100.0)
		} else if ch < pl {
			candles.DataRows[i].Set(ret, 0.0)
		} else if cl > ph {
			candles.DataRows[i].Set(ret, 0.0)
		} else {
			ov := m.Max(0, m.Min(ph, ch)-m.Max(pl, cl))
			candles.DataRows[i].Set(ret, ov/(ph-pl)*100.0)
		}
	}
	return ret
}

func myMax(values ...float64) float64 {
	m := values[0]
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

func myMin(values ...float64) float64 {
	m := values[0]
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

// https://alpaca.markets/learn/andean-oscillator-a-new-technical-indicator-based-on-an-online-algorithm-for-trend-analysis/
func Andean(candles *Matrix, period, signal int) int {
	// 0 : Bull 1 : Bear 2 : Signal
	bu := candles.AddNamedColumn("Bull")
	be := candles.AddNamedColumn("Bear")
	sig := candles.AddNamedColumn("Signal")
	//length     = input(50)
	//sig_length = input(9,'Signal Length')
	alpha := 2.0 / (float64(period) + 1.0)

	up1 := candles.AddColumn()
	up2 := candles.AddColumn()
	dn1 := candles.AddColumn()
	dn2 := candles.AddColumn()
	c := candles.DataRows[0]
	candles.DataRows[0].Set(up1, candles.CLOSE(0))
	candles.DataRows[0].Set(up2, candles.CLOSE(0)*candles.CLOSE(0))
	candles.DataRows[0].Set(dn1, candles.CLOSE(0))
	candles.DataRows[0].Set(dn2, candles.CLOSE(0)*candles.CLOSE(0))
	for i := 1; i < candles.Rows; i++ {
		c = candles.DataRows[i]
		p := candles.DataRows[i-1]
		//C = close
		//O = open
		// up1 :=  	nz(math.max(C, O, up1[1] - (up1[1] - C) * alpha), C)
		candles.DataRows[i].Set(up1, myMax(candles.CLOSE(i), candles.OPEN(i), p.Get(up1)-(p.Get(up1)-candles.CLOSE(i))*alpha))
		//up2 := nz(math.max(C*C, O*O, up2[1]-(up2[1]-C*C)*alpha), C*C)
		candles.DataRows[i].Set(up2, myMax(candles.CLOSE(i)*candles.CLOSE(i), candles.OPEN(i)*candles.OPEN(i), p.Get(up2)-(p.Get(up2)-candles.CLOSE(i)*candles.CLOSE(i))*alpha))

		//dn1 := nz(math.min(C, O, dn1[1]+(C-dn1[1])*alpha), C)
		candles.DataRows[i].Set(dn1, myMin(candles.CLOSE(i), candles.OPEN(i), p.Get(dn1)+(candles.CLOSE(i)-p.Get(dn1))*alpha))
		//dn2 := nz(math.min(C*C, O*O, dn2[1]+(C*C-dn2[1])*alpha), C*C)
		candles.DataRows[i].Set(dn2, myMin(candles.CLOSE(i)*candles.CLOSE(i), candles.OPEN(i)*candles.OPEN(i), p.Get(dn2)+(candles.CLOSE(i)*candles.CLOSE(i)-p.Get(dn2))*alpha))
	}
	for i := 0; i < candles.Rows; i++ {
		c = candles.DataRows[i]
		//Components
		v1 := c.Get(dn2) - c.Get(dn1)*c.Get(dn1)
		if v1 < 0.0 {
			v1 = 0.0
		}
		v2 := c.Get(up2) - c.Get(up1)*c.Get(up1)
		if v2 < 0.0 {
			v2 = 0.0
		}
		candles.DataRows[i].Set(bu, math.Sqrt(v1))
		candles.DataRows[i].Set(be, math.Sqrt(v2))
	}
	tmp := candles.AddColumn()
	for i := 0; i < candles.Rows; i++ {
		c = candles.DataRows[i]
		if c.Get(bu) > c.Get(be) {
			candles.DataRows[i].Set(tmp, c.Get(bu))
		} else {
			candles.DataRows[i].Set(tmp, c.Get(be))
		}
	}
	ei := EMA(candles, tmp, signal)
	candles.CopyColumn(ei, sig)
	candles.RemoveColumns(6)
	//
	return bu
}

/*
study("Ehlers Distance Coefficient Filter", shorttitle="EDCF", overlay=true)

length = input(title="Length", type=integer, defval=14)
src = input(title="Source", type=source, defval=hl2)

srcSum = 0.0
coefSum = 0.0

for count = 0 to length - 1

	distance = 0.0

	for lookback = 1 to length - 1
		distance := distance + pow(src[count] - src[count + lookback], 2)

	srcSum := srcSum + distance * src[count]
	coefSum := coefSum + distance

dcf = coefSum != 0 ? srcSum / coefSum : 0.0

plot(dcf, title="EDCF", linewidth=2, color=#6d1e7f, transp=0)
*/
func EDCF(candles *Matrix, period int) int {
	ret := candles.AddNamedColumn("EDCF")
	src := candles.Apply(func(mr MatrixRow) float64 {
		return (mr.Get(1) + mr.Get(2)) / 2.0
	})
	for i := period; i < candles.Rows; i++ {

		srcSum := 0.0
		coefSum := 0.0

		start := i - period

		for count := 0; count < period-1; count++ {
			distance := 0.0
			for lookback := 1; lookback < period-1; lookback++ {
				distance = distance + math.Pow(candles.DataRows[start+count].Get(src)-candles.DataRows[start+count-lookback].Get(src), 2)
			}
			srcSum = srcSum + distance*candles.DataRows[start+count].Get(src)
			coefSum = coefSum + distance
		}
		dcf := 0.0
		if coefSum != 0.0 {
			dcf = srcSum / coefSum
		}
		candles.DataRows[i].Set(ret, dcf)
	}
	candles.RemoveColumns(1)
	return ret
}

// -----------------------------------------------------------------------
// ATR
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/a/atr.asp
func EfficiencyRatio(m *Matrix, days int) int {
	// 0 = ER
	ret := m.AddNamedColumn("ER")
	for i := days; i < m.Rows; i++ {
		direction := math.Abs(m.DataRows[i].Get(4) - m.DataRows[i-days].Get(4))
		volatility := 0.0
		for j := 0; j < days; j++ {
			if i-j-1 >= 0 {
				volatility += math.Abs(m.DataRows[i-j].Get(4) - m.DataRows[i-j-1].Get(4))
			}
		}
		//volatility *= float64(days)
		er := 0.0
		if volatility != 0.0 {
			er = direction / volatility
		}
		m.DataRows[i].Set(ret, er)
	}
	return ret
}

// -----------------------------------------------------------------------
// WaveTrend
// -----------------------------------------------------------------------
// https://www.tradingview.com/script/2KE8wTuF-Indicator-WaveTrend-Oscillator-WT/
// https://aryatrader.medium.com/what-is-wavetrend-indicator-a-multifunctional-tool-d8d1e5153843
// 10,21
func WaveTrend(m *Matrix, n1, n2 int) int {
	// 0 = WT 1 = WT-Signal 2 = Histo
	ret := m.AddNamedColumn("WT")
	sig := m.AddNamedColumn("WT-Sig")
	histo := m.AddNamedColumn("WT-Hist")
	// ap = hlc3
	hlc := m.Apply(func(mr MatrixRow) float64 {
		return (mr.Get(1) + mr.Get(2) + mr.Get(4)) / 3.0
	})
	// esa = ema(ap, n1)
	esa := EMA(m, n1, hlc)
	// d = ema(abs(ap - esa), n1)
	tmp := m.Apply(func(mr MatrixRow) float64 {
		return math.Abs(mr.Get(hlc) - mr.Get(esa))
	})
	d := EMA(m, n1, tmp)
	// ci = (ap - esa) / (0.015 * d)
	ci := m.Apply(func(mr MatrixRow) float64 {
		if mr.Get(d) != 0.0 {
			return (mr.Get(hlc) - mr.Get(esa)) / (0.015 * mr.Get(d))
		}
		return 0.0
	})
	// tci = ema(ci, n2)
	tci := EMA(m, n2, ci)
	m.CopyColumn(tci, ret)
	wt2 := SMA(m, 4, ret)
	m.CopyColumn(wt2, sig)
	m.ApplyRow(histo, func(mr MatrixRow) float64 {
		return mr.Get(ret) - mr.Get(sig)
	})
	m.RemoveColumns(6)
	return ret
}

// https://kaabar-sofien.medium.com/the-lbr-oscillator-for-trading-727da3530808
// 3,10,16
func LBR(m *Matrix, fast, slow, signal int) int {
	ret := m.AddNamedColumn("LBR-Line")
	sig := m.AddNamedColumn("LBR-Signal")
	hist := m.AddNamedColumn("LBR-Hist")
	//fast_length = input(title='Fast Length', defval=3)
	//slow_length = input(title='Slow Length', defval=10)
	//src = input(title='Source', defval=close)
	//stdev = input(title='St. dev.', defval=1.5)
	//stdev_length = input(title='St. dev. length', defval=100)
	//signal_length = input.int(title='Signal Smoothing', minval=1, maxval=50, defval=16)
	//sma_source = input(title='Simple MA(Oscillator)', defval=true)
	//sma_signal = input(title='Simple MA(Signal Line)', defval=true)
	//show_std_bands = input(title='Show std. bands', defval=false)
	//volatility_normalized = input(title='Volatility normalized', defval=false)

	// Calculating
	//fast_ma = sma_source ? ta.sma(src, fast_length) : ta.ema(src, fast_length)
	s1 := SMA(m, fast, 4)
	//slow_ma = sma_source ? ta.sma(src, slow_length) : ta.ema(src, slow_length)
	s2 := SMA(m, slow, 4)
	//macd = volatility_normalized ? ((fast_ma - slow_ma) / ta.atr(slow_length)) * 100 : (fast_ma - slow_ma)
	m.ApplyRow(ret, func(mr MatrixRow) float64 {
		return mr.Get(s1) - mr.Get(s2)
	})
	sg := SMA(m, signal, ret)
	m.CopyColumn(sg, sig)
	//signal = sma_signal ? ta.sma(macd, signal_length) : ta.ema(macd, signal_length)
	//hist = macd - signal
	m.ApplyRow(hist, func(mr MatrixRow) float64 {
		return mr.Get(ret) - mr.Get(sig)
	})
	m.RemoveColumns(3)
	return ret
}

func LBRNormalized(m *Matrix, slow, fast, signal int) int {
	ret := m.AddNamedColumn("LBR-Line")
	sig := m.AddNamedColumn("LBR-Signal")
	hist := m.AddNamedColumn("LBR-Hist")
	//fast_length = input(title='Fast Length', defval=3)
	//slow_length = input(title='Slow Length', defval=10)
	//src = input(title='Source', defval=close)
	//stdev = input(title='St. dev.', defval=1.5)
	//stdev_length = input(title='St. dev. length', defval=100)
	//signal_length = input.int(title='Signal Smoothing', minval=1, maxval=50, defval=16)
	//sma_source = input(title='Simple MA(Oscillator)', defval=true)
	//sma_signal = input(title='Simple MA(Signal Line)', defval=true)
	//show_std_bands = input(title='Show std. bands', defval=false)
	//volatility_normalized = input(title='Volatility normalized', defval=false)

	// Calculating
	//fast_ma = sma_source ? ta.sma(src, fast_length) : ta.ema(src, fast_length)
	s1 := SMA(m, fast, 4)
	//slow_ma = sma_source ? ta.sma(src, slow_length) : ta.ema(src, slow_length)
	s2 := SMA(m, slow, 4)
	ai := ATR(m, slow)
	//macd = volatility_normalized ? ((fast_ma - slow_ma) / ta.atr(slow_length)) * 100 : (fast_ma - slow_ma)
	m.ApplyRow(ret, func(mr MatrixRow) float64 {
		if mr.Get(ai) != 0.0 {
			return (mr.Get(s1) - mr.Get(s2)) / mr.Get(ai) * 100.0
		}
		return 0.0
	})
	sg := SMA(m, signal, ret)
	m.CopyColumn(sg, sig)
	//signal = sma_signal ? ta.sma(macd, signal_length) : ta.ema(macd, signal_length)
	//hist = macd - signal
	m.ApplyRow(hist, func(mr MatrixRow) float64 {
		return mr.Get(ret) - mr.Get(sig)
	})
	m.RemoveColumns(4)
	return ret
}

/*
study("Stochastics Momentum Index", shorttitle = "Stoch_MTM")
a = input(10, "Percent K Length")
b = input(3, "Percent D Length")
ob = input(40, "Overbought")
os = input(-40, "Oversold")
// Range Calculation
ll = lowest (low, a)
hh = highest (high, a)
diff = hh - ll
rdiff = close - (hh+ll)/2

avgrel = ema(ema(rdiff,b),b)
avgdiff = ema(ema(diff,b),b)
// SMI calculations
SMI = avgdiff != 0 ? (avgrel/(avgdiff/2)*100) : 0
SMIsignal = ema(SMI,b)
emasignal = ema(SMI, 10)
*/
func SMI(m *Matrix, k, d int) int {
	// 0 = SMI 1 = Signal 2 = Histogram
	//a = input(10, "Percent K Length")
	//b = input(3, "Percent D Length")
	//ob = input(40, "Overbought")
	//os = input(-40, "Oversold")
	ret := m.AddNamedColumn("SMI")
	sig := m.AddNamedColumn("SMI-Signal")
	histo := m.AddNamedColumn("SMI-Hist")
	diff := m.AddColumn()
	rdiff := m.AddColumn()
	for i := k; i < m.Rows; i++ {
		// Range Calculation
		hh := m.FindMaxBetween(1, i-k, k)
		ll := m.FindMinBetween(2, i-k, k)
		med := (hh + ll) / 2.0
		//ll = lowest (low, a)
		//hh = highest (high, a)
		//diff = hh - ll
		m.DataRows[i].Set(diff, hh-ll)
		//rdiff = close - (hh+ll)/2
		m.DataRows[i].Set(rdiff, m.DataRows[i].Get(4)-med)
	}
	e1 := EMA(m, d, rdiff)
	avgrel := EMA(m, d, e1)
	//avgrel = ema(ema(rdiff,b),b)
	//avgdiff = ema(ema(diff,b),b)
	e2 := EMA(m, d, diff)
	avgdiff := EMA(m, d, e2)
	// SMI calculations
	for i := k; i < m.Rows; i++ {
		//SMI = avgdiff != 0 ? (avgrel/(avgdiff/2)*100) : 0
		if m.DataRows[i].Get(avgdiff) != 0.0 {
			m.DataRows[i].Set(ret, (m.DataRows[i].Get(avgrel)/(m.DataRows[i].Get(avgdiff)/2.0))*100.0)
		}
	}
	e3 := EMA(m, d, ret)
	m.CopyColumn(e3, sig)
	m.ApplyRow(histo, func(mr MatrixRow) float64 {
		return mr.Get(ret) - mr.Get(sig)
	})
	//SMIsignal = ema(SMI,b)
	//emasignal = ema(SMI, 10)
	m.RemoveColumns(7)
	return ret
}

/*
study(title="Normalized smoothed MACD", shorttitle = "NSM", overlay=false)
//
inpFastPeriod   = input(defval=12, title="MACD fast period", minval=1, type=input.integer)
inpSlowPeriod   = input(defval=26, title="MACD slow period", minval=1, type=input.integer)
inpMacdSignal   = input(defval=9, title="Signal period", minval=1, type=input.integer)
inpSmoothPeriod = input(defval=5, title="Smoothing period", minval=1, type=input.integer)
inpNormPeriod   = input(defval=20, title="Normalization period", minval=1, type=input.integer)
price           = input(close, title="Price Source",type=input.source)
//
emaf = 0.0
emas = 0.0
val  = 0.0
nval = 0.0
sig  = 0.0
//
red  =color.new(#FF0000, 0)
green=color.new(#32CD32, 0)
black=color.new(#000000, 0)
//
if bar_index > inpSlowPeriod
    alphaf   = 2.0/(1.0+max(inpFastPeriod,1))
    alphas   = 2.0/(1.0+max(inpSlowPeriod,1))
    alphasig = 2.0/(1.0+max(inpMacdSignal,1))
    alphasm  = 2.0/(1.0+max(inpSmoothPeriod,1))

    emaf := emaf[1]+alphaf*(price-emaf[1])
    emas := emas[1]+alphas*(price-emas[1])
    imacd = emaf-emas

    mmax = highest(imacd,inpNormPeriod)
    mmin = lowest(imacd,inpNormPeriod)
    if mmin != mmax
        nval := 2.0*(imacd-mmin)/(mmax-mmin)-1.0
    else
        nval := 0

    val := val[1] + alphasm*(nval-val[1])
    sig := sig[1] + alphasig*(val-sig[1])
//
plot(val, color=val>val[1]?green:red, style=plot.style_line, linewidth=2, title="Reg smooth MACD")
plot(sig, color=black, style=plot.style_cross, linewidth=1, title="Signal line")
hline(0, title='0', color=color.gray, linestyle=hline.style_dotted, linewidth=1)
//
alertcondition(crossunder(val,sig),title="Sell",message="Sell")
alertcondition(crossover(val,sig),title="Buy",message="Buy")
alertcondition(crossunder(val,sig) or crossover(val,sig) ,title="Sell/Buy",message="Sell/Buy")
*/

/*

study(title="Zero Lag MACD Enhanced - Version 1.2", shorttitle="Zero Lag MACD Enhanced 1.2")
source = close

fastLength = input(12, title="Fast MM period", minval=1)
slowLength = input(26,title="Slow MM period",  minval=1)
signalLength =input(9,title="Signal MM period",  minval=1)
MacdEmaLength =input(9, title="MACD EMA period", minval=1)
useEma = input(true, title="Use EMA (otherwise SMA)")
useOldAlgo = input(false, title="Use Glaz algo (otherwise 'real' original zero lag)")
showDots = input(true, title="Show symbols to indicate crossing")
dotsDistance = input(1.5, title="Symbols distance factor", minval=0.1)

// Fast line
ma1= useEma ? ema(source, fastLength) : sma(source, fastLength)
ma2 = useEma ?  ema(ma1,fastLength) :  sma(ma1,fastLength)
zerolagEMA = ((2 * ma1) - ma2)

// Slow line
mas1=  useEma ? ema(source , slowLength) :  sma(source , slowLength)
mas2 =  useEma ? ema(mas1 , slowLength): sma(mas1 , slowLength)
zerolagslowMA = ((2 * mas1) - mas2)

// MACD line
ZeroLagMACD = zerolagEMA - zerolagslowMA

// Signal line
emasig1 = ema(ZeroLagMACD, signalLength)
emasig2 = ema(emasig1, signalLength)
signal = useOldAlgo ? sma(ZeroLagMACD, signalLength) : (2 * emasig1) - emasig2

hist = ZeroLagMACD - signal

upHist = (hist > 0) ? hist : 0
downHist = (hist <= 0) ? hist : 0


p1 = plot(upHist, color=green, transp=40, style=columns, title='Positive delta')
p2 = plot(downHist, color=purple, transp=40, style=columns, title='Negative delta')

zeroLine = plot(ZeroLagMACD, color=black, transp=0, linewidth=2, title='MACD line')
signalLine = plot(signal, color=gray, transp=0, linewidth=2, title='Signal')

ribbonDiff = hist > 0 ? green : purple
fill(zeroLine, signalLine, color=ribbonDiff)

circleYPosition = signal*dotsDistance
plot(ema(ZeroLagMACD,MacdEmaLength) , color=red, transp=0, linewidth=2, title='EMA on MACD line')

ribbonDiff2 = hist > 0 ? green : purple
plot(showDots and cross(ZeroLagMACD, signal) ? circleYPosition : na,style=circles, linewidth=4, color=ribbonDiff2, title='Dots')
*/
// -----------------------------------------------------------------------
// SSLChannel
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/a/atr.asp
func SSLChannel(m *Matrix, period int) int {
	// 0 = SSL Down 1 = SSL Up
	ret := m.AddNamedColumn("SSL Down")
	upi := m.AddNamedColumn("SSL Up")
	hlv := m.AddColumn()
	sh := SMA(m, period, 1)
	sl := SMA(m, period, 2)
	for i := 1; i < m.Rows; i++ {
		if m.DataRows[i].Get(4) > m.DataRows[i].Get(sh) {
			m.DataRows[i].Set(hlv, 1.0)
		} else if m.DataRows[i].Get(4) < m.DataRows[i].Get(sl) {
			m.DataRows[i].Set(hlv, -1.0)
		} else {
			m.DataRows[i].Set(hlv, m.DataRows[i].Get(hlv))
		}

	}
	for i := 1; i < m.Rows; i++ {
		if m.DataRows[i].Get(hlv) < 0.0 {
			m.DataRows[i].Set(ret, m.DataRows[i].Get(sh))
			m.DataRows[i].Set(upi, m.DataRows[i].Get(sl))
		} else {
			m.DataRows[i].Set(ret, m.DataRows[i].Get(sl))
			m.DataRows[i].Set(upi, m.DataRows[i].Get(sh))
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// Waddah Attar Explosion
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/a/atr.asp
func WAE(m *Matrix, sensitivity, fast, slow, length int, multiplier float64) int {
	// 0 = Trend Up 1 = Trend Down 2 = Explosion line
	ret := m.AddNamedColumn("TrendUp")
	dwi := m.AddNamedColumn("TrendDown")
	upi := m.AddNamedColumn("ExplosionLine")
	t1 := m.AddColumn()
	e1 := EMA(m, fast, 4)
	e2 := EMA(m, slow, 4)
	di := m.Apply(func(mr MatrixRow) float64 {
		return mr.Get(e1) - mr.Get(e2)
	})
	sens := float64(sensitivity)
	bb := BollingerBand(m, length, multiplier, multiplier)
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(t1, (m.DataRows[i].Get(di)-m.DataRows[i-1].Get(di))*sens)
	}
	m.ApplyRow(upi, func(mr MatrixRow) float64 {
		return mr.Get(bb) - mr.Get(bb+1)
	})
	for i := 0; i < m.Rows; i++ {
		t := m.DataRows[i].Get(t1)
		if t >= 0.0 {
			m.DataRows[i].Set(ret, t)
		}
		if t < 0.0 {
			m.DataRows[i].Set(dwi, t*-1.0)
		}
	}
	m.RemoveColumns(6)
	return ret
}

func SmoothedHeikinAshi(cn *Matrix, len1, len2 int) int {
	// 0 = Open 1 = High 2 = Low 3 = Close
	or := cn.AddNamedColumn("SHA-Open")
	hr := cn.AddNamedColumn("SHA-High")
	lr := cn.AddNamedColumn("SHA-Low")
	cr := cn.AddNamedColumn("SHA-Close")
	o := EMA(cn, len1, 0)
	c := EMA(cn, len1, 4)
	h := EMA(cn, len1, 1)
	l := EMA(cn, len1, 2)
	haclose := cn.AddColumn()
	haopen := cn.AddColumn()
	hahigh := cn.AddColumn()
	halow := cn.AddColumn()
	for i := 1; i < cn.Rows; i++ {
		cur := &cn.DataRows[i]
		p := cn.DataRows[i-1]
		// = (o+h+l+c)/4
		cur.Set(haclose, (cur.Get(o)+cur.Get(c)+cur.Get(h)+cur.Get(l))/4.0)
		//na(haopen[1]) ? (o + c)/2 : (haopen[1] + haclose[1]) / 2
		if i == 1 {
			cur.Set(haopen, (cur.Get(o)+cur.Get(c))/2.0)
		} else {
			cur.Set(haopen, (p.Get(haopen)+p.Get(haclose))/2.0)
		}
		// max (h, max(haopen,haclose))
		cur.Set(hahigh, m.Max(cur.Get(h), m.Max(cur.Get(haopen), cur.Get(haclose))))
		//min (l, min(haopen,haclose))
		cur.Set(halow, m.Min(cur.Get(l), m.Min(cur.Get(haopen), cur.Get(haclose))))
	}

	o2 := EMA(cn, len2, haopen)
	c2 := EMA(cn, len2, haclose)
	h2 := EMA(cn, len2, hahigh)
	l2 := EMA(cn, len2, halow)
	cn.CopyColumn(o2, or)
	cn.CopyColumn(c2, cr)
	cn.CopyColumn(h2, hr)
	cn.CopyColumn(l2, lr)
	//cn.RemoveColumns(12)
	return or

}

func rateCDV(uw, lw, body float64, cond bool) float64 {
	ip := 2.0 * body
	if !cond {
		ip = 0.0
	}
	ret := 0.5 * (uw + lw + ip) / (uw + lw + body)
	if ret == 0.0 {
		ret = 0.5
	}
	//_rate(cond) =>
	//  ret = 0.5 * (tw + bw + (cond ? 2 * body : 0)) / (tw + bw + body)
	//    ret := nz(ret) == 0 ? 0.5 : ret
	return ret
}

func CDV(cn *Matrix) int {
	//ret := cn.AddColumn()
	tmp := cn.AddColumn()
	delta := cn.AddNamedColumn("Delta")
	for i := 0; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		//tw = high - max(open, close)
		uw := c.Get(1) - m.Max(c.Get(0), c.Get(4))
		//bw = min(open, close) - low
		lw := m.Min(c.Get(0), c.Get(4)) - c.Get(2)
		//body = abs(close - open)
		body := m.Abs(c.Get(0) - c.Get(4))
		//deltaup =  volume * _rate(open <= close)
		du := c.Get(5) * rateCDV(uw, lw, body, c.Get(0) <= c.Get(4))
		//deltadown = volume * _rate(open > close)
		dd := c.Get(5) * rateCDV(uw, lw, body, c.Get(0) > c.Get(4))
		//delta = close >= open ? deltaup : -deltadown
		de := du
		if c.Get(0) < c.Get(4) {
			de = -1.0 * dd
		}
		c.Set(tmp, de)
		//cumdelta = cum(delta)
	}
	for i := 0; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		if i > 0 {
			c.Set(delta, cn.DataRows[i-1].Get(tmp)+c.Get(tmp))
		} else {
			c.Set(delta, c.Get(tmp))
		}
	}
	oi := cn.AddColumn()
	hi := cn.AddColumn()
	li := cn.AddColumn()
	ci := cn.AddColumn()
	aci := cn.AddColumn()

	for i := 1; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		p := cn.DataRows[i-1]
		c.Set(oi, p.Get(delta))
		c.Set(hi, m.Max(c.Get(delta), p.Get(delta)))
		c.Set(li, m.Min(c.Get(delta), p.Get(delta)))
		c.Set(ci, c.Get(delta))
		c.Set(aci, c.Get(delta))
	}
	return oi
}

// volume weighted candles
func VWC(cn *Matrix) int {
	si := SMA(cn, 10, 5)
	oi := cn.AddColumn()
	hi := cn.AddColumn()
	li := cn.AddColumn()
	ci := cn.AddColumn()
	aci := cn.AddColumn()
	ri := cn.AddColumn()
	ti := cn.AddColumn()
	for i := 1; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		v := 0.0
		if c.Get(si) != 0.0 {
			v = c.Get(5) / c.Get(si)
		}
		c.Set(ti, v)
		c.Set(hi, c.Get(1)+c.Get(1)*v)
		c.Set(li, c.Get(2)-c.Get(2)*v)
		if c.Get(4) < c.Get(0) {
			c.Set(oi, c.Get(0)+c.Get(0)*v)
			c.Set(ci, c.Get(3)-c.Get(3)*v)
			c.Set(aci, c.Get(4)-c.Get(4)*v)
		} else {
			c.Set(oi, c.Get(0)-c.Get(0)*v)
			c.Set(ci, c.Get(3)+c.Get(3)*v)
			c.Set(aci, c.Get(4)+c.Get(4)*v)
		}
		c.Set(ri, (c.Get(1)-c.Get(2))*v)
	}
	return oi
}

func CandleSentiments(cn *Matrix, atr, vma int) int {
	// 0 = Range/ATR 1 = Body 2 = B/R 3 = Body/ATR 4 = Resistance 5 = P-E 6 = V/SMA
	//si := SMA(cn, vma, 5)
	raIdx := cn.AddNamedColumn("R/A")
	bIdx := cn.AddNamedColumn("Body")
	brIdx := cn.AddNamedColumn("B/R")
	baIdx := cn.AddNamedColumn("B/A")
	rsIdx := cn.AddNamedColumn("Resistance")
	peIdx := cn.AddNamedColumn("P-E")
	vsaIdx := cn.AddNamedColumn("V/SMA")
	ai := ATR(cn, atr)
	ei := EMA(cn, vma, 4)
	vi := SMA(cn, vma, 5)
	for i := 0; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		uw := c.Get(1) - m.Max(c.Get(0), c.Get(4))
		lw := m.Min(c.Get(0), c.Get(4)) - c.Get(2)
		body := m.Abs(c.Get(0) - c.Get(4))
		rng := c.Get(1) - c.Get(2)
		if c.Get(ai) != 0.0 {
			c.Set(raIdx, rng/c.Get(ai))
			c.Set(bIdx, body)
			c.Set(brIdx, body/rng)
			c.Set(baIdx, body/c.Get(ai))
		}
		if c.Get(0) < c.Get(4) {
			c.Set(rsIdx, uw/rng)
		} else {
			c.Set(rsIdx, lw/rng)
		}
		c.Set(peIdx, c.Get(4)-c.Get(ei))
		if c.Get(vi) != 0.0 {
			c.Set(vsaIdx, c.Get(5)/c.Get(vi))
		}
	}
	cn.RemoveColumns(4)
	return raIdx
}

// PVT = [((CurrentClose - PreviousClose) / PreviousClose) x Volume] + PreviousPVT
func PVT(cn *Matrix, signal int) int {
	ret := cn.AddNamedColumn("PVT")
	for i := 1; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		p := cn.DataRows[i-1]
		d := ((c.Get(4)-p.Get(4))/p.Get(4))*c.Get(5) + p.Get(ret)
		c.Set(ret, d)
	}
	EMA(cn, signal, ret)
	return ret
}

// Rolling VWAP - not anchored
func RVWAP(candles *Matrix, period int) int {
	ret := candles.AddNamedColumn("RWAP")
	tp := HLC3(candles)
	//typical_price_volume = typical_price * volume
	tpv := candles.Apply(func(mr MatrixRow) float64 {
		return mr.Get(tp) * mr.Get(5)
	})
	//cumulative_price_volume = math.sum(typical_price_volume, cumulative_period)
	cpv := candles.Sum(tpv, period)
	//cumulative_volume = math.sum(volume, cumulative_period)
	cv := candles.Sum(5, period)
	//vwap = cumulative_price_volume / cumulative_volume
	candles.ApplyRow(ret, func(mr MatrixRow) float64 {
		if mr.Get(cv) != 0.0 {
			return mr.Get(cpv) / mr.Get(cv)
		}
		return 0.0
	})
	candles.RemoveColumns(4)
	return ret
}

// PVT = [((CurrentClose - PreviousClose) / PreviousClose) x Volume] + PreviousPVT
func PVR(cn *Matrix) int {
	// 1 = strong uptrend 0.5 = weak uptrend -0.5 = weak downtrend -1 = strong downtrend
	ret := cn.AddNamedColumn("PVR")
	for i := 1; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		p := cn.DataRows[i-1]
		d := 0.0
		if p.Get(4) > c.Get(4) && p.Get(5) < c.Get(5) {
			d = -1.0
		}
		if p.Get(4) > c.Get(4) && p.Get(5) > c.Get(5) {
			d = -0.5
		}
		if p.Get(4) < c.Get(4) && p.Get(5) < c.Get(5) {
			d = 1.0
		}
		if p.Get(4) < c.Get(4) && p.Get(5) > c.Get(5) {
			d = 0.5
		}
		c.Set(ret, d)
	}
	return ret
}

//-----------------------------------------------------------------------------}
//Augmented RSI
//-----------------------------------------------------------------------------{
/*
length = input.int(14, minval = 2)
smoType1 = input.string('RMA', 'Method', options = ['EMA', 'SMA', 'RMA', 'TMA'])
src = input(close, 'Source')

arsiCss = input(color.silver, 'Color', inline = 'rsicss')
autoCss = input(true, 'Auto', inline = 'rsicss')

//Signal Line
smooth = input.int(14, minval = 1, group = 'Signal Line')
*/
func UltimateRSI(cn *Matrix, length int) int {
	ret := cn.AddColumn()
	//upper = ta.highest(src, length)
	upper := Highest(cn, length, 4)
	//lower = ta.lowest(src, length)
	lower := Lowest(cn, length, 4)
	//r = upper - lower
	r := cn.Apply(func(mr MatrixRow) float64 {
		return mr.Get(upper) - mr.Get(lower)
	})
	//d = src - src[1]
	d := ROC(cn, 1, 4)

	//diff = upper > upper[1] ? r
	//: lower < lower[1] ? -r
	//: d
	diff := cn.AddColumn()
	for i := 1; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		p := cn.DataRows[i-1]
		if c.Get(upper) > p.Get(upper) {
			c.Set(diff, c.Get(r))
		} else if c.Get(lower) < p.Get(lower) {
			c.Set(diff, c.Get(r)*-1.0)
		} else {
			c.Set(diff, c.Get(d))
		}
	}
	//num = ma(diff, length, smoType1)
	num := EMA(cn, diff, length)
	ad := cn.Apply(func(mr MatrixRow) float64 {
		return m.Abs(mr.Get(diff))
	})
	//den = ma(math.abs(diff), length, smoType1)
	den := EMA(cn, ad, length)
	//arsi = num / den * 50 + 50
	cn.ApplyRow(ret, func(mr MatrixRow) float64 {
		if mr.Get(den) != 0.0 {
			return mr.Get(num)/mr.Get(den)*50.0 + 50.0
		}
		return 0.0
	})
	//signal = ma(arsi, smooth, smoType2)
	cn.RemoveColumns(8)
	return ret
}

func Range(cn *Matrix, ema int) int {
	// 0 = Range 1 = EMA 2 = Relation
	ret := cn.AddNamedColumn("Range")
	for i := 0; i < cn.Rows; i++ {
		c := &cn.DataRows[i]
		c.Set(ret, c.High()-c.Low())
	}
	e := EMA(cn, ema, ret)
	cn.Apply(func(mr MatrixRow) float64 {
		if mr.Get(e) != 0.0 {
			return mr.Get(ret) / mr.Get(e)
		}
		return 0.0
	})
	return ret
}

/*
func SmoothedPivotPoints(cn *Matrix, lookback int) int {
	fs := cn.AddNamedColumn("FirstSupport")
	fr := cn.AddNamedColumn("FirstResistance")

	// Adjusted highs
	adjusted_high = ta.highest(high, lookback_piv)
	// Adjusted lows
	adjusted_low = ta.lowest(low, lookback_piv)
	// Adjusted close
	adjusted_close = ta.sma(close, lookback_piv)
	// Pivot point
	pivot_point = (adjusted_high + adjusted_low + adjusted_close) / 3
	first_support = (ta.lowest(pivot_point, 12) * 2) - adjusted_high
	first_resistance = (ta.highest(pivot_point, 12) * 2) - adjusted_low
	return fs

}
*/
