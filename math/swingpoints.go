package math

func (m *Matrix) FindSwingPoints() SwingPoints {
	var tmp SwingPoints
	lv := 0.0
	hv := 0.0
	for i := 2; i < m.Rows-2; i++ {
		p1 := m.DataRows[i-2]
		p2 := m.DataRows[i-1]
		pc := m.DataRows[i]
		p3 := m.DataRows[i+1]
		p4 := m.DataRows[i+2]
		if p1.Get(1) < pc.Get(1) && p2.Get(1) < pc.Get(1) && p3.Get(1) < pc.Get(1) && p4.Get(1) < pc.Get(1) {
			sp := SwingPoint{
				Timestamp: pc.Key,
				Type:      High,
				Value:     pc.Get(1),
				Price:     pc.Get(4),
				Index:     i,
				BaseType:  High,
			}
			if sp.Value > hv {
				sp.Type = HigherHigh
				sp.Trend = 1
				hv = sp.Value
			} else {
				sp.Trend = -1
				sp.Type = LowerHigh
				hv = sp.Value
			}
			tmp = append(tmp, sp)
		}
		if p1.Get(2) > pc.Get(2) && p2.Get(2) > pc.Get(2) && p3.Get(2) > pc.Get(2) && p4.Get(2) > pc.Get(2) {
			sp := SwingPoint{
				Timestamp: pc.Key,
				Type:      Low,
				Value:     pc.Get(2),
				Price:     pc.Get(4),
				Index:     i,
				BaseType:  Low,
			}
			if sp.Value < lv {
				sp.Trend = -1
				sp.Type = LowerLow
				lv = sp.Value
			} else {
				sp.Trend = 1
				sp.Type = HigherLow
				lv = sp.Value
			}
			tmp = append(tmp, sp)
		}
	}
	for i := 1; i < len(tmp); i++ {
		c := &tmp[i]
		p := tmp[i-1]
		c.Delta = c.Value - p.Value
	}
	for i := 0; i < len(tmp); i++ {
		c := &tmp[i]
		for j := c.Index; j < m.Rows; j++ {
			if c.BaseType == High && m.DataRows[j].High() > c.Value {
				c.Broken = true
			}
			if c.BaseType == Low && m.DataRows[j].Low() < c.Value {
				c.Broken = true
			}
		}
	}
	return tmp
}

func (m *Matrix) FindTurningPoints(field int) SwingPoints {
	var tmp SwingPoints
	lv := 0.0
	hv := 0.0
	for i := 1; i < m.Rows-1; i++ {
		p := m.DataRows[i-1]
		c := m.DataRows[i]
		n := m.DataRows[i+1]
		if p.Get(field) < c.Get(field) && n.Get(field) < c.Get(field) {
			sp := SwingPoint{
				Timestamp: c.Key,
				Type:      High,
				Value:     c.Get(field),
				Price:     c.Get(field),
				Index:     i,
				BaseType:  High,
			}
			if sp.Value > hv {
				sp.Type = HigherHigh
				hv = sp.Value
			} else {
				sp.Type = LowerHigh
				hv = sp.Value
			}
			tmp = append(tmp, sp)
		}
		if p.Get(field) > c.Get(field) && n.Get(field) > c.Get(field) {
			sp := SwingPoint{
				Timestamp: c.Key,
				Type:      Low,
				Value:     c.Get(field),
				Price:     c.Get(field),
				Index:     i,
				BaseType:  Low,
			}
			if sp.Value < lv {
				sp.Type = LowerLow
				lv = sp.Value
			} else {
				sp.Type = HigherLow
				lv = sp.Value
			}
			tmp = append(tmp, sp)
		}
	}
	return tmp
}

func (m *Matrix) Fractals(size int) SwingPoints {
	var tmp SwingPoints
	dist := size / 2
	lv := 0.0
	hv := 0.0
	for i := dist; i < m.Rows-dist; i++ {
		c := m.DataRows[i]
		ch := 0
		cl := 0
		for j := 1; j <= dist; j++ {
			idx := i - j
			if m.DataRows[idx].High() > c.High() {
				ch++
			}
			if m.DataRows[idx].Low() < c.Low() {
				cl++
			}
			idx = i + j
			if m.DataRows[idx].High() > c.High() {
				ch++
			}
			if m.DataRows[idx].Low() < c.Low() {
				cl++
			}
		}
		if ch == 0 {
			sp := SwingPoint{
				Timestamp: c.Key,
				Type:      High,
				Value:     c.Get(1),
				Price:     c.Get(4),
				Index:     i,
				BaseType:  High,
			}
			if sp.Value > hv {
				sp.Type = HigherHigh
				sp.Trend = 1
				hv = sp.Value
			} else {
				sp.Trend = -1
				sp.Type = LowerHigh
				hv = sp.Value
			}
			tmp = append(tmp, sp)
		}
		if cl == 0 {
			sp := SwingPoint{
				Timestamp: c.Key,
				Type:      Low,
				Value:     c.Get(2),
				Price:     c.Get(4),
				Index:     i,
				BaseType:  Low,
			}
			if sp.Value < lv {
				sp.Trend = -1
				sp.Type = LowerLow
				lv = sp.Value
			} else {
				sp.Trend = 1
				sp.Type = HigherLow
				lv = sp.Value
			}
			tmp = append(tmp, sp)
		}
	}
	for i := 1; i < len(tmp); i++ {
		c := &tmp[i]
		p := tmp[i-1]
		c.Delta = c.Value - p.Value
	}
	for i := 0; i < len(tmp); i++ {
		c := &tmp[i]
		for j := c.Index; j < m.Rows; j++ {
			if c.BaseType == High && m.DataRows[j].High() > c.Value {
				c.Broken = true
			}
			if c.BaseType == Low && m.DataRows[j].Low() < c.Value {
				c.Broken = true
			}
		}
	}
	return tmp
}

func (m *Matrix) FindSwingPointsByField(field int) SwingPoints {
	var tmp SwingPoints
	lv := 0.0
	hv := 0.0
	for i := 2; i < m.Rows-2; i++ {
		p1 := m.DataRows[i-2]
		p2 := m.DataRows[i-1]
		pc := m.DataRows[i]
		p3 := m.DataRows[i+1]
		p4 := m.DataRows[i+2]
		if p1.Get(field) < pc.Get(field) && p2.Get(field) < pc.Get(field) && p3.Get(field) < pc.Get(field) && p4.Get(field) < pc.Get(field) {
			sp := SwingPoint{
				Timestamp: pc.Key,
				Type:      High,
				Value:     pc.Get(field),
				Price:     pc.Get(4),
				Index:     i,
				BaseType:  High,
			}
			if sp.Value > hv {
				sp.Type = HigherHigh
				hv = sp.Value
			} else {
				sp.Type = LowerHigh
				hv = sp.Value
			}
			tmp = append(tmp, sp)
		}
		if p1.Get(field) > pc.Get(field) && p2.Get(field) > pc.Get(field) && p3.Get(field) > pc.Get(field) && p4.Get(field) > pc.Get(field) {
			sp := SwingPoint{
				Timestamp: pc.Key,
				Type:      Low,
				Value:     pc.Get(field),
				Price:     pc.Get(4),
				Index:     i,
				BaseType:  Low,
			}
			if sp.Value < lv {
				sp.Type = LowerLow
				lv = sp.Value
			} else {
				sp.Type = HigherLow
				lv = sp.Value
			}
			tmp = append(tmp, sp)
		}
	}
	return tmp
}
