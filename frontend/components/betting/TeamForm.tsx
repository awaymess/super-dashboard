'use client';

import { GlassCard, Badge } from '@/components/ui';

interface TeamFormProps {
  teamName: string;
  form: Array<'W' | 'D' | 'L'>;
  stats: {
    goalsScored: number;
    goalsConceded: number;
    cleanSheets: number;
    failedToScore: number;
  };
}

export function TeamForm({ teamName, form, stats }: TeamFormProps) {
  const formColors = {
    W: 'success',
    D: 'warning',
    L: 'danger',
  } as const;

  const wins = form.filter(r => r === 'W').length;
  const draws = form.filter(r => r === 'D').length;
  const losses = form.filter(r => r === 'L').length;
  const points = wins * 3 + draws;

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">{teamName} Form</h3>
      
      <div className="flex gap-2 mb-6">
        {form.slice(0, 5).map((result, i) => (
          <Badge key={i} variant={formColors[result]} className="w-8 h-8 flex items-center justify-center">
            {result}
          </Badge>
        ))}
      </div>

      <div className="grid grid-cols-3 gap-4 mb-6 text-center">
        <div>
          <p className="text-2xl font-bold text-success">{wins}</p>
          <p className="text-xs text-gray-400">Wins</p>
        </div>
        <div>
          <p className="text-2xl font-bold text-warning">{draws}</p>
          <p className="text-xs text-gray-400">Draws</p>
        </div>
        <div>
          <p className="text-2xl font-bold text-danger">{losses}</p>
          <p className="text-xs text-gray-400">Losses</p>
        </div>
      </div>

      <div className="space-y-3 pt-4 border-t border-white/10">
        <div className="flex justify-between">
          <span className="text-gray-400">Goals Scored</span>
          <span className="font-medium text-white">{stats.goalsScored}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-400">Goals Conceded</span>
          <span className="font-medium text-white">{stats.goalsConceded}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-400">Clean Sheets</span>
          <span className="font-medium text-white">{stats.cleanSheets}</span>
        </div>
        <div className="flex justify-between">
          <span className="text-gray-400">Failed to Score</span>
          <span className="font-medium text-white">{stats.failedToScore}</span>
        </div>
        <div className="flex justify-between pt-3 border-t border-white/10">
          <span className="text-gray-400">Points (Last {form.length})</span>
          <span className="font-bold text-primary">{points}</span>
        </div>
      </div>
    </GlassCard>
  );
}
