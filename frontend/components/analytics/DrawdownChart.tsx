'use client';

import { GlassCard } from '@/components/ui';
import { AreaChart } from '@/components/charts';
import { TrendingDown } from 'lucide-react';

interface DrawdownChartProps {
  data: number[];
  labels: string[];
}

export function DrawdownChart({ data, labels }: DrawdownChartProps) {
  const maxDrawdown = Math.min(...data);
  const currentDrawdown = data[data.length - 1] || 0;

  return (
    <GlassCard>
      <div className="flex items-start justify-between mb-4">
        <div className="flex items-center gap-2">
          <TrendingDown className="w-5 h-5 text-danger" />
          <h3 className="font-bold text-white">Drawdown</h3>
        </div>
        <div className="text-right">
          <div className="flex items-center gap-4">
            <div>
              <p className="text-xs text-gray-400">Current</p>
              <p className={`font-bold ${currentDrawdown < 0 ? 'text-danger' : 'text-success'}`}>
                {currentDrawdown.toFixed(2)}%
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-400">Max</p>
              <p className="font-bold text-danger">{maxDrawdown.toFixed(2)}%</p>
            </div>
          </div>
        </div>
      </div>
      
      <AreaChart 
        data={data} 
        labels={labels} 
        gradientFrom="#ef4444" 
        height={200} 
      />
    </GlassCard>
  );
}
