'use client';

import { useState } from 'react';
import { Trash2, Calculator } from 'lucide-react';
import { GlassCard, GlassButton, GlassInput } from '@/components/ui';
import type { Bet } from '@/types/betting';

interface BetSlipProps {
  bets: Bet[];
  onRemoveBet: (id: string) => void;
  onPlaceBet: (stake: number) => void;
  onClear: () => void;
}

export function BetSlip({ bets, onRemoveBet, onPlaceBet, onClear }: BetSlipProps) {
  const [stake, setStake] = useState<string>('10');
  const stakeNum = parseFloat(stake) || 0;

  const totalOdds = bets.reduce((acc, bet) => acc * bet.odds, 1);
  const potentialWin = stakeNum * totalOdds;

  return (
    <GlassCard className="sticky top-20">
      <div className="flex items-center justify-between mb-4">
        <h3 className="font-bold text-white">Bet Slip</h3>
        {bets.length > 0 && (
          <button onClick={onClear} className="text-xs text-gray-400 hover:text-white transition-colors">
            Clear All
          </button>
        )}
      </div>

      {bets.length === 0 ? (
        <div className="text-center py-8">
          <Calculator className="w-12 h-12 mx-auto text-gray-600 mb-3" />
          <p className="text-gray-400">No selections added</p>
          <p className="text-xs text-gray-500 mt-1">Click on odds to add selections</p>
        </div>
      ) : (
        <>
          <div className="space-y-3 mb-4">
            {bets.map(bet => (
              <div key={bet.id} className="bg-white/5 rounded-lg p-3">
                <div className="flex items-start justify-between">
                  <div>
                    <p className="font-medium text-white text-sm">{bet.betType}</p>
                    <p className="text-xs text-gray-400">{bet.match.homeTeam.name} vs {bet.match.awayTeam.name}</p>
                  </div>
                  <div className="flex items-center gap-2">
                    <span className="font-bold text-primary">{bet.odds.toFixed(2)}</span>
                    <button
                      onClick={() => onRemoveBet(bet.id)}
                      className="text-gray-400 hover:text-danger transition-colors"
                    >
                      <Trash2 className="w-4 h-4" />
                    </button>
                  </div>
                </div>
              </div>
            ))}
          </div>

          {bets.length > 1 && (
            <div className="bg-primary/10 rounded-lg p-3 mb-4">
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">Accumulator Odds</span>
                <span className="font-bold text-primary">{totalOdds.toFixed(2)}</span>
              </div>
            </div>
          )}

          <div className="mb-4">
            <label className="text-sm text-gray-400 mb-2 block">Stake ($)</label>
            <GlassInput
              type="number"
              value={stake}
              onChange={e => setStake(e.target.value)}
              min={0}
            />
          </div>

          <div className="bg-success/10 rounded-lg p-3 mb-4">
            <div className="flex justify-between">
              <span className="text-gray-300">Potential Win</span>
              <span className="font-bold text-success text-lg">${potentialWin.toFixed(2)}</span>
            </div>
          </div>

          <GlassButton
            variant="primary"
            className="w-full"
            onClick={() => onPlaceBet(stakeNum)}
            disabled={stakeNum <= 0}
          >
            Place Bet
          </GlassButton>
        </>
      )}
    </GlassCard>
  );
}
