package math

import (
	m "math"
)

// Vortex indicator
// https://www.investopedia.com/terms/v/vortex-indicator-vi.asp
func Vortex(candles *Matrix, period int) int {
	vp := candles.AddNamedColumn("VP")
	vm := candles.AddNamedColumn("VM")
	t1 := candles.AddColumn()
	t2 := candles.AddColumn()
	for i := 1; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		p := candles.DataRows[i-1]
		c.Set(t1, m.Abs(c.Get(1)-p.Get(2)))
		c.Set(t2, m.Abs(c.Get(2)-p.Get(1)))

	}
	//VMP = math.sum( math.abs( high - low[1]), period_ )
	//VMM = math.sum( math.abs( low - high[1]), period_ )
	s1 := candles.Sum(t1, period)
	s2 := candles.Sum(t2, period)
	ai := ATR(candles, 1)
	str := candles.Sum(ai, period)
	//STR = math.sum( ta.atr(1), period_ )
	for i := 1; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		//VIP = VMP / STR
		c.Set(vp, c.Get(s1)/c.Get(str))
		//VIM = VMM / STR
		c.Set(vm, c.Get(s2)/c.Get(str))
	}
	candles.RemoveColumns(7)
	return vp
}

type Reversal struct {
	Index   int
	Trend   int
	End     int
	Closest float64
	Value   float64
}

func Reversals(candles *Matrix, period int, threshold float64) []Reversal {
	indices := make([]Reversal, 0)
	ai := ATR(candles, period)
	for i := 0; i < candles.Rows-1; i++ {
		c := candles.DataRows[i]
		bs := m.Abs(c.Get(0) - c.Get(4))
		rbs := 0.0
		if c.Get(ai) != 0.0 {
			rbs = bs / c.Get(ai) * 100.0
		}
		if rbs >= threshold {
			tr := 1
			if c.Get(0) > c.Get(4) {
				tr = -1
			}
			val := c.Get(2)
			if tr == -1 {
				val = c.Get(1)
			}
			indices = append(indices, Reversal{
				Index: i,
				Trend: tr,
				End:   -1,
				Value: val,
			})
		}
	}
	for i := 0; i < len(indices); i++ {
		cin := &indices[i]
		for j := cin.Index + 1; j < candles.Rows; j++ {
			if cin.Trend == 1 {
				if candles.DataRows[j].Get(2) <= candles.DataRows[cin.Index].Get(2) && cin.End == -1 {
					cin.End = j

				}
			} else {
				if candles.DataRows[j].Get(1) >= candles.DataRows[cin.Index].Get(1) && cin.End == -1 {
					cin.End = j
				}
			}
		}
		if cin.End == -1 {
			d := candles.Last().Get(4) - candles.DataRows[cin.Index].Get(2)
			if cin.Trend == -1 {
				d = candles.DataRows[cin.Index].Get(1) - candles.Last().Get(4)
			}
			cin.Closest = d
		}
	}
	return indices
}
