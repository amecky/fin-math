package math

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	m "math"
	"os"
	"sort"
	"strconv"
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

type DateIndex struct {
	Date  string
	Index int
}
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
	Trend     int
	Delta     float64
	Broken    bool
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

func (sps SwingPoints) Contains(index int) (bool, *SwingPoint) {
	for i := range sps {
		s := sps[i]
		if s.Index == index {
			return true, &sps[i]
		}
	}
	return false, nil
}

func (sps SwingPoints) FindTwoHH(start int) int {
	ht := 0
	for i := start; i < len(sps); i++ {
		cur := sps[i]
		if cur.BaseType == High {
			if cur.Type == HigherHigh {
				ht++
			} else {
				ht = 0
			}
		}
		if ht == 2 {
			return i
		}
	}
	return -1
}

func (sps SwingPoints) FindTwo(start int, base, tp SwingPointType) int {
	ht := 0
	for i := start; i < len(sps); i++ {
		cur := sps[i]
		if cur.BaseType == base {
			if cur.Type == tp {
				ht++
			} else {
				ht = 0
			}
		}
		if ht == 2 {
			return i
		}
	}
	return -1
}

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

// PriceData represents a single data point in the market
type PriceData struct {
	Timestamp string
	Open      float64
	High      float64
	Low       float64
	Close     float64
}

func (m *Matrix) GetAsPriceData() []PriceData {
	ret := make([]PriceData, m.Rows)
	for i := 0; i < m.Rows; i++ {
		c := m.DataRows[i]
		ret[i] = PriceData{
			Timestamp: c.Key,
			Open:      c.Open(),
			High:      c.High(),
			Low:       c.Low(),
			Close:     c.Close(),
		}
	}
	return ret
}

func (m *Matrix) GetPriceData(start, count int) []PriceData {
	end := start + count
	if end > m.Rows {
		end = m.Rows
	}
	ret := make([]PriceData, end-start)
	for i := start; i < end; i++ {
		c := m.DataRows[i]
		ret[i-start] = PriceData{
			Timestamp: c.Key,
			Open:      c.Open(),
			High:      c.High(),
			Low:       c.Low(),
			Close:     c.Close(),
		}
	}
	return ret
}
func (m *Matrix) ExtractDates() []DateIndex {
	ret := make([]DateIndex, 0)
	mapping := make(map[string]int)
	for i := 0; i < m.Rows; i++ {
		dt := m.DataRows[i].Key
		dt = dt[0:strings.Index(dt, " ")]
		if _, ok := mapping[dt]; !ok {
			mapping[dt] = i
		}
	}
	for v, k := range mapping {
		ret = append(ret, DateIndex{
			Date:  v,
			Index: k,
		})
	}
	sort.Slice(ret, func(i, j int) bool {
		return ret[i].Date < ret[j].Date
	})
	return ret
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

func (m *Matrix) ForcedAddRow(key string) *MatrixRow {
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

func (m *Matrix) SetHeader(index int, header string) {
	if index >= 0 && index < len(m.Headers) {
		m.Headers[index+1] = header
	}
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

func (m *Matrix) Set(x, y int, value float64) *MatrixRow {
	mr := m.Row(y)
	mr.Set(x, value)
	return mr
}

func (m *Matrix) ForEach(fn func(index int, row *MatrixRow)) {
	for i := range m.Rows {
		fn(i, m.Row(i))
	}
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

func (m Matrix) SearchRowIndex(key string) int {
	ln := len(key)
	for i, r := range m.DataRows {
		cur := r.Key
		if len(cur) > ln {
			cur = cur[0:ln]
		}
		if cur == key {
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

func (m *Matrix) Row(index int) *MatrixRow {
	if index < 0 {
		index = m.Rows + index
	}
	if m.Rows > index {
		return &m.DataRows[index]
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

func (m *Matrix) SetComment(row int, cmt string) {
	if m.Rows > row {
		m.DataRows[row].Comment = cmt
	}
}

func (m *Matrix) Get(col, row int) float64 {
	if row >= 0 && row < m.Rows && col >= 0 && col < m.Cols {
		return m.DataRows[row].Get(col)
	}
	return 0.0
}

func (m *Matrix) GetColumn(col int) []float64 {
	ret := make([]float64, 0)
	for i := 0; i < m.Rows; i++ {
		ret = append(ret, m.DataRows[i].Get(col))
	}
	return ret
}

func (m *Matrix) GetPartialColumn(col, start, end int) []float64 {
	ret := make([]float64, end-start)
	for i := start; i < end; i++ {
		ret[i-start] = m.DataRows[i].Get(col)
	}
	return ret
}

func (m *Matrix) BuildColumn(conv func(m *Matrix, index int) float64) []float64 {
	ret := make([]float64, 0)
	for i := 0; i < m.Rows; i++ {
		ret = append(ret, conv(m, i))
	}
	return ret
}

func (m *Matrix) GetIntColumn(col int) []int {
	ret := make([]int, 0)
	for i := 0; i < m.Rows; i++ {
		ret = append(ret, int(m.DataRows[i].Get(col)))
	}
	return ret
}

func (m *Matrix) GetCommentColumn() []string {
	ret := make([]string, m.Rows)
	for i := 0; i < m.Rows; i++ {
		ret[i] = m.DataRows[i].Comment
	}
	return ret
}

func (m *Matrix) GetKeys() []string {
	ret := make([]string, 0)
	for _, c := range m.DataRows {
		ret = append(ret, c.Key)
	}
	return ret
}

func (m *Matrix) GetConvertedKeys(conv func(key string) string) []string {
	ret := make([]string, 0)
	for _, c := range m.DataRows {
		ret = append(ret, conv(c.Key))
	}
	return ret
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
	if end > m.Rows {
		end = m.Rows
	}
	low := m.DataRows[start].Get(2)
	high := m.DataRows[start].Get(1)
	hIdx := start
	lIdx := start
	for i := start + 1; i < end; i++ {
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

func (m *Matrix) FindHighestIndex(index, count int) int {
	if m.Rows == 0 {
		return 0.0
	}
	end := index + count
	if end > m.Rows {
		end = m.Rows
	}
	start := index
	hi := start
	high := m.DataRows[hi].High()
	for i := start + 1; i < end; i++ {
		h := m.DataRows[i].High()
		if h > high {
			high = h
			hi = i
		}
	}
	return hi
}

func (m *Matrix) FindLowestIndex(index, count int) int {
	if m.Rows == 0 {
		return 0.0
	}
	end := index + count
	if end > m.Rows {
		end = m.Rows
	}
	start := index
	hi := start
	low := m.DataRows[hi].Low()
	for i := start + 1; i < end; i++ {
		h := m.DataRows[i].Low()
		if h < low {
			low = h
			hi = i
		}
	}
	return hi
}

func (m *Matrix) FindHighestHigh(index, count int) float64 {
	if m.Rows == 0 {
		return 0.0
	}
	start := index - count
	if start < 0 {
		start = 0
	}
	end := index
	if end > m.Rows {
		end = m.Rows
	}
	high := m.DataRows[start].Get(1)
	for i := start + 1; i < end; i++ {
		h := m.DataRows[i].Get(1)
		if h >= high {
			high = h
		}
	}
	return high
}

func (m *Matrix) FindLowestLow(index, count int) float64 {
	if m.Rows == 0 {
		return 0.0
	}
	start := index - count
	if start < 0 {
		start = 0
	}
	end := index
	if end > m.Rows {
		end = m.Rows
	}
	low := m.DataRows[start].Get(2)
	for i := start + 1; i < end; i++ {
		l := m.DataRows[i].Get(2)
		if l < low {
			low = l
		}
	}
	return low
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

func (m *Matrix) FindMinMax(field, start, count int) (float64, float64) {
	if m.Rows == 0 {
		return 0.0, 0.0
	}
	if start < 0 {
		start = 0
	}
	end := min(start+count, m.Rows)
	low := m.DataRows[start].Get(field)
	high := m.DataRows[start].Get(field)
	for i := start + 1; i < end; i++ {
		c := m.DataRows[i].Get(field)
		if c >= high {
			high = c
		}
		if c <= low {
			low = c
		}
	}
	return low, high
}

func (m *Matrix) Shift(field, period int) int {
	ret := m.AddColumn()
	if period > 0 {
		for i := m.Rows - 1; i >= period; i-- {
			idx := i - period
			m.DataRows[i].Set(ret, m.DataRows[idx].Get(field))
		}
		for i := 0; i < period; i++ {
			m.DataRows[i].Set(ret, 0.0)
		}
	} else if period < 0 {
		period *= -1
		for i := period; i < m.Rows; i++ {
			idx := i - period
			m.DataRows[idx].Set(ret, m.DataRows[i].Get(field))
		}
		for i := m.Rows - period; i < m.Rows; i++ {
			m.DataRows[i].Set(ret, 0.0)
		}
	}
	return ret
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

func (m *Matrix) CrossUp(first, second, index int) bool {
	if index > 0 {
		f := m.DataRows[index-1]
		s := m.DataRows[index]
		if f.Get(first) < f.Get(second) && s.Get(first) > s.Get(second) {
			return true
		}
	}
	return false
}

func (m *Matrix) CrossDown(first, second, index int) bool {
	if index > 0 {
		f := m.DataRows[index-1]
		s := m.DataRows[index]
		if f.Get(first) > f.Get(second) && s.Get(first) < s.Get(second) {
			return true
		}
	}
	return false
}

func (m *Matrix) CrossOver(first, second int) int {
	ret := m.AddColumn()
	for index := 1; index < m.Rows; index++ {
		f := m.DataRows[index-1]
		s := &m.DataRows[index]
		if f.Get(first) > f.Get(second) && s.Get(first) < s.Get(second) {
			s.Set(ret, -1.0)
		}
		if f.Get(first) < f.Get(second) && s.Get(first) > s.Get(second) {
			s.Set(ret, 1.0)
		}
	}
	return ret
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

func (m *Matrix) Sort(field int) {
	sort.Slice(m.DataRows, func(i, j int) bool {
		return m.DataRows[i].Get(field) > m.DataRows[j].Get(field)
	})
}

func (m *Matrix) SortByKey() {
	sort.Slice(m.DataRows, func(i, j int) bool {
		return m.DataRows[i].Key > m.DataRows[j].Key
	})
}

func (m *Matrix) SortReverse(field int) {
	sort.Slice(m.DataRows, func(i, j int) bool {
		return m.DataRows[i].Get(field) < m.DataRows[j].Get(field)
	})
}

func (m *Matrix) Sample(field int, steps []float64) int {
	ret := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		for j, s := range steps {
			m.DataRows[i].Set(ret, float64(len(steps)+1))
			if m.DataRows[i].Get(field) <= s {
				m.DataRows[i].Set(ret, float64(j)+1.0)
				break
			}
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

type Evaluator struct {
	Name  string
	Start int
	Run   func(m *Matrix, index int) bool
}

func (m *Matrix) Evaluate(eval Evaluator) int {
	ret := m.AddNamedColumn(eval.Name)
	for i := eval.Start; i < m.Rows; i++ {
		if i >= eval.Start {
			if eval.Run(m, i) {
				m.DataRows[i].Set(ret, 1.0)
			}
		}
	}
	return ret
}

type CompareFn func(mat *Matrix, index int) bool

func (m *Matrix) CompareAll(functions []CompareFn) int {
	ret := m.AddColumn()
	total := len(functions)
	for i := 0; i < m.Rows; i++ {
		cnt := 0
		for _, f := range functions {
			if f(m, i) {
				cnt++
			}
		}
		if cnt == total {
			m.DataRows[i].Set(ret, 1.0)
		} else {
			m.DataRows[i].Set(ret, -1.0)
		}
	}
	return ret
}

func (m *Matrix) Apply(fn func(mr MatrixRow) float64) int {
	ret := m.AddColumn()
	for i, s := range m.DataRows {
		v := fn(s)
		m.DataRows[i].Set(ret, v)
	}
	return ret
}

func (m *Matrix) ApplyCurrentPrevious(fn func(prev, cur MatrixRow) float64) int {
	ret := m.AddColumn()
	for i, s := range m.DataRows {
		if i > 0 {
			v := fn(m.DataRows[i-1], s)
			m.DataRows[i].Set(ret, v)
		}
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

func (m *Matrix) OPEN(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(OPEN)
	}
	return 0.0
}

func (m *Matrix) HIGH(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(HIGH)
	}
	return 0.0
}

func (m *Matrix) LOW(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(LOW)
	}
	return 0.0
}

func (m *Matrix) CLOSE(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(ADJ_CLOSE)
	}
	return 0.0
}

func (m *Matrix) CopyColumn(source, destination int) {
	for j := 0; j < m.Rows; j++ {
		m.DataRows[j].Values[destination] = m.DataRows[j].Values[source]
	}
}

func (m *Matrix) Sublist(start, count int) *Matrix {
	ret := NewMatrix(m.Cols)
	ret.Headers = m.Headers
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

func (m *Matrix) Recent(count int) *Matrix {
	ret := NewMatrixWithHeaders(m.Cols, m.Headers)
	start := m.Rows - count
	if start < 0 {
		start = 0
	}
	for i := start; i < m.Rows; i++ {
		ret.DataRows = append(ret.DataRows, m.DataRows[i])
		ret.Rows++
	}
	return ret
}

func (m *Matrix) Subset(start, end int) *Matrix {
	ret := NewMatrixWithHeaders(m.Cols, m.Headers)
	if start < 0 {
		start = 0
	}
	if end > m.Rows {
		end = m.Rows
	}
	for i := start; i < end; i++ {
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

func (m *Matrix) Filter(compare func(m *Matrix, index int) bool) *Matrix {
	ret := NewMatrix(m.Cols)
	for i := 0; i < m.Rows; i++ {
		if compare(m, i) {
			ret.DataRows = append(ret.DataRows, m.DataRows[i])
			ret.Rows++
		}
	}
	return ret
}

func (m *Matrix) Copy() *Matrix {
	ret := NewMatrix(m.Cols)
	for i := 0; i < m.Rows; i++ {
		ret.DataRows = append(ret.DataRows, m.DataRows[i])
		ret.Rows++
	}
	return ret
}

func (m *Matrix) FilterByKeysStraight(start, end string) *Matrix {
	ret := NewMatrix(m.Cols)
	for i := 0; i < m.Rows; i++ {
		ts := m.DataRows[i].Key
		if ts >= start && ts <= end {
			ret.DataRows = append(ret.DataRows, m.DataRows[i])
			ret.Rows++
		}
	}
	return ret
}

func CalculateMean(numbers []float64) float64 {
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}

// CalculateStandardDeviation calculates the standard deviation of a slice of numbers.
func CalculateStandardDeviation(numbers []float64) float64 {
	mean := CalculateMean(numbers)
	var variance float64
	for _, number := range numbers {
		variance += math.Pow(number-mean, 2)
	}
	variance /= float64(len(numbers))
	return math.Sqrt(variance)
}

func (m *Matrix) StdDev(field, period int) int {
	ret := m.AddColumn()
	vol := m.GetColumn(field)
	for j := 0; j < len(vol)-period; j++ {
		window := vol[j : j+period]
		stdDev := CalculateStandardDeviation(window)
		m.DataRows[j+period].Set(ret, stdDev)
	}
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

func SaveMatrix(m *Matrix, fileName string) error {
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	// first line are headers
	for i, h := range m.Headers {
		if i != 0 {
			f.WriteString(";")
		}
		f.WriteString(h)
	}
	f.WriteString("\n")
	for _, p := range m.DataRows {
		_, err2 := f.WriteString(p.Key)
		if err2 != nil {
			fmt.Println(err2)
			return err2
		}
		for i := 0; i < p.Num; i++ {
			_, err2 = f.WriteString(fmt.Sprintf(";%.2f", p.Get(i)))
			if err2 != nil {
				fmt.Println(err2)
				return err2
			}
		}
		f.WriteString("\n")

	}
	return nil
}

func readMatrixFile(path string) ([]string, error) {
	inFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer inFile.Close()
	r := bufio.NewReader(inFile)
	bytes := []byte{}
	lines := []string{}
	for {
		line, isPrefix, err := r.ReadLine()
		if err != nil {
			break
		}
		bytes = append(bytes, line...)
		if !isPrefix {
			str := strings.TrimSpace(string(bytes))
			if len(str) > 0 {
				lines = append(lines, str)
				bytes = []byte{}
			}
		}
	}
	if len(bytes) > 0 {
		lines = append(lines, string(bytes))
	}
	return lines, nil
}

func LoadMatrix(fileName string) (*Matrix, error) {
	lines, err := readMatrixFile(fileName)
	if err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return nil, errors.New("empty file")
	}
	headerLine := lines[0]
	headers := strings.Split(headerLine, ";")
	max := len(headers)
	m := NewMatrixWithHeaders(len(headers), headers)
	for i, str := range lines {
		if i > 0 {
			entries := strings.Split(str, ";")
			cnt := len(entries)
			if cnt > max {
				cnt = max
			}
			r := m.AddRow(entries[0])
			for i := 0; i < cnt-1; i++ {
				open, _ := strconv.ParseFloat(entries[i+1], 64)
				r.Set(i, open)
			}
		}
	}
	return m, nil
}

/*
type ValueMapEntry struct {
	Key   string
	Value float64
}
type ValueMap struct {
	data  []ValueMapEntry
	names []string
}

func NewValueMap() *ValueMap {
	ret := ValueMap{}
	ret.data = make([]ValueMapEntry, 0)
	ret.names = make([]string, 0)
	return &ret
}

func (vm *ValueMap) Get(key string, defaultValue float64) float64 {
	idx := vm.IndexOf(key)
	if idx != -1 {
		return vm.data[idx].Value
	}
	return defaultValue
}

func (vm *ValueMap) Set(key string, value float64) {
	idx := vm.IndexOf(key)
	if idx == -1 {
		vm.data = append(vm.data, ValueMapEntry{
			Key:   key,
			Value: value,
		})
		vm.names = append(vm.names, key)
	} else {
		vm.data[idx].Value = value
	}
}

func (vm *ValueMap) Keys() []string {
	return vm.names
}

func (vm *ValueMap) Contains(key string) bool {
	for _, k := range vm.data {
		if k.Key == key {
			return true
		}
	}
	return false
}

func (vm *ValueMap) IndexOf(key string) int {
	for i, k := range vm.data {
		if k.Key == key {
			return i
		}
	}
	return -1
}
*/
