'use client';

import { useState } from 'react';
import { GlassCard, GlassInput } from '@/components/ui';
import { calculatePoissonProbabilities, calculateGoalProbabilities } from '@/lib/calculations';

export function PoissonCalculator() {
  const [homeGoals, setHomeGoals] = useState('1.5');
  const [awayGoals, setAwayGoals] = useState('1.2');

  const homeExpected = parseFloat(homeGoals) || 0;
  const awayExpected = parseFloat(awayGoals) || 0;

  const homeProbs = calculatePoissonProbabilities(homeExpected, 5);
  const awayProbs = calculatePoissonProbabilities(awayExpected, 5);
  const goalProbs = calculateGoalProbabilities(homeExpected, awayExpected);

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Poisson Goal Predictor</h3>
      
      <div className="grid grid-cols-2 gap-4 mb-6">
        <div>
          <label className="text-sm text-gray-400 mb-2 block">Home xG</label>
          <GlassInput
            type="number"
            value={homeGoals}
            onChange={e => setHomeGoals(e.target.value)}
            min={0}
            step={0.1}
          />
        </div>
        <div>
          <label className="text-sm text-gray-400 mb-2 block">Away xG</label>
          <GlassInput
            type="number"
            value={awayGoals}
            onChange={e => setAwayGoals(e.target.value)}
            min={0}
            step={0.1}
          />
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 mb-6">
        <div className="space-y-2">
          <p className="text-sm text-gray-400">Home Goals Distribution</p>
          {homeProbs.map((prob, i) => (
            <div key={i} className="flex items-center gap-2">
              <span className="w-6 text-xs text-gray-400">{i}</span>
              <div className="flex-1 h-4 bg-white/5 rounded overflow-hidden">
                <div 
                  className="h-full bg-primary transition-all"
                  style={{ width: `${prob * 100}%` }}
                />
              </div>
              <span className="text-xs text-gray-300 w-12 text-right">{(prob * 100).toFixed(1)}%</span>
            </div>
          ))}
        </div>
        <div className="space-y-2">
          <p className="text-sm text-gray-400">Away Goals Distribution</p>
          {awayProbs.map((prob, i) => (
            <div key={i} className="flex items-center gap-2">
              <span className="w-6 text-xs text-gray-400">{i}</span>
              <div className="flex-1 h-4 bg-white/5 rounded overflow-hidden">
                <div 
                  className="h-full bg-secondary transition-all"
                  style={{ width: `${prob * 100}%` }}
                />
              </div>
              <span className="text-xs text-gray-300 w-12 text-right">{(prob * 100).toFixed(1)}%</span>
            </div>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-3 gap-4 pt-4 border-t border-white/10">
        <div className="text-center p-3 bg-primary/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Home Win</p>
          <p className="text-xl font-bold text-primary">{(goalProbs.homeWin * 100).toFixed(1)}%</p>
        </div>
        <div className="text-center p-3 bg-white/5 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Draw</p>
          <p className="text-xl font-bold text-gray-300">{(goalProbs.draw * 100).toFixed(1)}%</p>
        </div>
        <div className="text-center p-3 bg-secondary/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Away Win</p>
          <p className="text-xl font-bold text-secondary">{(goalProbs.awayWin * 100).toFixed(1)}%</p>
        </div>
      </div>

      <div className="grid grid-cols-2 gap-4 mt-4 pt-4 border-t border-white/10">
        <div className="text-center p-3 bg-success/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Over 2.5</p>
          <p className="text-xl font-bold text-success">{(goalProbs.over25 * 100).toFixed(1)}%</p>
        </div>
        <div className="text-center p-3 bg-danger/10 rounded-lg">
          <p className="text-xs text-gray-400 mb-1">Under 2.5</p>
          <p className="text-xl font-bold text-danger">{(goalProbs.under25 * 100).toFixed(1)}%</p>
        </div>
      </div>
    </GlassCard>
  );
}
