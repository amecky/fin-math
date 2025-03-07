package math

import (
	"math"
	m "math"
)

func TranslatePattern(v float64) string {
	switch v {
	case 1.0:
		return "Hammer"
	case 2.0:
		return "Shooting Star"
	case 3.0:
		return "Bullish Engulfing"
	case 4.0:
		return "Bearish Engulfing"
	case 5.0:
		return "Bullish Inside Bar"
	case 6.0:
		return "Bearish Inside Bar"
	case 7.0:
		return "Harami Bearish"
	case 8.0:
		return "Harami Bullish"
	case 9.0:
		return "Doji"
	case 10.0:
		return "Inverted Hammer"
	case 11.0:
		return "Hanging Man"
	case 12.0:
		return "Bullish Marubozu"
	case 13.0:
		return "Bearish Marubozu"
	case 14.0:
		return "Piercing Line"
	case 15.0:
		return "Dark Cloud Cover"
	case 16.0:
		return "Tweezer Top"
	case 17.0:
		return "Tweezer Bottom"
	case 18.0:
		return "Three Bar BU Reversal"
	case 19.0:
		return "Three Bar BE Reversal"
	case 20.0:
		return "Inside Bar"
	}
	return "-"
}

func TranslatePatternShort(v float64) string {
	switch v {
	case 1.0:
		return "HAM"
	case 2.0:
		return "SS"
	case 3.0:
		return "BUE"
	case 4.0:
		return "BEE"
	case 5.0:
		return "IBU"
	case 6.0:
		return "IBE"
	case 7.0:
		return "HBE"
	case 8.0:
		return "HBU"
	case 9.0:
		return "DJI"
	case 10.0:
		return "IHA"
	case 11.0:
		return "HM"
	case 12.0:
		return "BUM"
	case 13.0:
		return "BEM"
	case 14.0:
		return "PL"
	case 15.0:
		return "DC"
	case 16.0:
		return "TWT"
	case 17.0:
		return "TWB"
	case 18.0:
		return "3BUR"
	case 19.0:
		return "3BER"
	case 20.0:
		return "IB"
	}
	return "-"
}

const (
	BULLISH = 1.0
	BEARISH = -1.0
)

func NearlyEquals(first, second, threshold float64) bool {
	return m.Abs(ChangePercentage(first, second)) < threshold
}

func FindCandleStickPatterns(prices *Matrix) int {

	ret := prices.AddColumn()
	// use EMA13 to determine trend
	// 0 = BodySize 1 = BodyPos 2 = Mid 3 = RelBodySize 4 = RelAvg 5 = Upper 6 = Lower 7 = Trend 8 = Spread 9 = RelSpread
	//ci := Candles(prices, 21)
	//const (
	//	MID           = 2
	//	REL_BODY_SIZE = 3
	//	UPPER         = 5
	//	LOWER         = 6
	//	TREND         = 7
	//)
	//
	// Single Bar patterns
	//
	for i := 0; i < prices.Rows; i++ {
		cur := prices.DataRows[i]
		cd := ConvertMatrixRow(cur)
		//fmt.Println(first.Key, "=", cd, "RB", cd.Body/cd.Range)
		// Hammer
		ht := cur.Get(HIGH) - 0.382*cd.Range
		if cur.Get(OPEN) > ht && cur.Get(ADJ_CLOSE) > ht {
			cur.Set(ret, 1.0)
		}
		// DOJI
		if cd.Body/cd.Range < 0.05 {
			cur.Set(ret, 9.0)
		}
		// Shooting star
		st := cur.Get(LOW) + 0.382*(cur.Get(HIGH)-cur.Get(LOW))
		if cur.Get(OPEN) < st && cur.Get(ADJ_CLOSE) < st {
			cur.Set(ret, 2.0)
		}
		// Bullish Marubozu
		if cd.IsBullish() && cd.Upper/cd.Range < 0.01 && cd.Lower/cd.Range < 0.01 {
			cur.Set(ret, 12.0)
		}
		// Bearish Marubozu
		if cd.IsBearish() && cd.Upper/cd.Range < 0.01 && cd.Lower/cd.Range < 0.01 {
			cur.Set(ret, 13.0)
		}

		// Hanging Man
		if cd.IsBearish() && cd.Lower/cd.Range >= 0.75 {
			cur.Set(ret, 11.0)
		}
		// Inverted Hammer
		if cd.IsBullish() && cd.Upper/cd.Range >= 0.75 {
			cur.Set(ret, 10.0)
		}
		if i > 0 {
			prev := prices.DataRows[i-1]
			// Inside Bar
			if cur.High() <= prev.High() && cur.Low() >= prev.Low() {
				prices.DataRows[i].Set(ret, 20.0)
			}
			// Bullish Engulfing
			if cur.IsGreen() && prev.IsRed() && cur.Get(ADJ_CLOSE) > prev.Get(OPEN) && cur.Open() < prev.Close() {
				prices.DataRows[i].Set(ret, 3.0)
			}
			if cur.IsRed() && prev.IsGreen() && cur.Get(ADJ_CLOSE) < prev.Get(OPEN) && cur.Get(OPEN) > prev.Get(ADJ_CLOSE) {
				prices.DataRows[i].Set(ret, 4.0)
			}
			pd := ConvertMatrixRow(prev)
			// Tweezer Top
			if cd.IsBullish() && pd.IsBearish() && NearlyEquals(cur.Get(HIGH), prev.Get(HIGH), 0.01) {
				prices.DataRows[i].Set(ret, 16.0)
			}
			// Tweezer Bottom
			if cd.IsBearish() && pd.IsBullish() && NearlyEquals(cur.Get(LOW), prev.Get(LOW), 0.01) {
				prices.DataRows[i].Set(ret, 17.0)
			}

			/*
				// Bullish Inside Bar
				if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && first.Get(HIGH) > second.Get(HIGH) && first.Get(2) < second.Get(2) {
					prices.DataRows[i].Set(ret, 5.0)
				}
				// Bearish Inside Bar
				if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && first.Get(HIGH) > second.Get(HIGH) && first.Get(2) < second.Get(2) {
					prices.DataRows[i].Set(ret, 6.0)
				}
				// Harami Bearish
				if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && first.Get(ADJ_CLOSE) > second.Get(OPEN) && first.Get(OPEN) < second.Get(ADJ_CLOSE) {
					prices.DataRows[i].Set(ret, 7.0)
				}
				// Harami Bullish
				if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && first.Get(OPEN) > second.Get(ADJ_CLOSE) && first.Get(ADJ_CLOSE) < second.Get(OPEN) {
					prices.DataRows[i].Set(ret, 8.0)
				}
				// Piercing line
				if first.Get(ci+TREND) == BEARISH && second.Get(ci+TREND) == BULLISH && first.Get(ADJ_CLOSE) < second.Get(OPEN) && first.Get(ci+MID) < second.Get(ADJ_CLOSE) {
					prices.DataRows[i].Set(ret, 14.0)
				}
				// Dark Cloud Cover
				if first.Get(ci+TREND) == BULLISH && second.Get(ci+TREND) == BEARISH && first.Get(ADJ_CLOSE) < second.Get(OPEN) && first.Get(ci+MID) > second.Get(ADJ_CLOSE) {
					prices.DataRows[i].Set(ret, 15.0)
				}
			*/
		}
	}
	/*
		// Double Bar patterns

		for i := 1; i < prices.Rows; i++ {
			first := prices.DataRows[i-1]
			second := prices.DataRows[i]
			fd := ConvertMatrixRow(&first)
			sd := ConvertMatrixRow(&second)



		}
	*/
	for i := 2; i < prices.Rows; i++ {
		first := prices.DataRows[i-2]
		second := prices.DataRows[i-1]
		third := prices.DataRows[i]

		// Three Bar Bullish reversal
		if first.IsRed() &&
			second.Get(LOW) < first.Get(LOW) &&
			second.Get(LOW) < third.Get(LOW) &&
			third.Get(ADJ_CLOSE) > second.Get(HIGH) &&
			third.Get(ADJ_CLOSE) > first.Get(HIGH) &&
			second.Get(ADJ_CLOSE) < first.Get(OPEN) {
			// three bar reversal
			prices.DataRows[i].Set(ret, 18.0)
		}
		// Three Bar Bearish reversal
		if first.IsGreen() &&
			second.Get(HIGH) > first.Get(HIGH) &&
			second.Get(HIGH) > third.Get(HIGH) &&
			third.Get(ADJ_CLOSE) < second.Get(LOW) &&
			third.Get(ADJ_CLOSE) < first.Get(LOW) &&
			second.Get(ADJ_CLOSE) > first.Get(ADJ_CLOSE) {
			// three bar reversal
			prices.DataRows[i].Set(ret, 19.0)
		}

	}
	prices.RemoveColumns(11)

	return ret
}

func FindSingleBarCandleStickPattern(row *MatrixRow) string {
	//
	// Single Bar patterns
	//
	top := math.Max(row.Get(OPEN), row.Get(ADJ_CLOSE))
	bottom := math.Min(row.Get(OPEN), row.Get(ADJ_CLOSE))
	uwAbs := row.Get(HIGH) - top
	lwAbs := bottom - row.Get(LOW)

	td := row.Get(HIGH) - row.Get(LOW)

	uw := uwAbs / td * 100.0
	lw := lwAbs / td * 100.0

	if uw > 45.0 && lw > 45.0 {
		return "DJI"
	}
	if lw >= 66.0 {
		return "HAM"
	}
	if uw > 66.0 {
		return "SHT"
	}
	return ""
}
