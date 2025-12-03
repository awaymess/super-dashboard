import React, { forwardRef } from 'react';
import { cn } from '@/lib/cn';

interface GlassButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'default' | 'primary' | 'secondary' | 'success' | 'danger' | 'ghost' | 'outline';
  size?: 'sm' | 'md' | 'lg' | 'icon';
  loading?: boolean;
  icon?: React.ReactNode;
  iconPosition?: 'left' | 'right';
}

export const GlassButton = forwardRef<HTMLButtonElement, GlassButtonProps>(
  (
    {
      className,
      variant = 'default',
      size = 'md',
      loading = false,
      icon,
      iconPosition = 'left',
      disabled,
      children,
      ...props
    },
    ref
  ) => {
    return (
      <button
        ref={ref}
        disabled={disabled || loading}
        className={cn(
          'relative inline-flex items-center justify-center gap-2 font-medium transition-all duration-200',
          'rounded-xl border backdrop-blur-xl',
          'focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-offset-transparent',
          'disabled:opacity-50 disabled:cursor-not-allowed',
          {
            'bg-white/8 border-white/15 text-white hover:bg-white/12 hover:border-white/20 focus:ring-white/30':
              variant === 'default',
            'bg-primary/80 border-primary/50 text-white hover:bg-primary hover:border-primary focus:ring-primary/50':
              variant === 'primary',
            'bg-secondary/80 border-secondary/50 text-white hover:bg-secondary hover:border-secondary focus:ring-secondary/50':
              variant === 'secondary',
            'bg-success/80 border-success/50 text-white hover:bg-success hover:border-success focus:ring-success/50':
              variant === 'success',
            'bg-danger/80 border-danger/50 text-white hover:bg-danger hover:border-danger focus:ring-danger/50':
              variant === 'danger',
            'bg-transparent border-transparent text-white hover:bg-white/8 focus:ring-white/30':
              variant === 'ghost',
            'bg-transparent border-white/20 text-white hover:bg-white/8 hover:border-white/30 focus:ring-white/30':
              variant === 'outline',
            'px-3 py-1.5 text-sm': size === 'sm',
            'px-4 py-2 text-sm': size === 'md',
            'px-6 py-3 text-base': size === 'lg',
            'p-2 aspect-square': size === 'icon',
          },
          className
        )}
        {...props}
      >
        {loading && (
          <svg
            className="animate-spin h-4 w-4"
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
          >
            <circle
              className="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              strokeWidth="4"
            />
            <path
              className="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            />
          </svg>
        )}
        {!loading && icon && iconPosition === 'left' && icon}
        {children}
        {!loading && icon && iconPosition === 'right' && icon}
      </button>
    );
  }
);

GlassButton.displayName = 'GlassButton';

export default GlassButton;
