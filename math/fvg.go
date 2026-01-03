package math

import (
	"fmt"
	m "math"
)

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

func FindFairValueGaps3(candles *Matrix, bodyMultiplier float64) []FairValueGap {
	var gaps = make([]FairValueGap, 0)
	avgBody := candles.Apply(func(mr MatrixRow) float64 {
		return m.Abs(mr.Close() - mr.Open())
	})
	for i := 2; i < candles.Rows; i++ {
		c := candles.Row(i)
		mid := candles.Row(i - 1)
		p := candles.Row(i - 2)
		if c.Get(HIGH) < p.Get(LOW) && m.Abs(mid.Close()-mid.Open()) > c.Get(avgBody)*bodyMultiplier {
			gaps = append(gaps, FairValueGap{
				Timestamp: p.Key,
				Filled:    0,
				Upper:     p.Get(LOW),
				Lower:     c.Get(HIGH),
				Index:     i - 2,
				Type:      -1,
			})
		}
		if c.Get(LOW) > p.Get(HIGH) && m.Abs(mid.Close()-mid.Open()) > c.Get(avgBody)*bodyMultiplier {
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

func FindFairValueGaps(candles *Matrix, bodyMultiplier float64) []Gap {
	avgBody := candles.Apply(func(mr MatrixRow) float64 {
		return m.Abs(mr.Close() - mr.Open())
	})
	var gaps = make([]Gap, 0)
	for i := 2; i < candles.Rows; i++ {
		c := candles.Row(i)
		mid := candles.Row(i - 1)
		p := candles.Row(i - 2)
		if c.Get(HIGH) < p.Get(LOW) && m.Abs(mid.Close()-mid.Open()) > c.Get(avgBody)*bodyMultiplier {
			gaps = append(gaps, Gap{
				Timestamp: c.Key,
				Filled:    0,
				Upper:     p.Get(LOW),
				Lower:     c.Get(HIGH),
				Index:     i - 2,
			})
		}
		if c.Get(LOW) > p.Get(HIGH) && m.Abs(mid.Close()-mid.Open()) > c.Get(avgBody)*bodyMultiplier {
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
