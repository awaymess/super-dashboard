'use client';

import React, { useState, useRef, useEffect } from 'react';
import { cn } from '@/lib/cn';
import { ChevronDown } from 'lucide-react';

interface DropdownItem {
  id: string;
  label: string;
  icon?: React.ReactNode;
  onClick?: () => void;
  disabled?: boolean;
  danger?: boolean;
  divider?: boolean;
}

interface DropdownProps {
  trigger: React.ReactNode;
  items: DropdownItem[];
  align?: 'left' | 'right';
  className?: string;
}

export function Dropdown({ trigger, items, align = 'left', className }: DropdownProps) {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  return (
    <div className={cn('relative inline-block', className)} ref={dropdownRef}>
      <div onClick={() => setIsOpen(!isOpen)}>{trigger}</div>
      {isOpen && (
        <div
          className={cn(
            'absolute z-50 mt-2 min-w-[180px] py-1',
            'bg-surface border border-white/10 rounded-xl shadow-glass-lg backdrop-blur-xl',
            'animate-fade-in',
            {
              'left-0': align === 'left',
              'right-0': align === 'right',
            }
          )}
        >
          {items.map((item) =>
            item.divider ? (
              <div key={item.id} className="my-1 border-t border-white/10" />
            ) : (
              <button
                key={item.id}
                onClick={() => {
                  item.onClick?.();
                  setIsOpen(false);
                }}
                disabled={item.disabled}
                className={cn(
                  'w-full flex items-center gap-2 px-4 py-2 text-sm text-left transition-colors',
                  item.disabled
                    ? 'opacity-50 cursor-not-allowed text-white/50'
                    : item.danger
                    ? 'text-danger hover:bg-danger/10'
                    : 'text-white hover:bg-white/5'
                )}
              >
                {item.icon}
                {item.label}
              </button>
            )
          )}
        </div>
      )}
    </div>
  );
}

interface DropdownButtonProps {
  label: string;
  items: DropdownItem[];
  variant?: 'default' | 'primary';
  className?: string;
}

export function DropdownButton({ label, items, variant = 'default', className }: DropdownButtonProps) {
  return (
    <Dropdown
      trigger={
        <button
          className={cn(
            'flex items-center gap-2 px-4 py-2 rounded-xl border font-medium transition-colors',
            variant === 'default'
              ? 'bg-white/5 border-white/10 text-white hover:bg-white/10'
              : 'bg-primary border-primary text-white hover:bg-primary-600'
          )}
        >
          {label}
          <ChevronDown className="w-4 h-4" />
        </button>
      }
      items={items}
      className={className}
    />
  );
}

export default Dropdown;
