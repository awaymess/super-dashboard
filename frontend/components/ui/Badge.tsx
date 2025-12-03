import React from 'react';
import { cn } from '@/lib/cn';

interface BadgeProps {
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'warning' | 'danger' | 'outline';
  size?: 'sm' | 'md' | 'lg';
  dot?: boolean;
  icon?: React.ReactNode;
  children: React.ReactNode;
  className?: string;
}

export function Badge({
  variant = 'default',
  size = 'md',
  dot = false,
  icon,
  children,
  className,
}: BadgeProps) {
  return (
    <span
      className={cn(
        'inline-flex items-center gap-1 font-medium rounded-full',
        {
          'bg-white/10 text-white': variant === 'default',
          'bg-primary/20 text-primary': variant === 'primary',
          'bg-secondary/20 text-secondary': variant === 'secondary',
          'bg-success/20 text-success': variant === 'success',
          'bg-warning/20 text-warning': variant === 'warning',
          'bg-danger/20 text-danger': variant === 'danger',
          'border border-white/20 bg-transparent text-white': variant === 'outline',
          'px-2 py-0.5 text-xs': size === 'sm',
          'px-2.5 py-1 text-xs': size === 'md',
          'px-3 py-1.5 text-sm': size === 'lg',
        },
        className
      )}
    >
      {dot && (
        <span
          className={cn('w-1.5 h-1.5 rounded-full', {
            'bg-white': variant === 'default' || variant === 'outline',
            'bg-primary': variant === 'primary',
            'bg-secondary': variant === 'secondary',
            'bg-success': variant === 'success',
            'bg-warning': variant === 'warning',
            'bg-danger': variant === 'danger',
          })}
        />
      )}
      {icon}
      {children}
    </span>
  );
}

export default Badge;
