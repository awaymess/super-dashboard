'use client';

import { motion } from 'framer-motion';
import { Zap, TrendingUp, Percent } from 'lucide-react';
import { GlassCard, Badge } from '@/components/ui';
import type { ValueBet } from '@/types/betting';

interface ValueBetCardProps {
  valueBet: ValueBet;
  onClick?: () => void;
}

export function ValueBetCard({ valueBet, onClick }: ValueBetCardProps) {
  const edgePercent = ((valueBet.value - 1) * 100).toFixed(1);

  return (
    <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
      <GlassCard className="cursor-pointer border-l-4 border-l-success" onClick={onClick}>
        <div className="flex items-start justify-between mb-3">
          <div className="flex items-center gap-2">
            <Zap className="w-5 h-5 text-warning" />
            <Badge variant="success">+{edgePercent}% Edge</Badge>
          </div>
          <span className="text-2xl font-bold text-success">{valueBet.odds.toFixed(2)}</span>
        </div>

        <h3 className="font-semibold text-white mb-1">{valueBet.selection}</h3>
        <p className="text-sm text-gray-400 mb-4">{valueBet.match}</p>

        <div className="grid grid-cols-3 gap-4 pt-4 border-t border-white/10">
          <div className="text-center">
            <p className="text-xs text-gray-400 mb-1">True Prob</p>
            <p className="font-bold text-primary">{(valueBet.probability * 100).toFixed(1)}%</p>
          </div>
          <div className="text-center">
            <p className="text-xs text-gray-400 mb-1">Kelly Stake</p>
            <p className="font-bold text-secondary">{valueBet.kellyStake.toFixed(1)}%</p>
          </div>
          <div className="text-center">
            <p className="text-xs text-gray-400 mb-1">Confidence</p>
            <p className="font-bold text-warning">{valueBet.confidence}</p>
          </div>
        </div>
      </GlassCard>
    </motion.div>
  );
}
