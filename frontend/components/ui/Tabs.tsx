'use client';

import React from 'react';
import { cn } from '@/lib/cn';
import { TabItem } from '@/types/common';

interface TabsProps {
  tabs: TabItem[];
  activeTab: string;
  onChange: (tabId: string) => void;
  variant?: 'default' | 'pills' | 'underline';
  fullWidth?: boolean;
  className?: string;
}

export function Tabs({
  tabs,
  activeTab,
  onChange,
  variant = 'default',
  fullWidth = false,
  className,
}: TabsProps) {
  return (
    <div
      className={cn(
        'flex',
        {
          'gap-1 p-1 rounded-xl bg-white/5 border border-white/10': variant === 'default' || variant === 'pills',
          'gap-6 border-b border-white/10': variant === 'underline',
          'w-full': fullWidth,
        },
        className
      )}
    >
      {tabs.map((tab) => (
        <button
          key={tab.id}
          onClick={() => !tab.disabled && onChange(tab.id)}
          disabled={tab.disabled}
          className={cn(
            'flex items-center gap-2 font-medium transition-all duration-200',
            {
              'px-4 py-2 rounded-lg text-sm': variant === 'default' || variant === 'pills',
              'pb-3 text-sm': variant === 'underline',
              'flex-1': fullWidth && (variant === 'default' || variant === 'pills'),
            },
            tab.disabled && 'opacity-50 cursor-not-allowed',
            activeTab === tab.id
              ? {
                  'bg-primary text-white': variant === 'default',
                  'bg-white/10 text-white': variant === 'pills',
                  'text-white border-b-2 border-primary -mb-px': variant === 'underline',
                }
              : {
                  'text-white/60 hover:text-white hover:bg-white/5': variant === 'default' || variant === 'pills',
                  'text-white/60 hover:text-white': variant === 'underline',
                }
          )}
        >
          {tab.icon}
          <span>{tab.label}</span>
        </button>
      ))}
    </div>
  );
}

interface TabPanelProps {
  id: string;
  activeTab: string;
  children: React.ReactNode;
  className?: string;
}

export function TabPanel({ id, activeTab, children, className }: TabPanelProps) {
  if (id !== activeTab) return null;

  return <div className={cn('animate-fade-in', className)}>{children}</div>;
}

export default Tabs;
