'use client';

import React, { useState } from 'react';
import { cn } from '@/lib/cn';

interface TooltipProps {
  content: React.ReactNode;
  children: React.ReactNode;
  position?: 'top' | 'bottom' | 'left' | 'right';
  delay?: number;
  className?: string;
}

export function Tooltip({
  content,
  children,
  position = 'top',
  delay = 200,
  className,
}: TooltipProps) {
  const [isVisible, setIsVisible] = useState(false);
  const [timeoutId, setTimeoutId] = useState<NodeJS.Timeout | null>(null);

  const showTooltip = () => {
    const id = setTimeout(() => setIsVisible(true), delay);
    setTimeoutId(id);
  };

  const hideTooltip = () => {
    if (timeoutId) {
      clearTimeout(timeoutId);
    }
    setIsVisible(false);
  };

  return (
    <div
      className="relative inline-flex"
      onMouseEnter={showTooltip}
      onMouseLeave={hideTooltip}
      onFocus={showTooltip}
      onBlur={hideTooltip}
    >
      {children}
      {isVisible && (
        <div
          className={cn(
            'absolute z-50 px-3 py-1.5 text-xs font-medium text-white rounded-lg',
            'bg-surface border border-white/10 shadow-glass-sm backdrop-blur-xl',
            'animate-fade-in whitespace-nowrap',
            {
              'bottom-full left-1/2 -translate-x-1/2 mb-2': position === 'top',
              'top-full left-1/2 -translate-x-1/2 mt-2': position === 'bottom',
              'right-full top-1/2 -translate-y-1/2 mr-2': position === 'left',
              'left-full top-1/2 -translate-y-1/2 ml-2': position === 'right',
            },
            className
          )}
        >
          {content}
          <div
            className={cn(
              'absolute w-2 h-2 bg-surface border-white/10 rotate-45',
              {
                'top-full left-1/2 -translate-x-1/2 -mt-1 border-b border-r': position === 'top',
                'bottom-full left-1/2 -translate-x-1/2 mb-[-4px] border-t border-l': position === 'bottom',
                'left-full top-1/2 -translate-y-1/2 -ml-1 border-t border-r': position === 'left',
                'right-full top-1/2 -translate-y-1/2 mr-[-4px] border-b border-l': position === 'right',
              }
            )}
          />
        </div>
      )}
    </div>
  );
}

export default Tooltip;
