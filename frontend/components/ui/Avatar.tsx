import React from 'react';
import { cn } from '@/lib/cn';
import { User } from 'lucide-react';

interface AvatarProps {
  src?: string;
  alt?: string;
  name?: string;
  size?: 'xs' | 'sm' | 'md' | 'lg' | 'xl';
  status?: 'online' | 'offline' | 'away' | 'busy';
  className?: string;
}

const sizeClasses = {
  xs: 'w-6 h-6 text-xs',
  sm: 'w-8 h-8 text-sm',
  md: 'w-10 h-10 text-base',
  lg: 'w-12 h-12 text-lg',
  xl: 'w-16 h-16 text-xl',
};

const statusClasses = {
  online: 'bg-success',
  offline: 'bg-gray-400',
  away: 'bg-warning',
  busy: 'bg-danger',
};

function getInitials(name: string): string {
  const parts = name.split(' ');
  if (parts.length >= 2) {
    return `${parts[0][0]}${parts[1][0]}`.toUpperCase();
  }
  return name.slice(0, 2).toUpperCase();
}

export function Avatar({ src, alt, name, size = 'md', status, className }: AvatarProps) {
  const initials = name ? getInitials(name) : null;

  return (
    <div className={cn('relative inline-flex', className)}>
      <div
        className={cn(
          'relative rounded-full overflow-hidden flex items-center justify-center',
          'bg-gradient-to-br from-primary/50 to-secondary/50',
          'border-2 border-white/10',
          sizeClasses[size]
        )}
      >
        {src ? (
          <img src={src} alt={alt || name || 'Avatar'} className="w-full h-full object-cover" />
        ) : initials ? (
          <span className="font-semibold text-white">{initials}</span>
        ) : (
          <User className="w-1/2 h-1/2 text-white/60" />
        )}
      </div>
      {status && (
        <span
          className={cn(
            'absolute bottom-0 right-0 block rounded-full ring-2 ring-surface',
            statusClasses[status],
            {
              'w-2 h-2': size === 'xs' || size === 'sm',
              'w-2.5 h-2.5': size === 'md',
              'w-3 h-3': size === 'lg',
              'w-4 h-4': size === 'xl',
            }
          )}
        />
      )}
    </div>
  );
}

export default Avatar;
