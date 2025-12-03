'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { Trophy, Zap, Calendar, Clock, Filter } from 'lucide-react';
import { GlassCard, GlassButton, Badge, Tabs } from '@/components/ui';
import { MatchCard, ValueBetCard, BetSlip, KellyCalculator, PoissonCalculator } from '@/components/betting';
import { matches, mockValueBets } from '@/lib/mock-data';
import type { Bet } from '@/types/betting';

export default function BettingPage() {
  const [selectedTab, setSelectedTab] = useState('matches');
  const [bets, setBets] = useState<Bet[]>([]);
  const [statusFilter, setStatusFilter] = useState<'all' | 'live' | 'upcoming'>('all');

  const filteredMatches = matches.filter(m => 
    statusFilter === 'all' || m.status === statusFilter
  );

  const addBet = (match: typeof matches[0], selection: string, odds: number) => {
    const newBet: Bet = {
      id: `${match.id}-${selection}`,
      matchId: match.id,
      match,
      betType: selection,
      odds,
      stake: 10,
      potentialWin: 10 * odds,
      status: 'pending',
      placedAt: new Date().toISOString(),
    };
    setBets(prev => [...prev.filter(b => b.id !== newBet.id), newBet]);
  };

  //  const addBet = (match: typeof mockMatches[0], selection: string, odds: number) => {
  //   const newBet: Bet = {
  //     id: `${match.id}-${selection}`,
  //     match: `${match.homeTeam.name} vs ${match.awayTeam.name}`,
  //     selection,
  //     odds,
  //     stake: 10,
  //   };
  //   setBets(prev => [...prev.filter(b => b.id !== newBet.id), newBet]);
  // };

  const removeBet = (id: string) => {
    setBets(prev => prev.filter(b => b.id !== id));
  };

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <Trophy className="w-8 h-8 text-warning" />
            Betting
          </h1>
          <p className="text-gray-400 mt-1">Find value bets and analyze matches</p>
        </div>
      </div>

      <Tabs
        tabs={[
          { id: 'matches', label: 'Matches', icon: <Calendar /> },
          { id: 'value-bets', label: 'Value Bets', icon: <Zap /> },
          { id: 'calculators', label: 'Calculators', icon: <Clock /> },
        ]}
        activeTab={selectedTab}
        onChange={setSelectedTab}
      />

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2">
          {selectedTab === 'matches' && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
              <div className="flex gap-2 mb-4">
                {(['all', 'live', 'upcoming'] as const).map((filter) => (
                  <GlassButton
                    key={filter}
                    variant={statusFilter === filter ? 'primary' : 'ghost'}
                    size="sm"
                    onClick={() => setStatusFilter(filter)}
                  >
                    {filter.charAt(0).toUpperCase() + filter.slice(1)}
                  </GlassButton>
                ))}
              </div>
              <div className="grid gap-4">
                {filteredMatches.slice(0, 10).map((match) => (
                  <MatchCard key={match.id} match={match} />
                ))}
              </div>
            </motion.div>
          )}

          {selectedTab === 'value-bets' && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
              <div className="grid gap-4">
                {mockValueBets.map((valueBet) => (
                  <ValueBetCard key={valueBet.matchId} valueBet={valueBet} />
                ))}
              </div>
            </motion.div>
          )}

          {selectedTab === 'calculators' && (
            <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
              <div className="grid gap-6">
                <KellyCalculator />
                <PoissonCalculator />
              </div>
            </motion.div>
          )}
        </div>

        <div>
          <BetSlip
            bets={bets}
            onRemoveBet={removeBet}
            onPlaceBet={(stake) => {
              console.log('Placing bet with stake:', stake);
              setBets([]);
            }}
            onClear={() => setBets([])}
          />
        </div>
      </div>
    </div>
  );
}
