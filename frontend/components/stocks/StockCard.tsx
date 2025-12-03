'use client';

import { motion } from 'framer-motion';
import { TrendingUp, TrendingDown } from 'lucide-react';
import { GlassCard, Badge } from '@/components/ui';
import { SparkLine } from '@/components/charts';
import type { Stock } from '@/types/stocks';

interface StockCardProps {
  stock: Stock;
  onClick?: () => void;
}

export function StockCard({ stock, onClick }: StockCardProps) {
  const isPositive = stock.change >= 0;
  const priceHistory = stock.priceHistory || [];

  return (
    <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
      <GlassCard className="cursor-pointer" onClick={onClick}>
        <div className="flex items-start justify-between mb-3">
          <div>
            <h3 className="font-bold text-white">{stock.symbol}</h3>
            <p className="text-sm text-gray-400">{stock.name}</p>
          </div>
          <Badge variant={isPositive ? 'success' : 'danger'}>
            {stock.sector}
          </Badge>
        </div>

        <div className="flex items-end justify-between">
          <div>
            <p className="text-2xl font-bold text-white">${stock.price.toFixed(2)}</p>
            <div className={`flex items-center gap-1 text-sm ${isPositive ? 'text-success' : 'text-danger'}`}>
              {isPositive ? <TrendingUp className="w-4 h-4" /> : <TrendingDown className="w-4 h-4" />}
              <span>{isPositive ? '+' : ''}{stock.change.toFixed(2)}</span>
              <span>({isPositive ? '+' : ''}{stock.changePercent.toFixed(2)}%)</span>
            </div>
          </div>

          {priceHistory.length > 0 && (
            <SparkLine data={priceHistory} height={40} width={80} />
          )}
        </div>

        <div className="grid grid-cols-3 gap-4 mt-4 pt-4 border-t border-white/10 text-center">
          <div>
            <p className="text-xs text-gray-400">Volume</p>
            <p className="font-medium text-white text-sm">
              {(stock.volume / 1e6).toFixed(2)}M
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-400">Market Cap</p>
            <p className="font-medium text-white text-sm">
              {(stock.marketCap / 1e9).toFixed(2)}B
            </p>
          </div>
          <div>
            <p className="text-xs text-gray-400">P/E</p>
            <p className="font-medium text-white text-sm">
              {stock.peRatio?.toFixed(2) || 'N/A'}
            </p>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );
}
