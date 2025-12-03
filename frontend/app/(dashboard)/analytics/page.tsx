'use client';

import { motion } from 'framer-motion';
import { Target, TrendingUp, Trophy, DollarSign, Percent, BarChart3 } from 'lucide-react';
import { StatsCard, PerformanceChart, GoalTracker, DrawdownChart } from '@/components/analytics';
import { DoughnutChart, BarChart } from '@/components/charts';
import { GlassCard } from '@/components/ui';

export default function AnalyticsPage() {
  const performanceData = [10000, 10500, 10200, 11000, 10800, 11500, 12000, 11800, 12500, 13000, 12800, 14000];
  const labels = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
  const drawdownData = [0, -2, -1.5, 0, -3, -2, 0, -1, 0, -4, -2, 0];

  const goals = [
    { id: '1', title: 'Monthly Profit Target', target: 5000, current: 3200, unit: 'USD', deadline: 'Dec 31' },
    { id: '2', title: 'Win Rate Goal', target: 70, current: 67, unit: '%', deadline: 'Dec 31' },
    { id: '3', title: 'Trade Volume', target: 100, current: 85, unit: 'trades', deadline: 'Dec 31' },
    { id: '4', title: 'ROI Target', target: 25, current: 18.5, unit: '%' },
  ];

  const monthlyReturns = [5.2, 3.1, -1.8, 4.5, 2.3, 6.1, -0.5, 3.8, 4.2, -2.1, 5.5, 8.2];

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <Target className="w-8 h-8 text-warning" />
            Analytics
          </h1>
          <p className="text-gray-400 mt-1">Track your performance and goals</p>
        </div>
      </div>

      {/* Stats Overview */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }}>
          <StatsCard title="Total Profit" value="$14,250" change={18.5} icon={DollarSign} color="success" />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.1 }}>
          <StatsCard title="ROI" value="28.5%" change={5.2} icon={Percent} color="primary" />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.2 }}>
          <StatsCard title="Win Rate" value="67%" change={3.1} icon={Trophy} color="warning" />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.3 }}>
          <StatsCard title="Total Trades" value="156" change={12} icon={BarChart3} color="secondary" />
        </motion.div>
      </div>

      {/* Charts */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.4 }}>
          <PerformanceChart title="Portfolio Growth" data={performanceData} labels={labels} />
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.5 }}>
          <DrawdownChart data={drawdownData} labels={labels} />
        </motion.div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.6 }} className="lg:col-span-2">
          <GlassCard>
            <h3 className="font-bold text-white mb-4">Monthly Returns</h3>
            <BarChart
              data={monthlyReturns}
              labels={labels}
              colors={monthlyReturns.map(v => v >= 0 ? '#10b981' : '#ef4444')}
              height={300}
            />
          </GlassCard>
        </motion.div>
        <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.7 }}>
          <GlassCard>
            <h3 className="font-bold text-white mb-4">Profit by Category</h3>
            <DoughnutChart
              data={[5200, 4800, 2500, 1750]}
              labels={['Stocks', 'Betting', 'Options', 'Crypto']}
              height={250}
            />
          </GlassCard>
        </motion.div>
      </div>

      {/* Goals */}
      <motion.div initial={{ opacity: 0, y: 20 }} animate={{ opacity: 1, y: 0 }} transition={{ delay: 0.8 }}>
        <GoalTracker goals={goals} />
      </motion.div>
    </div>
  );
}
