import { TechnicalIndicators } from '@/types/stocks';

export function calculateSMA(prices: number[], period: number): number {
  if (prices.length < period) return 0;
  const slice = prices.slice(-period);
  return slice.reduce((sum, price) => sum + price, 0) / period;
}

export function calculateEMA(prices: number[], period: number): number {
  if (prices.length < period) return 0;
  
  const multiplier = 2 / (period + 1);
  let ema = calculateSMA(prices.slice(0, period), period);
  
  for (let i = period; i < prices.length; i++) {
    ema = (prices[i] - ema) * multiplier + ema;
  }
  
  return ema;
}

export function calculateRSI(prices: number[], period: number = 14): number {
  if (prices.length < period + 1) return 50;

  const changes: number[] = [];
  for (let i = 1; i < prices.length; i++) {
    changes.push(prices[i] - prices[i - 1]);
  }

  let avgGain = 0;
  let avgLoss = 0;

  for (let i = 0; i < period; i++) {
    if (changes[i] >= 0) {
      avgGain += changes[i];
    } else {
      avgLoss += Math.abs(changes[i]);
    }
  }

  avgGain /= period;
  avgLoss /= period;

  for (let i = period; i < changes.length; i++) {
    if (changes[i] >= 0) {
      avgGain = (avgGain * (period - 1) + changes[i]) / period;
      avgLoss = (avgLoss * (period - 1)) / period;
    } else {
      avgGain = (avgGain * (period - 1)) / period;
      avgLoss = (avgLoss * (period - 1) + Math.abs(changes[i])) / period;
    }
  }

  if (avgLoss === 0) return 100;
  const rs = avgGain / avgLoss;
  return 100 - (100 / (1 + rs));
}

export function calculateMACD(
  prices: number[],
  fastPeriod: number = 12,
  slowPeriod: number = 26,
  signalPeriod: number = 9
): { macd: number; signal: number; histogram: number } {
  if (prices.length < slowPeriod) {
    return { macd: 0, signal: 0, histogram: 0 };
  }

  const ema12 = calculateEMA(prices, fastPeriod);
  const ema26 = calculateEMA(prices, slowPeriod);
  const macdLine = ema12 - ema26;

  const macdHistory: number[] = [];
  for (let i = slowPeriod; i <= prices.length; i++) {
    const slice = prices.slice(0, i);
    const e12 = calculateEMA(slice, fastPeriod);
    const e26 = calculateEMA(slice, slowPeriod);
    macdHistory.push(e12 - e26);
  }

  const signal = macdHistory.length >= signalPeriod 
    ? calculateEMA(macdHistory, signalPeriod)
    : macdLine;

  return {
    macd: macdLine,
    signal,
    histogram: macdLine - signal,
  };
}

export function calculateBollingerBands(
  prices: number[],
  period: number = 20,
  stdDev: number = 2
): { upper: number; middle: number; lower: number } {
  if (prices.length < period) {
    const lastPrice = prices[prices.length - 1] || 0;
    return { upper: lastPrice, middle: lastPrice, lower: lastPrice };
  }

  const slice = prices.slice(-period);
  const sma = slice.reduce((sum, p) => sum + p, 0) / period;
  
  const squaredDiffs = slice.map(p => Math.pow(p - sma, 2));
  const variance = squaredDiffs.reduce((sum, d) => sum + d, 0) / period;
  const std = Math.sqrt(variance);

  return {
    upper: sma + stdDev * std,
    middle: sma,
    lower: sma - stdDev * std,
  };
}

export function calculateATR(
  highs: number[],
  lows: number[],
  closes: number[],
  period: number = 14
): number {
  if (highs.length < period + 1) return 0;

  const trueRanges: number[] = [];
  
  for (let i = 1; i < highs.length; i++) {
    const tr = Math.max(
      highs[i] - lows[i],
      Math.abs(highs[i] - closes[i - 1]),
      Math.abs(lows[i] - closes[i - 1])
    );
    trueRanges.push(tr);
  }

  return calculateSMA(trueRanges, period);
}

export function calculateStochastic(
  highs: number[],
  lows: number[],
  closes: number[],
  kPeriod: number = 14,
  dPeriod: number = 3
): { k: number; d: number } {
  if (highs.length < kPeriod) return { k: 50, d: 50 };

  const kValues: number[] = [];
  
  for (let i = kPeriod - 1; i < closes.length; i++) {
    const periodHighs = highs.slice(i - kPeriod + 1, i + 1);
    const periodLows = lows.slice(i - kPeriod + 1, i + 1);
    const highest = Math.max(...periodHighs);
    const lowest = Math.min(...periodLows);
    
    const k = highest !== lowest 
      ? ((closes[i] - lowest) / (highest - lowest)) * 100 
      : 50;
    kValues.push(k);
  }

  const k = kValues[kValues.length - 1];
  const d = kValues.length >= dPeriod 
    ? calculateSMA(kValues.slice(-dPeriod), dPeriod)
    : k;

  return { k, d };
}

export function calculateAllIndicators(
  prices: number[],
  highs: number[],
  lows: number[],
  closes: number[]
): TechnicalIndicators {
  const rsi = calculateRSI(closes);
  const macdData = calculateMACD(closes);
  const bollinger = calculateBollingerBands(closes);
  const stoch = calculateStochastic(highs, lows, closes);
  const currentPrice = closes[closes.length - 1];

  let rsiSignal: 'oversold' | 'neutral' | 'overbought' = 'neutral';
  if (rsi < 30) rsiSignal = 'oversold';
  else if (rsi > 70) rsiSignal = 'overbought';

  let macdTrend: 'bullish' | 'bearish' | 'neutral' = 'neutral';
  if (macdData.histogram > 0) macdTrend = 'bullish';
  else if (macdData.histogram < 0) macdTrend = 'bearish';

  let bollingerPosition: 'above' | 'within' | 'below' = 'within';
  if (currentPrice > bollinger.upper) bollingerPosition = 'above';
  else if (currentPrice < bollinger.lower) bollingerPosition = 'below';

  return {
    rsi,
    rsiSignal,
    macd: {
      value: macdData.macd,
      signal: macdData.signal,
      histogram: macdData.histogram,
      trend: macdTrend,
    },
    bollingerBands: {
      upper: bollinger.upper,
      middle: bollinger.middle,
      lower: bollinger.lower,
      position: bollingerPosition,
    },
    sma20: calculateSMA(closes, 20),
    sma50: calculateSMA(closes, 50),
    sma200: calculateSMA(closes, 200),
    ema12: calculateEMA(closes, 12),
    ema26: calculateEMA(closes, 26),
    atr: calculateATR(highs, lows, closes),
    stochastic: stoch,
  };
}

export function getTechnicalSignal(indicators: TechnicalIndicators): {
  signal: 'buy' | 'sell' | 'hold';
  strength: number;
  reasons: string[];
} {
  let bullishSignals = 0;
  let bearishSignals = 0;
  const reasons: string[] = [];

  if (indicators.rsi < 30) {
    bullishSignals += 2;
    reasons.push('RSI oversold');
  } else if (indicators.rsi > 70) {
    bearishSignals += 2;
    reasons.push('RSI overbought');
  }

  if (indicators.macd.histogram > 0 && indicators.macd.value > indicators.macd.signal) {
    bullishSignals += 2;
    reasons.push('MACD bullish crossover');
  } else if (indicators.macd.histogram < 0 && indicators.macd.value < indicators.macd.signal) {
    bearishSignals += 2;
    reasons.push('MACD bearish crossover');
  }

  if (indicators.bollingerBands.position === 'below') {
    bullishSignals += 1;
    reasons.push('Price below lower Bollinger Band');
  } else if (indicators.bollingerBands.position === 'above') {
    bearishSignals += 1;
    reasons.push('Price above upper Bollinger Band');
  }

  if (indicators.stochastic.k < 20) {
    bullishSignals += 1;
    reasons.push('Stochastic oversold');
  } else if (indicators.stochastic.k > 80) {
    bearishSignals += 1;
    reasons.push('Stochastic overbought');
  }

  const totalSignals = bullishSignals + bearishSignals;
  const strength = totalSignals > 0 ? Math.abs(bullishSignals - bearishSignals) / totalSignals : 0;

  let signal: 'buy' | 'sell' | 'hold' = 'hold';
  if (bullishSignals > bearishSignals + 1) signal = 'buy';
  else if (bearishSignals > bullishSignals + 1) signal = 'sell';

  return { signal, strength: strength * 100, reasons };
}

/**
 * Calculate Fibonacci Retracement levels
 * @param high - Swing high price
 * @param low - Swing low price
 * @param isUptrend - Whether the trend is up (retracement from high) or down (retracement from low)
 * @returns Fibonacci retracement levels
 */
export function calculateFibonacciRetracement(
  high: number,
  low: number,
  isUptrend: boolean = true
): {
  level0: number;
  level236: number;
  level382: number;
  level500: number;
  level618: number;
  level786: number;
  level1000: number;
  extension1272: number;
  extension1618: number;
} {
  const diff = high - low;
  
  if (isUptrend) {
    // Retracement from high to low (price pulled back from high)
    // Extensions project above the high (beyond the original move)
    return {
      level0: high,
      level236: high - diff * 0.236,
      level382: high - diff * 0.382,
      level500: high - diff * 0.500,
      level618: high - diff * 0.618,
      level786: high - diff * 0.786,
      level1000: low,
      // Extensions: 127.2% = high + (diff * 0.272), 161.8% = high + (diff * 0.618)
      extension1272: high + diff * 0.272,
      extension1618: high + diff * 0.618,
    };
  } else {
    // Retracement from low to high (price bounced from low in downtrend)
    // Extensions project below the low (beyond the original move down)
    return {
      level0: low,
      level236: low + diff * 0.236,
      level382: low + diff * 0.382,
      level500: low + diff * 0.500,
      level618: low + diff * 0.618,
      level786: low + diff * 0.786,
      level1000: high,
      // Extensions project below low for downtrend continuation
      extension1272: low - diff * 0.272,
      extension1618: low - diff * 0.618,
    };
  }
}

/**
 * Detect support and resistance levels using volume profile
 * @param prices - Array of close prices
 * @param volumes - Array of volumes
 * @param bins - Number of price bins
 * @returns Volume profile with support/resistance levels
 */
export function calculateVolumeProfile(
  prices: number[],
  volumes: number[],
  bins: number = 20
): {
  levels: { price: number; volume: number; type: 'support' | 'resistance' | 'neutral' }[];
  pointOfControl: number;
  valueAreaHigh: number;
  valueAreaLow: number;
} {
  if (prices.length === 0 || prices.length !== volumes.length) {
    return {
      levels: [],
      pointOfControl: 0,
      valueAreaHigh: 0,
      valueAreaLow: 0,
    };
  }

  const minPrice = Math.min(...prices);
  const maxPrice = Math.max(...prices);
  const binSize = (maxPrice - minPrice) / bins;
  
  // Initialize bins
  const volumeByPrice: { price: number; volume: number }[] = [];
  for (let i = 0; i < bins; i++) {
    volumeByPrice.push({
      price: minPrice + (i + 0.5) * binSize,
      volume: 0,
    });
  }
  
  // Aggregate volume by price level
  for (let i = 0; i < prices.length; i++) {
    const binIndex = Math.min(
      Math.floor((prices[i] - minPrice) / binSize),
      bins - 1
    );
    volumeByPrice[binIndex].volume += volumes[i];
  }
  
  // Find Point of Control (highest volume level)
  const poc = volumeByPrice.reduce((max, level) => 
    level.volume > max.volume ? level : max
  , volumeByPrice[0]);
  
  // Calculate value area (70% of volume)
  const totalVolume = volumeByPrice.reduce((sum, l) => sum + l.volume, 0);
  const valueAreaThreshold = totalVolume * 0.7;
  
  const sortedByVolume = [...volumeByPrice].sort((a, b) => b.volume - a.volume);
  let cumulativeVolume = 0;
  const valueAreaLevels: number[] = [];
  
  for (const level of sortedByVolume) {
    if (cumulativeVolume >= valueAreaThreshold) break;
    valueAreaLevels.push(level.price);
    cumulativeVolume += level.volume;
  }
  
  const valueAreaHigh = Math.max(...valueAreaLevels);
  const valueAreaLow = Math.min(...valueAreaLevels);
  const currentPrice = prices[prices.length - 1];
  
  // Determine support/resistance for each level
  const levels = volumeByPrice.map(level => ({
    ...level,
    type: level.price < currentPrice ? 'support' as const : 
          level.price > currentPrice ? 'resistance' as const : 'neutral' as const,
  }));
  
  return {
    levels,
    pointOfControl: poc.price,
    valueAreaHigh,
    valueAreaLow,
  };
}

/**
 * Detect volume spike (unusual volume)
 * @param currentVolume - Current volume
 * @param avgVolume - Average volume
 * @param threshold - Spike threshold multiplier (default 2x)
 * @returns Volume spike analysis
 */
export function detectVolumeSpike(
  currentVolume: number,
  avgVolume: number,
  threshold: number = 2
): {
  isSpike: boolean;
  multiplier: number;
  significance: 'normal' | 'elevated' | 'high' | 'extreme';
} {
  const multiplier = avgVolume > 0 ? currentVolume / avgVolume : 0;
  
  let significance: 'normal' | 'elevated' | 'high' | 'extreme' = 'normal';
  if (multiplier >= 4) significance = 'extreme';
  else if (multiplier >= 3) significance = 'high';
  else if (multiplier >= 2) significance = 'elevated';
  
  return {
    isSpike: multiplier >= threshold,
    multiplier,
    significance,
  };
}

/**
 * Detect breakout or breakdown patterns
 * @param prices - Array of close prices
 * @param period - Lookback period for range
 * @returns Breakout/breakdown detection
 */
export function detectBreakout(
  prices: number[],
  period: number = 20
): {
  type: 'breakout' | 'breakdown' | 'none';
  level: number;
  strength: number;
  confirmed: boolean;
} {
  if (prices.length < period + 1) {
    return { type: 'none', level: 0, strength: 0, confirmed: false };
  }
  
  const range = prices.slice(-period - 1, -1);
  const high = Math.max(...range);
  const low = Math.min(...range);
  const currentPrice = prices[prices.length - 1];
  const previousPrice = prices[prices.length - 2];
  
  const rangeSize = high - low;
  
  if (currentPrice > high) {
    const strength = rangeSize > 0 ? ((currentPrice - high) / rangeSize) * 100 : 0;
    return {
      type: 'breakout',
      level: high,
      strength: Math.min(strength, 100),
      confirmed: previousPrice > high, // Confirmed if previous candle also above
    };
  }
  
  if (currentPrice < low) {
    const strength = rangeSize > 0 ? ((low - currentPrice) / rangeSize) * 100 : 0;
    return {
      type: 'breakdown',
      level: low,
      strength: Math.min(strength, 100),
      confirmed: previousPrice < low,
    };
  }
  
  return { type: 'none', level: 0, strength: 0, confirmed: false };
}

/**
 * Detect divergence between price and indicator
 * @param prices - Array of prices
 * @param indicator - Array of indicator values (e.g., RSI, MACD)
 * @param lookback - Lookback period for divergence detection
 * @returns Divergence detection result
 */
export function detectDivergence(
  prices: number[],
  indicator: number[],
  lookback: number = 14
): {
  type: 'bullish' | 'bearish' | 'none';
  strength: 'weak' | 'moderate' | 'strong';
  description: string;
} {
  if (prices.length < lookback || indicator.length < lookback) {
    return { type: 'none', strength: 'weak', description: 'Insufficient data' };
  }
  
  const recentPrices = prices.slice(-lookback);
  const recentIndicator = indicator.slice(-lookback);
  
  const priceStart = recentPrices[0];
  const priceEnd = recentPrices[recentPrices.length - 1];
  const indicatorStart = recentIndicator[0];
  const indicatorEnd = recentIndicator[recentIndicator.length - 1];
  
  const priceChange = (priceEnd - priceStart) / priceStart;
  const indicatorChange = indicatorStart !== 0 ? (indicatorEnd - indicatorStart) / Math.abs(indicatorStart) : 0;
  
  // Bullish divergence: price making lower lows, indicator making higher lows
  if (priceChange < -0.02 && indicatorChange > 0.05) {
    const strength = Math.abs(indicatorChange) > 0.15 ? 'strong' : 
                     Math.abs(indicatorChange) > 0.10 ? 'moderate' : 'weak';
    return {
      type: 'bullish',
      strength,
      description: 'Price falling while indicator rising - potential reversal up',
    };
  }
  
  // Bearish divergence: price making higher highs, indicator making lower highs
  if (priceChange > 0.02 && indicatorChange < -0.05) {
    const strength = Math.abs(indicatorChange) > 0.15 ? 'strong' : 
                     Math.abs(indicatorChange) > 0.10 ? 'moderate' : 'weak';
    return {
      type: 'bearish',
      strength,
      description: 'Price rising while indicator falling - potential reversal down',
    };
  }
  
  return { type: 'none', strength: 'weak', description: 'No divergence detected' };
}

/**
 * Detect Golden Cross or Death Cross
 * @param prices - Array of close prices
 * @returns Cross detection result
 */
export function detectMovingAverageCross(
  prices: number[]
): {
  type: 'golden_cross' | 'death_cross' | 'none';
  ma50: number;
  ma200: number;
  daysAgo: number;
} {
  if (prices.length < 200) {
    return { type: 'none', ma50: 0, ma200: 0, daysAgo: -1 };
  }
  
  const ma50Current = calculateSMA(prices, 50);
  const ma200Current = calculateSMA(prices, 200);
  
  // Check for recent cross (within last 5 days)
  for (let i = 1; i <= 5; i++) {
    const pastPrices = prices.slice(0, -i);
    if (pastPrices.length < 200) continue;
    
    const ma50Past = calculateSMA(pastPrices, 50);
    const ma200Past = calculateSMA(pastPrices, 200);
    
    // Golden Cross: MA50 crosses above MA200
    if (ma50Current > ma200Current && ma50Past <= ma200Past) {
      return {
        type: 'golden_cross',
        ma50: ma50Current,
        ma200: ma200Current,
        daysAgo: i - 1,
      };
    }
    
    // Death Cross: MA50 crosses below MA200
    if (ma50Current < ma200Current && ma50Past >= ma200Past) {
      return {
        type: 'death_cross',
        ma50: ma50Current,
        ma200: ma200Current,
        daysAgo: i - 1,
      };
    }
  }
  
  return {
    type: 'none',
    ma50: ma50Current,
    ma200: ma200Current,
    daysAgo: -1,
  };
}
