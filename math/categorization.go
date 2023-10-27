package math

func TripleEMACategorization(m *Matrix, e1, e2, e3 int) int {
	ema1 := EMA(m, e1, ADJ_CLOSE)
	ema2 := EMA(m, e2, ADJ_CLOSE)
	ema3 := EMA(m, e3, ADJ_CLOSE)
	res := m.AddColumn()
	for i := 1; i < m.Rows; i++ {
		cnt := 0.0
		c := &m.DataRows[i]
		p := m.DataRows[i-1]
		if c.Get(ema1) > p.Get(ema1) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(ema2) > p.Get(ema2) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(ema3) > p.Get(ema3) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(ema1) > c.Get(ema2) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(ema2) > c.Get(ema3) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(4) > c.Get(ema1) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(4) > c.Get(ema2) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(4) > c.Get(ema3) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(ema1)-p.Get(ema1) > c.Get(ema2)-p.Get(ema2) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if c.Get(ema2)-p.Get(ema2) > c.Get(ema3)-p.Get(ema3) {
			cnt += 1.0
		} else {
			cnt -= 1.0
		}
		if i > 5 {
			if c.Get(ema1) > m.DataRows[i-5].Get(ema1) {
				cnt += 1.0
			} else {
				cnt -= 1.0
			}
			if c.Get(ema2) > m.DataRows[i-5].Get(ema2) {
				cnt += 1.0
			} else {
				cnt -= 1.0
			}
			if c.Get(ema3) > m.DataRows[i-5].Get(ema3) {
				cnt += 1.0
			} else {
				cnt -= 1.0
			}
			c.Set(res, cnt/13.0)
		} else {
			c.Set(res, cnt/10.0)
		}
	}
	return res
}
