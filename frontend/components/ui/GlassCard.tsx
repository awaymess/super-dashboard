import React from 'react';
import { cn } from '@/lib/cn';

interface GlassCardProps extends React.HTMLAttributes<HTMLDivElement> {
  variant?: 'default' | 'elevated' | 'outlined';
  padding?: 'none' | 'sm' | 'md' | 'lg';
  hover?: boolean;
  glow?: boolean;
  glowColor?: 'primary' | 'success' | 'danger' | 'warning';
}

export function GlassCard({
  className,
  variant = 'default',
  padding = 'md',
  hover = false,
  glow = false,
  glowColor = 'primary',
  children,
  ...props
}: GlassCardProps) {
  return (
    <div
      className={cn(
        'rounded-2xl border transition-all duration-200',
        'bg-white/5 backdrop-blur-xl',
        'border-white/10',
        {
          'shadow-glass': variant === 'default',
          'shadow-glass-lg': variant === 'elevated',
          'border-white/20': variant === 'outlined',
          'p-0': padding === 'none',
          'p-3': padding === 'sm',
          'p-5': padding === 'md',
          'p-7': padding === 'lg',
          'hover:bg-white/8 hover:border-white/15 hover:-translate-y-0.5': hover,
          'shadow-glow': glow && glowColor === 'primary',
          'shadow-glow-success': glow && glowColor === 'success',
          'shadow-glow-danger': glow && glowColor === 'danger',
        },
        className
      )}
      {...props}
    >
      {children}
    </div>
  );
}

export default GlassCard;
