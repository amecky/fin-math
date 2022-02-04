package math

func FindCandleStickPatterns(prices *Matrix) int {
	ret := prices.AddColumn()
	// 0 = BodySize 1 = BodyPos 2 = Mid 3 = RelBodySize 4 = RelAvg 5 = Upper 6 = Lower 7 = Trend 8 = Spread 9 = RelSpread
	ci := Candles(prices, 21)
	for i := 1; i < prices.Rows; i++ {
		first := prices.DataRows[i-1]
		if first.Get(ci+7) == 1.0 && first.Get(ci+6) >= 66.0 {
			//pattern.Type = model.HAMMER
			prices.DataRows[i].Set(ret, 1.0)

		}
		if first.Get(ci+7) == -1 && first.Get(ci+5) >= 66.0 {
			//pattern.Type = model.SHOOTING_STAR
			prices.DataRows[i].Set(ret, 2.0)

		}

		second := prices.DataRows[i]
		if first.Get(ci+7) != second.Get(ci+7) {
			if second.Get(ci+7) == 1 {
				if second.Get(4) > first.Get(0) && second.Get(0) < first.Get(4) {
					//pattern.Type = model.BULLISH_ENGULFING
					prices.DataRows[i].Set(ret, 3.0)
				}
			}
			if second.Get(ci+7) == -1 {
				if second.Get(4) < first.Get(0) && second.Get(0) > first.Get(4) {
					//pattern.Type = model.BEARISH_ENGULFING
					prices.DataRows[i].Set(ret, 4.0)

				}
			}
		}

	}
	for i := 0; i < 10; i++ {
		prices.RemoveColumn()
	}
	return ret
}
