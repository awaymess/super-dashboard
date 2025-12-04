import type { Meta, StoryObj } from '@storybook/react';
import { TradeJournal } from './TradeJournal';

const trades = [
  {
    id: 'trade-1',
    symbol: 'NVDA',
    side: 'buy' as const,
    orderType: 'limit' as const,
    quantity: 20,
    price: 480.5,
    pnl: 420.3,
    notes: 'Bought on breakout above 52-week high.',
    timestamp: '2025-11-29 09:42',
  },
  {
    id: 'trade-2',
    symbol: 'TSLA',
    side: 'sell' as const,
    orderType: 'market' as const,
    quantity: 15,
    price: 202.1,
    pnl: -185.9,
    notes: 'Stopped out after failed bounce near 200 MA.',
    timestamp: '2025-11-30 14:05',
  },
  {
    id: 'trade-3',
    symbol: 'AAPL',
    side: 'buy' as const,
    orderType: 'market' as const,
    quantity: 30,
    price: 192.45,
    pnl: 120.0,
    timestamp: '2025-12-01 10:12',
  },
];

const meta = {
  title: 'PaperTrading/TradeJournal',
  component: TradeJournal,
  parameters: {
    layout: 'padded',
  },
  args: {
    trades,
  },
} satisfies Meta<typeof TradeJournal>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const LongJournal: Story = {
  args: {
    trades: Array.from({ length: 10 }).map((_, i) => ({
      ...trades[i % trades.length],
      id: `trade-${i}`,
      timestamp: `2025-11-${10 + i} ${(9 + i) % 24}:15`,
    })),
  },
};
