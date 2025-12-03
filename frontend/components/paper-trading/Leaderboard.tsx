'use client';

import { Trophy, TrendingUp, Medal } from 'lucide-react';
import { GlassCard, Avatar, Badge } from '@/components/ui';
import type { LeaderboardEntry } from '@/types/paper-trading';

interface LeaderboardProps {
  entries: LeaderboardEntry[];
  currentUser?: string;
}

export function Leaderboard({ entries, currentUser }: LeaderboardProps) {
  const getRankIcon = (rank: number) => {
    switch (rank) {
      case 1:
        return <Trophy className="w-5 h-5 text-yellow-400" />;
      case 2:
        return <Medal className="w-5 h-5 text-gray-300" />;
      case 3:
        return <Medal className="w-5 h-5 text-amber-600" />;
      default:
        return <span className="w-5 text-center text-gray-400">{rank}</span>;
    }
  };

  return (
    <GlassCard>
      <div className="flex items-center justify-between mb-4">
        <div className="flex items-center gap-2">
          <Trophy className="w-5 h-5 text-warning" />
          <h3 className="font-bold text-white">Leaderboard</h3>
        </div>
        <Badge>This Month</Badge>
      </div>

      <div className="space-y-2">
        {entries.map((entry) => {
          const isCurrentUser = entry.username === currentUser;
          const isPositive = entry.totalReturn >= 0;
          
          return (
            <div
              key={entry.rank}
              className={`flex items-center gap-3 p-3 rounded-lg transition-colors ${
                isCurrentUser ? 'bg-primary/20 border border-primary/30' : 'bg-white/5 hover:bg-white/10'
              }`}
            >
              <div className="w-8 flex justify-center">
                {getRankIcon(entry.rank)}
              </div>
              
              <Avatar name={entry.username} src={entry.avatar} size="sm" />
              
              <div className="flex-1 min-w-0">
                <p className={`font-medium truncate ${isCurrentUser ? 'text-primary' : 'text-white'}`}>
                  {entry.username}
                </p>
                <p className="text-xs text-gray-400">{entry.totalTrades} trades</p>
              </div>
              
              <div className="text-right">
                <p className={`font-bold ${isPositive ? 'text-success' : 'text-danger'}`}>
                  {isPositive ? '+' : ''}{entry.totalReturnPercent.toFixed(2)}%
                </p>
                <p className="text-xs text-gray-400">{entry.winRate.toFixed(0)}% win rate</p>
              </div>
            </div>
          );
        })}
      </div>
    </GlassCard>
  );
}
