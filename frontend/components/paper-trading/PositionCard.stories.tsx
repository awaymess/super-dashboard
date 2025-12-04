import type { Meta, StoryObj } from '@storybook/react';
import { PositionCard } from './PositionCard';
import type { Position } from '@/types/paper-trading';
import { action } from '@storybook/addon-actions';

const position: Position = {
  id: 'pos-1',
  symbol: 'NVDA',
  name: 'NVIDIA Corp',
  quantity: 42,
  avgCost: 412.35,
  currentPrice: 505.8,
  marketValue: 21243.6,
  unrealizedPL: 3920.7,
  unrealizedPLPercent: 22.6,
  dayChange: 210.4,
  dayChangePercent: 1.02,
  weight: 0.18,
  openedAt: '2025-08-19T00:00:00Z',
};

const negativePosition: Position = {
  ...position,
  id: 'pos-2',
  symbol: 'TSLA',
  name: 'Tesla Inc',
  currentPrice: 198.3,
  avgCost: 245.2,
  marketValue: 10411.6,
  unrealizedPL: -1975.4,
  unrealizedPLPercent: -15.9,
};

const meta = {
  title: 'PaperTrading/PositionCard',
  component: PositionCard,
  parameters: {
    layout: 'centered',
  },
  args: {
    position,
    onClose: action('close-position'),
  },
} satisfies Meta<typeof PositionCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Positive: Story = {};

export const Negative: Story = {
  args: {
    position: negativePosition,
  },
};
