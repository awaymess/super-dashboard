'use client';

import { motion } from 'framer-motion';
import { DollarSign, TrendingUp, Trophy, Target, ArrowUp, ArrowDown } from 'lucide-react';
import { GlassCard, Badge } from '@/components/ui';
import { StatsCard, PerformanceChart } from '@/components/analytics';
import { AreaChart, DoughnutChart } from '@/components/charts';
import { mockMatches, mockStocks, mockPortfolio } from '@/lib/mock-data';

export default function DashboardPage() {
  const performanceData = [10000, 10500, 10200, 11000, 10800, 11500, 12000, 11800, 12500, 13000];
  const labels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct'];

  const topStocks = mockStocks.slice(0, 5);
  const liveMatches = mockMatches.filter(m => m.status === 'live').slice(0, 3);
  const upcomingMatches = mockMatches.filter(m => m.status === 'upcoming').slice(0, 3);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white">Dashboard</h1>
          <p className="text-gray-400 mt-1">Welcome back! Here&apos;s your overview.</p>
        </div>
        <Badge variant="success">All Systems Operational</Badge>
      </div>

      {/* Stats Cards */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
          <StatsCard
            title="Portfolio Value"
            value={`$${mockPortfolio.totalValue.toLocaleString()}`}
            change={mockPortfolio.totalPnLPercent}
            icon={DollarSign}
            color="primary"
          />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
          <StatsCard
            title="Betting Profit"
            value="$2,450"
            change={12.5}
            icon={Trophy}
            color="success"
          />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2 }}>
          <StatsCard
            title="Stock Returns"
            value="$4,320"
            change={8.3}
            icon={TrendingUp}
            color="secondary"
          />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3 }}>
          <StatsCard
            title="Win Rate"
            value="67%"
            change={5.2}
            icon={Target}
            color="warning"
          />
        </motion.div>
      </div>

      {/* Performance Chart */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.4 }}
          className="lg:col-span-2"
        >
          <PerformanceChart
            title="Portfolio Performance"
            data={performanceData}
            labels={labels}
          />
        </motion.div>

        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.5 }}
        >
          <GlassCard className="h-full">
            <h3 className="font-bold text-white mb-4">Asset Allocation</h3>
            <DoughnutChart
              data={[45, 30, 15, 10]}
              labels={['Stocks', 'Betting', 'Cash', 'Crypto']}
              height={250}
            />
          </GlassCard>
        </motion.div>
      </div>

      {/* Quick Access */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Top Stocks */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.6 }}
        >
          <GlassCard>
            <h3 className="font-bold text-white mb-4">Top Movers</h3>
            <div className="space-y-3">
              {topStocks.map((stock) => (
                <div key={stock.symbol} className="flex items-center justify-between p-3 bg-white/5 rounded-lg">
                  <div>
                    <p className="font-semibold text-white">{stock.symbol}</p>
                    <p className="text-xs text-gray-400">{stock.name}</p>
                  </div>
                  <div className="text-right">
                    <p className="font-semibold text-white">${stock.price.toFixed(2)}</p>
                    <p className={`text-sm flex items-center justify-end gap-1 ${stock.change >= 0 ? 'text-success' : 'text-danger'}`}>
                      {stock.change >= 0 ? <ArrowUp className="w-3 h-3" /> : <ArrowDown className="w-3 h-3" />}
                      {stock.change >= 0 ? '+' : ''}{stock.changePercent.toFixed(2)}%
                    </p>
                  </div>
                </div>
              ))}
            </div>
          </GlassCard>
        </motion.div>

        {/* Upcoming Matches */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.7 }}
        >
          <GlassCard>
            <h3 className="font-bold text-white mb-4">Upcoming Matches</h3>
            <div className="space-y-3">
              {upcomingMatches.map((match) => (
                <div key={match.id} className="flex items-center justify-between p-3 bg-white/5 rounded-lg">
                  <div className="flex items-center gap-3">
                    <Badge variant="warning">{match.league}</Badge>
                    <div>
                      <p className="font-medium text-white text-sm">{match.homeTeam.name} vs {match.awayTeam.name}</p>
                      <p className="text-xs text-gray-400">{match.startTime}</p>
                    </div>
                  </div>
                  {match.odds && (
                    <div className="flex gap-2 text-xs">
                      <span className="px-2 py-1 bg-primary/20 text-primary rounded">{match.odds.home.toFixed(2)}</span>
                      <span className="px-2 py-1 bg-white/10 text-gray-300 rounded">{match.odds.draw.toFixed(2)}</span>
                      <span className="px-2 py-1 bg-secondary/20 text-secondary rounded">{match.odds.away.toFixed(2)}</span>
                    </div>
                  )}
                </div>
              ))}
            </div>
          </GlassCard>
        </motion.div>
      </div>
    </div>
  );
}
