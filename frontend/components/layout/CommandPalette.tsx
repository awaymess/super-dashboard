'use client';

import { useState, useEffect, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { motion, AnimatePresence } from 'framer-motion';
import {
  Search,
  LayoutDashboard,
  Trophy,
  LineChart,
  Wallet,
  Target,
  Settings,
  TrendingUp,
  TrendingDown,
  ArrowRight,
} from 'lucide-react';
import { GlassInput } from '@/components/ui';

interface CommandPaletteProps {
  isOpen: boolean;
  onClose: () => void;
}

interface CommandItem {
  id: string;
  title: string;
  description?: string;
  icon: React.ReactNode;
  action: () => void;
  category: string;
}

export function CommandPalette({ isOpen, onClose }: CommandPaletteProps) {
  const [query, setQuery] = useState('');
  const [selectedIndex, setSelectedIndex] = useState(0);
  const router = useRouter();

  const commands: CommandItem[] = [
    // Navigation
    { id: 'nav-dashboard', title: 'Go to Dashboard', icon: <LayoutDashboard className="w-4 h-4" />, action: () => router.push('/dashboard'), category: 'Navigation' },
    { id: 'nav-betting', title: 'Go to Betting', icon: <Trophy className="w-4 h-4" />, action: () => router.push('/betting'), category: 'Navigation' },
    { id: 'nav-stocks', title: 'Go to Stocks', icon: <LineChart className="w-4 h-4" />, action: () => router.push('/stocks'), category: 'Navigation' },
    { id: 'nav-paper', title: 'Go to Paper Trading', icon: <Wallet className="w-4 h-4" />, action: () => router.push('/paper-trading'), category: 'Navigation' },
    { id: 'nav-analytics', title: 'Go to Analytics', icon: <Target className="w-4 h-4" />, action: () => router.push('/analytics'), category: 'Navigation' },
    { id: 'nav-settings', title: 'Go to Settings', icon: <Settings className="w-4 h-4" />, action: () => router.push('/settings'), category: 'Navigation' },
    // Actions
    { id: 'action-valuebets', title: 'Find Value Bets', description: 'Scan for value betting opportunities', icon: <TrendingUp className="w-4 h-4" />, action: () => router.push('/betting/value-bets'), category: 'Actions' },
    { id: 'action-screener', title: 'Stock Screener', description: 'Filter stocks by criteria', icon: <TrendingDown className="w-4 h-4" />, action: () => router.push('/stocks/screener'), category: 'Actions' },
  ];

  const filteredCommands = query
    ? commands.filter(
        cmd =>
          cmd.title.toLowerCase().includes(query.toLowerCase()) ||
          cmd.description?.toLowerCase().includes(query.toLowerCase())
      )
    : commands;

  const groupedCommands = filteredCommands.reduce((acc, cmd) => {
    if (!acc[cmd.category]) acc[cmd.category] = [];
    acc[cmd.category].push(cmd);
    return acc;
  }, {} as Record<string, CommandItem[]>);

  const handleKeyDown = useCallback(
    (e: KeyboardEvent) => {
      if (!isOpen) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex(i => (i + 1) % filteredCommands.length);
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex(i => (i - 1 + filteredCommands.length) % filteredCommands.length);
          break;
        case 'Enter':
          e.preventDefault();
          if (filteredCommands[selectedIndex]) {
            filteredCommands[selectedIndex].action();
            onClose();
          }
          break;
        case 'Escape':
          onClose();
          break;
      }
    },
    [isOpen, filteredCommands, selectedIndex, onClose]
  );

  useEffect(() => {
    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [handleKeyDown]);

  useEffect(() => {
    if (!isOpen) {
      setQuery('');
      setSelectedIndex(0);
    }
  }, [isOpen]);

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            onClick={onClose}
            className="fixed inset-0 bg-black/60 backdrop-blur-sm z-50"
          />
          <motion.div
            initial={{ opacity: 0, scale: 0.95, y: -20 }}
            animate={{ opacity: 1, scale: 1, y: 0 }}
            exit={{ opacity: 0, scale: 0.95, y: -20 }}
            className="fixed top-1/4 left-1/2 -translate-x-1/2 w-full max-w-xl bg-surface border border-white/10 rounded-2xl shadow-2xl overflow-hidden z-50"
          >
            {/* Search Input */}
            <div className="p-4 border-b border-white/10">
              <div className="flex items-center gap-3">
                <Search className="w-5 h-5 text-gray-400" />
                <input
                  type="text"
                  value={query}
                  onChange={e => setQuery(e.target.value)}
                  placeholder="Type a command or search..."
                  className="flex-1 bg-transparent text-white placeholder-gray-400 outline-none"
                  autoFocus
                />
                <kbd className="px-2 py-1 bg-white/10 rounded text-xs text-gray-400">ESC</kbd>
              </div>
            </div>

            {/* Results */}
            <div className="max-h-80 overflow-y-auto p-2">
              {Object.entries(groupedCommands).map(([category, items]) => (
                <div key={category} className="mb-4">
                  <p className="text-xs text-gray-500 px-3 py-2">{category}</p>
                  {items.map((cmd, idx) => {
                    const globalIndex = filteredCommands.indexOf(cmd);
                    return (
                      <button
                        key={cmd.id}
                        onClick={() => {
                          cmd.action();
                          onClose();
                        }}
                        className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg transition-colors ${
                          globalIndex === selectedIndex
                            ? 'bg-primary/20 text-primary'
                            : 'text-gray-300 hover:bg-white/5'
                        }`}
                      >
                        <span className="p-2 bg-white/5 rounded-lg">{cmd.icon}</span>
                        <div className="flex-1 text-left">
                          <p className="font-medium">{cmd.title}</p>
                          {cmd.description && (
                            <p className="text-xs text-gray-500">{cmd.description}</p>
                          )}
                        </div>
                        <ArrowRight className="w-4 h-4 opacity-50" />
                      </button>
                    );
                  })}
                </div>
              ))}
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  );
}
