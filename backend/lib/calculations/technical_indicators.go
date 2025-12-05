package calculations

import (
	"math"
)

// PriceData represents OHLCV data for a single period.
type PriceData struct {
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// SMA calculates Simple Moving Average.
func SMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return []float64{}
	}

	result := make([]float64, len(prices)-period+1)
	for i := range result {
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += prices[i+j]
		}
		result[i] = sum / float64(period)
	}

	return result
}

// EMA calculates Exponential Moving Average.
func EMA(prices []float64, period int) []float64 {
	if len(prices) < period {
		return []float64{}
	}

	result := make([]float64, len(prices))
	multiplier := 2.0 / float64(period+1)

	// Start with SMA for first value
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += prices[i]
	}
	result[period-1] = sum / float64(period)

	// Calculate EMA for remaining values
	for i := period; i < len(prices); i++ {
		result[i] = (prices[i]-result[i-1])*multiplier + result[i-1]
	}

	return result
}

// RSI calculates Relative Strength Index.
func RSI(prices []float64, period int) []float64 {
	if len(prices) < period+1 {
		return []float64{}
	}

	changes := make([]float64, len(prices)-1)
	for i := 1; i < len(prices); i++ {
		changes[i-1] = prices[i] - prices[i-1]
	}

	gains := make([]float64, len(changes))
	losses := make([]float64, len(changes))
	for i, change := range changes {
		if change > 0 {
			gains[i] = change
		} else {
			losses[i] = -change
		}
	}

	result := make([]float64, len(prices)-period)

	// Calculate initial averages
	avgGain := 0.0
	avgLoss := 0.0
	for i := 0; i < period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Calculate RSI for each point
	for i := period; i < len(changes); i++ {
		avgGain = ((avgGain * float64(period-1)) + gains[i]) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + losses[i]) / float64(period)

		rs := avgGain / avgLoss
		rsi := 100 - (100 / (1 + rs))

		result[i-period] = rsi
	}

	return result
}

// MACD calculates Moving Average Convergence Divergence.
type MACDResult struct {
	MACD      float64
	Signal    float64
	Histogram float64
}

func MACD(prices []float64, fastPeriod, slowPeriod, signalPeriod int) []MACDResult {
	if len(prices) < slowPeriod {
		return []MACDResult{}
	}

	fastEMA := EMA(prices, fastPeriod)
	slowEMA := EMA(prices, slowPeriod)

	// Calculate MACD line
	macdLine := make([]float64, len(slowEMA))
	for i := range slowEMA {
		macdLine[i] = fastEMA[i+len(fastEMA)-len(slowEMA)] - slowEMA[i]
	}

	// Calculate signal line
	signalLine := EMA(macdLine, signalPeriod)

	// Calculate histogram
	results := make([]MACDResult, len(signalLine))
	for i := range results {
		macdVal := macdLine[i+len(macdLine)-len(signalLine)]
		results[i] = MACDResult{
			MACD:      macdVal,
			Signal:    signalLine[i],
			Histogram: macdVal - signalLine[i],
		}
	}

	return results
}

// BollingerBands calculates Bollinger Bands.
type BollingerBandsResult struct {
	Upper  float64
	Middle float64
	Lower  float64
}

func BollingerBands(prices []float64, period int, stdDevMultiplier float64) []BollingerBandsResult {
	if len(prices) < period {
		return []BollingerBandsResult{}
	}

	sma := SMA(prices, period)
	results := make([]BollingerBandsResult, len(sma))

	for i := range sma {
		// Calculate standard deviation
		sum := 0.0
		for j := 0; j < period; j++ {
			diff := prices[i+j] - sma[i]
			sum += diff * diff
		}
		stdDev := math.Sqrt(sum / float64(period))

		results[i] = BollingerBandsResult{
			Upper:  sma[i] + (stdDevMultiplier * stdDev),
			Middle: sma[i],
			Lower:  sma[i] - (stdDevMultiplier * stdDev),
		}
	}

	return results
}

// ATR calculates Average True Range.
func ATR(data []PriceData, period int) []float64 {
	if len(data) < period+1 {
		return []float64{}
	}

	trueRanges := make([]float64, len(data)-1)
	for i := 1; i < len(data); i++ {
		highLow := data[i].High - data[i].Low
		highClose := math.Abs(data[i].High - data[i-1].Close)
		lowClose := math.Abs(data[i].Low - data[i-1].Close)

		trueRanges[i-1] = math.Max(highLow, math.Max(highClose, lowClose))
	}

	// Calculate initial ATR as simple average
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += trueRanges[i]
	}
	initialATR := sum / float64(period)

	result := make([]float64, len(trueRanges)-period+1)
	result[0] = initialATR

	// Calculate subsequent ATR values using smoothing
	for i := 1; i < len(result); i++ {
		result[i] = ((result[i-1] * float64(period-1)) + trueRanges[period+i-1]) / float64(period)
	}

	return result
}

// Stochastic calculates Stochastic Oscillator.
type StochasticResult struct {
	K float64 // Fast line
	D float64 // Slow line (SMA of K)
}

func Stochastic(data []PriceData, kPeriod, dPeriod int) []StochasticResult {
	if len(data) < kPeriod {
		return []StochasticResult{}
	}

	kValues := make([]float64, len(data)-kPeriod+1)

	for i := 0; i < len(kValues); i++ {
		high := data[i].High
		low := data[i].Low

		// Find highest high and lowest low in period
		for j := 1; j < kPeriod; j++ {
			if data[i+j].High > high {
				high = data[i+j].High
			}
			if data[i+j].Low < low {
				low = data[i+j].Low
			}
		}

		currentClose := data[i+kPeriod-1].Close
		kValues[i] = ((currentClose - low) / (high - low)) * 100
	}

	// Calculate D (SMA of K)
	dValues := SMA(kValues, dPeriod)

	results := make([]StochasticResult, len(dValues))
	for i := range results {
		results[i] = StochasticResult{
			K: kValues[i+len(kValues)-len(dValues)],
			D: dValues[i],
		}
	}

	return results
}

// ADX calculates Average Directional Index.
func ADX(data []PriceData, period int) []float64 {
	if len(data) < period+1 {
		return []float64{}
	}

	// Calculate +DM and -DM
	plusDM := make([]float64, len(data)-1)
	minusDM := make([]float64, len(data)-1)

	for i := 1; i < len(data); i++ {
		highDiff := data[i].High - data[i-1].High
		lowDiff := data[i-1].Low - data[i].Low

		if highDiff > lowDiff && highDiff > 0 {
			plusDM[i-1] = highDiff
		}
		if lowDiff > highDiff && lowDiff > 0 {
			minusDM[i-1] = lowDiff
		}
	}

	// Calculate ATR
	atr := ATR(data, period)

	// Calculate +DI and -DI
	plusDI := make([]float64, len(atr))
	minusDI := make([]float64, len(atr))

	for i := range atr {
		plusDI[i] = (plusDM[i+period-1] / atr[i]) * 100
		minusDI[i] = (minusDM[i+period-1] / atr[i]) * 100
	}

	// Calculate DX
	dx := make([]float64, len(plusDI))
	for i := range dx {
		diDiff := math.Abs(plusDI[i] - minusDI[i])
		diSum := plusDI[i] + minusDI[i]
		if diSum > 0 {
			dx[i] = (diDiff / diSum) * 100
		}
	}

	// Calculate ADX (smoothed DX)
	return SMA(dx, period)
}

// CCI calculates Commodity Channel Index.
func CCI(data []PriceData, period int) []float64 {
	if len(data) < period {
		return []float64{}
	}

	typicalPrices := make([]float64, len(data))
	for i := range data {
		typicalPrices[i] = (data[i].High + data[i].Low + data[i].Close) / 3
	}

	sma := SMA(typicalPrices, period)
	result := make([]float64, len(sma))

	for i := range sma {
		// Calculate mean deviation
		sum := 0.0
		for j := 0; j < period; j++ {
			sum += math.Abs(typicalPrices[i+j] - sma[i])
		}
		meanDeviation := sum / float64(period)

		result[i] = (typicalPrices[i+period-1] - sma[i]) / (0.015 * meanDeviation)
	}

	return result
}

// OBV calculates On-Balance Volume.
func OBV(data []PriceData) []float64 {
	if len(data) < 2 {
		return []float64{}
	}

	result := make([]float64, len(data))
	result[0] = data[0].Volume

	for i := 1; i < len(data); i++ {
		if data[i].Close > data[i-1].Close {
			result[i] = result[i-1] + data[i].Volume
		} else if data[i].Close < data[i-1].Close {
			result[i] = result[i-1] - data[i].Volume
		} else {
			result[i] = result[i-1]
		}
	}

	return result
}

// WilliamsR calculates Williams %R.
func WilliamsR(data []PriceData, period int) []float64 {
	if len(data) < period {
		return []float64{}
	}

	result := make([]float64, len(data)-period+1)

	for i := 0; i < len(result); i++ {
		high := data[i].High
		low := data[i].Low

		for j := 1; j < period; j++ {
			if data[i+j].High > high {
				high = data[i+j].High
			}
			if data[i+j].Low < low {
				low = data[i+j].Low
			}
		}

		currentClose := data[i+period-1].Close
		result[i] = ((high - currentClose) / (high - low)) * -100
	}

	return result
}

// VWAP calculates Volume Weighted Average Price.
func VWAP(data []PriceData) []float64 {
	result := make([]float64, len(data))
	cumulativeTPV := 0.0
	cumulativeVolume := 0.0

	for i := range data {
		typicalPrice := (data[i].High + data[i].Low + data[i].Close) / 3
		cumulativeTPV += typicalPrice * data[i].Volume
		cumulativeVolume += data[i].Volume

		if cumulativeVolume > 0 {
			result[i] = cumulativeTPV / cumulativeVolume
		}
	}

	return result
}

// ParabolicSAR calculates Parabolic SAR.
func ParabolicSAR(data []PriceData, acceleration, maxAcceleration float64) []float64 {
	if len(data) < 2 {
		return []float64{}
	}

	result := make([]float64, len(data))
	result[0] = data[0].Low

	isUptrend := true
	af := acceleration
	ep := data[0].High
	sar := data[0].Low

	for i := 1; i < len(data); i++ {
		sar = sar + af*(ep-sar)

		if isUptrend {
			if data[i].Low < sar {
				isUptrend = false
				sar = ep
				ep = data[i].Low
				af = acceleration
			} else {
				if data[i].High > ep {
					ep = data[i].High
					af = math.Min(af+acceleration, maxAcceleration)
				}
			}
		} else {
			if data[i].High > sar {
				isUptrend = true
				sar = ep
				ep = data[i].High
				af = acceleration
			} else {
				if data[i].Low < ep {
					ep = data[i].Low
					af = math.Min(af+acceleration, maxAcceleration)
				}
			}
		}

		result[i] = sar
	}

	return result
}

// IchimokuCloud calculates Ichimoku Cloud components.
type IchimokuResult struct {
	TenkanSen   float64 // Conversion Line
	KijunSen    float64 // Base Line
	SenkouSpanA float64 // Leading Span A
	SenkouSpanB float64 // Leading Span B
	ChikouSpan  float64 // Lagging Span
}

func IchimokuCloud(data []PriceData, conversionPeriod, basePeriod, spanBPeriod, displacement int) []IchimokuResult {
	if len(data) < spanBPeriod {
		return []IchimokuResult{}
	}

	result := make([]IchimokuResult, len(data))

	for i := 0; i < len(data); i++ {
		// Tenkan-sen (Conversion Line)
		if i >= conversionPeriod-1 {
			high := data[i-conversionPeriod+1].High
			low := data[i-conversionPeriod+1].Low
			for j := i - conversionPeriod + 2; j <= i; j++ {
				if data[j].High > high {
					high = data[j].High
				}
				if data[j].Low < low {
					low = data[j].Low
				}
			}
			result[i].TenkanSen = (high + low) / 2
		}

		// Kijun-sen (Base Line)
		if i >= basePeriod-1 {
			high := data[i-basePeriod+1].High
			low := data[i-basePeriod+1].Low
			for j := i - basePeriod + 2; j <= i; j++ {
				if data[j].High > high {
					high = data[j].High
				}
				if data[j].Low < low {
					low = data[j].Low
				}
			}
			result[i].KijunSen = (high + low) / 2
		}

		// Senkou Span A (Leading Span A)
		if i >= basePeriod-1 {
			result[i].SenkouSpanA = (result[i].TenkanSen + result[i].KijunSen) / 2
		}

		// Senkou Span B (Leading Span B)
		if i >= spanBPeriod-1 {
			high := data[i-spanBPeriod+1].High
			low := data[i-spanBPeriod+1].Low
			for j := i - spanBPeriod + 2; j <= i; j++ {
				if data[j].High > high {
					high = data[j].High
				}
				if data[j].Low < low {
					low = data[j].Low
				}
			}
			result[i].SenkouSpanB = (high + low) / 2
		}

		// Chikou Span (Lagging Span)
		result[i].ChikouSpan = data[i].Close
	}

	return result
}
