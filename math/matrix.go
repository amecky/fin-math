package math

import (
	"fmt"
	m "math"
	ma "math"
	"strings"
)

const (
	OPEN      int = 0
	HIGH      int = 1
	LOW       int = 2
	CLOSE     int = 3
	ADJ_CLOSE int = 4
	VOLUME    int = 5
)

type MajorLevel struct {
	Value float64
	Count int
	Date  string
	Index int
}

type MajorLevels struct {
	Levels    []MajorLevel
	Threshold float64
}

func NewMajorLevels(threshold float64) *MajorLevels {
	ret := MajorLevels{
		Threshold: threshold,
	}
	ret.Levels = make([]MajorLevel, 0)
	return &ret
}

func (ml *MajorLevels) Add(value float64, count int, date string) {
	cnt := 0
	for _, cm := range ml.Levels {
		d := (cm.Value/value - 1.0) * 100.0
		if cnt == 0 && m.Abs(d) < ml.Threshold {
			cnt = 1
		}
	}
	if cnt == 0 {
		ml.Levels = append(ml.Levels, MajorLevel{
			Value: value,
			Count: count,
			Date:  date,
		})
	}
}

type SwingPointType int

const (
	High       SwingPointType = 1
	Low        SwingPointType = -1
	LowerHigh  SwingPointType = 2
	HigherHigh SwingPointType = 3
	HigherLow  SwingPointType = -2
	LowerLow   SwingPointType = -3
)

type ClusterPoint struct {
	Value      float64
	Count      int
	Timestamps []string
	Values     []float64
	Average    float64
}

type SwingPoint struct {
	Timestamp string
	BaseType  SwingPointType
	Type      SwingPointType
	Value     float64
	Price     float64
	Index     int
}

func (s SwingPointType) String() string {
	switch s {
	case LowerLow:
		return "LL"
	case HigherLow:
		return "HL"
	case Low:
		return "L"
	case 0:
		return "-"
	case High:
		return "H"
	case LowerHigh:
		return "LH"
	case HigherHigh:
		return "HH"
	}
	return "-"
}

type SwingPoints []SwingPoint

func (sps *SwingPoints) Add(timestamp string, trend SwingPointType, value, price float64) {
	*sps = append(*sps, SwingPoint{
		Timestamp: timestamp,
		Type:      trend,
		Value:     value,
		Price:     price,
	})
}

func (sps SwingPoints) FindByDate(timestamp string) *SwingPoint {
	for _, s := range sps {
		if s.Timestamp == timestamp {
			return &s
		}
	}
	return nil
}

func (sps SwingPoints) FilterByType(baseType SwingPointType) SwingPoints {
	var ret SwingPoints
	for _, s := range sps {
		if s.BaseType == baseType {
			ret = append(ret, s)
		}
	}
	return ret
}

func (sps SwingPoints) FindRecentByBaseType(baseType SwingPointType) *SwingPoint {
	idx := -1
	for i, s := range sps {
		if s.BaseType == baseType {
			idx = i
		}
	}
	if idx == -1 {
		return nil
	} else {
		return &sps[idx]
	}
}

func (sps SwingPoints) GetAngle(first, second int) float64 {
	steps := float64(sps[second].Index - sps[first].Index)
	y1 := sps[first].Value
	y2 := sps[second].Value
	return (y2 - y1) / steps //m.Atan((y2-y1)/steps) * 180.0 / m.Pi
}

type MatrixRow struct {
	Key     string
	Num     int
	Values  []float64
	Comment string
}
type Matrix struct {
	Info     string
	Rows     int
	Cols     int
	DataRows []MatrixRow
	Headers  []string
}

func NewMatrix(cols int) *Matrix {
	m := Matrix{
		Cols:    cols,
		Headers: []string{"Key"},
	}
	for i := 1; i < cols; i++ {
		m.Headers = append(m.Headers, "")
	}
	return &m
}

func NewMatrixWithHeaders(cols int, headers []string) *Matrix {
	m := Matrix{
		Cols:    0,
		Headers: []string{"Key"},
	}
	for _, h := range headers {
		m.Headers = append(m.Headers, h)
		m.Cols++
	}
	return &m
}

func (m *Matrix) AddRow(key string) *MatrixRow {
	r := m.FindRow(key)
	if r == nil {
		mr := MatrixRow{
			Key:    key,
			Values: make([]float64, 0),
			Num:    m.Cols,
		}
		for i := 0; i < m.Cols; i++ {
			mr.Values = append(mr.Values, 0.0)
		}
		m.DataRows = append(m.DataRows, mr)
		m.Rows++
		return &m.DataRows[m.Rows-1]
	}
	return r
}

func (m *Matrix) AddNamedColumn(header string) int {
	m.Cols++
	m.Headers = append(m.Headers, header)
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Num++
		m.DataRows[i].Values = append(m.DataRows[i].Values, 0.0)
	}
	return m.Cols - 1
}

func (m *Matrix) AddColumn() int {
	m.Cols++
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Num++
		m.DataRows[i].Values = append(m.DataRows[i].Values, 0.0)
	}
	m.Headers = append(m.Headers, "")
	return m.Cols - 1
}

func (m Matrix) FindRow(key string) *MatrixRow {
	for _, r := range m.DataRows {
		if r.Key == key {
			return &r
		}
	}
	return nil
}

func (m Matrix) FindRowIndex(key string) int {
	for i, r := range m.DataRows {
		if r.Key == key {
			return i
		}
	}
	return -1
}

func (m Matrix) Last() *MatrixRow {
	if m.Rows > 0 {
		return &m.DataRows[m.Rows-1]
	}
	return nil
}

func (m *Matrix) PartialSum(field, start, count int) float64 {
	sum := 0.0
	if start >= m.Rows {
		return 0.0
	}
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	for i := start; i < end; i++ {
		sum += m.DataRows[i].Values[field]
	}
	return sum
}

func (m *Matrix) String() string {
	var builder strings.Builder
	for _, d := range m.DataRows {
		builder.WriteString(d.String())
		builder.WriteString("\n")
	}
	return builder.String()
}

func (m *MatrixRow) Set(index int, value float64) *MatrixRow {
	if m != nil {
		if index < m.Num {
			m.Values[index] = value
		}
		return m
	}
	return nil
}

func (m *MatrixRow) SetComment(cmt string) *MatrixRow {
	if m != nil {
		m.Comment = cmt
		return m
	}
	return nil
}

func (m *Matrix) SetComment(row int, cmt string) {
	if m.Rows > row {
		m.DataRows[row].Comment = cmt
	}
}

func (m *MatrixRow) Get(index int) float64 {
	if index >= 0 && index < m.Num {
		return m.Values[index]
	}
	return 0.0
}

func (m *MatrixRow) ShortKey() string {
	idx := strings.Index(m.Key, " ")
	if idx != -1 {
		return m.Key[0:idx]
	}
	return m.Key
}

func (m *MatrixRow) Sum() float64 {
	sum := 0.0
	for i := 0; i < m.Num; i++ {
		sum += m.Values[i]
	}
	return sum
}

func (m *MatrixRow) PartialSum(start, count int) float64 {
	sum := 0.0
	if start >= m.Num {
		return 0.0
	}
	end := start + count
	if end > m.Num {
		end = m.Num
	}
	for i := start; i < end; i++ {
		sum += m.Values[i]
	}
	return sum
}

func (mr MatrixRow) String() string {
	var builder strings.Builder
	builder.WriteString("Key: ")
	builder.WriteString(mr.Key)
	for _, n := range mr.Values {
		builder.WriteString(fmt.Sprintf("%.2f ", n))
	}
	builder.WriteString("C: ")
	builder.WriteString(mr.Comment)
	return builder.String()
}

func (m *Matrix) FindMinMaxIndexBetween(field, start, count int) (int, int) {
	if m.Rows < 1 {
		return -1, -1
	}
	if start < 0 {
		start = 0
	}
	if start >= m.Rows {
		start = m.Rows - 1
	}
	rma := 0
	rmi := 0
	max := m.DataRows[start].Get(field)
	min := m.DataRows[start].Get(field)
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	for i := start + 1; i < end; i++ {
		cur := m.DataRows[i].Get(field)
		if cur > max {
			max = cur
			rma = i
		}
		if cur < min {
			min = cur
			rmi = i
		}
	}
	return rmi, rma
}

func (m *Matrix) FindMinMaxBetween(field, start, count int) (float64, float64) {
	if m.Rows < 1 {
		return 0.0, 0.0
	}
	if start < 0 {
		start = 0
	}
	if start >= m.Rows {
		start = m.Rows - 1
	}
	max := m.DataRows[start].Get(field)
	min := m.DataRows[start].Get(field)
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	for i := start + 1; i < end; i++ {
		cur := m.DataRows[i].Get(field)
		if cur > max {
			max = cur
		}
		if cur < min {
			min = cur
		}
	}
	return min, max
}

func (m *Matrix) FindHighLowIndex(start, count int) (int, int) {
	end := start + count
	low := 100000.0
	high := 0.0
	hIdx := start
	lIdx := start
	for i := start; i < end; i++ {
		h := m.DataRows[i].Get(1)
		if h >= high {
			hIdx = i
			high = h
		}
		l := m.DataRows[i].Get(2)
		if l <= low {
			lIdx = i
			low = l
		}
	}
	return hIdx, lIdx
}

func (m *Matrix) FindHighestHighLowestLow(start, count int) (float64, float64) {
	if m.Rows == 0 {
		return 0.0, 0.0
	}
	if start < 0 {
		start = 0
	}
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	low := m.DataRows[start].Get(2)
	high := m.DataRows[start].Get(1)
	for i := start; i < end; i++ {
		h := m.DataRows[i].Get(1)
		if h >= high {
			high = h
		}
		l := m.DataRows[i].Get(2)
		if l <= low {
			low = l
		}
	}
	return high, low
}

func (m *Matrix) FindMinBetween(field, start, count int) float64 {
	if start >= m.Rows {
		start = m.Rows - 1
	}
	min := m.DataRows[start].Get(field)
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	for i := start + 1; i < end; i++ {
		cur := m.DataRows[i].Get(field)
		if cur < min {
			min = cur
		}
	}
	return min
}

func (m *Matrix) FindMaxBetween(field, start, count int) float64 {
	if start >= m.Rows {
		start = m.Rows - 1
	}
	max := m.DataRows[start].Get(field)
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	for i := start + 1; i < end; i++ {
		cur := m.DataRows[i].Get(field)
		if cur > max {
			max = cur
		}
	}
	return max
}

func (m *Matrix) Stochastic(days int, field int) int {
	ret := m.AddColumn()
	total := m.Rows
	if total > days {

		for i := days; i < m.Rows; i++ {
			low, high := m.FindMinMaxBetween(field, i-days+1, days)
			value := (m.DataRows[i].Get(field) - low) / (high - low) * 100.0
			m.DataRows[i].Set(ret, value)
		}
	}
	return ret
}

func (m *Matrix) Add(field int, other *Matrix, otherField int) int {
	ret := m.AddColumn()
	for _, s := range other.DataRows {
		row := m.FindRow(s.Key)
		if row != nil {
			row.Set(ret, row.Get(field)+s.Get(otherField))
		}
	}
	return ret
}

func (m *Matrix) Subtract(field, otherField int) int {
	ret := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(field)-m.DataRows[i].Get(otherField))
	}
	return ret
}

func (m *Matrix) Sum(field, length int) int {
	ret := m.AddColumn()
	for i := length; i < m.Rows; i++ {
		s := 0.0
		for j := 0; j < length; j++ {
			s += m.DataRows[i-j].Get(field)
		}
		m.DataRows[i].Set(ret, s)
	}
	return ret
}

func (m *Matrix) SlopePercentage(field, lookback int) int {
	ret := m.AddColumn()
	for i := lookback; i < m.Rows; i++ {
		cur := m.DataRows[i].Get(field)
		prev := m.DataRows[i-lookback].Get(field)
		if prev != 0.0 {
			s := (cur/prev - 1.0) * 100.0
			m.DataRows[i].Set(ret, s)
		}
	}
	return ret
}

func (m *Matrix) Categorize(field, target int, check func(mr MatrixRow) float64) {
	for i, s := range m.DataRows {
		v := check(s)
		m.DataRows[i].Set(target, v)
	}
}

func (m *Matrix) Categorize2(target int, check func(mr MatrixRow) float64) {
	for i, s := range m.DataRows {
		v := check(s)
		m.DataRows[i].Set(target, v)
	}
}

func (m *Matrix) Apply(fn func(mr MatrixRow) float64) int {
	ret := m.AddColumn()
	for i, s := range m.DataRows {
		v := fn(s)
		m.DataRows[i].Set(ret, v)
	}
	return ret
}

func (m *Matrix) ApplyRow(row int, fn func(mr MatrixRow) float64) {
	for i, s := range m.DataRows {
		v := fn(s)
		m.DataRows[i].Set(row, v)
	}
}

func (m *Matrix) RemoveColumn() {
	for j := 0; j < m.Rows; j++ {
		m.DataRows[j].Values = m.DataRows[j].Values[:m.DataRows[j].Num-1]
		m.DataRows[j].Num--
	}
	m.Headers = m.Headers[:len(m.Headers)-1]
	m.Cols--
}

func (m *Matrix) RemoveColumns(cnt int) {
	for i := 0; i < cnt; i++ {
		m.RemoveColumn()
	}
}

func (m *Matrix) CopyColumn(source, destination int) {
	for j := 0; j < m.Rows; j++ {
		m.DataRows[j].Values[destination] = m.DataRows[j].Values[source]
	}
}
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
				hv = sp.Value
			} else {
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
				Value:     c.Get(1),
				Price:     c.Get(4),
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
				Value:     c.Get(2),
				Price:     c.Get(4),
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

func (m *Matrix) Sublist(start, count int) *Matrix {
	ret := NewMatrix(m.Cols)
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	for i := start; i < end; i++ {
		ret.DataRows = append(ret.DataRows, m.DataRows[i])
		ret.Rows++
	}
	return ret
}

func (m *Matrix) Recent(start int) *Matrix {
	ret := NewMatrix(m.Cols)
	for i := start; i < m.Rows; i++ {
		ret.DataRows = append(ret.DataRows, m.DataRows[i])
		ret.Rows++
	}
	return ret
}

func (m *Matrix) FilterByKeys(start, end string) *Matrix {
	ret := NewMatrix(m.Cols)
	idx := strings.Index(start, " ")
	if idx == -1 {
		start = start + " 00:00"
	}
	idx = strings.Index(end, " ")
	if idx == -1 {
		end = end + " 00:00"
	}
	for i := 0; i < m.Rows; i++ {
		ts := m.DataRows[i].Key
		idx := strings.Index(ts, " ")
		if idx != -1 {
			ts = ts[:idx] + " 00:00"
		}
		if ts >= start && ts <= end {
			ret.DataRows = append(ret.DataRows, m.DataRows[i])
			ret.Rows++
		}
	}
	return ret
}

func (m *Matrix) StdDev(field, period int) int {
	ret := m.AddColumn()
	avg := SMA(m, period, field)
	for j := period; j < m.Rows; j++ {
		sum := 0.0
		cavg := m.DataRows[j].Get(avg)
		for i := 0; i < period; i++ {
			idx := j - i
			d := m.DataRows[idx].Get(field) - cavg
			sum += d * d
		}
		m.DataRows[j].Set(ret, ma.Sqrt(sum/float64(period)))
	}
	m.RemoveColumn()
	return ret
}

func FindMax(values []float64) float64 {
	if len(values) == 1 {
		return values[0]
	}
	max := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] > max {
			max = values[i]
		}
	}
	return max
}

func FindMin(values []float64) float64 {
	if len(values) == 1 {
		return values[0]
	}
	min := values[0]
	for i := 1; i < len(values); i++ {
		if values[i] < min {
			min = values[i]
		}
	}
	return min
}

func (m *Matrix) ChangePercentage(first, second, field int) float64 {
	return (m.DataRows[first].Get(field)/m.DataRows[second].Get(field) - 1.0) * 100.0
}

func (m *Matrix) FibonacciLevels(index int) *MajorLevels {
	// top bottom: 0 (Low) 23.6%, 38.2%, 50%, 61.8%, and 78.6%. 100 (High)
	levels := NewMajorLevels(0.0)
	cur := m.DataRows[index]
	hs := cur.Get(HIGH)
	ls := cur.Get(LOW)
	diff := hs - ls
	levels.Add(hs, 1, cur.Key)
	levels.Add(hs-diff*0.236, 1, cur.Key)
	levels.Add(hs-diff*0.382, 1, cur.Key)
	levels.Add(hs-diff*0.5, 1, cur.Key)
	levels.Add(hs-diff*0.618, 1, cur.Key)
	levels.Add(hs-diff*0.786, 1, cur.Key)
	levels.Add(ls, 1, cur.Key)
	return levels
}

func (m *Matrix) InverseFibonacciLevels(index int) *MajorLevels {
	// bottom to top: 0 (Low) 23.6%, 38.2%, 50%, 61.8%, and 78.6%. 100 (High)
	levels := NewMajorLevels(0.0)
	cur := m.DataRows[index]
	hs := cur.Get(HIGH)
	ls := cur.Get(LOW)
	diff := hs - ls
	levels.Add(ls, 1, cur.Key)
	levels.Add(ls+diff*0.236, 1, cur.Key)
	levels.Add(ls+diff*0.382, 1, cur.Key)
	levels.Add(ls+diff*0.5, 1, cur.Key)
	levels.Add(ls+diff*0.618, 1, cur.Key)
	levels.Add(ls+diff*0.786, 1, cur.Key)
	levels.Add(hs, 1, cur.Key)
	return levels
}

func (ma *Matrix) FindSupportResistance(lookback int, threshold float64) []float64 {
	sp := ma.FindSwingPoints()
	var points = make([]float64, 0, lookback)
	s := len(sp) - lookback
	if s < 0 {
		s = 0
	}
	for i := s; i < len(sp); i++ {
		d := 1
		for _, p := range points {
			cp := m.Abs(ChangePercentage(sp[i].Value, p))
			if cp < 5.0 {
				d = 0
			}
		}

		if d != 0 {
			points = append(points, sp[i].Value)
		}
	}
	return points
}
