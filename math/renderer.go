package math

import "fmt"

type IndicatorValueRenderer interface {
	Convert(v float64) (string, int)
}

type DefaultRenderer struct{}

func (r *DefaultRenderer) Convert(v float64) (string, int) {
	return fmt.Sprintf("%.2f", v), 0
}

type LowerUpperThresholdRenderer struct {
	Lower float64
	Upper float64
}

func (r *LowerUpperThresholdRenderer) Convert(v float64) (string, int) {
	tr := 0
	if v >= r.Upper {
		tr = 1
	}
	if v <= r.Lower {
		tr = -1
	}
	return fmt.Sprintf("%.2f", v), tr
}

type PercentageRenderer struct{}

func (r *PercentageRenderer) Convert(v float64) (string, int) {
	tr := 0
	if v > 0.0 {
		tr = 1
	}
	if v < 0.0 {
		tr = -1
	}
	return fmt.Sprintf("%.2f%%", v*100.0), tr
}

type PercentageRangeRenderer struct{}

func (r *PercentageRangeRenderer) Convert(v float64) (string, int) {
	tr := 0
	if v >= 0.8 {
		tr = 6
	} else if v >= 0.6 {
		tr = 5
	} else if v >= 0.4 {
		tr = 4
	} else if v >= 0.2 {
		tr = 3
	} else {
		tr = 2
	}
	return fmt.Sprintf("%.2f%%", v*100.0), tr
}

type UpDownRenderer struct{}

func (r *UpDownRenderer) Convert(v float64) (string, int) {
	if v < 0.0 {
		return "■", -1
	} else {
		return "■", 1
	}
}
