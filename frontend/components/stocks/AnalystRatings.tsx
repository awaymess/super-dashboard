'use client';

import { GlassCard, Progress } from '@/components/ui';
import { TrendingUp, TrendingDown, Minus, Target } from 'lucide-react';

interface AnalystRatingsProps {
  buy: number;
  hold: number;
  sell: number;
  priceTarget: {
    low: number;
    average: number;
    high: number;
    current: number;
  };
}

export function AnalystRatings({ buy, hold, sell, priceTarget }: AnalystRatingsProps) {
  const total = buy + hold + sell;
  const buyPercent = (buy / total) * 100;
  const holdPercent = (hold / total) * 100;
  const sellPercent = (sell / total) * 100;

  const consensus = buyPercent > 60 ? 'Strong Buy' : 
                    buyPercent > 40 ? 'Buy' : 
                    sellPercent > 40 ? 'Sell' : 'Hold';
  
  const consensusColor = consensus.includes('Buy') ? 'text-success' : 
                         consensus.includes('Sell') ? 'text-danger' : 'text-warning';

  const upside = ((priceTarget.average - priceTarget.current) / priceTarget.current) * 100;

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Analyst Ratings</h3>

      <div className="text-center mb-6">
        <p className={`text-3xl font-bold ${consensusColor}`}>{consensus}</p>
        <p className="text-sm text-gray-400">Based on {total} analysts</p>
      </div>

      <div className="flex gap-1 mb-4 h-3 rounded-full overflow-hidden">
        <div className="bg-success" style={{ width: `${buyPercent}%` }} />
        <div className="bg-warning" style={{ width: `${holdPercent}%` }} />
        <div className="bg-danger" style={{ width: `${sellPercent}%` }} />
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="text-center">
          <div className="flex items-center justify-center gap-1 text-success mb-1">
            <TrendingUp className="w-4 h-4" />
            <span className="font-bold">{buy}</span>
          </div>
          <p className="text-xs text-gray-400">Buy</p>
        </div>
        <div className="text-center">
          <div className="flex items-center justify-center gap-1 text-warning mb-1">
            <Minus className="w-4 h-4" />
            <span className="font-bold">{hold}</span>
          </div>
          <p className="text-xs text-gray-400">Hold</p>
        </div>
        <div className="text-center">
          <div className="flex items-center justify-center gap-1 text-danger mb-1">
            <TrendingDown className="w-4 h-4" />
            <span className="font-bold">{sell}</span>
          </div>
          <p className="text-xs text-gray-400">Sell</p>
        </div>
      </div>

      <div className="pt-4 border-t border-white/10">
        <div className="flex items-center gap-2 mb-4">
          <Target className="w-5 h-5 text-primary" />
          <span className="text-gray-400">Price Target</span>
        </div>
        
        <div className="relative mb-4">
          <div className="flex justify-between text-xs text-gray-400 mb-2">
            <span>${priceTarget.low.toFixed(2)}</span>
            <span>${priceTarget.high.toFixed(2)}</span>
          </div>
          <div className="h-2 bg-white/10 rounded-full relative">
            <div 
              className="absolute top-0 h-full w-1 bg-primary rounded-full"
              style={{ 
                left: `${((priceTarget.current - priceTarget.low) / (priceTarget.high - priceTarget.low)) * 100}%` 
              }}
            />
            <div 
              className="absolute top-0 h-full w-1 bg-success rounded-full"
              style={{ 
                left: `${((priceTarget.average - priceTarget.low) / (priceTarget.high - priceTarget.low)) * 100}%` 
              }}
            />
          </div>
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div className="p-3 bg-white/5 rounded-lg text-center">
            <p className="text-xs text-gray-400 mb-1">Average Target</p>
            <p className="text-xl font-bold text-success">${priceTarget.average.toFixed(2)}</p>
          </div>
          <div className="p-3 bg-white/5 rounded-lg text-center">
            <p className="text-xs text-gray-400 mb-1">Upside</p>
            <p className={`text-xl font-bold ${upside >= 0 ? 'text-success' : 'text-danger'}`}>
              {upside >= 0 ? '+' : ''}{upside.toFixed(1)}%
            </p>
          </div>
        </div>
      </div>
    </GlassCard>
  );
}
