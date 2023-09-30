
<!--- Awesome Oscilator -->

# Awesome Oscilator

## Description

[Tradingview](https://de.tradingview.com/scripts/awesomeoscillator/)

## Notes

## Parameters

| Name  | Description         |
|-------|---------------------|
| short | Period of short SMA |
| long  | Period of long SMA  |

## Example
```
AO(8,21)
```

## Provider

The provider shows 1 column. 

| Name |Description        |
|------|-------------------|
| AO   |The current value  |

<!--- Elder-Bars -->

# Elder-Bars

## Description

The Impulse System is based on two indicators, a 13-day exponential moving average and the MACD-Histogram. The moving average identifies the trend, while the MACD-Histogram measures momentum. As a result, the Impulse System combines trend following and momentum to identify tradable impulses. This unique indicator combination is color coded into the price bars for easy reference.

* [Medium](https://kaabar-sofien.medium.com/the-elder-ray-index-for-trading-b54c9b1741aa)
* [Stockcharts](https://school.stockcharts.com/doku.php?id=chart_analysis:elder_impulse_system)

## Notes

The indicator only has 3 modes 1,0,-1 correlating to the bullish, bearish or mixed mode. Change in value can show change in trend

## Parameters

None

## Example
```
Elder-Bars
```

## Provider

The provider shows 1 column. 

| Name |Description                                      |
|------|-------------------------------------------------|
| EB   | The elder bar as text only if there is a change |

<!--- Market-Regime -->

# Market-Regime

## Description

Compares current price to EMA 5,8,13 steps before to detect the market regime

## Notes

Ranges between -1 and 1. Look for crossovers of the zero line. Shows als the strength of trend

## Parameters

| Name   |Description      |
|--------|-----------------|
| period | The period EMA  |

## Example
```
Market-Regime(10)
```

## Provider

The provider shows 1 column. 

| Name |Description                   |
|------|--------------------------------|
| MR   |The market regime as histogram  |

<!--- Mean Breakout -->

# Mean Breakout

## Description

The Mean Breakout (MBO) compares the difference between the closing price of a candle and a moving average over N periods to the difference between the min and max value of the closing price over the same N periods. 

[Medium](https://medium.com/superalgos/all-in-one-indicator-for-exponential-moving-average-crossover-strategy-b6b4b0da957e)

## Notes

The MBO indicator is centered around 0. A crossing to the positive zone indicates and upward trend and will give a +1 to the MBO component of the all-in-one indicator and a -1 will be attributed at a downward trend signal, MBO crossing to the negative region

## Parameters

| Name   |Description                  |
|--------|-----------------------------|
| period | The period for caluclation  |

## Example
```
MBO(14)
```

## Provider

The provider shows 1 column. 

| Name |Description                |
|------|---------------------------|
| MBO  |The MBO value as histogram |


<!--- Vortex -->

# Vortex

## Description

A vortex indicator (VI) is an indicator composed of two lines - an uptrend line (VI+) and a downtrend line (VI-). These lines are typically colored green and red respectively. A vortex indicator is used to spot trend reversals and confirm current trends.

[Investopedia](https://www.investopedia.com/terms/v/vortex-indicator-vi.asp)

## Notes

Look for crossing of the two line. This can mark a reversal of the trend.
Also look for reversals of each line as it shows reversals.

## Parameters

| Name   |Description              |
|--------|-------------------------|
| period | The period for the sum  |

## Example
```
Vortex(10)
```

## Provider

The provider shows 3 columns. 

| Name |Description                  |
|------|-----------------------------|
| VP   |The Vortex plus values       |
| VM   | The Vortex minus values     |
| Diff | The difference as histogram |

<!--- Template -->

# Template

## Description

DESC

[Investopedia](https://www.investopedia.com/terms/v/vortex-indicator-vi.asp)

## Notes

NOTES

## Parameters

| Name   |Description              |
|--------|-------------------------|
| period | The period for the sum  |

## Example
```
XXXX(10)
```

## Provider

The provider shows 3 columns. 

| Name |Description                  |
|------|-----------------------------|
| VP   |The Vortex plus values       |


