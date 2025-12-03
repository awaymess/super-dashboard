import React from 'react';
import { cn } from '@/lib/cn';

interface SkeletonProps {
  variant?: 'text' | 'circular' | 'rectangular' | 'rounded';
  width?: string | number;
  height?: string | number;
  className?: string;
}

export function Skeleton({ variant = 'text', width, height, className }: SkeletonProps) {
  return (
    <div
      className={cn(
        'animate-pulse bg-white/10',
        {
          'h-4 rounded': variant === 'text',
          'rounded-full aspect-square': variant === 'circular',
          'rounded-none': variant === 'rectangular',
          'rounded-lg': variant === 'rounded',
        },
        className
      )}
      style={{
        width: width,
        height: variant !== 'circular' ? height : undefined,
      }}
    />
  );
}

interface SkeletonCardProps {
  lines?: number;
  hasImage?: boolean;
  className?: string;
}

export function SkeletonCard({ lines = 3, hasImage = false, className }: SkeletonCardProps) {
  return (
    <div className={cn('rounded-xl border border-white/10 bg-white/5 p-5', className)}>
      {hasImage && <Skeleton variant="rounded" className="w-full h-40 mb-4" />}
      <Skeleton className="w-3/4 h-5 mb-3" />
      {[...Array(lines)].map((_, i) => (
        <Skeleton key={i} className={cn('h-3 mb-2', i === lines - 1 && 'w-1/2')} />
      ))}
    </div>
  );
}

interface SkeletonTableProps {
  rows?: number;
  columns?: number;
  className?: string;
}

export function SkeletonTable({ rows = 5, columns = 4, className }: SkeletonTableProps) {
  return (
    <div className={cn('rounded-xl border border-white/10 bg-white/5 overflow-hidden', className)}>
      <div className="border-b border-white/10 p-4 flex gap-4">
        {[...Array(columns)].map((_, i) => (
          <Skeleton key={i} className="flex-1 h-4" />
        ))}
      </div>
      {[...Array(rows)].map((_, rowIndex) => (
        <div key={rowIndex} className="border-b border-white/5 p-4 flex gap-4">
          {[...Array(columns)].map((_, colIndex) => (
            <Skeleton key={colIndex} className="flex-1 h-4" />
          ))}
        </div>
      ))}
    </div>
  );
}

export default Skeleton;
