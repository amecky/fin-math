package math

import (
	"fmt"
	"math"
	m "math"
)

// -----------------------------------------------------------------------
// Normaizes Value into -1 / 1 range
// -----------------------------------------------------------------------
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
//
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
	mid := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(mid, (m.DataRows[i].Get(HIGH)+m.DataRows[i].Get(2))/2.0)
	}
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
//  Momentum
// -----------------------------------------------------------------------
func Momentum(m *Matrix, days int) int {
	// 0 = Momentum 1 = Momentum Percentage
	ret := m.AddNamedColumn("Momentum")
	per := m.AddNamedColumn("Momentum-Pct")
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(ADJ_CLOSE) - m.DataRows[i-days].Get(ADJ_CLOSE)))
		m.DataRows[i].Set(per, (m.DataRows[i].Get(ADJ_CLOSE)-m.DataRows[i-days].Get(ADJ_CLOSE))/m.DataRows[i-days].Get(ADJ_CLOSE)*100.0)
	}
	return ret
}

// -----------------------------------------------------------------------
//  Daily Percentage Change
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
func MeanBreakout(m *Matrix, period int) int {
	// 0 = MBO
	ret := m.AddColumn()
	sma := SMA(m, period, 4)
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
//  RSI
// -----------------------------------------------------------------------
func RSI(m *Matrix, days, field int) int {
	// 0 = RSI
	ret := m.AddNamedColumn("RSI")
	sub := m.AddColumn()
	up := m.AddColumn()
	down := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		m.DataRows[i].Set(sub, m.DataRows[i].Get(field)-m.DataRows[i-1].Get(field))
	}

	for i := 1; i < m.Rows; i++ {
		if m.DataRows[i].Get(sub) > 0.0 {
			m.DataRows[i].Set(up, m.DataRows[i].Get(sub))
		} else {
			m.DataRows[i].Set(down, -1.0*m.DataRows[i].Get(sub))
		}
	}
	ue := RMA(m, days, up)
	de := RMA(m, days, down)
	for i := 1; i < m.Rows; i++ {
		if m.DataRows[i].Get(de) != 0.0 {
			rs := m.DataRows[i].Get(ue) / m.DataRows[i].Get(de)
			m.DataRows[i].Set(ret, (100.0 - 100.0/(1.0+rs)))
		}
	}
	m.RemoveColumns(5)
	return ret
}

// -----------------------------------------------------------------------
//  RSI
// -----------------------------------------------------------------------
func RSI_BB(m *Matrix, days, field int) int {
	// 0 = RSI 1 = Upper 2 = Lower 3 = Mid
	ri := RSI(m, days, 4)
	BollingerBandExt(m, ri, days, 2.0, 2.0)
	return ri
}

// -----------------------------------------------------------------------
//  RSI Momentum
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
//
func ADR(m *Matrix, days int) int {
	// 0 = ADR 1 = Real ADR
	ret := m.AddColumn()
	rel := m.AddColumn()
	relTmp := m.AddColumn()
	sh := SMA(m, days, 1)
	sl := SMA(m, days, 2)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(sh) - m.DataRows[i].Get(sl)))
	}
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(relTmp, m.DataRows[i].Get(HIGH)/m.DataRows[i].Get(2))
	}
	sr := SMA(m, days, relTmp)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(rel, 100.0*(m.DataRows[i].Get(sr)-1.0))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
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
		m.DataRows[i].Set(ret, 100.0*(m.DataRows[i].Get(si)-1.0))
	}
	return ret
}

// -----------------------------------------------------------------------
// ROC
// -----------------------------------------------------------------------
func ROC(m *Matrix, days, field int) int {
	// 0 = ROC
	ret := m.AddNamedColumn(fmt.Sprintf("ROC%d", days))
	start := days
	for i := days; i < m.Rows; i++ {
		if i >= start {
			current := m.DataRows[i].Get(field)
			prev := m.DataRows[i-days].Get(field)
			v := 0.0
			if prev != 0.0 {
				v = (current - prev) / prev * 100.0
			}
			m.DataRows[i].Set(ret, v)
		}
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
func StochasticRSI(m *Matrix, days, smoothK, smoothD int) int {
	// 0 = K 1 = D
	ki := m.AddColumn()
	di := m.AddColumn()
	rsi := RSI(m, days, 4)
	sr := StochasticExt(m, 14, 3, rsi, rsi, rsi)
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
//  DonchianChannel
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
//  RSI - ATR - RSI
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
//  WilliamsRange
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
	ret := m.AddColumn()
	trIdx := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		curH := m.DataRows[i].Get(HIGH)
		curL := m.DataRows[i].Get(2)
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
	rma := RMA(m, days, trIdx)
	stoch := StochasticExt(m, days, rma, rma, rma, rma)
	m.CopyColumn(stoch, ret)
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
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
	// 0 = BodySize 1 = Upper Wick 2 = Lower Wick 3 = BodyPos
	relBS := m.AddNamedColumn("RelBodySize")
	upIdx := m.AddNamedColumn("Upper")
	lowIdx := m.AddNamedColumn("Lower")
	trendIdx := m.AddNamedColumn("Trend")
	for i := 0; i < m.Rows; i++ {
		p := m.DataRows[i]
		top := math.Max(p.Get(OPEN), p.Get(ADJ_CLOSE))
		bottom := math.Min(p.Get(OPEN), p.Get(ADJ_CLOSE))

		d := p.Get(OPEN) - p.Get(ADJ_CLOSE)
		if d < 0.0 {
			d *= -1.0
		}
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
//  Relative Volume
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
//  VO
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
//  AverageVolume
// -----------------------------------------------------------------------
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/volume-oscillator
// https://commodity.com/technical-analysis/volume-oscillator/
func AverageVolume(m *Matrix, lookback int) int {
	// 0 = SMA of volume
	sma := SMA(m, lookback, 5)
	return sma
}

// -----------------------------------------------------------------------
//  Ichimoku
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/i/ichimoku-cloud.asp
func Ichimoku(m *Matrix, short, mid, long int) int {
	// 0 = Conversion line 1 = Base Line 2 = Leading Span A 3 = Leading Span B 4 = Lagging Span
	ret := m.AddColumn()
	midIdx := m.AddColumn()
	lsaIdx := m.AddColumn()
	lsbIdx := m.AddColumn()
	lsIdx := m.AddColumn()
	for i := short; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-short, short)
		high := m.FindMaxBetween(1, i-short, short)
		m.DataRows[i].Set(ret, (high+low)/2.0)
	}
	for i := mid; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-mid, mid)
		high := m.FindMaxBetween(1, i-mid, mid)
		m.DataRows[i].Set(midIdx, (high+low)/2.0)
	}
	for i := mid; i < m.Rows; i++ {
		m.DataRows[i].Set(lsaIdx, (m.DataRows[i].Get(ret)+m.DataRows[i].Get(midIdx))/2.0)
	}
	for i := long; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-long, long)
		high := m.FindMaxBetween(1, i-long, long)
		m.DataRows[i].Set(lsbIdx, (high+low)/2.0)
	}
	for i := mid; i < m.Rows; i++ {
		m.DataRows[i].Set(lsIdx, m.DataRows[i-mid].Get(ADJ_CLOSE))
	}
	return ret
}

// -----------------------------------------------------------------------
//  Weighted Trend Intensity
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
		value := ChangePercentage(sl, cp.Get(atrIdx))
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
		value := ChangePercentage(cp.Get(HIGH)-cp.Get(LOW), cp.Get(atrIdx))
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
//  STD
// -----------------------------------------------------------------------
func STD(prices *Matrix, days int) int {
	// 0 = STD
	s := prices.StdDev(4, days)
	return s
}

// -----------------------------------------------------------------------
//  STDChannel
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
//  STD
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
	// 0 = Bullish Bearish Count
	ret := m.AddColumn()
	for i := period; i < m.Rows; i++ {
		bu := 0
		for j := 0; j < period; j++ {
			idx := i - j
			if m.DataRows[idx].Get(0) <= m.DataRows[idx].Get(ADJ_CLOSE) {
				bu++
			}
		}
		d := float64(bu) / float64(period) * 100.0
		m.DataRows[i].Set(ret, d)
	}
	return ret
}

// -----------------------------------------------------------------------
//  OBV
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
//  Aroon
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
//  TrendIntensity
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
//  Chaikin A/D Line
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
//  TSI
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
//  Divergence
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
//  Divergence
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
//  ADX
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

// -----------------------------------------------------------------------
//  MFI
// -----------------------------------------------------------------------
/*
Typical Price: (High + Low + Close) / 3

If TP(t) > TP(t-1) then MFP = MFP(t-1) + TP(t) * V(t)
If TP(t) < TP(t-1) then MFN = MFN(t-1) + TP(t) * V(t)

MFI(t) = 100 - 100 / ( (1 + MFP(t)) / MFN(t))

*/
func MFI(m *Matrix, days int) int {
	// 0 = MFI
	ret := m.AddColumn()
	tp := m.AddColumn()
	rmf := m.AddColumn()
	mfpi := m.AddColumn()
	mfni := m.AddColumn()
	mr := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		v := (m.DataRows[i].Get(HIGH) + m.DataRows[i].Get(2) + m.DataRows[i].Get(ADJ_CLOSE)) / 3.0
		m.DataRows[i].Set(tp, v)
		m.DataRows[i].Set(rmf, m.DataRows[i].Get(5)*v)
	}
	// money ratio
	for i := days + 1; i < m.Rows; i++ {
		mfp := 0.0
		mfn := 0.0
		for j := 0; j < days; j++ {
			cur := m.DataRows[i-j].Get(tp)
			prev := m.DataRows[i-1-j].Get(tp)
			if cur > prev {
				mfp += m.DataRows[i].Get(rmf)
			} else {
				mfn += m.DataRows[i].Get(rmf)
			}
		}

		m.DataRows[i].Set(mfpi, mfp)
		m.DataRows[i].Set(mfni, mfn)
		if mfn != 0.0 {
			m.DataRows[i].Set(mr, mfp/mfn)
			m.DataRows[i].Set(ret, 100.0-100.0/(1.0+m.DataRows[i].Get(mr)))
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
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
//  HeikinAshi
// -----------------------------------------------------------------------
func HeikinAshi(m *Matrix) int {
	// 0 = Open 1 = High 2 = Low 3 = Close 4 = AdjClose 5 = Volume
	oi := m.AddColumn()
	hi := m.AddColumn()
	li := m.AddColumn()
	ci := m.AddColumn()
	aci := m.AddColumn()
	vi := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		c := m.DataRows[i]
		prev := m.DataRows[i-1]
		m.DataRows[i].Set(ci, 0.25*(c.Get(0)+c.Get(HIGH)+c.Get(2)+c.Get(ADJ_CLOSE)))
		m.DataRows[i].Set(oi, 0.5*(prev.Get(0)+prev.Get(ADJ_CLOSE)))
		m.DataRows[i].Set(hi, FindMax([]float64{c.Get(HIGH), c.Get(0), c.Get(ADJ_CLOSE)}))
		m.DataRows[i].Set(li, FindMin([]float64{c.Get(2), c.Get(0), c.Get(ADJ_CLOSE)}))
		m.DataRows[i].Set(aci, c.Get(ADJ_CLOSE))
		m.DataRows[i].Set(vi, float64(c.Get(5)))
	}
	return oi
}

// -----------------------------------------------------------------------
//  Hull Moving Average
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
// CMF
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
func STC(prices *Matrix, short, long, stoch int) int {
	// 0 = STC
	ret := prices.AddColumn()
	// EMA1 = EMA (Close, Short Length);
	// EMA2 = EMA (Close, Long Length);
	es := EMA(prices, short, 4)
	el := EMA(prices, long, 4)
	// MACD = EMA1 – EMA2.
	macd := prices.AddColumn()
	for i := 0; i < prices.Rows; i++ {
		prices.DataRows[i].Set(macd, prices.DataRows[i].Get(es)-prices.DataRows[i].Get(el))
	}
	// Second, the 10-period Stochastic from the MACD values is calculated:
	// %K (MACD) = %KV (MACD, 10);
	// %D (MACD) = %DV (MACD, 10);
	tmp := StochasticExt(prices, stoch, 3, macd, macd, macd)
	for i := 0; i < prices.Rows; i++ {
		// Schaff = 100 x (MACD – %K (MACD)) / (%D (MACD) – %K (MACD)).
		d := prices.DataRows[i].Get(tmp+1) - prices.DataRows[i].Get(tmp)
		if d != 0.0 {
			prices.DataRows[i].Set(ret, 100.0*(prices.DataRows[i].Get(macd)-prices.DataRows[i].Get(tmp))/d)
		}
	}
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
	prices.RemoveColumn()
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
//  Choppiness
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
//  Guppy GMMA
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
//  Volatility
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
//  z-Score
// -----------------------------------------------------------------------
// https://kaabar-sofien.medium.com/using-z-score-in-trading-a-python-study-5f4b21b41aa0
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
//  z-Normalization Bollinger Bands
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
		} else {
			sum -= 1.0
		}
		if p.Get(e2) < c.Get(e2) {
			sum += 1.0
		} else {
			sum -= 1.0
		}
		if p.Get(e3) < c.Get(e3) {
			sum += 1.0
		} else {
			sum -= 1.0
		}

		if c.Get(e1) > c.Get(e2) {
			sum += 1.0
		} else {
			sum -= 1.0
		}
		if c.Get(e2) > c.Get(e3) {
			sum += 1.0
		} else {
			sum -= 1.0
		}
		m.DataRows[i].Set(ret, sum/5.0)
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
func SqueezeMomentum(prices *Matrix, lookback int) int {
	// 0 = State 1 = Line
	ret := prices.AddColumn()
	li := prices.AddColumn()
	bb := BollingerBand(prices, lookback, 2.0, 2.0)
	kc := Keltner(prices, lookback, lookback, 1.5)
	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		if c.Get(bb+1) > c.Get(kc+1) && c.Get(bb) < c.Get(kc) {
			prices.DataRows[i].Set(ret, 1.0)
		}
		if c.Get(bb+1) < c.Get(kc+1) && c.Get(bb) > c.Get(kc) {
			prices.DataRows[i].Set(ret, -1.0)
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
	for i := lookback; i < prices.Rows; i++ {
		h, l := prices.FindHighLowIndex(i-lookback, lookback)
		d := (prices.DataRows[h].Get(ADJ_CLOSE) + prices.DataRows[l].Get(ADJ_CLOSE)) / 2.0
		prices.DataRows[i].Set(di, d)

	}

	si := SMA(prices, lookback, 4)
	for i := lookback; i < prices.Rows; i++ {
		d := prices.DataRows[i].Get(ADJ_CLOSE) - (prices.DataRows[i].Get(di)+prices.DataRows[i].Get(si))/2.0
		prices.DataRows[i].Set(li, d)

	}

	/*
		prices.RemoveColumn()
		prices.RemoveColumn()
		prices.RemoveColumn()
		prices.RemoveColumn()
		prices.RemoveColumn()
		prices.RemoveColumn()
		prices.RemoveColumn()
		prices.RemoveColumn()
	*/
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

func FisherTransform(prices *Matrix, period, signal int) int {
	// 0 = FT
	ret := prices.AddNamedColumn("FT")
	tmp := prices.AddColumn()
	for i := period; i < prices.Rows; i++ {
		l, h := prices.FindMinMaxBetween(ADJ_CLOSE, i-period, period)
		prices.DataRows[i].Set(tmp, NormalizeRange(prices.DataRows[i].Get(ADJ_CLOSE), l, h))
	}
	for i := 0; i < prices.Rows; i++ {
		c := prices.DataRows[i]
		if c.Get(tmp) != 1.0 {
			d := (1 + c.Get(tmp)) / (1.0 - c.Get(tmp))
			if d != 0.0 {
				c.Set(ret, 0.5*m.Log(d))
			}
		}
	}
	prices.RemoveColumns(1)
	SMA(prices, signal, ret)
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
	// 0 = TWAP
	ret := prices.AddColumn()
	sum := prices.Apply(func(mr MatrixRow) float64 {
		return (mr.Get(OPEN) + mr.Get(HIGH) + mr.Get(LOW) + mr.Get(ADJ_CLOSE)) / 4.0
	})
	si := SMA(prices, period, sum)
	prices.CopyColumn(si, ret)
	prices.RemoveColumns(2)
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
	for i := cnt; i < prices.Rows; i++ {
		cn := 1
		cur := prices.DataRows[i].Get(HIGH)
		for j := i - 1; j >= 0; j-- {
			if prices.DataRows[j].Get(HIGH) < cur {
				cn++
				cur = prices.DataRows[j].Get(HIGH)
			} else {
				break
			}
		}
		if cn > 4 {
			c := &prices.DataRows[i]
			c.Set(r, 1.0)
			c.SetComment(fmt.Sprintf("PMG UP %d", cn))
		}
	}
	for i := cnt; i < prices.Rows; i++ {
		cn := 1
		cur := prices.DataRows[i].Get(LOW)
		for j := i - 1; j >= 0; j-- {
			if prices.DataRows[j].Get(LOW) > cur {
				cn++
				cur = prices.DataRows[j].Get(LOW)
			} else {
				break
			}
		}
		if cn > 4 {
			c := &prices.DataRows[i]
			c.Set(r, -1.0)
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
				Timestamp: c.Key,
				Filled:    0,
				Upper:     p.Get(LOW),
				Lower:     c.Get(HIGH),
				Index:     i - 2,
			})
		}
		if c.Get(LOW) > p.Get(HIGH) {
			gaps = append(gaps, FairValueGap{
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
