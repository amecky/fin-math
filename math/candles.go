package math

import (
	m "math"
)

type CandleDesc struct {
	Body  float64
	Upper float64
	Lower float64
	Range float64
	Trend int
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
	} else {
		cd.Trend = 1
	}
	return cd
}

func (cd CandleDesc) IsBullish() bool {
	return cd.Trend == 1
}

func (cd CandleDesc) IsBearish() bool {
	return cd.Trend == -1
}
