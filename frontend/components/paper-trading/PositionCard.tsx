'use client';

import { TrendingUp, TrendingDown, X } from 'lucide-react';
import { GlassCard, GlassButton, Badge } from '@/components/ui';
import type { Position } from '@/types/paper-trading';

interface PositionCardProps {
  position: Position;
  onClose?: () => void;
}

export function PositionCard({ position, onClose }: PositionCardProps) {
  const isPositive = position.pnl >= 0;
  const pnlPercent = (position.pnl / (position.entryPrice * position.quantity)) * 100;

  return (
    <GlassCard>
      <div className="flex items-start justify-between mb-4">
        <div>
          <div className="flex items-center gap-2">
            <h4 className="font-bold text-white">{position.symbol}</h4>
            <Badge variant={position.side === 'long' ? 'success' : 'danger'}>
              {position.side.toUpperCase()}
            </Badge>
          </div>
          <p className="text-sm text-gray-400">{position.quantity} shares</p>
        </div>
        <GlassButton size="sm" variant="ghost" onClick={onClose}>
          <X className="w-4 h-4" />
        </GlassButton>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-4">
        <div>
          <p className="text-xs text-gray-400 mb-1">Entry Price</p>
          <p className="font-semibold text-white">${position.entryPrice.toFixed(2)}</p>
        </div>
        <div>
          <p className="text-xs text-gray-400 mb-1">Current Price</p>
          <p className="font-semibold text-white">${position.currentPrice.toFixed(2)}</p>
        </div>
      </div>

      <div className={`p-3 rounded-lg ${isPositive ? 'bg-success/10' : 'bg-danger/10'}`}>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-2">
            {isPositive ? (
              <TrendingUp className="w-5 h-5 text-success" />
            ) : (
              <TrendingDown className="w-5 h-5 text-danger" />
            )}
            <span className="text-gray-300">P&L</span>
          </div>
          <div className="text-right">
            <p className={`font-bold ${isPositive ? 'text-success' : 'text-danger'}`}>
              {isPositive ? '+' : ''}${position.pnl.toFixed(2)}
            </p>
            <p className={`text-sm ${isPositive ? 'text-success' : 'text-danger'}`}>
              {isPositive ? '+' : ''}{pnlPercent.toFixed(2)}%
            </p>
          </div>
        </div>
      </div>

      <div className="flex gap-2 mt-4">
        <GlassButton variant="primary" size="sm" className="flex-1">
          Add to Position
        </GlassButton>
        <GlassButton variant="danger" size="sm" className="flex-1" onClick={onClose}>
          Close Position
        </GlassButton>
      </div>
    </GlassCard>
  );
}
