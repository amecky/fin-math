package math

import "math"

// -----------------------------------------------------------------------
// SMA - Simple moving average
// -----------------------------------------------------------------------
func SMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
	n := float64(days)
	for i := days - 1; i < m.Rows; i++ {
		sum := 0.0
		current := m.DataRows[i]
		for j := 0; j < days; j++ {
			idx := i - days + 1 + j
			sum += m.DataRows[idx].Get(field)
		}
		avg := sum / n
		row := m.FindRow(current.Key)
		row.Set(ret, avg)
	}
	return ret
}

// -----------------------------------------------------------------------
// EMA
// -----------------------------------------------------------------------
func EMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
	if m.Rows >= days {
		n := float64(days)
		sma := SMA(m, days, field)
		multiplier := 2.0 / (n + 1)
		m.DataRows[days].Set(ret, m.DataRows[days-1].Get(sma))
		for i := days + 1; i < m.Rows; i++ {
			v := m.DataRows[i-1].Get(ret)*(1.0-multiplier) + m.DataRows[i].Get(field)*multiplier
			m.DataRows[i].Set(ret, v)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
// RMA
// -----------------------------------------------------------------------
func RMA(m *Matrix, days, field int) int {
	ret := m.AddColumn()
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
	ret := m.AddColumn()
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
	sma := SMA(prices, days, 4)
	for _, p := range prices.DataRows {
		p.Set(ret, (p.Get(4)-p.Get(sma))/p.Get(sma)*100.0)
	}
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
		m.DataRows[i].Set(mid, (m.DataRows[i].Get(1)+m.DataRows[i].Get(2))/2.0)
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
// MACD
// -----------------------------------------------------------------------
func MACD(m *Matrix, short, long, signal int) int {
	// 0 = Line 1 = Signal 2 = Diff
	ret := m.AddColumn()
	sig := m.AddColumn()
	diff := m.AddColumn()

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
//  Momentum
// -----------------------------------------------------------------------
func Momentum(m *Matrix, days int) int {
	// 0 = Momentum 1 = Momentum Percentage
	ret := m.AddColumn()
	per := m.AddColumn()
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(4) - m.DataRows[i-days].Get(4)))
		m.DataRows[i].Set(per, (m.DataRows[i].Get(4)-m.DataRows[i-days].Get(4))/m.DataRows[i-days].Get(4)*100.0)
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
			d := (cp.Get(4) - cs) / (max - min)
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
		price := cp.Get(4)
		for j := i + 1; j < m.Rows; j++ {
			if counting == true {
				p := m.DataRows[j]
				if p.Get(4) > pp.Get(2) && p.Get(4) < cp.Get(1) {
					cnt++
					price += p.Get(4)
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
			m.DataRows[i].Set(ret, (m.DataRows[i].Get(4)-m.DataRows[i].Get(sub))/m.DataRows[i].Get(sub)*100.0)
		}
	}
	return ret
}

// -----------------------------------------------------------------------
//  RSI
// -----------------------------------------------------------------------
func RSI(m *Matrix, days, field int) int {
	// 0 = RSI
	ret := m.AddColumn()
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
		rs := m.DataRows[i].Get(ue) / m.DataRows[i].Get(de)
		m.DataRows[i].Set(ret, (100.0 - 100.0/(1.0+rs)))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
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
	ret := m.AddColumn()
	smoothed := m.AddColumn()
	trIdx := m.AddColumn()
	if m.Rows < 1 {
		return -1
	}
	for i := 1; i < m.Rows; i++ {
		curH := m.DataRows[i].Get(1)
		curL := m.DataRows[i].Get(2)
		prevC := m.DataRows[i-1].Get(4)
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
	m.RemoveColumn()
	m.RemoveColumn()
	m.RemoveColumn()
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
	// 0 = ADR
	ret := m.AddColumn()
	sh := SMA(m, days, 1)
	sl := SMA(m, days, 2)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, (m.DataRows[i].Get(sh) - m.DataRows[i].Get(sl)))
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
// ROC
// -----------------------------------------------------------------------
func ROC(m *Matrix, days, field int) int {
	// 0 = ROC
	ret := m.AddColumn()
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
// Stochastic
// -----------------------------------------------------------------------
func Stochastic(m *Matrix, days, ema int) int {
	// 0 = K 1 = D
	k := m.AddColumn()
	d := m.AddColumn()
	v := m.AddColumn()
	total := m.Rows
	if total < days {
		return -1
	}
	for i := days; i < m.Rows; i++ {
		low := m.FindMinBetween(2, i-days+1, days)
		high := m.FindMaxBetween(1, i-days+1, days)
		m.DataRows[i].Set(v, (m.DataRows[i].Get(4)-low)/(high-low)*100.0)
	}
	slowData := SMA(m, ema, v)
	dData := SMA(m, 3, slowData)
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(k, m.DataRows[i].Get(slowData))
		m.DataRows[i].Set(d, m.DataRows[i].Get(dData))
	}
	m.RemoveColumn()
	m.RemoveColumn()
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
	return k
}

func StochasticSMA(m *Matrix, sma, days, ema int) int {
	smaIdx := SMA(m, sma, 4)
	return StochasticExt(m, days, ema, smaIdx, smaIdx, smaIdx)
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

func standardAbbreviation(m *Matrix, v float64, offset int, cnt int) float64 {
	sum := 0.0
	end := offset + cnt
	if end > m.Rows {
		end = m.Rows
	}
	for i := offset; i < end; i++ {
		d := m.DataRows[i].Get(4) - v
		sum += d * d
	}
	return math.Sqrt(sum / float64(cnt))
}

// -----------------------------------------------------------------------
// BollingerBand
// -----------------------------------------------------------------------
func BollingerBand(m *Matrix, ema int, upper, lower float64) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	midIdx := m.AddColumn()
	sma := SMA(m, ema, 4)
	for i := 0; i < m.Rows; i++ {
		sa := standardAbbreviation(m, m.DataRows[i].Get(sma), i, ema)
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(sma)+sa*upper)
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(sma)-sa*lower)
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(sma))
	}
	m.RemoveColumn()
	return upIdx
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
		m.DataRows[i].Set(midIdx, m.DataRows[i].Get(4))
	}
	return ret
}

// -----------------------------------------------------------------------
// KeltnerChannel
// -----------------------------------------------------------------------
// https://www.investopedia.com/terms/k/keltnerchannel.asp
func Keltner(m *Matrix, ema, atr int) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	ed := EMA(m, ema, 4)
	ad := ATR(m, atr)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(upIdx, m.DataRows[i].Get(ed)+2.0*m.DataRows[i].Get(ad))
		m.DataRows[i].Set(lowIdx, m.DataRows[i].Get(ed)-2.0*m.DataRows[i].Get(ad))
	}
	m.RemoveColumn()
	return upIdx
}

// -----------------------------------------------------------------------
//  DonchianChannel
// -----------------------------------------------------------------------
func DonchianChannel(m *Matrix, days int) int {
	// 0 = Upper 1 = Lower 2 = Mid
	upIdx := m.AddColumn()
	lowIdx := m.AddColumn()
	midIdx := m.AddColumn()
	for i := days; i < m.Rows; i++ {
		m.DataRows[i].Set(upIdx, m.FindMaxBetween(1, i-days, days))
		m.DataRows[i].Set(lowIdx, m.FindMinBetween(2, i-days, days))
		m.DataRows[i].Set(midIdx, (m.DataRows[i].Get(upIdx) + m.DataRows[i].Get(lowIdx)/2.0))
	}
	return upIdx
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
			value = (high - m.DataRows[i].Get(4)) / dl * -100.0
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
		m.DataRows[i].Set(d, m.DataRows[i].Get(4)-m.DataRows[i].Get(sd))
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
		m.DataRows[i].Set(ret, m.DataRows[i].Get(4)-m.DataRows[i].Get(ed))
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
		curH := m.DataRows[i].Get(1)
		curL := m.DataRows[i].Get(2)
		prevC := m.DataRows[i-1].Get(4)
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
		top := math.Max(p.Get(0), p.Get(4))
		bottom := math.Min(p.Get(0), p.Get(4))

		d := p.Get(0) - p.Get(4)
		if d < 0.0 {
			d *= -1.0
		}
		td := p.Get(1) - p.Get(2)
		uwAbs := p.Get(1) - top
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
		if p.Get(0) > p.Get(4) {
			m.DataRows[i].Set(trendIdx, -1.0)
		} else {
			m.DataRows[i].Set(trendIdx, 1.0)
		}
		m.DataRows[i].Set(bodyPos, (1.0-((p.Get(1)-top)/td))*100.0)
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
//  VO
// -----------------------------------------------------------------------
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/volume-oscillator
// https://commodity.com/technical-analysis/volume-oscillator/
func VO(m *Matrix, fast, slow int) int {
	ret := m.AddColumn()
	emaSlow := EMA(m, slow, 5)
	emaFast := EMA(m, fast, 5)
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
		m.DataRows[i].Set(lsIdx, m.DataRows[i-mid].Get(4))
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
			if m.DataRows[i-j].Get(4) > m.DataRows[i-j].Get(0) {
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
		m.DataRows[i].Set(buIdx, (cp.Get(1)+cp.Get(2))/2.0+multiplier*cp.Get(atrIdx))
		m.DataRows[i].Set(blIdx, (cp.Get(1)+cp.Get(2))/2.0-multiplier*cp.Get(atrIdx))
	}

	//lowerBand := lowerBand > prevLowerBand or close[1] < prevLowerBand ? lowerBand : prevLowerBand
	//upperBand := upperBand < prevUpperBand or close[1] > prevUpperBand ? upperBand : prevUpperBand
	for i := 1; i < m.Rows; i++ {
		cp := m.DataRows[i]
		pp := m.DataRows[i-1]

		if cp.Get(buIdx) < pp.Get(buIdx) || pp.Get(4) > pp.Get(buIdx) {
			m.DataRows[i].Set(bfuIdx, cp.Get(buIdx))
		} else {
			m.DataRows[i].Set(bfuIdx, pp.Get(bfuIdx))
		}

		if cp.Get(blIdx) > pp.Get(blIdx) || pp.Get(4) < pp.Get(blIdx) {
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
			if cp.Get(4) > cp.Get(bfuIdx) {
				dir = 1.0
			} else {
				dir = -1.0
			}
		} else {
			if cp.Get(4) < cp.Get(bflIdx) {
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

func GAP(m *Matrix) int {
	ret := m.AddColumn()
	atrIdx := ATR(m, 14)
	for i := 15; i < m.Rows; i++ {
		cp := m.DataRows[i]
		prev := m.DataRows[i-1]
		value := (cp.Get(0) - prev.Get(4)) / cp.Get(atrIdx) * 100.0
		m.DataRows[i].Set(ret, value)
	}
	return ret
}
