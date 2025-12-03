'use client';

import { motion } from 'framer-motion';
import { Clock, Trophy, TrendingUp } from 'lucide-react';
import { GlassCard, Badge } from '@/components/ui';
import type { Match } from '@/types/betting';

interface MatchCardProps {
  match: Match;
  onClick?: () => void;
}

export function MatchCard({ match, onClick }: MatchCardProps) {
  const statusColors = {
    upcoming: 'warning',
    live: 'danger',
    finished: 'default',
  } as const;

  return (
    <motion.div whileHover={{ scale: 1.02 }} whileTap={{ scale: 0.98 }}>
      <GlassCard className="cursor-pointer" onClick={onClick}>
        <div className="flex items-center justify-between mb-4">
          <Badge variant={statusColors[match.status]}>{match.status.toUpperCase()}</Badge>
          <span className="text-xs text-gray-400">{match.league}</span>
        </div>

        <div className="flex items-center justify-between">
          <div className="flex-1 text-center">
            <div className="w-12 h-12 mx-auto bg-white/10 rounded-full flex items-center justify-center mb-2">
              <Trophy className="w-6 h-6 text-primary" />
            </div>
            <p className="font-medium text-white text-sm">{match.homeTeam.name}</p>
            <p className="text-xs text-gray-400">Home</p>
          </div>

          <div className="px-4 text-center">
            {match.status === 'finished' ? (
              <p className="text-2xl font-bold text-white">
                {match.homeTeam.score} - {match.awayTeam.score}
              </p>
            ) : (
              <p className="text-lg font-medium text-gray-400">VS</p>
            )}
            <div className="flex items-center gap-1 text-xs text-gray-400 mt-2">
              <Clock className="w-3 h-3" />
              {match.startTime}
            </div>
          </div>

          <div className="flex-1 text-center">
            <div className="w-12 h-12 mx-auto bg-white/10 rounded-full flex items-center justify-center mb-2">
              <Trophy className="w-6 h-6 text-secondary" />
            </div>
            <p className="font-medium text-white text-sm">{match.awayTeam.name}</p>
            <p className="text-xs text-gray-400">Away</p>
          </div>
        </div>

        {match.odds && (
          <div className="mt-4 pt-4 border-t border-white/10 flex justify-center gap-4">
            <div className="text-center">
              <p className="text-xs text-gray-400">1</p>
              <p className="font-bold text-primary">{match.odds.home.toFixed(2)}</p>
            </div>
            <div className="text-center">
              <p className="text-xs text-gray-400">X</p>
              <p className="font-bold text-gray-300">{match.odds.draw.toFixed(2)}</p>
            </div>
            <div className="text-center">
              <p className="text-xs text-gray-400">2</p>
              <p className="font-bold text-secondary">{match.odds.away.toFixed(2)}</p>
            </div>
          </div>
        )}
      </GlassCard>
    </motion.div>
  );
}
