'use client';

import { useState } from 'react';
import { Search, Bell, Command, Sun, Moon, Globe, User } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import { useTheme } from '@/hooks/useTheme';
import { Avatar, Badge, Dropdown } from '@/components/ui';

interface HeaderProps {
  onOpenCommandPalette: () => void;
}

export function Header({ onOpenCommandPalette }: HeaderProps) {
  const { theme, toggleTheme } = useTheme();
  const [showNotifications, setShowNotifications] = useState(false);

  const notifications = [
    { id: 1, title: 'Value Bet Found', message: 'Arsenal vs Chelsea - Over 2.5 Goals at 2.10', time: '2m ago', unread: true },
    { id: 2, title: 'Stock Alert', message: 'AAPL reached target price of $190', time: '15m ago', unread: true },
    { id: 3, title: 'Trade Executed', message: 'Bought 100 shares of MSFT at $380', time: '1h ago', unread: false },
  ];

  const unreadCount = notifications.filter(n => n.unread).length;

  return (
    <header className="h-16 bg-surface/80 backdrop-blur-xl border-b border-white/10 sticky top-0 z-30">
      <div className="h-full flex items-center justify-between px-6">
        {/* Search */}
        <button
          onClick={onOpenCommandPalette}
          className="flex items-center gap-3 px-4 py-2 rounded-xl bg-white/5 hover:bg-white/10 transition-colors text-gray-400 group"
        >
          <Search className="w-4 h-4" />
          <span className="text-sm">Search...</span>
          <div className="flex items-center gap-1 ml-8">
            <kbd className="px-2 py-0.5 bg-white/10 rounded text-xs group-hover:bg-primary/20 group-hover:text-primary transition-colors">
              <Command className="w-3 h-3 inline" />
            </kbd>
            <kbd className="px-2 py-0.5 bg-white/10 rounded text-xs group-hover:bg-primary/20 group-hover:text-primary transition-colors">
              K
            </kbd>
          </div>
        </button>

        {/* Right Section */}
        <div className="flex items-center gap-4">
          {/* Theme Toggle */}
          <button
            onClick={toggleTheme}
            className="p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-gray-400 hover:text-white"
          >
            {theme === 'dark' ? <Sun className="w-5 h-5" /> : <Moon className="w-5 h-5" />}
          </button>

          {/* Language Toggle */}
          <button className="p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-gray-400 hover:text-white">
            <Globe className="w-5 h-5" />
          </button>

          {/* Notifications */}
          <div className="relative">
            <button
              onClick={() => setShowNotifications(!showNotifications)}
              className="p-2 rounded-lg bg-white/5 hover:bg-white/10 transition-colors text-gray-400 hover:text-white relative"
            >
              <Bell className="w-5 h-5" />
              {unreadCount > 0 && (
                <span className="absolute -top-1 -right-1 w-5 h-5 bg-danger rounded-full text-xs flex items-center justify-center text-white">
                  {unreadCount}
                </span>
              )}
            </button>

            <AnimatePresence>
              {showNotifications && (
                <motion.div
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, y: 10 }}
                  className="absolute right-0 top-full mt-2 w-80 bg-surface border border-white/10 rounded-xl shadow-xl overflow-hidden"
                >
                  <div className="p-4 border-b border-white/10">
                    <h3 className="font-semibold text-white">Notifications</h3>
                  </div>
                  <div className="max-h-80 overflow-y-auto">
                    {notifications.map(notification => (
                      <div
                        key={notification.id}
                        className={`p-4 border-b border-white/5 hover:bg-white/5 transition-colors cursor-pointer ${
                          notification.unread ? 'bg-primary/5' : ''
                        }`}
                      >
                        <div className="flex items-start gap-3">
                          {notification.unread && (
                            <span className="w-2 h-2 bg-primary rounded-full mt-2" />
                          )}
                          <div className="flex-1">
                            <p className="font-medium text-white text-sm">{notification.title}</p>
                            <p className="text-gray-400 text-xs mt-1">{notification.message}</p>
                            <p className="text-gray-500 text-xs mt-2">{notification.time}</p>
                          </div>
                        </div>
                      </div>
                    ))}
                  </div>
                </motion.div>
              )}
            </AnimatePresence>
          </div>

          {/* User Menu */}
          <div className="flex items-center gap-3 pl-4 border-l border-white/10">
            <Avatar name="John Doe" size="sm" />
            <div className="hidden md:block">
              <p className="text-sm font-medium text-white">John Doe</p>
              <p className="text-xs text-gray-400">Pro Member</p>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
}
