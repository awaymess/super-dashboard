import type { Meta, StoryObj } from '@storybook/react';
import { StockQuote } from './StockQuote';
import type { Stock } from '@/types/stocks';

const meta = {
  title: 'Stocks/StockQuote',
  component: StockQuote,
  parameters: {
    layout: 'padded',
  },
  args: {
    stock: {
      symbol: 'NVDA',
      name: 'NVIDIA Corporation',
      exchange: 'NASDAQ',
      sector: 'Technology',
      industry: 'Semiconductors',
      marketCap: 1_250_000_000_000,
      price: 505.8,
      change: 7.6,
      changePercent: 1.53,
      volume: 34_120_000,
      avgVolume: 39_500_000,
      high52Week: 570.0,
      low52Week: 390.0,
      pe: 75.2,
      eps: 6.72,
      dividend: 0.16,
      dividendYield: 0.03,
      beta: 1.7,
    } satisfies Stock,
    showDetails: true,
  },
} satisfies Meta<typeof StockQuote>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const IntradayDown: Story = {
  args: {
    stock: {
      symbol: 'NVDA',
      name: 'NVIDIA Corporation',
      exchange: 'NASDAQ',
      sector: 'Technology',
      industry: 'Semiconductors',
      marketCap: 1_250_000_000_000,
      price: 492.15,
      change: -6.05,
      changePercent: -1.21,
      volume: 28_450_000,
      avgVolume: 39_500_000,
      high52Week: 570.0,
      low52Week: 390.0,
      pe: 75.2,
      eps: 6.72,
      dividend: 0.16,
      dividendYield: 0.03,
      beta: 1.7,
    } satisfies Stock,
  },
};
