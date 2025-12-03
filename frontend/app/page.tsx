'use client';

import { useRouter } from 'next/navigation';
import { motion } from 'framer-motion';
import { TrendingUp, Trophy, LineChart, Wallet, Target, ArrowRight, Sparkles } from 'lucide-react';
import { GlassButton, GlassCard } from '@/components/ui';
import { LiquidBackground } from '@/components/layout';

export default function HomePage() {
  const router = useRouter();

  const features = [
    { icon: Trophy, title: 'Betting Analytics', description: 'Value bets, Poisson models, Kelly criterion' },
    { icon: LineChart, title: 'Stock Monitoring', description: 'Real-time quotes, technical analysis' },
    { icon: Wallet, title: 'Paper Trading', description: 'Risk-free trading simulation' },
    { icon: Target, title: 'Performance Analytics', description: 'Track your progress and goals' },
  ];

  return (
    <div className="min-h-screen relative overflow-hidden">
      <LiquidBackground />
      
      <div className="relative z-10 container mx-auto px-4 py-12">
        {/* Hero Section */}
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="text-center max-w-4xl mx-auto mb-16"
        >
          <div className="flex items-center justify-center gap-2 mb-6">
            <div className="w-16 h-16 rounded-2xl bg-gradient-to-br from-primary to-secondary flex items-center justify-center">
              <TrendingUp className="w-8 h-8 text-white" />
            </div>
          </div>
          
          <h1 className="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-white via-gray-200 to-gray-400 bg-clip-text text-transparent">
            Super Dashboard
          </h1>
          
          <p className="text-xl text-gray-400 mb-8 max-w-2xl mx-auto">
            Your integrated platform for sports betting analytics and stock market monitoring. 
            Make data-driven decisions with powerful tools and real-time insights.
          </p>
          
          <div className="flex items-center justify-center gap-4">
            <GlassButton variant="primary" size="lg" onClick={() => router.push('/dashboard')}>
              Get Started
              <ArrowRight className="w-5 h-5 ml-2" />
            </GlassButton>
            <GlassButton variant="ghost" size="lg" onClick={() => router.push('/login')}>
              Sign In
            </GlassButton>
          </div>
        </motion.div>

        {/* Features Grid */}
        <motion.div
          initial={{ opacity: 0, y: 40 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ delay: 0.2 }}
          className="grid md:grid-cols-2 lg:grid-cols-4 gap-6 max-w-6xl mx-auto"
        >
          {features.map((feature, i) => (
            <motion.div
              key={feature.title}
              initial={{ opacity: 0, y: 20 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ delay: 0.3 + i * 0.1 }}
            >
              <GlassCard className="h-full hover:border-primary/50 transition-colors cursor-pointer group">
                <div className="p-3 bg-primary/20 rounded-xl w-fit mb-4 group-hover:bg-primary/30 transition-colors">
                  <feature.icon className="w-6 h-6 text-primary" />
                </div>
                <h3 className="font-bold text-white mb-2">{feature.title}</h3>
                <p className="text-sm text-gray-400">{feature.description}</p>
              </GlassCard>
            </motion.div>
          ))}
        </motion.div>

        {/* Stats Section */}
        <motion.div
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          transition={{ delay: 0.6 }}
          className="mt-20 grid grid-cols-2 md:grid-cols-4 gap-8 max-w-4xl mx-auto text-center"
        >
          {[
            { value: '50+', label: 'Matches Tracked' },
            { value: '100+', label: 'Stocks Monitored' },
            { value: '5+', label: 'Analysis Tools' },
            { value: '24/7', label: 'Real-time Data' },
          ].map((stat, i) => (
            <div key={stat.label}>
              <p className="text-4xl font-bold bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                {stat.value}
              </p>
              <p className="text-gray-400 mt-1">{stat.label}</p>
            </div>
          ))}
        </motion.div>
      </div>
    </div>
  );
}
