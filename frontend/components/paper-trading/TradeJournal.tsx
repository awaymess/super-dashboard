'use client';

import { GlassCard, Badge } from '@/components/ui';
import { Calendar, TrendingUp, TrendingDown, FileText } from 'lucide-react';
import type { Trade } from '@/types/paper-trading';

interface TradeJournalProps {
  trades: Trade[];
}

export function TradeJournal({ trades }: TradeJournalProps) {
  return (
    <GlassCard>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <FileText className="w-5 h-5 text-primary" />
          <h3 className="font-bold text-white">Trade Journal</h3>
        </div>
        <Badge>{trades.length} trades</Badge>
      </div>

      <div className="space-y-3">
        {trades.map((trade) => {
          const isPositive = trade.pnl && trade.pnl >= 0;
          return (
            <div
              key={trade.id}
              className="p-4 bg-white/5 rounded-lg hover:bg-white/10 transition-colors cursor-pointer"
            >
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center gap-2">
                  <span className="font-bold text-white">{trade.symbol}</span>
                  <Badge variant={trade.side === 'buy' ? 'success' : 'danger'}>
                    {trade.side.toUpperCase()}
                  </Badge>
                  <Badge variant="default">{trade.orderType}</Badge>
                </div>
                {trade.pnl !== undefined && (
                  <div className={`flex items-center gap-1 ${isPositive ? 'text-success' : 'text-danger'}`}>
                    {isPositive ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />}
                    <span className="font-bold">
                      {isPositive ? '+' : ''}${trade.pnl.toFixed(2)}
                    </span>
                  </div>
                )}
              </div>
              
              <div className="grid grid-cols-3 gap-4 text-sm">
                <div>
                  <p className="text-gray-400">Quantity</p>
                  <p className="text-white">{trade.quantity}</p>
                </div>
                <div>
                  <p className="text-gray-400">Price</p>
                  <p className="text-white">${trade.price.toFixed(2)}</p>
                </div>
                <div>
                  <p className="text-gray-400">Total</p>
                  <p className="text-white">${(trade.quantity * trade.price).toFixed(2)}</p>
                </div>
              </div>

              {trade.notes && (
                <p className="text-sm text-gray-400 mt-3 pt-3 border-t border-white/10">
                  {trade.notes}
                </p>
              )}

              <div className="flex items-center gap-1 text-xs text-gray-500 mt-3">
                <Calendar className="w-3 h-3" />
                {trade.timestamp}
              </div>
            </div>
          );
        })}
      </div>
    </GlassCard>
  );
}
