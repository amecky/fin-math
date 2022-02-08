package math

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
	}
	return "-"
}

const (
	BULLISH = 1.0
	BEARISH = -1.0
)

func FindCandleStickPatterns(prices *Matrix) int {

	ret := prices.AddColumn()
	// use EMA13 to determine trend
	//e1 := EMA(prices, 13, 4)
	// 0 = BodySize 1 = BodyPos 2 = Mid 3 = RelBodySize 4 = RelAvg 5 = Upper 6 = Lower 7 = Trend 8 = Spread 9 = RelSpread
	ci := Candles(prices, 21)

	//
	// Single Bar patterns
	//
	for i := 0; i < prices.Rows; i++ {
		first := prices.DataRows[i]
		// Hammer
		if first.Get(ci+7) == BULLISH && first.Get(ci+6) >= 66.0 {
			prices.DataRows[i].Set(ret, 1.0)
		}
		// DOJI
		if first.Get(ci+3) < 1.0 {
			prices.DataRows[i].Set(ret, 9.0)
		}
		// Hanging Man
		if first.Get(ci+7) == BEARISH && first.Get(ci+6) >= 66.0 {
			prices.DataRows[i].Set(ret, 11.0)
		}
		// Inverted Hammer
		if first.Get(ci+7) == BULLISH && first.Get(ci+5) >= 66.0 {
			prices.DataRows[i].Set(ret, 10.0)
		}
		// Shooting star
		if first.Get(ci+7) == BEARISH && first.Get(ci+5) >= 66.0 {
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
	}

	for i := 1; i < prices.Rows; i++ {
		first := prices.DataRows[i-1]
		second := prices.DataRows[i]
		if first.Get(ci+7) != second.Get(ci+7) {
			// Bullish Engulfing
			if second.Get(ci+7) == BULLISH {
				if second.Get(4) > first.Get(0) && second.Get(0) < first.Get(4) {
					prices.DataRows[i].Set(ret, 3.0)
				}
			}
			// Bearish Engulfing
			if second.Get(ci+7) == BEARISH {
				if second.Get(4) < first.Get(0) && second.Get(0) > first.Get(4) {
					//pattern.Type = model.BEARISH_ENGULFING
					prices.DataRows[i].Set(ret, 4.0)

				}
			}
		}
		// Bullish Inside Bar
		if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && first.Get(1) > second.Get(1) && first.Get(2) < second.Get(2) {
			prices.DataRows[i].Set(ret, 5.0)
		}
		// Bearish Inside Bar
		if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && first.Get(1) > second.Get(1) && first.Get(2) < second.Get(2) {
			prices.DataRows[i].Set(ret, 6.0)
		}
		// Harami Bearish
		if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && first.Get(4) > second.Get(0) && first.Get(0) < second.Get(4) {
			prices.DataRows[i].Set(ret, 7.0)
		}
		// Harami Bullish
		if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && first.Get(0) > second.Get(4) && first.Get(4) < second.Get(0) {
			prices.DataRows[i].Set(ret, 8.0)
		}
		// Piercing line
		if first.Get(ci+7) == BEARISH && second.Get(ci+7) == BULLISH && first.Get(2) > second.Get(0) && first.Get(ci+2) < second.Get(4) {
			prices.DataRows[i].Set(ret, 14.0)
		}
		// Dark Cloud Cover
		if first.Get(ci+7) == BULLISH && second.Get(ci+7) == BEARISH && first.Get(1) < second.Get(0) && first.Get(ci+2) > second.Get(4) {
			prices.DataRows[i].Set(ret, 15.0)
		}

	}
	for i := 0; i < 10; i++ {
		prices.RemoveColumn()
	}
	return ret
}
