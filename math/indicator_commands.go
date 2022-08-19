package math

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type IndicatorCmd struct {
	Name        string
	CountParams int
	Run         func(candles *Matrix, params []string) int
}

var INDICATOR_COMMANDS = []*IndicatorCmd{
	closeCmd,
	highCmd,
	openCmd,
	lowCmd,
	volumeCmd,
	emaCmd,
	rsiCmd,
	stochasticCmd,
	twapCmd,
	smaCmd,
	adxCmd,
	rocCmd,
	swmaCmd,
	rmaCmd,
	wmaCmd,
	temaCmd,
	demaCmd,
	zlemaCmd,
	zlsmaCmd,
	disparityCmd,
	aoCmd,
	accCmd,
	macdCmd,
	macdzlCmd,
	macdextCmd,
	momentumCmd,
	dpcCmd,
	meanbreakoutCmd,
	consolidatedpricedifferenceCmd,
	rsi_bbCmd,
	rsimomentumCmd,
	atrCmd,
	adrCmd,
	dailyrangeCmd,
	tdrocCmd,
	stochasticextCmd,
	stochasticsmaCmd,
	stochasticrsiCmd,
	rssCmd,
	ppoCmd,
	bollingerbandCmd,
	bollingerband_price_relationCmd,
	esdbandCmd,
	bollingerbandextCmd,
	bollingerbandsqueezeCmd,
	bollingerbandwidthCmd,
	kenvelopeCmd,
	keltnerCmd,
	donchianchannelCmd,
	rarCmd,
	williamsrangeCmd,
	meandistanceCmd,
	perCmd,
	stochasticatrCmd,
	relativevolumeCmd,
	averagepriceCmd,
	voCmd,
	averagevolumeCmd,
	ichimokuCmd,
	weightedtrendintensityCmd,
	supertrendCmd,
	gap_atrCmd,
	gapCmd,
	priceatrCmd,
	kriCmd,
	stdCmd,
	stdchannelCmd,
	stdstochasticCmd,
	demarkCmd,
	demarkerCmd,
	bullishbearishCmd,
	obvCmd,
	aroonCmd,
	trendintensityCmd,
	adCmd,
	tsiCmd,
	divergenceCmd,
	highlowchannelCmd,
	highlowemachannelCmd,
	highestlowestchannelCmd,
	rviCmd,
	rvistochasticCmd,
	kdCmd,
	dpoCmd,
	mfiCmd,
	cciCmd,
	doscCmd,
	hmaCmd,
	cogCmd,
	griCmd,
	cmfCmd,
	stcCmd,
	choppinessCmd,
	spreadCmd,
	gmmaCmd,
	volatilityCmd,
	trixCmd,
	squeezemomentumCmd,
	spreadrangerelationCmd,
	historicalvolatilityCmd,
	minerviniscoreCmd,
	percentrankCmd,
	tsvCmd,
	laguerrefilterCmd,
	apzCmd,
	almaCmd,
	radCmd,
	vptCmd,
	parabolicsarCmd,
	chandelierexitCmd,
	stratclassificationCmd,
	stratpmgCmd,
}

var closeCmd = &IndicatorCmd{
	Name:        "Close",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return ADJ_CLOSE
	},
}

var highCmd = &IndicatorCmd{
	Name:        "High",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return HIGH
	},
}

var openCmd = &IndicatorCmd{
	Name:        "Open",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return OPEN
	},
}

var lowCmd = &IndicatorCmd{
	Name:        "Low",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return LOW
	},
}

var volumeCmd = &IndicatorCmd{
	Name:        "Volume",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return VOLUME
	},
}

var emaCmd = &IndicatorCmd{
	Name:        "EMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return EMA(candles, days, field)
	},
}

var rsiCmd = &IndicatorCmd{
	Name:        "RSI",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return RSI(candles, days, ADJ_CLOSE)
	},
}

var stochasticCmd = &IndicatorCmd{
	Name:        "Stochastic",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		smooth, _ := strconv.Atoi(params[1])
		return Stochastic(candles, days, smooth)
	},
}

var twapCmd = &IndicatorCmd{
	Name:        "TWAP",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return TWAP(candles, days)
	},
}

var smaCmd = &IndicatorCmd{
	Name:        "SMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return SMA(candles, days, field)
	},
}

var adxCmd = &IndicatorCmd{
	Name:        "ADX",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return ADX(candles, days)
	},
}

var rocCmd = &IndicatorCmd{
	Name:        "ROC",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return ROC(candles, days, ADJ_CLOSE)
	},
}

var swmaCmd = &IndicatorCmd{
	Name:        "SWMA",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		field, _ := strconv.Atoi(params[0])
		return SWMA(candles, field)
	},
}

var rmaCmd = &IndicatorCmd{
	Name:        "RMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return RMA(candles, days, field)
	},
}

var wmaCmd = &IndicatorCmd{
	Name:        "WMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return WMA(candles, days, field)
	},
}

var temaCmd = &IndicatorCmd{
	Name:        "TEMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return TEMA(candles, days, field)
	},
}

var demaCmd = &IndicatorCmd{
	Name:        "DEMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return DEMA(candles, days, field)
	},
}

var zlemaCmd = &IndicatorCmd{
	Name:        "ZLEMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return ZLEMA(candles, days, field)
	},
}

func FindIndicatorCmd(name string) bool {
	for _, ic := range INDICATOR_COMMANDS {
		if ic.Name == name {
			return true
		}
	}
	return false
}

var zlsmaCmd = &IndicatorCmd{
	Name:        "ZLSMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return ZLSMA(candles, days, field)
	},
}

var disparityCmd = &IndicatorCmd{
	Name:        "Disparity",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return Disparity(candles, days)
	},
}

var aoCmd = &IndicatorCmd{
	Name:        "AO",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		return AO(candles, short, long)
	},
}

var accCmd = &IndicatorCmd{
	Name:        "ACC",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		s, _ := strconv.Atoi(params[2])
		return ACC(candles, short, long, s)
	},
}

var macdCmd = &IndicatorCmd{
	Name:        "MACD",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		signal, _ := strconv.Atoi(params[2])
		return MACD(candles, short, long, signal)
	},
}

var macdzlCmd = &IndicatorCmd{
	Name:        "MACDZL",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		signal, _ := strconv.Atoi(params[2])
		return MACDZL(candles, short, long, signal)
	},
}

var macdextCmd = &IndicatorCmd{
	Name:        "MACDExt",
	CountParams: 4,
	Run: func(candles *Matrix, params []string) int {
		field, _ := strconv.Atoi(params[0])
		short, _ := strconv.Atoi(params[1])
		long, _ := strconv.Atoi(params[2])
		signal, _ := strconv.Atoi(params[3])
		return MACDExt(candles, field, short, long, signal)
	},
}

var momentumCmd = &IndicatorCmd{
	Name:        "Momentum",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return Momentum(candles, days)
	},
}

var dpcCmd = &IndicatorCmd{
	Name:        "DPC",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return DPC(candles)
	},
}

var meanbreakoutCmd = &IndicatorCmd{
	Name:        "MeanBreakout",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return MeanBreakout(candles, period)
	},
}

var consolidatedpricedifferenceCmd = &IndicatorCmd{
	Name:        "ConsolidatedPriceDifference",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		min, _ := strconv.Atoi(params[0])
		return ConsolidatedPriceDifference(candles, min)
	},
}

var rsi_bbCmd = &IndicatorCmd{
	Name:        "RSI_BB",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return RSI_BB(candles, days, field)
	},
}

var rsimomentumCmd = &IndicatorCmd{
	Name:        "RSIMomentum",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		field, _ := strconv.Atoi(params[2])
		return RSIMomentum(candles, short, long, field)
	},
}

var atrCmd = &IndicatorCmd{
	Name:        "ATR",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return ATR(candles, days)
	},
}

var adrCmd = &IndicatorCmd{
	Name:        "ADR",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return ADR(candles, days)
	},
}

var dailyrangeCmd = &IndicatorCmd{
	Name:        "DailyRange",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return DailyRange(candles, days)
	},
}

var tdrocCmd = &IndicatorCmd{
	Name:        "TDROC",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return TDROC(candles, days, field)
	},
}

var stochasticextCmd = &IndicatorCmd{
	Name:        "StochasticExt",
	CountParams: 5,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		ema, _ := strconv.Atoi(params[1])
		highField, _ := strconv.Atoi(params[2])
		lowField, _ := strconv.Atoi(params[3])
		priceField, _ := strconv.Atoi(params[4])
		return StochasticExt(candles, days, ema, highField, lowField, priceField)
	},
}

var stochasticsmaCmd = &IndicatorCmd{
	Name:        "StochasticSMA",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		sma, _ := strconv.Atoi(params[0])
		days, _ := strconv.Atoi(params[1])
		ema, _ := strconv.Atoi(params[2])
		return StochasticSMA(candles, sma, days, ema)
	},
}

var stochasticrsiCmd = &IndicatorCmd{
	Name:        "StochasticRSI",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		smoothK, _ := strconv.Atoi(params[1])
		smoothD, _ := strconv.Atoi(params[2])
		return StochasticRSI(candles, days, smoothK, smoothD)
	},
}

var rssCmd = &IndicatorCmd{
	Name:        "RSS",
	CountParams: 4,
	Run: func(candles *Matrix, params []string) int {
		slow, _ := strconv.Atoi(params[0])
		fast, _ := strconv.Atoi(params[1])
		rsi, _ := strconv.Atoi(params[2])
		smoothing, _ := strconv.Atoi(params[3])
		return RSS(candles, slow, fast, rsi, smoothing)
	},
}

var ppoCmd = &IndicatorCmd{
	Name:        "PPO",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		signal, _ := strconv.Atoi(params[2])
		return PPO(candles, short, long, signal)
	},
}

var bollingerbandCmd = &IndicatorCmd{
	Name:        "BollingerBand",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		upper, _ := strconv.ParseFloat(params[1], 64)
		lower, _ := strconv.ParseFloat(params[2], 64)
		return BollingerBand(candles, ema, upper, lower)
	},
}

var bollingerband_price_relationCmd = &IndicatorCmd{
	Name:        "BollingerBand_Price_Relation",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		upper, _ := strconv.ParseFloat(params[1], 64)
		lower, _ := strconv.ParseFloat(params[2], 64)
		return BollingerBand_Price_Relation(candles, ema, upper, lower)
	},
}

var esdbandCmd = &IndicatorCmd{
	Name:        "ESDBand",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		upper, _ := strconv.ParseFloat(params[1], 64)
		lower, _ := strconv.ParseFloat(params[2], 64)
		return ESDBand(candles, ema, upper, lower)
	},
}

var bollingerbandextCmd = &IndicatorCmd{
	Name:        "BollingerBandExt",
	CountParams: 4,
	Run: func(candles *Matrix, params []string) int {
		field, _ := strconv.Atoi(params[0])
		ema, _ := strconv.Atoi(params[1])
		upper, _ := strconv.ParseFloat(params[2], 64)
		lower, _ := strconv.ParseFloat(params[3], 64)
		return BollingerBandExt(candles, field, ema, upper, lower)
	},
}

var bollingerbandsqueezeCmd = &IndicatorCmd{
	Name:        "BollingerBandSqueeze",
	CountParams: 4,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		upper, _ := strconv.ParseFloat(params[1], 64)
		lower, _ := strconv.ParseFloat(params[2], 64)
		period, _ := strconv.Atoi(params[3])
		return BollingerBandSqueeze(candles, ema, upper, lower, period)
	},
}

var bollingerbandwidthCmd = &IndicatorCmd{
	Name:        "BollingerBandWidth",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		upper, _ := strconv.ParseFloat(params[1], 64)
		lower, _ := strconv.ParseFloat(params[2], 64)
		return BollingerBandWidth(candles, ema, upper, lower)
	},
}

var kenvelopeCmd = &IndicatorCmd{
	Name:        "KEnvelope",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return KEnvelope(candles, days)
	},
}

var keltnerCmd = &IndicatorCmd{
	Name:        "Keltner",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		atrLength, _ := strconv.Atoi(params[1])
		multiplier, _ := strconv.ParseFloat(params[2], 64)
		return Keltner(candles, ema, atrLength, multiplier)
	},
}

var donchianchannelCmd = &IndicatorCmd{
	Name:        "DonchianChannel",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return DonchianChannel(candles, days)
	},
}

var rarCmd = &IndicatorCmd{
	Name:        "RAR",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return RAR(candles, days)
	},
}

var williamsrangeCmd = &IndicatorCmd{
	Name:        "WilliamsRange",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return WilliamsRange(candles, days)
	},
}

var meandistanceCmd = &IndicatorCmd{
	Name:        "MeanDistance",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return MeanDistance(candles, lookback)
	},
}

var perCmd = &IndicatorCmd{
	Name:        "PER",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		ema, _ := strconv.Atoi(params[0])
		smoothing, _ := strconv.Atoi(params[1])
		return PER(candles, ema, smoothing)
	},
}

var stochasticatrCmd = &IndicatorCmd{
	Name:        "StochasticATR",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return StochasticATR(candles, days)
	},
}

var relativevolumeCmd = &IndicatorCmd{
	Name:        "RelativeVolume",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return RelativeVolume(candles, period)
	},
}

var averagepriceCmd = &IndicatorCmd{
	Name:        "AveragePrice",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return AveragePrice(candles)
	},
}

var voCmd = &IndicatorCmd{
	Name:        "VO",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		fast, _ := strconv.Atoi(params[0])
		slow, _ := strconv.Atoi(params[1])
		return VO(candles, fast, slow)
	},
}

var averagevolumeCmd = &IndicatorCmd{
	Name:        "AverageVolume",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return AverageVolume(candles, lookback)
	},
}

var ichimokuCmd = &IndicatorCmd{
	Name:        "Ichimoku",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		mid, _ := strconv.Atoi(params[1])
		long, _ := strconv.Atoi(params[2])
		return Ichimoku(candles, short, mid, long)
	},
}

var weightedtrendintensityCmd = &IndicatorCmd{
	Name:        "WeightedTrendIntensity",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return WeightedTrendIntensity(candles, period)
	},
}

var supertrendCmd = &IndicatorCmd{
	Name:        "Supertrend",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		multiplier, _ := strconv.ParseFloat(params[1], 64)
		return Supertrend(candles, period, multiplier)
	},
}

var gap_atrCmd = &IndicatorCmd{
	Name:        "GAP_ATR",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return GAP_ATR(candles)
	},
}

var gapCmd = &IndicatorCmd{
	Name:        "GAP",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return GAP(candles)
	},
}

var priceatrCmd = &IndicatorCmd{
	Name:        "PriceATR",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return PriceATR(candles, period)
	},
}

var kriCmd = &IndicatorCmd{
	Name:        "KRI",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return KRI(candles, period)
	},
}

var stdCmd = &IndicatorCmd{
	Name:        "STD",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return STD(candles, days)
	},
}

var stdchannelCmd = &IndicatorCmd{
	Name:        "STDChannel",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		std, _ := strconv.ParseFloat(params[1], 64)
		return STDChannel(candles, days, std)
	},
}

var stdstochasticCmd = &IndicatorCmd{
	Name:        "STDStochastic",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return STDStochastic(candles, days)
	},
}

var demarkCmd = &IndicatorCmd{
	Name:        "DeMark",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		return DeMark(candles)
	},
}

var demarkerCmd = &IndicatorCmd{
	Name:        "DeMarker",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return DeMarker(candles, days)
	},
}

var bullishbearishCmd = &IndicatorCmd{
	Name:        "BullishBearish",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return BullishBearish(candles, period)
	},
}

var obvCmd = &IndicatorCmd{
	Name:        "OBV",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		scale, _ := strconv.ParseFloat(params[0], 64)
		return OBV(candles, scale)
	},
}

var aroonCmd = &IndicatorCmd{
	Name:        "Aroon",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return Aroon(candles, days)
	},
}

var trendintensityCmd = &IndicatorCmd{
	Name:        "TrendIntensity",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return TrendIntensity(candles, days)
	},
}

var adCmd = &IndicatorCmd{
	Name:        "AD",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return AD(candles)
	},
}

var tsiCmd = &IndicatorCmd{
	Name:        "TSI",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		signal, _ := strconv.Atoi(params[2])
		return TSI(candles, short, long, signal)
	},
}

var divergenceCmd = &IndicatorCmd{
	Name:        "Divergence",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		first, _ := strconv.Atoi(params[0])
		second, _ := strconv.Atoi(params[1])
		period, _ := strconv.Atoi(params[2])
		return Divergence(candles, first, second, period)
	},
}

var highlowchannelCmd = &IndicatorCmd{
	Name:        "HighLowChannel",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		highPeriod, _ := strconv.Atoi(params[0])
		lowPeriod, _ := strconv.Atoi(params[1])
		return HighLowChannel(candles, highPeriod, lowPeriod)
	},
}

var highlowemachannelCmd = &IndicatorCmd{
	Name:        "HighLowEMAChannel",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		highPeriod, _ := strconv.Atoi(params[0])
		lowPeriod, _ := strconv.Atoi(params[0])
		return HighLowEMAChannel(candles, highPeriod, lowPeriod)
	},
}

var highestlowestchannelCmd = &IndicatorCmd{
	Name:        "HighestLowestChannel",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return HighestLowestChannel(candles, period)
	},
}

var rviCmd = &IndicatorCmd{
	Name:        "RVI",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		signal, _ := strconv.Atoi(params[1])
		return RVI(candles, lookback, signal)
	},
}

var rvistochasticCmd = &IndicatorCmd{
	Name:        "RVIStochastic",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		signal, _ := strconv.Atoi(params[1])
		return RVIStochastic(candles, lookback, signal)
	},
}

var kdCmd = &IndicatorCmd{
	Name:        "KD",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return KD(candles, period)
	},
}

var dpoCmd = &IndicatorCmd{
	Name:        "DPO",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return DPO(candles, period)
	},
}

var mfiCmd = &IndicatorCmd{
	Name:        "MFI",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return MFI(candles, days)
	},
}

var cciCmd = &IndicatorCmd{
	Name:        "CCI",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		smoothed, _ := strconv.Atoi(params[1])
		return CCI(candles, days, smoothed)
	},
}

var doscCmd = &IndicatorCmd{
	Name:        "DOSC",
	CountParams: 5,
	Run: func(candles *Matrix, params []string) int {
		r, _ := strconv.Atoi(params[0])
		e1, _ := strconv.Atoi(params[1])
		e2, _ := strconv.Atoi(params[2])
		s, _ := strconv.Atoi(params[3])
		sl, _ := strconv.Atoi(params[4])
		return DOSC(candles, r, e1, e2, s, sl)
	},
}

var hmaCmd = &IndicatorCmd{
	Name:        "HMA",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return HMA(candles, period, field)
	},
}

var cogCmd = &IndicatorCmd{
	Name:        "COG",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return COG(candles, period)
	},
}

var griCmd = &IndicatorCmd{
	Name:        "GRI",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return GRI(candles, period)
	},
}

var cmfCmd = &IndicatorCmd{
	Name:        "CMF",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		return CMF(candles, period)
	},
}

var stcCmd = &IndicatorCmd{
	Name:        "STC",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		short, _ := strconv.Atoi(params[0])
		long, _ := strconv.Atoi(params[1])
		stoch, _ := strconv.Atoi(params[2])
		return STC(candles, short, long, stoch)
	},
}

var choppinessCmd = &IndicatorCmd{
	Name:        "Choppiness",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		days, _ := strconv.Atoi(params[0])
		return Choppiness(candles, days)
	},
}

var spreadCmd = &IndicatorCmd{
	Name:        "Spread",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return Spread(candles, lookback)
	},
}

var gmmaCmd = &IndicatorCmd{
	Name:        "GMMA",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return GMMA(candles)
	},
}

var volatilityCmd = &IndicatorCmd{
	Name:        "Volatility",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return Volatility(candles, lookback, field)
	},
}

var trixCmd = &IndicatorCmd{
	Name:        "TRIX",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return TRIX(candles, lookback)
	},
}

var squeezemomentumCmd = &IndicatorCmd{
	Name:        "SqueezeMomentum",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return SqueezeMomentum(candles, lookback)
	},
}

var spreadrangerelationCmd = &IndicatorCmd{
	Name:        "SpreadRangeRelation",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return SpreadRangeRelation(candles, lookback)
	},
}

var historicalvolatilityCmd = &IndicatorCmd{
	Name:        "HistoricalVolatility",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		lookback, _ := strconv.Atoi(params[0])
		return HistoricalVolatility(candles, lookback)
	},
}

var minerviniscoreCmd = &IndicatorCmd{
	Name:        "MinerviniScore",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return MinerviniScore(candles)
	},
}

var percentrankCmd = &IndicatorCmd{
	Name:        "PercentRank",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		field, _ := strconv.Atoi(params[1])
		return PercentRank(candles, period, field)
	},
}

var tsvCmd = &IndicatorCmd{
	Name:        "TSV",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return TSV(candles)
	},
}

var laguerrefilterCmd = &IndicatorCmd{
	Name:        "LaguerreFilter",
	CountParams: 1,
	Run: func(candles *Matrix, params []string) int {
		gamma, _ := strconv.ParseFloat(params[0], 64)
		return LaguerreFilter(candles, gamma)
	},
}

var apzCmd = &IndicatorCmd{
	Name:        "APZ",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		dev, _ := strconv.ParseFloat(params[1], 64)
		return APZ(candles, period, dev)
	},
}

var almaCmd = &IndicatorCmd{
	Name:        "ALMA",
	CountParams: 3,
	Run: func(candles *Matrix, params []string) int {
		windowSize, _ := strconv.Atoi(params[0])
		offset, _ := strconv.ParseFloat(params[1], 64)
		sigma, _ := strconv.ParseFloat(params[2], 64)
		return ALMA(candles, windowSize, offset, sigma)
	},
}

var radCmd = &IndicatorCmd{
	Name:        "RAD",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		ma, _ := strconv.Atoi(params[0])
		period, _ := strconv.Atoi(params[1])
		return RAD(candles, ma, period)
	},
}

var vptCmd = &IndicatorCmd{
	Name:        "VPT",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return VPT(candles)
	},
}

var parabolicsarCmd = &IndicatorCmd{
	Name:        "ParabolicSAR",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return ParabolicSAR(candles)
	},
}

var chandelierexitCmd = &IndicatorCmd{
	Name:        "ChandelierExit",
	CountParams: 2,
	Run: func(candles *Matrix, params []string) int {
		period, _ := strconv.Atoi(params[0])
		multiplier, _ := strconv.ParseFloat(params[1], 64)
		return ChandelierExit(candles, period, multiplier)
	},
}

var stratclassificationCmd = &IndicatorCmd{
	Name:        "StratClassification",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return StratClassification(candles)
	},
}

var stratpmgCmd = &IndicatorCmd{
	Name:        "StratPMG",
	CountParams: 0,
	Run: func(candles *Matrix, params []string) int {
		return StratPMG(candles)
	},
}

func RunIndicatorCmd(name string, candles *Matrix, params string) (int, error) {
	for _, ic := range INDICATOR_COMMANDS {
		if ic.Name == name {
			entries := strings.Split(params, ",")
			if len(entries) == ic.CountParams || ic.CountParams == 0 {
				return ic.Run(candles, entries), nil
			} else {
				return -1, errors.New(fmt.Sprintf("Not enough arguments for %s - expected: %d but got %s", name, ic.CountParams, params))
			}
		}
	}
	return -1, errors.New("No matching indicator found: " + name)
}
