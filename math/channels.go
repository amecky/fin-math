package math

import (
	"time"
)

type TrendLine struct {
	Start string
	End   string
	Value float64
	Angle float64
}

type TrendChannel struct {
	Trend int
	Upper TrendLine
	Lower TrendLine
}

func CalculateDaysBetween(first, second string) (int, error) {
	cnt := 0
	ind := 1
	start := first
	end := second
	if start > end {
		start = second
		end = first
		ind = -1
	}
	running := true
	ed, err := time.Parse("2006-01-02 15:04", start)
	if err != nil {
		return 0, err
	}
	for running {
		ed = ed.AddDate(0, 0, 1)
		wd := int(ed.Weekday())
		if wd != 0 && wd != 6 {
			cmp := ed.Format("2006-01-02") + " 00:00"
			if cmp > end {
				running = false
			}
			cnt++
		}
	}
	return ind * cnt, nil
}

func CalculateTrendChannel(prices *Matrix) int {
	// 0 = Upper 1 = Lower
	ui := prices.AddColumn()
	li := prices.AddColumn()
	points := prices.FindSwingPoints()
	highs := points.FilterByType(High)
	hc := len(highs)
	if hc > 1 {
		for i := 1; i < hc; i++ {
			ch := highs[i]
			ph := highs[i-1]
			m := highs.GetAngle(i-1, i)
			for j := ph.Index; j < ch.Index; j++ {
				s := j - ph.Index
				prices.DataRows[j].Set(ui, ph.Value+float64(s)*m)
			}
		}
		start := highs[hc-2].Index
		h1 := highs[hc-2]
		m := highs.GetAngle(hc-2, hc-1)
		for i := start; i < prices.Rows; i++ {
			s := i - start
			prices.DataRows[i].Set(ui, h1.Value+float64(s)*m)
		}

		lows := points.FilterByType(Low)
		lc := len(lows)
		for i := 1; i < lc; i++ {
			ch := lows[i]
			ph := lows[i-1]
			m := lows.GetAngle(i-1, i)
			for j := ph.Index; j < ch.Index; j++ {
				s := j - ph.Index
				prices.DataRows[j].Set(li, ph.Value+float64(s)*m)
			}
		}

		start = lows[lc-2].Index
		l1 := lows[lc-2]
		m = lows.GetAngle(lc-2, lc-1)
		for i := start; i < prices.Rows; i++ {
			s := i - start
			prices.DataRows[i].Set(li, l1.Value+float64(s)*m)
		}
	}
	return ui
}

/*
type MarketCycle struct {
	Start string
	End   string
	Days  int
	Trend int
	Type  int
}

type MarketCycles []MarketCycle

func (m MarketCycles) Contains(date string) bool {
	for _, cm := range m {
		if cm.Start <= date && cm.End >= date {
			return true
		}
	}
	return false
}

func FindContractionPhases(prices *model.DataFrame, threshold float64) []MarketCycle {
	var ret MarketCycles
	nh := prices.Normalize("High", "NH", 50)
	nl := prices.Normalize("Low", "NL", 50)
	for i, p := range nh.DataRows {
		cur := p.Get("NH")
		cl := nl.FindByTimestamp(p.Timestamp).Get("NL")
		hc := 0
		for j := i; j < nh.Rows; j++ {
			c := nh.DataRows[j].Get("NH")
			if m.Abs(c-cur) > threshold {
				break
			}
			hc++
		}
		lc := 0
		ed := ""
		for j := i; j < nl.Rows; j++ {
			c := nl.DataRows[j].Get("NL")
			if m.Abs(c-cl) > threshold {
				break
			}
			ed = nl.DataRows[j].Timestamp
			lc++
		}
		if hc > 2 && lc > 2 {
			sd := p.Timestamp
			if ret.Contains(sd) == false {
				mc := MarketCycle{
					Start: p.Timestamp,
					End:   ed,
					Days:  hc,
					Type:  1,
					Trend: 0,
				}
				ret = append(ret, mc)
			}
		}
	}
	return ret
}
*/
