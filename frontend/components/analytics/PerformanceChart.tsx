'use client';

import { GlassCard } from '@/components/ui';
import { AreaChart } from '@/components/charts';

interface PerformanceChartProps {
  title: string;
  data: number[];
  labels: string[];
  color?: string;
}

export function PerformanceChart({
  title,
  data,
  labels,
  color = '#3b82f6',
}: PerformanceChartProps) {
  const currentValue = data[data.length - 1] || 0;
  const startValue = data[0] || 0;
  const change = startValue !== 0 ? ((currentValue - startValue) / startValue) * 100 : 0;
  const isPositive = change >= 0;

  return (
    <GlassCard>
      <div className="flex items-start justify-between mb-4">
        <div>
          <h3 className="font-bold text-white">{title}</h3>
          <p className="text-sm text-gray-400">Performance over time</p>
        </div>
        <div className="text-right">
          <p className="text-2xl font-bold text-white">${currentValue.toLocaleString()}</p>
          <p className={`text-sm ${isPositive ? 'text-success' : 'text-danger'}`}>
            {isPositive ? '+' : ''}{change.toFixed(2)}%
          </p>
        </div>
      </div>
      
      <AreaChart data={data} labels={labels} gradientFrom={color} height={250} />
    </GlassCard>
  );
}
