'use client';

import { useState } from 'react';
import { motion } from 'framer-motion';
import { Settings as SettingsIcon, User, Bell, Shield, Palette, Globe, CreditCard } from 'lucide-react';
import { GlassCard, GlassInput, GlassButton, Toggle, Tabs, GlassSelect } from '@/components/ui';

export default function SettingsPage() {
  const [selectedTab, setSelectedTab] = useState('profile');
  const [notifications, setNotifications] = useState({
    email: true,
    push: true,
    valueBets: true,
    stockAlerts: true,
    tradeExecuted: false,
  });

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-3xl font-bold text-white flex items-center gap-3">
          <SettingsIcon className="w-8 h-8 text-gray-400" />
          Settings
        </h1>
        <p className="text-gray-400 mt-1">Manage your account and preferences</p>
      </div>

      <Tabs
        tabs={[
          { id: 'profile', label: 'Profile', icon: User },
          { id: 'notifications', label: 'Notifications', icon: Bell },
          { id: 'security', label: 'Security', icon: Shield },
          { id: 'appearance', label: 'Appearance', icon: Palette },
        ]}
        activeTab={selectedTab}
        onChange={setSelectedTab}
      />

      {selectedTab === 'profile' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <GlassCard>
            <h3 className="font-bold text-white mb-6">Profile Information</h3>
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Full Name</label>
                <GlassInput defaultValue="John Doe" />
              </div>
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Email</label>
                <GlassInput type="email" defaultValue="john@example.com" />
              </div>
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Username</label>
                <GlassInput defaultValue="johndoe" />
              </div>
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Timezone</label>
                <GlassSelect
                  options={[
                    { value: 'utc', label: 'UTC' },
                    { value: 'est', label: 'Eastern Time' },
                    { value: 'pst', label: 'Pacific Time' },
                    { value: 'gmt7', label: 'GMT+7 (Bangkok)' },
                  ]}
                  defaultValue="utc"
                />
              </div>
            </div>
            <div className="mt-6 pt-6 border-t border-white/10">
              <GlassButton variant="primary">Save Changes</GlassButton>
            </div>
          </GlassCard>
        </motion.div>
      )}

      {selectedTab === 'notifications' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <GlassCard>
            <h3 className="font-bold text-white mb-6">Notification Preferences</h3>
            <div className="space-y-6">
              <div className="flex items-center justify-between py-3 border-b border-white/10">
                <div>
                  <p className="font-medium text-white">Email Notifications</p>
                  <p className="text-sm text-gray-400">Receive notifications via email</p>
                </div>
                <Toggle
                  checked={notifications.email}
                  onChange={(checked) => setNotifications({ ...notifications, email: checked })}
                />
              </div>
              <div className="flex items-center justify-between py-3 border-b border-white/10">
                <div>
                  <p className="font-medium text-white">Push Notifications</p>
                  <p className="text-sm text-gray-400">Receive push notifications in browser</p>
                </div>
                <Toggle
                  checked={notifications.push}
                  onChange={(checked) => setNotifications({ ...notifications, push: checked })}
                />
              </div>
              <div className="flex items-center justify-between py-3 border-b border-white/10">
                <div>
                  <p className="font-medium text-white">Value Bet Alerts</p>
                  <p className="text-sm text-gray-400">Get notified when new value bets are found</p>
                </div>
                <Toggle
                  checked={notifications.valueBets}
                  onChange={(checked) => setNotifications({ ...notifications, valueBets: checked })}
                />
              </div>
              <div className="flex items-center justify-between py-3 border-b border-white/10">
                <div>
                  <p className="font-medium text-white">Stock Price Alerts</p>
                  <p className="text-sm text-gray-400">Get notified when stocks hit target prices</p>
                </div>
                <Toggle
                  checked={notifications.stockAlerts}
                  onChange={(checked) => setNotifications({ ...notifications, stockAlerts: checked })}
                />
              </div>
              <div className="flex items-center justify-between py-3">
                <div>
                  <p className="font-medium text-white">Trade Executed</p>
                  <p className="text-sm text-gray-400">Get notified when trades are executed</p>
                </div>
                <Toggle
                  checked={notifications.tradeExecuted}
                  onChange={(checked) => setNotifications({ ...notifications, tradeExecuted: checked })}
                />
              </div>
            </div>
          </GlassCard>
        </motion.div>
      )}

      {selectedTab === 'security' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }} className="space-y-6">
          <GlassCard>
            <h3 className="font-bold text-white mb-6">Change Password</h3>
            <div className="space-y-4 max-w-md">
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Current Password</label>
                <GlassInput type="password" placeholder="••••••••" />
              </div>
              <div>
                <label className="text-sm text-gray-400 mb-2 block">New Password</label>
                <GlassInput type="password" placeholder="••••••••" />
              </div>
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Confirm New Password</label>
                <GlassInput type="password" placeholder="••••••••" />
              </div>
              <GlassButton variant="primary">Update Password</GlassButton>
            </div>
          </GlassCard>

          <GlassCard>
            <h3 className="font-bold text-white mb-6">Two-Factor Authentication</h3>
            <p className="text-gray-400 mb-4">Add an extra layer of security to your account</p>
            <GlassButton variant="secondary">Enable 2FA</GlassButton>
          </GlassCard>
        </motion.div>
      )}

      {selectedTab === 'appearance' && (
        <motion.div initial={{ opacity: 0 }} animate={{ opacity: 1 }}>
          <GlassCard>
            <h3 className="font-bold text-white mb-6">Appearance Settings</h3>
            <div className="space-y-6">
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Theme</label>
                <div className="flex gap-4">
                  <button className="p-4 bg-[#0a0a0f] border-2 border-primary rounded-xl">
                    <div className="w-16 h-10 bg-[#12121a] rounded" />
                    <p className="text-sm text-white mt-2">Dark</p>
                  </button>
                  <button className="p-4 bg-gray-100 border-2 border-transparent rounded-xl opacity-50">
                    <div className="w-16 h-10 bg-white rounded border" />
                    <p className="text-sm text-gray-800 mt-2">Light</p>
                  </button>
                </div>
              </div>
              <div>
                <label className="text-sm text-gray-400 mb-2 block">Language</label>
                <GlassSelect
                  options={[
                    { value: 'en', label: 'English' },
                    { value: 'th', label: 'ไทย (Thai)' },
                  ]}
                  defaultValue="en"
                />
              </div>
            </div>
          </GlassCard>
        </motion.div>
      )}
    </div>
  );
}
