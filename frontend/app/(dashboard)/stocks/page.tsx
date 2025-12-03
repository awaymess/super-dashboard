'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { LineChart as LineChartIcon, Search, Filter, Star, TrendingUp, TrendingDown } from 'lucide-react';
import { GlassCard, GlassInput, GlassButton, Badge, Tabs } from '@/components/ui';
import { StockCard, WatchlistCard, SectorHeatmap, NewsCard } from '@/components/stocks';
import { stocks, news } from '@/lib/mock-data';

export default function StocksPage() {
  const [selectedTab, setSelectedTab] = useState('overview');
  const [searchQuery, setSearchQuery] = useState('');
  const [watchlist, setWatchlist] = useState(stocks.slice(0, 5));

  const filteredStocks = stocks.filter(stock =>
    stock.symbol.toLowerCase().includes(searchQuery.toLowerCase()) ||
    stock.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  const sectors = [
    { name: 'Technology', change: 2.34 },
    { name: 'Healthcare', change: 1.12 },
    { name: 'Finance', change: -0.87 },
    { name: 'Energy', change: -1.23 },
    { name: 'Consumer', change: 0.45 },
    { name: 'Industrial', change: 0.78 },
    { name: 'Materials', change: -0.34 },
    { name: 'Utilities', change: 0.12 },
  ];

  const gainers = [...stocks].sort((a, b) => b.changePercent - a.changePercent).slice(0, 5);
  const losers = [...stocks].sort((a, b) => a.changePercent - b.changePercent).slice(0, 5);

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-white flex items-center gap-3">
            <LineChartIcon className="w-8 h-8 text-primary" />
            Stocks
          </h1>
          <p className="text-gray-400 mt-1">Monitor and analyze stock market</p>
        </div>
      </div>

      <div className="flex gap-4">
        <div className="flex-1 relative">
          <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-gray-400" />
          <GlassInput
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            placeholder="Search stocks..."
            className="pl-10"
          />
        </div>
        <GlassButton variant="ghost">
          <Filter className="w-5 h-5 mr-2" />
          Filter
        </GlassButton>
      </div>

      <Tabs
        tabs={[
          { id: 'overview', label: 'Overview' },
          { id: 'watchlist', label: 'Watchlist' },
          { id: 'news', label: 'News' },
        ]}
        activeTab={selectedTab}
        onChange={setSelectedTab}
      />

      {selectedTab === 'overview' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
          <SectorHeatmap sectors={sectors} />
          
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <GlassCard>
              <div className="flex items-center gap-2 mb-4">
                <TrendingUp className="w-5 h-5 text-success" />
                <h3 className="font-bold text-white">Top Gainers</h3>
              </div>
              <div className="space-y-3">
                {gainers.map((stock) => (
                  <div key={stock.symbol} className="flex items-center justify-between p-3 bg-white/5 rounded-lg">
                    <div>
                      <p className="font-semibold text-white">{stock.symbol}</p>
                      <p className="text-xs text-gray-400">{stock.name}</p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-white">${stock.price.toFixed(2)}</p>
                      <p className="text-sm text-success">+{stock.changePercent.toFixed(2)}%</p>
                    </div>
                  </div>
                ))}
              </div>
            </GlassCard>

            <GlassCard>
              <div className="flex items-center gap-2 mb-4">
                <TrendingDown className="w-5 h-5 text-danger" />
                <h3 className="font-bold text-white">Top Losers</h3>
              </div>
              <div className="space-y-3">
                {losers.map((stock) => (
                  <div key={stock.symbol} className="flex items-center justify-between p-3 bg-white/5 rounded-lg">
                    <div>
                      <p className="font-semibold text-white">{stock.symbol}</p>
                      <p className="text-xs text-gray-400">{stock.name}</p>
                    </div>
                    <div className="text-right">
                      <p className="font-semibold text-white">${stock.price.toFixed(2)}</p>
                      <p className="text-sm text-danger">{stock.changePercent.toFixed(2)}%</p>
                    </div>
                  </div>
                ))}
              </div>
            </GlassCard>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {filteredStocks.slice(0, 9).map((stock) => (
              <StockCard key={stock.symbol} stock={stock} />
            ))}
          </div>
        </motion.div>
      )}

      {selectedTab === 'watchlist' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <WatchlistCard stocks={watchlist} />
        </motion.div>
      )}

      {selectedTab === 'news' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="grid gap-4">
          {news.map((news) => (
            <NewsCard key={news.id} news={news} />
          ))}
        </motion.div>
      )}
    </div>
  );
}
