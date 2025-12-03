'use client';

import { GlassCard, Badge, Progress } from '@/components/ui';
import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

interface TechnicalIndicatorsProps {
  rsi: number;
  macd: { value: number; signal: number; histogram: number };
  sma20: number;
  sma50: number;
  sma200: number;
  currentPrice: number;
  bollingerBands: { upper: number; middle: number; lower: number };
}

export function TechnicalIndicators({
  rsi,
  macd,
  sma20,
  sma50,
  sma200,
  currentPrice,
  bollingerBands,
}: TechnicalIndicatorsProps) {
  const getRsiSignal = (rsi: number) => {
    if (rsi >= 70) return { label: 'Overbought', variant: 'danger' as const };
    if (rsi <= 30) return { label: 'Oversold', variant: 'success' as const };
    return { label: 'Neutral', variant: 'default' as const };
  };

  const getMacdSignal = (macd: { value: number; signal: number }) => {
    if (macd.value > macd.signal) return { label: 'Bullish', variant: 'success' as const };
    if (macd.value < macd.signal) return { label: 'Bearish', variant: 'danger' as const };
    return { label: 'Neutral', variant: 'default' as const };
  };

  const getSmaSignal = (price: number, sma: number) => {
    if (price > sma) return { icon: TrendingUp, color: 'text-success' };
    if (price < sma) return { icon: TrendingDown, color: 'text-danger' };
    return { icon: Minus, color: 'text-gray-400' };
  };

  const rsiSignal = getRsiSignal(rsi);
  const macdSignal = getMacdSignal(macd);

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Technical Indicators</h3>

      <div className="space-y-6">
        {/* RSI */}
        <div>
          <div className="flex items-center justify-between mb-2">
            <span className="text-gray-400">RSI (14)</span>
            <div className="flex items-center gap-2">
              <span className="font-bold text-white">{rsi.toFixed(2)}</span>
              <Badge variant={rsiSignal.variant}>{rsiSignal.label}</Badge>
            </div>
          </div>
          <div className="relative h-4 bg-white/5 rounded-full overflow-hidden">
            <div className="absolute inset-0 flex">
              <div className="w-[30%] bg-success/30" />
              <div className="w-[40%] bg-white/10" />
              <div className="w-[30%] bg-danger/30" />
            </div>
            <div 
              className="absolute top-0 h-full w-1 bg-white rounded-full"
              style={{ left: `${rsi}%` }}
            />
          </div>
          <div className="flex justify-between text-xs text-gray-500 mt-1">
            <span>0</span>
            <span>30</span>
            <span>70</span>
            <span>100</span>
          </div>
        </div>

        {/* MACD */}
        <div className="p-4 bg-white/5 rounded-lg">
          <div className="flex items-center justify-between mb-3">
            <span className="text-gray-400">MACD</span>
            <Badge variant={macdSignal.variant}>{macdSignal.label}</Badge>
          </div>
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <p className="text-xs text-gray-400 mb-1">MACD</p>
              <p className={`font-bold ${macd.value >= 0 ? 'text-success' : 'text-danger'}`}>
                {macd.value.toFixed(4)}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-400 mb-1">Signal</p>
              <p className="font-bold text-white">{macd.signal.toFixed(4)}</p>
            </div>
            <div>
              <p className="text-xs text-gray-400 mb-1">Histogram</p>
              <p className={`font-bold ${macd.histogram >= 0 ? 'text-success' : 'text-danger'}`}>
                {macd.histogram.toFixed(4)}
              </p>
            </div>
          </div>
        </div>

        {/* Moving Averages */}
        <div>
          <p className="text-gray-400 mb-3">Moving Averages</p>
          <div className="space-y-3">
            {[
              { label: 'SMA 20', value: sma20 },
              { label: 'SMA 50', value: sma50 },
              { label: 'SMA 200', value: sma200 },
            ].map(({ label, value }) => {
              const signal = getSmaSignal(currentPrice, value);
              return (
                <div key={label} className="flex items-center justify-between">
                  <span className="text-gray-400">{label}</span>
                  <div className="flex items-center gap-2">
                    <span className="font-medium text-white">${value.toFixed(2)}</span>
                    <signal.icon className={`w-4 h-4 ${signal.color}`} />
                  </div>
                </div>
              );
            })}
          </div>
        </div>

        {/* Bollinger Bands */}
        <div className="p-4 bg-white/5 rounded-lg">
          <p className="text-gray-400 mb-3">Bollinger Bands</p>
          <div className="grid grid-cols-3 gap-4 text-center">
            <div>
              <p className="text-xs text-gray-400 mb-1">Upper</p>
              <p className="font-bold text-danger">${bollingerBands.upper.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-xs text-gray-400 mb-1">Middle</p>
              <p className="font-bold text-white">${bollingerBands.middle.toFixed(2)}</p>
            </div>
            <div>
              <p className="text-xs text-gray-400 mb-1">Lower</p>
              <p className="font-bold text-success">${bollingerBands.lower.toFixed(2)}</p>
            </div>
          </div>
        </div>
      </div>
    </GlassCard>
  );
}
