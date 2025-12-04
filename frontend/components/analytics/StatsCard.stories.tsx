import type { Meta, StoryObj } from '@storybook/react';
import { TrendingUp, TrendingDown, Trophy } from 'lucide-react';
import { StatsCard } from './StatsCard';

const meta = {
  title: 'Analytics/StatsCard',
  component: StatsCard,
  parameters: {
    layout: 'padded',
  },
  argTypes: {
    icon: {
      control: false,
    },
    color: {
      control: 'select',
      options: ['primary', 'secondary', 'success', 'warning', 'danger'],
    },
  },
} satisfies Meta<typeof StatsCard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const PositiveGrowth: Story = {
  args: {
    title: 'Monthly Return',
    value: '$32,480',
    change: 8.43,
    changeLabel: 'vs last month',
    icon: TrendingUp,
    color: 'success',
  },
};

export const NegativeGrowth: Story = {
  args: {
    title: 'Daily P&L',
    value: '-$1,240',
    change: -3.21,
    changeLabel: 'vs previous day',
    icon: TrendingDown,
    color: 'danger',
  },
};

export const NeutralTrophy: Story = {
  args: {
    title: 'Best Streak',
    value: '12 days',
    change: 0,
    changeLabel: 'stable',
    icon: Trophy,
    color: 'secondary',
  },
};
