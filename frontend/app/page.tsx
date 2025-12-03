'use client';

import { GlassCard, GlassButton } from '@/components/ui';

export default function HomePage() {
  return (
    <main className="min-h-screen bg-gradient-to-br from-background via-surface to-background p-8">
      <div className="max-w-7xl mx-auto">
        <div className="text-center mb-12">
          <h1 className="text-5xl font-bold bg-gradient-to-r from-primary via-secondary to-primary bg-clip-text text-transparent mb-4">
            Super Dashboard
          </h1>
          <p className="text-xl text-gray-400">
            Integrated Sports Betting Analytics &amp; Stock Monitoring Platform
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <GlassCard className="p-6">
            <h2 className="text-2xl font-semibold mb-4 text-primary">🎯 Betting Analytics</h2>
            <p className="text-gray-400 mb-4">
              Advanced Poisson distribution, Kelly Criterion, and ELO-based predictions
            </p>
            <GlassButton variant="primary">Explore Bets</GlassButton>
          </GlassCard>

          <GlassCard className="p-6">
            <h2 className="text-2xl font-semibold mb-4 text-success">📈 Stock Monitoring</h2>
            <p className="text-gray-400 mb-4">
              Real-time quotes, technical indicators, and fundamental analysis
            </p>
            <GlassButton variant="success">View Stocks</GlassButton>
          </GlassCard>

          <GlassCard className="p-6">
            <h2 className="text-2xl font-semibold mb-4 text-secondary">💼 Paper Trading</h2>
            <p className="text-gray-400 mb-4">
              Practice trading with virtual portfolio and track performance
            </p>
            <GlassButton variant="secondary">Start Trading</GlassButton>
          </GlassCard>
        </div>

        <div className="mt-12 text-center">
          <GlassCard className="inline-block p-6">
            <p className="text-gray-400">
              Press <kbd className="px-2 py-1 bg-surface rounded border border-white/10">Ctrl+K</kbd> for Command Palette
            </p>
          </GlassCard>
        </div>
      </div>
    </main>
  );
}
