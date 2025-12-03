'use client';

import { TrendingUp, TrendingDown, Clock, BarChart2 } from 'lucide-react';
import { GlassCard } from '@/components/ui';
import type { Stock } from '@/types/stocks';

interface StockQuoteProps {
  stock: Stock;
  showDetails?: boolean;
}

export function StockQuote({ stock, showDetails = true }: StockQuoteProps) {
  const isPositive = stock.change >= 0;

  return (
    <GlassCard>
      <div className="flex items-start justify-between mb-6">
        <div>
          <div className="flex items-center gap-3">
            <h2 className="text-3xl font-bold text-white">{stock.symbol}</h2>
            <span className="px-3 py-1 bg-white/10 rounded-lg text-sm text-gray-400">
              {stock.exchange}
            </span>
          </div>
          <p className="text-gray-400 mt-1">{stock.name}</p>
        </div>
        <div className="text-right">
          <p className="text-4xl font-bold text-white">${stock.price.toFixed(2)}</p>
          <div className={`flex items-center justify-end gap-2 mt-1 ${isPositive ? 'text-success' : 'text-danger'}`}>
            {isPositive ? <TrendingUp className="w-5 h-5" /> : <TrendingDown className="w-5 h-5" />}
            <span className="text-lg font-semibold">
              {isPositive ? '+' : ''}{stock.change.toFixed(2)} ({isPositive ? '+' : ''}{stock.changePercent.toFixed(2)}%)
            </span>
          </div>
        </div>
      </div>

      {showDetails && (
        <>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div className="bg-white/5 p-4 rounded-lg">
              <p className="text-sm text-gray-400 mb-1">Open</p>
              <p className="text-lg font-semibold text-white">${stock.open?.toFixed(2)}</p>
            </div>
            <div className="bg-white/5 p-4 rounded-lg">
              <p className="text-sm text-gray-400 mb-1">High</p>
              <p className="text-lg font-semibold text-success">${stock.high?.toFixed(2)}</p>
            </div>
            <div className="bg-white/5 p-4 rounded-lg">
              <p className="text-sm text-gray-400 mb-1">Low</p>
              <p className="text-lg font-semibold text-danger">${stock.low?.toFixed(2)}</p>
            </div>
            <div className="bg-white/5 p-4 rounded-lg">
              <p className="text-sm text-gray-400 mb-1">Prev Close</p>
              <p className="text-lg font-semibold text-white">${stock.previousClose?.toFixed(2)}</p>
            </div>
          </div>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div className="flex items-center gap-3">
              <div className="p-2 bg-primary/20 rounded-lg">
                <BarChart2 className="w-5 h-5 text-primary" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Volume</p>
                <p className="font-semibold text-white">{(stock.volume / 1e6).toFixed(2)}M</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-secondary/20 rounded-lg">
                <TrendingUp className="w-5 h-5 text-secondary" />
              </div>
              <div>
                <p className="text-sm text-gray-400">Market Cap</p>
                <p className="font-semibold text-white">${(stock.marketCap / 1e9).toFixed(2)}B</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-success/20 rounded-lg">
                <Clock className="w-5 h-5 text-success" />
              </div>
              <div>
                <p className="text-sm text-gray-400">P/E Ratio</p>
                <p className="font-semibold text-white">{stock.peRatio?.toFixed(2) || 'N/A'}</p>
              </div>
            </div>
            <div className="flex items-center gap-3">
              <div className="p-2 bg-warning/20 rounded-lg">
                <TrendingUp className="w-5 h-5 text-warning" />
              </div>
              <div>
                <p className="text-sm text-gray-400">52W Range</p>
                <p className="font-semibold text-white">
                  ${stock.week52Low?.toFixed(0)} - ${stock.week52High?.toFixed(0)}
                </p>
              </div>
            </div>
          </div>
        </>
      )}
    </GlassCard>
  );
}
