'use client';

import { GlassCard, Badge } from '@/components/ui';
import { Trophy, Minus, X } from 'lucide-react';

interface H2HMatch {
  date: string;
  homeTeam: string;
  awayTeam: string;
  homeScore: number;
  awayScore: number;
  competition: string;
}

interface H2HHistoryProps {
  homeTeam: string;
  awayTeam: string;
  matches: H2HMatch[];
}

export function H2HHistory({ homeTeam, awayTeam, matches }: H2HHistoryProps) {
  const homeWins = matches.filter(m => 
    (m.homeTeam === homeTeam && m.homeScore > m.awayScore) ||
    (m.awayTeam === homeTeam && m.awayScore > m.homeScore)
  ).length;
  
  const awayWins = matches.filter(m => 
    (m.homeTeam === awayTeam && m.homeScore > m.awayScore) ||
    (m.awayTeam === awayTeam && m.awayScore > m.homeScore)
  ).length;
  
  const draws = matches.length - homeWins - awayWins;

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Head to Head</h3>
      
      <div className="grid grid-cols-3 gap-4 mb-6">
        <div className="text-center">
          <div className="w-12 h-12 mx-auto bg-primary/20 rounded-full flex items-center justify-center mb-2">
            <Trophy className="w-6 h-6 text-primary" />
          </div>
          <p className="text-2xl font-bold text-white">{homeWins}</p>
          <p className="text-xs text-gray-400">{homeTeam} Wins</p>
        </div>
        <div className="text-center">
          <div className="w-12 h-12 mx-auto bg-white/10 rounded-full flex items-center justify-center mb-2">
            <Minus className="w-6 h-6 text-gray-400" />
          </div>
          <p className="text-2xl font-bold text-white">{draws}</p>
          <p className="text-xs text-gray-400">Draws</p>
        </div>
        <div className="text-center">
          <div className="w-12 h-12 mx-auto bg-secondary/20 rounded-full flex items-center justify-center mb-2">
            <Trophy className="w-6 h-6 text-secondary" />
          </div>
          <p className="text-2xl font-bold text-white">{awayWins}</p>
          <p className="text-xs text-gray-400">{awayTeam} Wins</p>
        </div>
      </div>

      <div className="space-y-3">
        {matches.slice(0, 5).map((match, i) => {
          const isHomeTeamHome = match.homeTeam === homeTeam;
          const result = match.homeScore > match.awayScore 
            ? (isHomeTeamHome ? 'W' : 'L')
            : match.homeScore < match.awayScore 
            ? (isHomeTeamHome ? 'L' : 'W')
            : 'D';
          
          return (
            <div key={i} className="flex items-center justify-between p-3 bg-white/5 rounded-lg">
              <div className="flex items-center gap-3">
                <Badge 
                  variant={result === 'W' ? 'success' : result === 'L' ? 'danger' : 'default'}
                >
                  {result}
                </Badge>
                <div>
                  <p className="text-sm text-white">{match.homeTeam} vs {match.awayTeam}</p>
                  <p className="text-xs text-gray-400">{match.competition}</p>
                </div>
              </div>
              <div className="text-right">
                <p className="font-bold text-white">{match.homeScore} - {match.awayScore}</p>
                <p className="text-xs text-gray-400">{match.date}</p>
              </div>
            </div>
          );
        })}
      </div>
    </GlassCard>
  );
}
