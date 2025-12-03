'use client';

import { useState } from 'react';
import { GlassCard, GlassInput, GlassButton, Progress } from '@/components/ui';
import { calculateKelly, calculateFullKelly, calculateHalfKelly, calculateQuarterKelly } from '@/lib/calculations';

export function KellyCalculator() {
  const [probability, setProbability] = useState('55');
  const [odds, setOdds] = useState('2.00');

  const prob = parseFloat(probability) / 100 || 0;
  const oddsNum = parseFloat(odds) || 0;

  const fullKelly = calculateFullKelly(prob, oddsNum);
  const halfKelly = calculateHalfKelly(prob, oddsNum);
  const quarterKelly = calculateQuarterKelly(prob, oddsNum);
  const edge = calculateKelly(prob, oddsNum).edge;

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Kelly Criterion Calculator</h3>
      
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div>
          <label className="text-sm text-gray-400 mb-2 block">Win Probability (%)</label>
          <GlassInput
            type="number"
            value={probability}
            onChange={e => setProbability(e.target.value)}
            min={0}
            max={100}
          />
        </div>
        <div>
          <label className="text-sm text-gray-400 mb-2 block">Decimal Odds</label>
          <GlassInput
            type="number"
            value={odds}
            onChange={e => setOdds(e.target.value)}
            min={1}
            step={0.01}
          />
        </div>
      </div>

      <div className="mb-6">
        <div className="flex justify-between text-sm mb-2">
          <span className="text-gray-400">Edge</span>
          <span className={edge > 0 ? 'text-success' : 'text-danger'}>
            {edge > 0 ? '+' : ''}{(edge * 100).toFixed(2)}%
          </span>
        </div>
        <Progress value={Math.max(0, edge * 100)} max={50} color={edge > 0 ? 'success' : 'danger'} />
      </div>

      <div className="space-y-4">
        <div className="flex justify-between items-center p-3 bg-white/5 rounded-lg">
          <div>
            <p className="font-medium text-white">Full Kelly</p>
            <p className="text-xs text-gray-400">Maximum growth, high variance</p>
          </div>
          <span className="text-xl font-bold text-primary">{(fullKelly * 100).toFixed(2)}%</span>
        </div>
        
        <div className="flex justify-between items-center p-3 bg-success/10 rounded-lg border border-success/30">
          <div>
            <p className="font-medium text-white">Half Kelly</p>
            <p className="text-xs text-gray-400">Recommended for most bettors</p>
          </div>
          <span className="text-xl font-bold text-success">{(halfKelly * 100).toFixed(2)}%</span>
        </div>
        
        <div className="flex justify-between items-center p-3 bg-white/5 rounded-lg">
          <div>
            <p className="font-medium text-white">Quarter Kelly</p>
            <p className="text-xs text-gray-400">Conservative approach</p>
          </div>
          <span className="text-xl font-bold text-secondary">{(quarterKelly * 100).toFixed(2)}%</span>
        </div>
      </div>
    </GlassCard>
  );
}
