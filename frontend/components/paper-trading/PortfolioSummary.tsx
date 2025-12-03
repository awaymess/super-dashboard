'use client';

import { TrendingUp, TrendingDown, DollarSign, Percent, PieChart } from 'lucide-react';
import { GlassCard } from '@/components/ui';
import type { Portfolio } from '@/types/paper-trading';

interface PortfolioSummaryProps {
  portfolio: Portfolio;
}

export function PortfolioSummary({ portfolio }: PortfolioSummaryProps) {
  const isPositive = portfolio.totalPnL >= 0;

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-6">Portfolio Summary</h3>

      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
        <div className="p-4 bg-white/5 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            <DollarSign className="w-5 h-5 text-primary" />
            <span className="text-sm text-gray-400">Total Value</span>
          </div>
          <p className="text-2xl font-bold text-white">${portfolio.totalValue.toLocaleString()}</p>
        </div>

        <div className="p-4 bg-white/5 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            {isPositive ? (
              <TrendingUp className="w-5 h-5 text-success" />
            ) : (
              <TrendingDown className="w-5 h-5 text-danger" />
            )}
            <span className="text-sm text-gray-400">Total P&L</span>
          </div>
          <p className={`text-2xl font-bold ${isPositive ? 'text-success' : 'text-danger'}`}>
            {isPositive ? '+' : ''}${portfolio.totalPnL.toLocaleString()}
          </p>
        </div>

        <div className="p-4 bg-white/5 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            <Percent className="w-5 h-5 text-secondary" />
            <span className="text-sm text-gray-400">P&L %</span>
          </div>
          <p className={`text-2xl font-bold ${isPositive ? 'text-success' : 'text-danger'}`}>
            {isPositive ? '+' : ''}{portfolio.totalPnLPercent.toFixed(2)}%
          </p>
        </div>

        <div className="p-4 bg-white/5 rounded-lg">
          <div className="flex items-center gap-2 mb-2">
            <DollarSign className="w-5 h-5 text-warning" />
            <span className="text-sm text-gray-400">Cash</span>
          </div>
          <p className="text-2xl font-bold text-white">${portfolio.cash.toLocaleString()}</p>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4">
        <div className="text-center p-3 bg-success/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Win Rate</p>
          <p className="text-xl font-bold text-success">{portfolio.winRate.toFixed(1)}%</p>
        </div>
        <div className="text-center p-3 bg-primary/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Total Trades</p>
          <p className="text-xl font-bold text-primary">{portfolio.totalTrades}</p>
        </div>
        <div className="text-center p-3 bg-secondary/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Open Positions</p>
          <p className="text-xl font-bold text-secondary">{portfolio.openPositions}</p>
        </div>
      </div>
    </GlassCard>
  );
}
