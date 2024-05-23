package math

import (
	m "math"
)

type CandleDesc struct {
	Body      float64
	Upper     float64
	Lower     float64
	Range     float64
	Open      float64
	High      float64
	Low       float64
	Close     float64
	Rejection float64
	Trend     int
}

func NewCandleDesc(mr *MatrixRow) CandleDesc {
	cd := CandleDesc{
		Upper: mr.Get(1) - m.Max(mr.Get(0), mr.Get(4)),
		Lower: m.Min(mr.Get(0), mr.Get(4)) - mr.Get(2),
		Body:  m.Abs(mr.Get(0) - mr.Get(4)),
		Range: mr.Get(1) - mr.Get(2),
		Open:  mr.Get(0),
		High:  mr.Get(1),
		Low:   mr.Get(2),
		Close: mr.Get(4),
	}
	if mr.Get(0) > mr.Get(4) {
		cd.Trend = -1
	} else {
		cd.Trend = 1
	}
	return cd
}

func ConvertMatrixRow(mr *MatrixRow) CandleDesc {
	cd := CandleDesc{
		Upper: mr.Get(1) - m.Max(mr.Get(0), mr.Get(4)),
		Lower: m.Min(mr.Get(0), mr.Get(4)) - mr.Get(2),
		Body:  m.Abs(mr.Get(0) - mr.Get(4)),
		Range: mr.Get(1) - mr.Get(2),
	}
	if mr.Get(0) > mr.Get(4) {
		cd.Trend = -1
		cd.Rejection = cd.Lower
	} else {
		cd.Trend = 1
		cd.Rejection = cd.Upper
	}
	return cd
}

func (cd CandleDesc) IsBullish() bool {
	return cd.Trend == 1
}

func (cd CandleDesc) IsBearish() bool {
	return cd.Trend == -1
}

func KPivots(m *Matrix, period int) int {
	// 0 = Upper 1 = Lower 2 = MID
	upper := m.AddNamedColumn("KPUP")
	lower := m.AddNamedColumn("KPLO")
	mid := m.AddNamedColumn("KPMID")
	sa := SMA(m, period, 4)
	th := m.AddColumn()
	tl := m.AddColumn()
	for i := period; i < m.Rows; i++ {
		hi := m.FindMaxBetween(1, i-period, period)
		lo := m.FindMinBetween(2, i-period, period)
		m.DataRows[i].Set(mid, (hi+lo+m.DataRows[i].Get(sa))/3.0)
		m.DataRows[i].Set(th, hi)
		m.DataRows[i].Set(tl, lo)
	}
	step := period / 2
	//	To find the support, use this formula: Support = (Lowest K’s pivot point of the last 12 periods * 2) -Step 1.
	// To find the resistance, use this formula: Resistance = (Highest K’s pivot point of the last 12 periods * 2) — Step 2.
	for i := step; i < m.Rows; i++ {
		hi := m.FindMaxBetween(mid, i-step, step)
		m.DataRows[i].Set(upper, hi*2.0-m.DataRows[i].Get(tl))
		lo := m.FindMinBetween(mid, i-step, step)
		m.DataRows[i].Set(lower, lo*2.0-m.DataRows[i].Get(th))
	}
	m.RemoveColumns(3)
	return upper
}

func FVG(candles *Matrix) int {
	// 0 = Upper 1 = Lower 2 = Type 3 = Filled 4 = Gap
	// 0 = FVG 1 = Filled 2 = gap

	upper := candles.AddColumn()
	lower := candles.AddColumn()
	tp := candles.AddColumn()
	//bu := candles.AddColumn()
	filled := candles.AddColumn()
	gap := candles.AddColumn()
	for i := 0; i < candles.Rows; i++ {
		if i > 1 {
			c := candles.Row(i)
			p1 := candles.Row(i - 1)
			p2 := candles.Row(i - 2)
			// bull_fvg = low > high[2] and close[1] > high[2]
			if c.Low() > p2.High() && p1.Close() > p2.High() {
				c.Set(upper, c.Low())
				c.Set(lower, p2.High())
				c.Set(tp, 1.0)
				d := c.Low() - p2.High()
				c.Set(gap, d)
			} else if c.High() < p2.Low() && p1.Close() < p2.High() {
				// bear_fvg = high < low[2] and close[1] < low[2]
				c.Set(upper, p2.Low())
				c.Set(lower, c.High())
				c.Set(tp, -1.0)
				d := p2.Low() - c.High()
				c.Set(gap, d)
			}
		}
	}
	// check for filled FVG
	for i := 0; i < candles.Rows; i++ {
		c := candles.Row(i)
		if c.Get(tp) == 1 {
			for j := i + 1; j < candles.Rows; j++ {
				cur := candles.Row(j)
				if cur.Low() <= c.Get(upper) {
					c.Set(upper, cur.Low())
					if cur.Low() <= c.Get(lower) {
						c.Set(filled, 1.0)
						break
					}
				}
			}

		}
		if c.Get(tp) == -1 {
			for j := i + 1; j < candles.Rows; j++ {
				cur := candles.Row(j)
				if cur.High() > c.Get(lower) {
					c.Set(lower, cur.High())
					if cur.High() >= c.Get(upper) {
						c.Set(filled, 1.0)
						break
					}
				}
			}

		}
	}
	return upper
}

func RelativeCandleDescriptors(candles *Matrix) int {
	up := candles.AddNamedColumn("Upper")
	bs := candles.AddNamedColumn("Body")
	lo := candles.AddNamedColumn("Lower")
	rng := candles.AddNamedColumn("Range")
	trn := candles.AddNamedColumn("Trend")
	rj := candles.AddNamedColumn("Rejection")
	// 0 = Upper 1 = Body 2 = Lower 3 = Range 4 = Trend 5 = Rejection
	for i := 0; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		cd := ConvertMatrixRow(c)
		c.Set(up, cd.Upper/cd.Range*100.0)
		c.Set(bs, cd.Body/cd.Range*100.0)
		c.Set(lo, cd.Lower/cd.Range*100.0)
		c.Set(rng, cd.Range)
		c.Set(trn, float64(cd.Trend))
		c.Set(rj, cd.Rejection/cd.Range*100.0)
	}
	return up
}

// Internal Bar Strength
func IBS(candles *Matrix) int {
	ret := candles.AddNamedColumn("IBS")
	// 0 = Internal Bar Strength
	for i := 0; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		c.Set(ret, (c.Close()-c.Low())/(c.High()-c.Low()))
	}
	return ret
}
