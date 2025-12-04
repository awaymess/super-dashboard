import type { Meta, StoryObj } from '@storybook/react';
import { PortfolioSummary } from './PortfolioSummary';
import type { Portfolio } from '@/types/paper-trading';

const portfolio: Portfolio = {
  id: 'pf-1',
  name: 'Alpha Paper',
  initialBalance: 100000,
  currentBalance: 126400,
  totalValue: 133850,
  cashBalance: 24350,
  positions: [],
  transactions: Array.from({ length: 68 }).map((_, i) => ({
    id: `tx-${i}`,
    symbol: i % 2 ? 'AAPL' : 'NVDA',
    type: i % 3 ? 'buy' : 'sell',
    quantity: 10 + (i % 5) * 5,
    price: 240 + i,
    total: 2500,
    fees: 1.5,
    executedAt: '2025-11-25T00:00:00Z',
  })),
  performance: {
    totalReturn: 23850,
    totalReturnPercent: 23.85,
    dayReturn: 1450,
    dayReturnPercent: 1.12,
    weekReturn: 4320,
    weekReturnPercent: 3.3,
    monthReturn: 8120,
    monthReturnPercent: 6.7,
    yearReturn: 23850,
    yearReturnPercent: 23.85,
    sharpeRatio: 1.42,
    sortinoRatio: 1.71,
    maxDrawdown: -12.4,
    winRate: 61.5,
    avgWin: 450,
    avgLoss: -210,
    profitFactor: 2.1,
  },
  createdAt: '2025-03-01T00:00:00Z',
  updatedAt: '2025-12-01T00:00:00Z',
};

const meta = {
  title: 'PaperTrading/PortfolioSummary',
  component: PortfolioSummary,
  parameters: {
    layout: 'padded',
  },
  args: {
    portfolio,
  },
} satisfies Meta<typeof PortfolioSummary>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
