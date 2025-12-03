import React, { forwardRef, useId } from 'react';
import { cn } from '@/lib/cn';
import { SelectOption } from '@/types/common';

interface GlassSelectProps extends Omit<React.SelectHTMLAttributes<HTMLSelectElement>, 'size'> {
  label?: string;
  error?: string;
  hint?: string;
  size?: 'sm' | 'md' | 'lg';
  options: SelectOption[];
  placeholder?: string;
}

export const GlassSelect = forwardRef<HTMLSelectElement, GlassSelectProps>(
  (
    {
      className,
      label,
      error,
      hint,
      size = 'md',
      options,
      placeholder,
      id,
      ...props
    },
    ref
  ) => {
    const generatedId = useId();
    const selectId = id || generatedId;

    return (
      <div className="w-full">
        {label && (
          <label
            htmlFor={selectId}
            className="block text-sm font-medium text-white/70 mb-1.5"
          >
            {label}
          </label>
        )}
        <div className="relative">
          <select
            ref={ref}
            id={selectId}
            className={cn(
              'w-full rounded-xl border bg-white/5 backdrop-blur-xl appearance-none',
              'border-white/10 text-white',
              'transition-all duration-200',
              'focus:outline-none focus:bg-white/8 focus:border-primary',
              'focus:ring-2 focus:ring-primary/20',
              'disabled:opacity-50 disabled:cursor-not-allowed',
              {
                'border-danger focus:border-danger focus:ring-danger/20': error,
                'px-3 py-1.5 text-sm pr-8': size === 'sm',
                'px-4 py-2.5 text-sm pr-10': size === 'md',
                'px-4 py-3 text-base pr-10': size === 'lg',
              },
              className
            )}
            {...props}
          >
            {placeholder && (
              <option value="" className="bg-surface text-white/50">
                {placeholder}
              </option>
            )}
            {options.map((option) => (
              <option
                key={option.value}
                value={option.value}
                disabled={option.disabled}
                className="bg-surface text-white"
              >
                {option.label}
              </option>
            ))}
          </select>
          <div className="absolute right-3 top-1/2 -translate-y-1/2 pointer-events-none text-white/40">
            <svg
              className="w-4 h-4"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M19 9l-7 7-7-7"
              />
            </svg>
          </div>
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

GlassSelect.displayName = 'GlassSelect';

export default GlassSelect;
