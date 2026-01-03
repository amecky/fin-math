
# Matrix

### NewMatrix(cols int) *Matrix 

### NewMatrixWithHeaders(cols int, headers []string) *Matrix 

### (m *Matrix) GetAsPriceData() []PriceData 

### (m *Matrix) GetPriceData(start, count int) []PriceData 
	
### (m *Matrix) ExtractDates() []DateIndex 

### (m *Matrix) AddRow(key string) *MatrixRow 

### (m *Matrix) AddNamedColumn(header string) int 

### (m *Matrix) AddColumn() int 

### (m *Matrix) Set(x, y int, value float64) *MatrixRow 

### (m *Matrix) ForEach(fn func(index int, row *MatrixRow)) 

### (m Matrix) FindRow(key string) *MatrixRow 

### (m Matrix) FindRowIndex(key string) int 

### (m Matrix) SearchRowIndex(key string) int 

### (m Matrix) Last() *MatrixRow 

### (m *Matrix) Row(index int) *MatrixRow 

### (m *Matrix) PartialSum(field, start, count int) float64 

### (m *Matrix) SetComment(row int, cmt string)

### (m *Matrix) Get(col, row int) float64 

### (m *Matrix) GetColumn(col int) []float64 

### (m *Matrix) BuildColumn(conv func(m *Matrix, index int) float64) []float64 

### (m *Matrix) GetIntColumn(col int) []int 

### (m *Matrix) GetCommentColumn() []string 

### (m *Matrix) GetKeys() []string 

### (m *Matrix) GetConvertedKeys(conv func(key string) string) []string 

### (m *Matrix) FindMinMaxIndexBetween(field, start, count int) (int, int) 

### (m *Matrix) FindMinMaxBetween(field, start, count int) (float64, float64) 

### (m *Matrix) FindHighLowIndex(start, count int) (int, int) {
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

### (m *Matrix) FindHighestIndex(index, count int) int {
	if m.Rows == 0 {
		return 0.0
	}
	end := index
	if end > m.Rows {
		end = m.Rows
	}
	start := end - count
	if start < 0 {
		start = 0
	}
	high := m.DataRows[end].Get(1)
	hi := 0
	for i := end - 1; i >= start; i-- {
		h := m.DataRows[i].Get(1)
		if h > high {
			high = h
			hi = i - end
		}
	}
	return hi
}

### (m *Matrix) FindLowestIndex(index, count int) int {
	if m.Rows == 0 {
		return 0.0
	}
	end := index
	if end > m.Rows {
		end = m.Rows
	}
	start := end - count
	if start < 0 {
		start = 0
	}
	low := m.DataRows[end].Get(2)
	hi := 0
	for i := end; i > start; i-- {
		h := m.DataRows[i].Get(2)
		if h < low {
			low = h
			hi = i - end
		}
	}
	return hi
}

### (m *Matrix) FindHighestHigh(index, count int) float64 {
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

### (m *Matrix) FindLowestLow(index, count int) float64 {
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

### (m *Matrix) FindHighestHighLowestLow(start, count int) (float64, float64) {
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

### (m *Matrix) FindMinMax(field, start, count int) (float64, float64) {
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
	low := m.DataRows[start].Get(field)
	high := m.DataRows[start].Get(field)
	for i := start; i < end; i++ {
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

### (m *Matrix) Shift(field, period int) int {
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

### (m *Matrix) FindMinBetween(field, start, count int) float64 {
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

### (m *Matrix) FindMaxBetween(field, start, count int) float64 {
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

### (m *Matrix) CrossUp(first, second, index int) bool {
	if index > 0 {
		f := m.DataRows[index-1]
		s := m.DataRows[index]
		if f.Get(first) < f.Get(second) && s.Get(first) > s.Get(second) {
			return true
		}
	}
	return false
}

### (m *Matrix) CrossDown(first, second, index int) bool {
	if index > 0 {
		f := m.DataRows[index-1]
		s := m.DataRows[index]
		if f.Get(first) > f.Get(second) && s.Get(first) < s.Get(second) {
			return true
		}
	}
	return false
}

### (m *Matrix) CrossOver(first, second int) int {
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

### (m *Matrix) Stochastic(days int, field int) int {
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

### (m *Matrix) Add(field int, other *Matrix, otherField int) int {
	ret := m.AddColumn()
	for _, s := range other.DataRows {
		row := m.FindRow(s.Key)
		if row != nil {
			row.Set(ret, row.Get(field)+s.Get(otherField))
		}
	}
	return ret
}

### (m *Matrix) Subtract(field, otherField int) int {
	ret := m.AddColumn()
	for i := 0; i < m.Rows; i++ {
		m.DataRows[i].Set(ret, m.DataRows[i].Get(field)-m.DataRows[i].Get(otherField))
	}
	return ret
}

### (m *Matrix) Sum(field, length int) int {
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

### (m *Matrix) SlopePercentage(field, lookback int) int {
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

### (m *Matrix) Sort(field int) {
	sort.Slice(m.DataRows, func(i, j int) bool {
		return m.DataRows[i].Get(field) > m.DataRows[j].Get(field)
	})
}

### (m *Matrix) SortByKey() {
	sort.Slice(m.DataRows, func(i, j int) bool {
		return m.DataRows[i].Key > m.DataRows[j].Key
	})
}

### (m *Matrix) SortReverse(field int) {
	sort.Slice(m.DataRows, func(i, j int) bool {
		return m.DataRows[i].Get(field) < m.DataRows[j].Get(field)
	})
}

### (m *Matrix) Sample(field int, steps []float64) int {
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

### (m *Matrix) Categorize(field, target int, check func(mr MatrixRow) float64) {
	for i, s := range m.DataRows {
		v := check(s)
		m.DataRows[i].Set(target, v)
	}
}

### (m *Matrix) Categorize2(target int, check func(mr MatrixRow) float64) {
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

### (m *Matrix) Evaluate(eval Evaluator) int {
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

### (m *Matrix) CompareAll(functions []CompareFn) int {
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

### (m *Matrix) Apply(fn func(mr MatrixRow) float64) int {
	ret := m.AddColumn()
	for i, s := range m.DataRows {
		v := fn(s)
		m.DataRows[i].Set(ret, v)
	}
	return ret
}

### (m *Matrix) ApplyCurrentPrevious(fn func(prev, cur MatrixRow) float64) int {
	ret := m.AddColumn()
	for i, s := range m.DataRows {
		if i > 0 {
			v := fn(m.DataRows[i-1], s)
			m.DataRows[i].Set(ret, v)
		}
	}
	return ret
}

### (m *Matrix) ApplyRow(row int, fn func(mr MatrixRow) float64) {
	for i, s := range m.DataRows {
		v := fn(s)
		m.DataRows[i].Set(row, v)
	}
}

### (m *Matrix) RemoveColumn() {
	for j := 0; j < m.Rows; j++ {
		m.DataRows[j].Values = m.DataRows[j].Values[:m.DataRows[j].Num-1]
		m.DataRows[j].Num--
	}
	m.Headers = m.Headers[:len(m.Headers)-1]
	m.Cols--
}

### (m *Matrix) RemoveColumns(cnt int) {
	for i := 0; i < cnt; i++ {
		m.RemoveColumn()
	}
}

### (m *Matrix) OPEN(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(OPEN)
	}
	return 0.0
}

### (m *Matrix) HIGH(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(HIGH)
	}
	return 0.0
}

### (m *Matrix) LOW(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(LOW)
	}
	return 0.0
}

### (m *Matrix) CLOSE(row int) float64 {
	if row >= 0 && row < m.Rows {
		return m.DataRows[row].Get(ADJ_CLOSE)
	}
	return 0.0
}

### (m *Matrix) CopyColumn(source, destination int) {
	for j := 0; j < m.Rows; j++ {
		m.DataRows[j].Values[destination] = m.DataRows[j].Values[source]
	}
}

### (m *Matrix) Sublist(start, count int) *Matrix {
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

### (m *Matrix) Recent(count int) *Matrix {
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

### (m *Matrix) Subset(start, end int) *Matrix {
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

### (m *Matrix) FilterByKeys(start, end string) *Matrix {
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

### (m *Matrix) Filter(compare func(m *Matrix, index int) bool) *Matrix {
	ret := NewMatrix(m.Cols)
	for i := 0; i < m.Rows; i++ {
		if compare(m, i) {
			ret.DataRows = append(ret.DataRows, m.DataRows[i])
			ret.Rows++
		}
	}
	return ret
}

### (m *Matrix) FilterByKeysStraight(start, end string) *Matrix {
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

### CalculateMean(numbers []float64) float64 {
	sum := 0.0
	for _, number := range numbers {
		sum += number
	}
	return sum / float64(len(numbers))
}

// CalculateStandardDeviation calculates the standard deviation of a slice of numbers.
### CalculateStandardDeviation(numbers []float64) float64 {
	mean := CalculateMean(numbers)
	var variance float64
	for _, number := range numbers {
		variance += math.Pow(number-mean, 2)
	}
	variance /= float64(len(numbers))
	return math.Sqrt(variance)
}

### (m *Matrix) StdDev(field, period int) int {
	ret := m.AddColumn()
	vol := m.GetColumn(field)
	for j := 0; j < len(vol)-period; j++ {
		window := vol[j : j+period]
		stdDev := CalculateStandardDeviation(window)
		m.DataRows[j+period].Set(ret, stdDev)
	}
	return ret
}

### FindMax(values []float64) float64 {
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

### FindMin(values []float64) float64 {
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

### (m *Matrix) ChangePercentage(first, second, field int) float64 {
	return (m.DataRows[first].Get(field)/m.DataRows[second].Get(field) - 1.0) * 100.0
}

### (m *Matrix) FibonacciLevels(index int) *MajorLevels {
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

### (m *Matrix) InverseFibonacciLevels(index int) *MajorLevels {
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

### (ma *Matrix) FindSupportResistance(lookback int, threshold float64) []float64 {
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

### SaveMatrix(m *Matrix, fileName string) error {
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

### readMatrixFile(path string) ([]string, error) {
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

### LoadMatrix(fileName string) (*Matrix, error) {
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

### NewValueMap() *ValueMap {
	ret := ValueMap{}
	ret.data = make([]ValueMapEntry, 0)
	ret.names = make([]string, 0)
	return &ret
}

### (vm *ValueMap) Get(key string, defaultValue float64) float64 {
	idx := vm.IndexOf(key)
	if idx != -1 {
		return vm.data[idx].Value
	}
	return defaultValue
}

### (vm *ValueMap) Set(key string, value float64) {
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

### (vm *ValueMap) Keys() []string {
	return vm.names
}

### (vm *ValueMap) Contains(key string) bool {
	for _, k := range vm.data {
		if k.Key == key {
			return true
		}
	}
	return false
}

### (vm *ValueMap) IndexOf(key string) int {
	for i, k := range vm.data {
		if k.Key == key {
			return i
		}
	}
	return -1
}
*/
