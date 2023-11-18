package math

func TripleEMACategorization(candles *Matrix, e1, e2, e3 int) int {
	indices := []int{0, 1, 2, 4}
	emas := make([]int, 0)
	emas = append(emas, EMA(candles, e1, 4))
	emas = append(emas, EMA(candles, e2, 4))
	emas = append(emas, EMA(candles, e3, 4))
	ret := candles.AddColumn()
	for i := 1; i < candles.Rows; i++ {
		c := &candles.DataRows[i]
		cnt := 0.0
		// C > EMA
		for _, e := range emas {
			if c.Get(4) > c.Get(e) {
				cnt += 1.0
			}
		}
		// EMA > EMA[-1]
		for _, e := range emas {
			if candles.DataRows[i-1].Get(e) < c.Get(e) {
				cnt += 1.0
			}
		}
		// O,H,L,C > EMA
		for _, e := range emas {
			for _, in := range indices {
				if c.Get(in) > c.Get(e) {
					cnt += 1.0
				}
			}
		}
		for j := 1; j < 3; j++ {
			if c.Get(emas[j-1]) > c.Get(emas[j]) {
				cnt += 1.0
			}

		}
		c.Set(ret, cnt/20.0)
	}
	return ret
}
