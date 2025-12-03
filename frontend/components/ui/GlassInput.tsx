import React, { forwardRef } from 'react';
import { cn } from '@/lib/cn';

interface GlassInputProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'size'> {
  label?: string;
  error?: string;
  hint?: string;
  size?: 'sm' | 'md' | 'lg';
  icon?: React.ReactNode;
  iconPosition?: 'left' | 'right';
  onIconClick?: () => void;
}

export const GlassInput = forwardRef<HTMLInputElement, GlassInputProps>(
  (
    {
      className,
      label,
      error,
      hint,
      size = 'md',
      icon,
      iconPosition = 'left',
      onIconClick,
      id,
      ...props
    },
    ref
  ) => {
    const inputId = id || `input-${Math.random().toString(36).substr(2, 9)}`;

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={inputId}
            className="block text-sm font-medium text-white/70 mb-1.5"
          >
            {label}
          </label>
        )}
        <div className="relative">
          {icon && iconPosition === 'left' && (
            <div
              className={cn(
                'absolute left-3 top-1/2 -translate-y-1/2 text-white/40',
                onIconClick && 'cursor-pointer hover:text-white/60'
              )}
              onClick={onIconClick}
            >
              {icon}
            </div>
          )}
          <input
            ref={ref}
            id={inputId}
            className={cn(
              'w-full rounded-xl border bg-white/5 backdrop-blur-xl',
              'border-white/10 text-white placeholder:text-white/40',
              'transition-all duration-200',
              'focus:outline-none focus:bg-white/8 focus:border-primary',
              'focus:ring-2 focus:ring-primary/20',
              'disabled:opacity-50 disabled:cursor-not-allowed',
              {
                'border-danger focus:border-danger focus:ring-danger/20': error,
                'px-3 py-1.5 text-sm': size === 'sm',
                'px-4 py-2.5 text-sm': size === 'md',
                'px-4 py-3 text-base': size === 'lg',
                'pl-10': icon && iconPosition === 'left',
                'pr-10': icon && iconPosition === 'right',
              },
              className
            )}
            {...props}
          />
          {icon && iconPosition === 'right' && (
            <div
              className={cn(
                'absolute right-3 top-1/2 -translate-y-1/2 text-white/40',
                onIconClick && 'cursor-pointer hover:text-white/60'
              )}
              onClick={onIconClick}
            >
              {icon}
            </div>
          )}
        </div>
        {error && (
          <p className="mt-1.5 text-sm text-danger">{error}</p>
        )}
        {hint && !error && (
          <p className="mt-1.5 text-sm text-white/50">{hint}</p>
        )}
      </div>
    );
  }
);

GlassInput.displayName = 'GlassInput';

export default GlassInput;
