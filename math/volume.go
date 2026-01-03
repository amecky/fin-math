package math

// -----------------------------------------------------------------------
//
//	VO
//
// -----------------------------------------------------------------------
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/volume-oscillator
// https://commodity.com/technical-analysis/volume-oscillator/
func VO(m *Matrix, fast, slow int) int {
	ret := m.AddColumn()
	emaSlow := EMA(m, slow, VOLUME)
	emaFast := EMA(m, fast, VOLUME)
	for i := 0; i < m.Rows; i++ {
		if m.DataRows[i].Get(emaSlow) != 0.0 {
			vo := (m.DataRows[i].Get(emaFast) - m.DataRows[i].Get(emaSlow)) / m.DataRows[i].Get(emaSlow) * 100.0
			m.DataRows[i].Set(ret, vo)
		}
	}
	m.RemoveColumn()
	m.RemoveColumn()
	return ret
}

// -----------------------------------------------------------------------
//
//	AverageVolume
//
// -----------------------------------------------------------------------
// https://www.fidelity.com/learning-center/trading-investing/technical-analysis/technical-indicator-guide/volume-oscillator
// https://commodity.com/technical-analysis/volume-oscillator/
func AverageVolume(m *Matrix, lookback int) int {
	// 0 = SMA of volume
	sma := SMA(m, lookback, 5)
	return sma
}

/*
calcVolumes(OHLCV ohlcv) =>
    var VolumeData data = VolumeData.new()
    data.buyVol       := ohlcv.V * (ohlcv.C - ohlcv.L) / (ohlcv.H - ohlcv.L) // Calculate buy volume using the formula: volume * (close - low) / (high - low)
    data.sellVol      := ohlcv.V - data.buyVol                               // Calculate sell volume by subtracting buy volume from total volume
    data.pcBuy        := data.buyVol / ohlcv.V * 100                         // Calculate the percentage of buy volume
    data.pcSell       := 100 - data.pcBuy                                    // Calculate the percentage of sell volume (100% - buy percentage)
    data.isBuyGreater := data.buyVol > data.sellVol                          // Determine if buy volume is greater than sell volume
    data.higherVol    := data.isBuyGreater ? data.buyVol  : data.sellVol     // Assign the higher volume value based on the comparison
    data.lowerVol     := data.isBuyGreater ? data.sellVol : data.buyVol      // Assign the lower volume value based on the comparison
    data.higherCol    := data.isBuyGreater ? C_Up     : C_Down               // Assign the color for the higher volume bar based on the comparison
    data.lowerCol     := data.isBuyGreater ? C_Down   : C_Up                 // Assign the color for the lower volume bar based on the comparison
    data
*/
func DeltaVolume(m *Matrix) int {
	// 0 = Buy Volume Perentage 1 = Sell Volume Percentage
	buy := m.AddNamedColumn("BuyVolume")
	sell := m.AddNamedColumn("SellVolume")
	for i := range m.Rows {
		c := m.DataRows[i]
		bv := c.Get(5) * (c.Close() - c.Low()) / (c.High() - c.Low())
		bvp := bv / c.Get(5) * 100.0
		m.DataRows[i].Set(buy, bvp)
		m.DataRows[i].Set(sell, 100-bvp)
	}
	return buy
}
