'use client';

import { Line } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
} from 'chart.js';

ChartJS.register(CategoryScale, LinearScale, PointElement, LineElement);

interface SparkLineProps {
  data: number[];
  color?: string;
  width?: number;
  height?: number;
  showChange?: boolean;
}

export function SparkLine({
  data,
  color,
  width = 100,
  height = 32,
  showChange = false,
}: SparkLineProps) {
  const isPositive = data.length > 1 && data[data.length - 1] >= data[0];
  const lineColor = color || (isPositive ? '#10b981' : '#ef4444');

  const chartData = {
    labels: data.map((_, i) => i.toString()),
    datasets: [
      {
        data,
        borderColor: lineColor,
        borderWidth: 2,
        fill: false,
        tension: 0.4,
        pointRadius: 0,
      },
    ],
  };

  const options = {
    responsive: false,
    maintainAspectRatio: false,
    plugins: {
      legend: { display: false },
      tooltip: { enabled: false },
    },
    scales: {
      x: { display: false },
      y: { display: false },
    },
  };

  const changePercent = data.length > 1
    ? ((data[data.length - 1] - data[0]) / data[0]) * 100
    : 0;

  return (
    <div className="flex items-center gap-2">
      <div style={{ width, height }}>
        <Line data={chartData} options={options} width={width} height={height} />
      </div>
      {showChange && (
        <span className={`text-sm font-medium ${isPositive ? 'text-success' : 'text-danger'}`}>
          {isPositive ? '+' : ''}{changePercent.toFixed(2)}%
        </span>
      )}
    </div>
  );
}
