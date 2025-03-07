package math

func MarketRegime(candles *Matrix, period int) int {
	ret := candles.AddNamedColumn("MR")
	e1 := EMA(candles, period, ADJ_CLOSE)
	for i := 13; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		p1 := candles.DataRows[i-13]
		p2 := candles.DataRows[i-8]
		p3 := candles.DataRows[i-5]
		regime := 0.0
		if c.Get(4) > c.Get(e1) && c.Get(4) > p1.Get(4) && c.Get(1) > p2.Get(1) && c.Get(2) > p3.Get(2) {
			regime = 1
		}
		if c.Get(4) > c.Get(e1) && c.Get(4) > p1.Get(4) && c.Get(1) > p2.Get(1) && c.Get(2) < p3.Get(2) {
			regime = 0.75
		}
		if c.Get(4) > c.Get(e1) && c.Get(4) > p1.Get(4) && c.Get(1) < p2.Get(1) && c.Get(2) < p3.Get(2) {
			regime = 0.5
		}
		if c.Get(4) > c.Get(e1) && c.Get(4) < p1.Get(4) && c.Get(1) < p2.Get(1) && c.Get(2) < p3.Get(2) {
			regime = 0.25
		}
		if c.Get(4) < c.Get(e1) && c.Get(4) < p1.Get(4) && c.Get(1) < p2.Get(1) && c.Get(2) < p3.Get(2) {
			regime = -1
		}
		if c.Get(4) < c.Get(e1) && c.Get(4) < p1.Get(4) && c.Get(1) < p2.Get(1) && c.Get(2) > p3.Get(2) {
			regime = -0.75
		}
		if c.Get(4) < c.Get(e1) && c.Get(4) < p1.Get(4) && c.Get(1) > p2.Get(1) && c.Get(2) > p3.Get(2) {
			regime = -0.5
		}
		if c.Get(4) < c.Get(e1) && c.Get(4) > p1.Get(4) && c.Get(1) > p2.Get(1) && c.Get(2) > p3.Get(2) {
			regime = -0.25
		}
		c.Set(ret, regime)
	}
	candles.RemoveColumns(1)
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
	tp := HLC3(m)
	rmf := m.Apply(func(mr MatrixRow) float64 {
		return mr.Get(5) * mr.Get(tp)
	})

	// money ratio
	for i := days + 1; i < m.Rows; i++ {
		mfp := 0.0
		mfn := 0.0
		for j := 0; j < days; j++ {
			cur := m.DataRows[i-j]
			prev := m.DataRows[i-1-j]
			if cur.Get(tp) > prev.Get(tp) {
				mfp += cur.Get(rmf)
			} else {
				mfn += cur.Get(rmf)
			}
		}
		if mfn != 0.0 {
			mr := mfp / mfn
			m.DataRows[i].Set(ret, 100.0-100.0/(1.0+mr))
		}
	}
	m.RemoveColumns(2)
	return ret
}

// Larry Williams Trading Indicator
func LWTI(m *Matrix, days int) int {
	ret := m.AddNamedColumn("LWTI")
	//ma = ta.sma(close - nz(close[per]), per)
	d := m.AddColumn()
	for i := days; i < m.Rows; i++ {
		c := &m.DataRows[i]
		p := m.DataRows[i-days]
		c.Set(d, c.Get(4)-p.Get(4))
	}
	s := EMA(m, days, d)
	//atr = ta.atr(per)
	a := ATR(m, days)
	//out = ma/atr * 50 + 50
	for i := days; i < m.Rows; i++ {
		c := &m.DataRows[i]
		if c.Get(a) != 0.0 {
			c.Set(ret, c.Get(s)/c.Get(a)*50.0+50.0)
		} else {
			c.Set(ret, 50.0)
		}
	}
	return ret
}
