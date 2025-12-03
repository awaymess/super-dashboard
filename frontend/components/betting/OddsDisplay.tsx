'use client';

import { TrendingUp, TrendingDown, Minus } from 'lucide-react';

interface OddsDisplayProps {
  home: number;
  draw: number;
  away: number;
  previousHome?: number;
  previousDraw?: number;
  previousAway?: number;
  size?: 'sm' | 'md' | 'lg';
}

function OddTrend({ current, previous }: { current: number; previous?: number }) {
  if (!previous) return null;
  const diff = current - previous;
  if (Math.abs(diff) < 0.01) return <Minus className="w-3 h-3 text-gray-400" />;
  if (diff > 0) return <TrendingUp className="w-3 h-3 text-success" />;
  return <TrendingDown className="w-3 h-3 text-danger" />;
}

export function OddsDisplay({
  home,
  draw,
  away,
  previousHome,
  previousDraw,
  previousAway,
  size = 'md',
}: OddsDisplayProps) {
  const sizeClasses = {
    sm: 'text-sm px-3 py-1.5',
    md: 'text-base px-4 py-2',
    lg: 'text-lg px-6 py-3',
  };

  return (
    <div className="flex gap-2">
      <button className={`flex-1 bg-primary/10 hover:bg-primary/20 text-primary rounded-lg font-bold transition-colors flex items-center justify-center gap-1 ${sizeClasses[size]}`}>
        <span>1</span>
        <span>{home.toFixed(2)}</span>
        <OddTrend current={home} previous={previousHome} />
      </button>
      <button className={`flex-1 bg-white/5 hover:bg-white/10 text-gray-300 rounded-lg font-bold transition-colors flex items-center justify-center gap-1 ${sizeClasses[size]}`}>
        <span>X</span>
        <span>{draw.toFixed(2)}</span>
        <OddTrend current={draw} previous={previousDraw} />
      </button>
      <button className={`flex-1 bg-secondary/10 hover:bg-secondary/20 text-secondary rounded-lg font-bold transition-colors flex items-center justify-center gap-1 ${sizeClasses[size]}`}>
        <span>2</span>
        <span>{away.toFixed(2)}</span>
        <OddTrend current={away} previous={previousAway} />
      </button>
    </div>
  );
}
