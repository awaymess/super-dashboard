'use client';

interface HeatmapCell {
  x: number;
  y: number;
  value: number;
  label?: string;
}

interface HeatmapProps {
  data: HeatmapCell[];
  xLabels: string[];
  yLabels: string[];
  colorScale?: { min: string; mid: string; max: string };
  cellSize?: number;
}

function getColor(value: number, min: number, max: number, colorScale: { min: string; mid: string; max: string }) {
  const normalized = (value - min) / (max - min);
  
  if (normalized < 0.5) {
    return colorScale.min;
  } else if (normalized < 0.75) {
    return colorScale.mid;
  } else {
    return colorScale.max;
  }
}

export function Heatmap({
  data,
  xLabels,
  yLabels,
  colorScale = { min: '#ef4444', mid: '#f59e0b', max: '#10b981' },
  cellSize = 40,
}: HeatmapProps) {
  const values = data.map(d => d.value);
  const min = Math.min(...values);
  const max = Math.max(...values);

  return (
    <div className="overflow-auto">
      <div className="inline-block">
        <div className="flex">
          <div style={{ width: cellSize * 2 }} />
          {xLabels.map((label, i) => (
            <div
              key={i}
              className="text-center text-xs text-gray-400"
              style={{ width: cellSize }}
            >
              {label}
            </div>
          ))}
        </div>
        
        {yLabels.map((yLabel, y) => (
          <div key={y} className="flex items-center">
            <div
              className="text-xs text-gray-400 truncate"
              style={{ width: cellSize * 2 }}
            >
              {yLabel}
            </div>
            {xLabels.map((_, x) => {
              const cell = data.find(d => d.x === x && d.y === y);
              const value = cell?.value ?? 0;
              const color = getColor(value, min, max, colorScale);
              
              return (
                <div
                  key={x}
                  className="flex items-center justify-center text-xs font-medium rounded-sm m-0.5 hover:scale-110 transition-transform cursor-pointer"
                  style={{
                    width: cellSize - 4,
                    height: cellSize - 4,
                    backgroundColor: `${color}66`,
                    color: color,
                  }}
                  title={cell?.label || `${value}`}
                >
                  {value.toFixed(0)}
                </div>
              );
            })}
          </div>
        ))}
      </div>
    </div>
  );
}
