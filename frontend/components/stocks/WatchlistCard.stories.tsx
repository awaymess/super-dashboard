import type { Meta, StoryObj } from '@storybook/react';
import { WatchlistCard } from './WatchlistCard';
import type { Stock } from '@/types/stocks';

const meta = {
  title: 'Stocks/WatchlistCard',
  component: WatchlistCard,
  parameters: {
    layout: 'padded',
  },
  args: {
    stocks: [
      { symbol: 'NVDA', name: 'NVIDIA', exchange: 'NASDAQ', sector: 'Tech', industry: 'Semiconductors', marketCap: 1_250_000_000_000, price: 505.8, change: 7.6, changePercent: 1.53, volume: 34_120_000, avgVolume: 39_500_000, high52Week: 570.0, low52Week: 390.0, pe: 75.2, eps: 6.72, dividend: 0.16, dividendYield: 0.03, beta: 1.7 } as Stock,
      { symbol: 'AAPL', name: 'Apple', exchange: 'NASDAQ', sector: 'Tech', industry: 'Consumer Electronics', marketCap: 2_980_000_000_000, price: 192.5, change: 1.24, changePercent: 0.65, volume: 51_230_000, avgVolume: 58_400_000, high52Week: 199.62, low52Week: 164.08, pe: 32.4, eps: 6.05, dividend: 0.24, dividendYield: 0.55, beta: 1.2 } as Stock,
      { symbol: 'MSFT', name: 'Microsoft', exchange: 'NASDAQ', sector: 'Tech', industry: 'Software', marketCap: 2_650_000_000_000, price: 377.1, change: 1.58, changePercent: 0.42, volume: 22_100_000, avgVolume: 25_400_000, high52Week: 380.0, low52Week: 305.0, pe: 34.1, eps: 11.3, dividend: 2.72, dividendYield: 0.72, beta: 0.9 } as Stock,
      { symbol: 'AMD', name: 'AMD', exchange: 'NASDAQ', sector: 'Tech', industry: 'Semiconductors', marketCap: 210_000_000_000, price: 132.2, change: -0.5, changePercent: -0.38, volume: 18_900_000, avgVolume: 20_500_000, high52Week: 140.0, low52Week: 95.0, pe: 45.5, eps: 2.9, dividend: 0, dividendYield: 0, beta: 1.4 } as Stock,
    ],
  },
} satisfies Meta<typeof WatchlistCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const MixedPerformance: Story = {
  args: {
    stocks: [
      { symbol: 'TSLA', name: 'Tesla', exchange: 'NASDAQ', sector: 'Automotive', industry: 'EV', marketCap: 630_000_000_000, price: 198.3, change: -2.2, changePercent: -1.12, volume: 120_000_000, avgVolume: 150_000_000, high52Week: 275, low52Week: 165, pe: 70, eps: 2.8, dividend: 0, dividendYield: 0, beta: 2.0 } as Stock,
      { symbol: 'GOOGL', name: 'Alphabet', exchange: 'NASDAQ', sector: 'Tech', industry: 'Internet', marketCap: 1_750_000_000_000, price: 139.4, change: 0.31, changePercent: 0.22, volume: 25_000_000, avgVolume: 28_000_000, high52Week: 145, low52Week: 115, pe: 28, eps: 4.9, dividend: 0, dividendYield: 0, beta: 1.1 } as Stock,
      { symbol: 'META', name: 'Meta', exchange: 'NASDAQ', sector: 'Tech', industry: 'Social Media', marketCap: 900_000_000_000, price: 316.9, change: -1.3, changePercent: -0.41, volume: 19_000_000, avgVolume: 22_000_000, high52Week: 330, low52Week: 240, pe: 24, eps: 13.2, dividend: 0, dividendYield: 0, beta: 1.2 } as Stock,
      { symbol: 'AVGO', name: 'Broadcom', exchange: 'NASDAQ', sector: 'Tech', industry: 'Semiconductors', marketCap: 450_000_000_000, price: 987.4, change: 22.6, changePercent: 2.35, volume: 3_200_000, avgVolume: 3_800_000, high52Week: 1000, low52Week: 700, pe: 28, eps: 35.3, dividend: 18.4, dividendYield: 1.86, beta: 1.1 } as Stock,
    ],
  },
};
