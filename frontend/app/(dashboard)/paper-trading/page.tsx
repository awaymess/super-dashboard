'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { Wallet, Plus } from 'lucide-react';
import { GlassButton, Tabs } from '@/components/ui';
import { PortfolioSummary, PositionCard, TradeForm, TradeJournal, Leaderboard } from '@/components/paper-trading';
import { mockPortfolio, mockPositions, mockTrades, mockLeaderboard } from '@/lib/mock-data';

export default function PaperTradingPage() {
  const [selectedTab, setSelectedTab] = useState('portfolio');

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <Wallet className="w-8 h-8 text-secondary" />
            Paper Trading
          </h1>
          <p className="text-gray-400 mt-1">Practice trading without risk</p>
        </div>
        <GlassButton variant="primary">
          <Plus className="w-5 h-5 mr-2" />
          New Trade
        </GlassButton>
      </div>

      <PortfolioSummary portfolio={mockPortfolio} />

      <Tabs
        tabs={[
          { id: 'portfolio', label: 'Positions' },
          { id: 'trade', label: 'Trade' },
          { id: 'journal', label: 'Journal' },
          { id: 'leaderboard', label: 'Leaderboard' },
        ]}
        activeTab={selectedTab}
        onChange={setSelectedTab}
      />

      {selectedTab === 'portfolio' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {mockPositions.map((position) => (
              <PositionCard key={position.id} position={position} />
            ))}
          </div>
        </motion.div>
      )}

      {selectedTab === 'trade' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <div className="max-w-md">
            <TradeForm
              onSubmit={(trade) => console.log('Trade submitted:', trade)}
            />
          </div>
        </motion.div>
      )}

      {selectedTab === 'journal' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <TradeJournal trades={mockTrades} />
        </motion.div>
      )}

      {selectedTab === 'leaderboard' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <Leaderboard entries={mockLeaderboard} currentUser="JohnDoe" />
        </motion.div>
      )}
    </div>
  );
}
