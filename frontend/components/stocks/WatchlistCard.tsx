'use client';

import { Star, TrendingUp, TrendingDown, MoreVertical } from 'lucide-react';
import { GlassCard, Badge, Dropdown } from '@/components/ui';
import type { Stock } from '@/types/stocks';

interface WatchlistCardProps {
  stocks: Stock[];
  onRemove?: (symbol: string) => void;
  onView?: (symbol: string) => void;
}

export function WatchlistCard({ stocks, onRemove, onView }: WatchlistCardProps) {
  return (
    <GlassCard>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Star className="w-5 h-5 text-warning fill-warning" />
          <h3 className="font-bold text-white">Watchlist</h3>
        </div>
        <Badge>{stocks.length} stocks</Badge>
      </div>

      <div className="space-y-3">
        {stocks.map(stock => {
          const isPositive = stock.change >= 0;
          return (
            <div
              key={stock.symbol}
              className="flex items-center justify-between p-3 bg-white/5 rounded-lg hover:bg-white/10 transition-colors cursor-pointer"
              onClick={() => onView?.(stock.symbol)}
            >
              <div className="flex items-center gap-3">
                <div>
                  <p className="font-semibold text-white">{stock.symbol}</p>
                  <p className="text-xs text-gray-400">{stock.name}</p>
                </div>
              </div>
              
              <div className="flex items-center gap-4">
                <div className="text-right min-w-[100px]">
                  <p className="font-semibold text-white">${stock.price.toFixed(2)}</p>
                  <p className={`text-xs flex items-center justify-end gap-1 ${isPositive ? 'text-success' : 'text-danger'}`}>
                    {isPositive ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
                    {isPositive ? '+' : ''}{stock.changePercent.toFixed(2)}%
                  </p>
                </div>
              </div>
            </div>
          );
        })}
      </div>
    </GlassCard>
  );
}
