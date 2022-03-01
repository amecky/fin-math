package math

import (
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
	}
	return "-"
}

func TranslatePatternShort(v float64) string {
	switch v {
	case 1.0:
		return "HA"
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
		return "DJ"
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
	}
	return "-"
}

const (
	BULLISH = 1.0
	BEARISH = -1.0
)

func NearlyEquals(first, second float64) bool {
	return m.Abs(ChangePercentage(first, second)) < 0.1
}

func FindCandleStickPatterns(prices *Matrix) int {

	ret := prices.AddColumn()
	// use EMA13 to determine trend
	e1 := MASlope(prices, EMA, 20, 1)
	// 0 = BodySize 1 = BodyPos 2 = Mid 3 = RelBodySize 4 = RelAvg 5 = Upper 6 = Lower 7 = Trend 8 = Spread 9 = RelSpread
	ci := Candles(prices, 21)

	//
	// Single Bar patterns
	//
	for i := 0; i < prices.Rows; i++ {
		first := prices.DataRows[i]

		// Hammer
		ht := first.Get(HIGH) - 0.382*(first.Get(HIGH)-first.Get(LOW))
		if first.Get(OPEN) > ht && first.Get(ADJ_CLOSE) > ht {
			prices.DataRows[i].Set(ret, 1.0)
		}
		// DOJI
		if first.Get(ci+3) < 1.0 {
			prices.DataRows[i].Set(ret, 9.0)
		}
		// Shooting star
		st := first.Get(LOW) + 0.382*(first.Get(HIGH)-first.Get(LOW))
		if first.Get(OPEN) < st && first.Get(ADJ_CLOSE) < st {
			prices.DataRows[i].Set(ret, 2.0)
		}
		// Bullish Marubozu
		if first.Get(ci+7) == BULLISH && first.Get(ci+5) < 1.0 && first.Get(ci+6) < 1.0 {
			prices.DataRows[i].Set(ret, 12.0)
		}
		// Bearish Marubozu
		if first.Get(ci+7) == BEARISH && first.Get(ci+5) < 1.0 && first.Get(ci+6) < 1.0 {
			prices.DataRows[i].Set(ret, 13.0)
		}

		// Hanging Man
		if first.Get(ci+7) == BEARISH && first.Get(ci+6) >= 66.0 {
			prices.DataRows[i].Set(ret, 11.0)
		}
		// Inverted Hammer
		if first.Get(ci+7) == BULLISH && first.Get(ci+5) >= 66.0 {
			prices.DataRows[i].Set(ret, 10.0)
		}

	}

	// Double Bar patterns

	for i := 1; i < prices.Rows; i++ {
		first := prices.DataRows[i-1]
		second := prices.DataRows[i]

		if first.Get(ci+7) != second.Get(ci+7) {
			// Bullish Engulfing
			if second.Get(ci+7) == BULLISH && second.Get(ADJ_CLOSE) > first.Get(OPEN) && second.Get(OPEN) < first.Get(ADJ_CLOSE) {
				prices.DataRows[i].Set(ret, 3.0)
			}
			// Bearish Engulfing
			if second.Get(ci+7) == BEARISH && second.Get(ADJ_CLOSE) < first.Get(OPEN) && second.Get(OPEN) > first.Get(ADJ_CLOSE) {
				prices.DataRows[i].Set(ret, 4.0)
			}
		}

		// Tweezer Top
		if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && NearlyEquals(first.Get(HIGH), second.Get(HIGH)) == true && second.Get(e1) > 0.0 {
			prices.DataRows[i].Set(ret, 16.0)
		}
		// Tweezer Bottom
		if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && NearlyEquals(first.Get(LOW), second.Get(LOW)) == true && second.Get(e1) < 0.0 {
			prices.DataRows[i].Set(ret, 17.0)
		}

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
		if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && first.Get(2) > second.Get(OPEN) && first.Get(ci+2) < second.Get(ADJ_CLOSE) {
			prices.DataRows[i].Set(ret, 14.0)
		}
		// Dark Cloud Cover
		if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && first.Get(HIGH) < second.Get(OPEN) && first.Get(ci+2) > second.Get(ADJ_CLOSE) {
			prices.DataRows[i].Set(ret, 15.0)
		}

	}

	for i := 2; i < prices.Rows; i++ {
		first := prices.DataRows[i-2]
		second := prices.DataRows[i-1]
		third := prices.DataRows[i]

		// Three Bar Bullish reversal
		if first.Get(ci+7) == BEARISH &&
			second.Get(LOW) < first.Get(LOW) &&
			second.Get(LOW) < third.Get(LOW) &&
			third.Get(ADJ_CLOSE) > second.Get(HIGH) &&
			third.Get(ADJ_CLOSE) > first.Get(HIGH) &&
			second.Get(ADJ_CLOSE) < first.Get(OPEN) {
			// three bar reversal
			prices.DataRows[i].Set(ret, 18.0)
		}
		// Three Bar Bearish reversal
		if first.Get(ci+7) == BULLISH &&
			second.Get(HIGH) > first.Get(HIGH) &&
			second.Get(HIGH) > third.Get(HIGH) &&
			third.Get(ADJ_CLOSE) < second.Get(LOW) &&
			third.Get(ADJ_CLOSE) < first.Get(LOW) &&
			second.Get(ADJ_CLOSE) > first.Get(ADJ_CLOSE) {
			// three bar reversal
			prices.DataRows[i].Set(ret, 19.0)
		}

	}
	prices.RemoveColumns(10)
	return ret
}
