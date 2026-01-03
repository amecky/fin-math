package math

import (
	"fmt"
	"strings"
)

type Border struct {
	Size      int
	H_LINE    string
	V_LINE    string
	TL_CORNER string
	TR_CORNER string
	BR_CORNER string
	BL_CORNER string
	CROSS     string
	TOP_DEL   string
	BOT_DEL   string
	LEFT_DEL  string
	RIGHT_DEL string
}

var DefaultBorder = Border{
	Size:      1,
	H_LINE:    "│",
	V_LINE:    "─",
	TL_CORNER: "┌",
	TR_CORNER: "┐",
	BR_CORNER: "┘",
	BL_CORNER: "└",
	CROSS:     "┼",
	TOP_DEL:   "┬",
	BOT_DEL:   "┴",
	LEFT_DEL:  "├",
	RIGHT_DEL: "┤",
}

func AlignStrings(txt string, length int, align int) string {
	left := 0
	right := 0
	d := max(length-len(txt), 0)
	// left
	if align == 0 {
		right = d
	}
	// right
	if align == 2 {
		if d > 0 {
			left = d
		}
	}
	// center
	if align == 1 {
		if d > 0 {
			left = d / 2
			right = d - left
		}
	}
	return strings.Repeat(" ", left) + txt + strings.Repeat(" ", right)
}

type MatrixRenderer struct {
	m       *Matrix
	sizes   []int
	builder strings.Builder
}

func NewMatrixRenderer(m *Matrix) *MatrixRenderer {
	return &MatrixRenderer{
		m: m,
	}
}

func (mr *MatrixRenderer) addDelimier(left, mid, right string) {
	mr.builder.WriteString(left)
	for i, s := range mr.sizes {
		mr.builder.WriteString(strings.Repeat(DefaultBorder.V_LINE, s))
		if i < len(mr.sizes)-1 {
			mr.builder.WriteString(mid)
		}
	}
	mr.builder.WriteString(right)
	mr.builder.WriteString("\n")
}

func (mr *MatrixRenderer) addHeaders() {
	// headers
	for i, h := range mr.m.Headers {
		mr.builder.WriteString(DefaultBorder.H_LINE)
		mr.builder.WriteString(AlignStrings(" "+h+" ", mr.sizes[i], 1))
	}
	mr.builder.WriteString(DefaultBorder.H_LINE)
	mr.builder.WriteString(AlignStrings(" Comment ", mr.sizes[len(mr.sizes)-1], 1))
	mr.builder.WriteString(DefaultBorder.H_LINE)
	mr.builder.WriteString("\n")
}

func (mr *MatrixRenderer) addValues() {
	for j := range mr.m.Rows {
		//if j >= mr.m.Rows-20 {
		c := mr.m.DataRows[j]
		mr.builder.WriteString(DefaultBorder.H_LINE)
		mr.builder.WriteString(AlignStrings(" "+c.Key+" ", mr.sizes[0], 0))
		mr.builder.WriteString(DefaultBorder.H_LINE)
		for i := range len(c.Values) {
			mr.builder.WriteString(AlignStrings(fmt.Sprintf(" %.2f ", c.Get(i)), mr.sizes[i+1], 2))
			mr.builder.WriteString(DefaultBorder.H_LINE)
		}
		mr.builder.WriteString(AlignStrings(" "+c.Comment+" ", mr.sizes[len(mr.sizes)-1], 0))
		mr.builder.WriteString(DefaultBorder.H_LINE)
		mr.builder.WriteString("\n")
		//}
	}

}

func (mr *MatrixRenderer) calculateSizes() {
	mr.sizes = make([]int, len(mr.m.DataRows[0].Values)+2)
	for i, th := range mr.m.Headers {
		mr.sizes[i] = len(th) + 2
	}
	mr.sizes[len(mr.sizes)-1] = len("Comment") + 2
	mr.sizes[0] = len(mr.m.DataRows[0].Key)
	total := 0
	for _, s := range mr.sizes {
		total += s
	}
	for j := range mr.m.Rows {
		c := mr.m.DataRows[j]
		if len(c.Key)+2 > mr.sizes[0] {
			mr.sizes[0] = len(c.Key) + 2
		}
		for i := range len(c.Values) {
			cur := fmt.Sprintf("%.2f", c.Get(i))
			if len(cur)+2 > mr.sizes[i+1] {
				mr.sizes[i+1] = len(cur) + 2
			}
		}
	}
}

func (mr *MatrixRenderer) String() string {
	mr.calculateSizes()
	mr.addDelimier(DefaultBorder.TL_CORNER, DefaultBorder.TOP_DEL, DefaultBorder.TR_CORNER)
	mr.addHeaders()
	mr.addDelimier(DefaultBorder.LEFT_DEL, DefaultBorder.CROSS, DefaultBorder.RIGHT_DEL)
	mr.addValues()
	mr.addDelimier(DefaultBorder.BL_CORNER, DefaultBorder.BOT_DEL, DefaultBorder.BR_CORNER)
	return mr.builder.String()
}

func (m *Matrix) String() string {
	mr := NewMatrixRenderer(m)
	return mr.String()

}
