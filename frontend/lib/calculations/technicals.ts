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
