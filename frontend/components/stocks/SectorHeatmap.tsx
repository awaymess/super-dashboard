'use client';

import { GlassCard } from '@/components/ui';

interface SectorData {
  name: string;
  change: number;
}

interface SectorHeatmapProps {
  sectors: SectorData[];
}

export function SectorHeatmap({ sectors }: SectorHeatmapProps) {
  const maxChange = Math.max(...sectors.map(s => Math.abs(s.change)));

  const getColor = (change: number) => {
    const intensity = Math.abs(change) / maxChange;
    if (change >= 0) {
      return `rgba(16, 185, 129, ${0.2 + intensity * 0.6})`;
    }
    return `rgba(239, 68, 68, ${0.2 + intensity * 0.6})`;
  };

  const sortedSectors = [...sectors].sort((a, b) => b.change - a.change);

  return (
    <GlassCard>
      <h3 className="font-bold text-white mb-4">Sector Performance</h3>
      
      <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-2">
        {sortedSectors.map(sector => (
          <div
            key={sector.name}
            className="p-3 rounded-lg transition-transform hover:scale-105 cursor-pointer"
            style={{ backgroundColor: getColor(sector.change) }}
          >
            <p className="text-sm font-medium text-white truncate">{sector.name}</p>
            <p className={`text-lg font-bold ${sector.change >= 0 ? 'text-success' : 'text-danger'}`}>
              {sector.change >= 0 ? '+' : ''}{sector.change.toFixed(2)}%
            </p>
          </div>
        ))}
      </div>
    </GlassCard>
  );
}
