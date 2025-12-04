import type { Meta, StoryObj } from '@storybook/react';
import { StockCard } from './StockCard';
import type { Stock } from '@/types/stocks';

const meta = {
  title: 'Stocks/StockCard',
  component: StockCard,
  parameters: {
    layout: 'padded',
  },
  args: {
    stock: {
      symbol: 'AAPL',
      name: 'Apple Inc.',
      exchange: 'NASDAQ',
      sector: 'Technology',
      industry: 'Consumer Electronics',
      marketCap: 2_980_000_000_000,
      price: 192.45,
      change: 1.24,
      changePercent: 0.65,
      volume: 51_230_000,
      avgVolume: 58_400_000,
      high52Week: 199.62,
      low52Week: 164.08,
      pe: 32.4,
      eps: 6.05,
      dividend: 0.24,
      dividendYield: 0.55,
      beta: 1.2,
    } satisfies Stock,
  },
} satisfies Meta<typeof StockCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const NegativeChange: Story = {
  args: {
    stock: {
      symbol: 'AAPL',
      name: 'Apple Inc.',
      exchange: 'NASDAQ',
      sector: 'Technology',
      industry: 'Consumer Electronics',
      marketCap: 2_950_000_000_000,
      price: 188.12,
      change: -2.15,
      changePercent: -1.13,
      volume: 49_000_000,
      avgVolume: 58_400_000,
      high52Week: 199.62,
      low52Week: 164.08,
      pe: 32.1,
      eps: 6.05,
      dividend: 0.24,
      dividendYield: 0.55,
      beta: 1.2,
    } satisfies Stock,
  },
};
