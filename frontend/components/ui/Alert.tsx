import React from 'react';
import { cn } from '@/lib/cn';
import { X, CheckCircle, AlertCircle, AlertTriangle, Info } from 'lucide-react';

interface AlertProps {
  type?: 'success' | 'error' | 'warning' | 'info';
  variant?: 'success' | 'error' | 'warning' | 'info';
  title?: string;
  message?: string;
  children?: React.ReactNode;
  onDismiss?: () => void;
  className?: string;
}

const icons = {
  success: CheckCircle,
  error: AlertCircle,
  warning: AlertTriangle,
  info: Info,
};

const colors = {
  success: 'border-success/30 bg-success/10',
  error: 'border-danger/30 bg-danger/10',
  warning: 'border-warning/30 bg-warning/10',
  info: 'border-primary/30 bg-primary/10',
};

const iconColors = {
  success: 'text-success',
  error: 'text-danger',
  warning: 'text-warning',
  info: 'text-primary',
};

export function Alert({ type, variant, title, message, children, onDismiss, className }: AlertProps) {
  // Support both 'type' and 'variant' props for backwards compatibility
  const alertType = type || variant || 'info';
  const Icon = icons[alertType];

  return (
    <div
      className={cn(
        'flex items-start gap-3 p-4 rounded-xl border',
        colors[alertType],
        className
      )}
    >
      <Icon className={cn('w-5 h-5 mt-0.5 shrink-0', iconColors[alertType])} />
      <div className="flex-1 min-w-0">
        {title && (
          <h4 className="text-sm font-semibold text-white mb-1">{title}</h4>
        )}
        {message && <p className="text-sm text-white/80">{message}</p>}
        {children && <p className="text-sm text-white/80">{children}</p>}
      </div>
      {onDismiss && (
        <button
          onClick={onDismiss}
          className="p-1 rounded-lg text-white/50 hover:text-white hover:bg-white/10 transition-colors"
        >
          <X className="w-4 h-4" />
        </button>
      )}
    </div>
  );
}

export default Alert;
