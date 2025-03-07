package math

import (
	"fmt"
	"strings"
)

// BitSet
type BitSet struct {
	names []string
	n     int
	l     int
}

func NewBitSet(headers ...string) *BitSet {
	ret := &BitSet{}
	ret.names = append(ret.names, headers...)
	ret.l = len(headers)
	return ret
}

func (bs *BitSet) IsSet(pos int) bool {
	return (bs.n & (1 << pos)) != 0
}

func (bs *BitSet) AllSet() bool {
	cnt := 0
	for i := 0; i < bs.l; i++ {
		if bs.IsSet(i) {
			cnt++
		}
	}
	return cnt == bs.l
}

func (bs *BitSet) Clear(pos int) {
	bs.n = bs.n & ^(1 << pos)
}

func (bs *BitSet) ClearAll() {
	bs.n = 0
}

func (bs *BitSet) SetState(pos int, set bool) {
	if set {
		bs.n = bs.n | (1 << pos)
	} else {
		bs.n = bs.n & ^(1 << pos)
	}
}

func bitString(n, l int) string {
	sb := strings.Builder{}
	for i := l - 1; i >= 0; i-- {
		state := (n & (1 << i)) != 0
		if state {
			sb.WriteRune('1')
		} else {
			sb.WriteRune('0')
		}
	}
	return sb.String()
}

func (bs *BitSet) String() string {
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("%s - %d : ", bitString(bs.n, bs.l), bs.n))
	for i, n := range bs.names {
		if i > 0 {
			sb.WriteString(" | ")
		}
		if bs.IsSet(i) {
			sb.WriteString(fmt.Sprintf("%s = 1", n))
		} else {
			sb.WriteString(fmt.Sprintf("%s = 0", n))
		}
	}
	return sb.String()
}

// BitMap
type BitMap struct {
	names  []string
	l      int
	values []int
	keys   []string
}

func NewBitMap(headers ...string) *BitMap {
	ret := &BitMap{}
	ret.names = append(ret.names, headers...)
	ret.l = len(headers)
	return ret
}

func (bm *BitMap) NewLine(key string) int {
	bm.values = append(bm.values, 0)
	bm.keys = append(bm.keys, key)
	return len(bm.values) - 1
}

func (bm *BitMap) Lines() int {
	return len(bm.values)
}

func (bm *BitMap) Names() []string {
	return bm.names
}

func (bm *BitMap) NumValues() int {
	return bm.l
}

func (bm *BitMap) Key(line int) string {
	if line >= len(bm.values) {
		return ""
	}
	return bm.keys[line]
}

func (bm *BitMap) IsSet(line, pos int) bool {
	if line >= len(bm.values) {
		return false
	}
	return (bm.values[line] & (1 << pos)) != 0
}

func (bm *BitMap) AllSet(line int) bool {
	if line >= len(bm.values) {
		return false
	}
	cnt := 0
	for i := 0; i < bm.l; i++ {
		if bm.IsSet(line, i) {
			cnt++
		}
	}
	return cnt == bm.l
}

func (bm *BitMap) Clear(line, pos int) {
	if line < len(bm.values) {
		bm.values[line] = bm.values[line] & ^(1 << pos)
	}
}

func (bm *BitMap) ClearAll() {
	for i := 0; i < len(bm.values); i++ {
		bm.values[i] = 0
	}
}

func (bm *BitMap) SetState(line, pos int, set bool) {
	if line < len(bm.values) {
		if set {
			bm.values[line] = bm.values[line] | (1 << pos)
		} else {
			bm.values[line] = bm.values[line] & ^(1 << pos)
		}
	}
}

func (bm *BitMap) BitString(line int) string {
	return bitString(bm.values[line], bm.l)
}

func (bs *BitMap) String() string {
	sb := strings.Builder{}
	for i := 0; i < len(bs.values); i++ {
		sb.WriteString(fmt.Sprintf("%s %s - %d : ", bs.keys[i], bitString(bs.values[i], bs.l), bs.values[i]))
		for j, n := range bs.names {
			if j > 0 {
				sb.WriteString(" | ")
			}
			if bs.IsSet(i, j) {
				sb.WriteString(fmt.Sprintf("%s = 1", n))
			} else {
				sb.WriteString(fmt.Sprintf("%s = 0", n))
			}
		}
		sb.WriteRune('\n')
	}
	return sb.String()
}
