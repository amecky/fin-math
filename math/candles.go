package math

import (
	"fmt"
	"log"
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

func ConvertMatrixRow(mr MatrixRow) CandleDesc {
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

func (cd CandleDesc) UpperWickPer() float64 {
	return cd.Upper / cd.Range * 100.0
}

func (cd CandleDesc) LowerWickPer() float64 {
	return cd.Lower / cd.Range * 100.0
}

func (cd CandleDesc) BodyPer() float64 {
	return cd.Body / cd.Range * 100.0
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

type OrderBlock struct {
	Key    string
	Upper  float64
	Mid    float64
	Lower  float64
	Gap    float64
	Filled float64
	Type   int
	Index  int
}

func (ob OrderBlock) IsInside(value float64) bool {
	return value < ob.Upper && value > ob.Lower
}

func getOpenOrderBlocks(blocks []OrderBlock, candles *Matrix) []OrderBlock {
	filled := make([]float64, len(blocks))
	for _, ob := range blocks {
		filled = append(filled, ob.Gap)
	}
	for i := range len(blocks) {
		ob := &blocks[i]
		log.Println("==>", ob)
		for j := ob.Index + 1; j < candles.Rows; j++ {
			cur := candles.Row(j)
			if ob.IsInside(cur.Low()) {
				ob.Upper = cur.Low()
				log.Println("Low inside at", cur.Key, "new upper", ob.Upper)
			}
			if ob.IsInside(cur.High()) {
				ob.Lower = cur.High()
				log.Println("High inside at", cur.Key, "new lower", ob.Lower)
			}
			if cur.High() > ob.Upper && cur.Low() < ob.Lower {
				log.Println("completely filled", cur.Key)
				filled[i] = 0.0
				break
			} else {
				filled[i] = ob.Upper - ob.Lower
			}
		}
	}
	nr := make([]OrderBlock, 0)
	for i, o := range blocks {
		if filled[i] > 0.0 {
			o.Filled = (1.0 - filled[i]/o.Gap) * 100.0
			o.Mid = o.Lower + (o.Upper-o.Lower)/2.0
			nr = append(nr, o)
		}
	}
	return nr
}

func GetOrderBlocks(candles *Matrix) []OrderBlock {
	ret := make([]OrderBlock, 0)
	fi := FairValueGaps(candles, false, 1.5)
	log.Println("--------------- OB -----------------------")
	for _, f := range fi {
		log.Println(f)
		if f.Index > 2 {
			p2 := candles.DataRows[f.Index-2]
			tp := 1
			if p2.IsRed() {
				tp = -1
			}
			ret = append(ret, OrderBlock{
				Key:   p2.Key,
				Upper: p2.High(),
				Lower: p2.Low(),
				Mid:   p2.Low() + (p2.High()-p2.Low())/2.0,
				Type:  tp,
				Gap:   p2.High() - p2.Low(),
				Index: f.Index,
			})
		}
	}
	return getOpenOrderBlocks(ret, candles)
}

func FairValueGaps(candles *Matrix, open bool, bodyMultiplier float64) []OrderBlock {
	//threshold = auto ? ta.cum((high - low) / low) / bar_index : thresholdPer / 100
	//bull_fvg = low > high[2] and close[1] > high[2] and (low - high[2]) / high[2] > threshold
	//bear_fvg = high < low[2] and close[1] < low[2] and (low[2] - high) / high > threshold
	avgBody := candles.Apply(func(mr MatrixRow) float64 {
		return m.Abs(mr.Close() - mr.Open())
	})
	ret := make([]OrderBlock, 0)
	bs := NewBitSet("c.Low() > p2.High()", "p1.Close() > p2.High()", fmt.Sprintf("p1.body > %.2f*averageBody", bodyMultiplier))
	for i := 2; i < candles.Rows; i++ {
		bs.ClearAll()
		c := candles.Row(i)
		p1 := candles.Row(i - 1)
		p2 := candles.Row(i - 2)
		bs.SetState(0, c.Low() > p2.High())
		bs.SetState(1, p1.Close() > p2.High())
		bs.SetState(2, p1.Close() > p2.High())
		delta := c.Low() - p2.High()
		bs.SetState(2, m.Abs(p1.Open()-p1.Close()) > c.Get(avgBody)*bodyMultiplier)
		if bs.AllSet() {
			ret = append(ret, OrderBlock{
				Key:   c.Key,
				Upper: c.Low(),
				Lower: p2.High(),
				Mid:   p2.High() + (c.Low()-p2.High())/2.0,
				Gap:   delta,
				Index: i,
				Type:  1,
			})
		}

		bs.ClearAll()
		bs.SetState(0, c.High() < p2.Low())
		bs.SetState(1, p1.Close() < p2.High())
		delta = p2.Low() - c.High()
		bs.SetState(2, m.Abs(p1.Open()-p1.Close()) > c.Get(avgBody)*bodyMultiplier)
		if bs.AllSet() {
			delta := p2.Low() - c.High()
			ret = append(ret, OrderBlock{
				Key:   c.Key,
				Upper: p2.Low(),
				Lower: c.High(),
				Mid:   c.High() + delta/2.0, //c.High() + (p2.Low()-c.High())/2.0,
				Gap:   delta,
				Index: i,
				Type:  -1,
			})
		}
	}
	if open {
		return getOpenOrderBlocks(ret, candles)
	}
	return ret
}

func OrderBlocks(candles *Matrix) int {
	// 0 = Upper 1 = Lower 2 = Type 3 = Filled 4 = Gap
	upper := candles.AddColumn()
	lower := candles.AddColumn()
	tp := candles.AddColumn()
	filled := candles.AddColumn()
	gap := candles.AddColumn()
	fi := FVG(candles)
	for i := 2; i < candles.Rows-2; i++ {
		c := candles.DataRows[i]
		if (c.Get(fi+2) == 1.0 || c.Get(fi+2) == -1) && c.Get(fi+3) == 0.0 {
			p2 := candles.DataRows[i-2]
			p2.Set(upper, p2.High())
			p2.Set(lower, p2.Low())
			if p2.IsGreen() {
				p2.Set(tp, 1)
			} else {
				p2.Set(tp, -1)
			}
			d := p2.High() - p2.Low()
			p2.Set(gap, d)
			p2.Set(filled, 0)
		}
	}
	// check for filled OB
	for i := 0; i < candles.Rows; i++ {
		c := candles.Row(i)
		if c.Get(tp) == 1 {
			for j := i + 1; j < candles.Rows; j++ {
				cur := candles.Row(j)
				if cur.Low() < c.Get(upper) {
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
	candles.RemoveColumns(5)
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
	bodyMultiplier := 1.5
	avgBody := candles.Apply(func(mr MatrixRow) float64 {
		return m.Abs(mr.Close() - mr.Open())
	})
	for i := 2; i < candles.Rows; i++ {
		c := candles.Row(i)
		p1 := candles.Row(i - 1)
		p2 := candles.Row(i - 2)
		// bull_fvg = c.low > high[2] and close[1] > high[2] and body[1] > mul * averageBody
		if c.Low() > p2.High() && p1.Close() > p2.High() && m.Abs(p1.Close()-p1.Open()) > c.Get(avgBody)*bodyMultiplier {
			c.Set(upper, c.Low())
			c.Set(lower, p2.High())
			c.Set(tp, 1.0)
			d := c.Low() - p2.High()
			c.Set(gap, d)
		} else if c.High() < p2.Low() && p1.Close() < p2.High() && m.Abs(p1.Close()-p1.Open()) > c.Get(avgBody)*bodyMultiplier {
			// bear_fvg = high < low[2] and close[1] < low[2] and body[1] > mul * averageBody
			c.Set(upper, p2.Low())
			c.Set(lower, c.High())
			c.Set(tp, -1.0)
			d := p2.Low() - c.High()
			c.Set(gap, d)
		}
	}
	// check for filled FVG
	for i := 0; i < candles.Rows; i++ {
		c := candles.Row(i)
		if c.Get(tp) == 1 {
			for j := i + 1; j < candles.Rows; j++ {
				cur := candles.Row(j)
				if cur.Low() < c.Get(upper) {
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
		c := candles.DataRows[i]
		cd := ConvertMatrixRow(c)
		if cd.Range > 0.0 {
			c.Set(up, cd.Upper/cd.Range*100.0)
			c.Set(bs, cd.Body/cd.Range*100.0)
			c.Set(lo, cd.Lower/cd.Range*100.0)
			c.Set(rng, cd.Range)
			c.Set(trn, float64(cd.Trend))
			c.Set(rj, cd.Rejection/cd.Range*100.0)
		}
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
