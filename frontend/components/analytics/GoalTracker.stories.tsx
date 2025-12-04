import type { Meta, StoryObj } from '@storybook/react';
import { GoalTracker } from './GoalTracker';

const meta = {
  title: 'Analytics/GoalTracker',
  component: GoalTracker,
  parameters: {
    layout: 'padded',
  },
} satisfies Meta<typeof GoalTracker>;

export default meta;
type Story = StoryObj<typeof meta>;

const baseGoals = [
  {
    id: '1',
    title: 'Grow paper portfolio to $50k',
    target: 50000,
    current: 32000,
    unit: 'USD',
    deadline: 'Q2 2026',
  },
  {
    id: '2',
    title: 'Lock 200 value bets',
    target: 200,
    current: 145,
    unit: 'bets',
    deadline: 'Mar 2026',
  },
  {
    id: '3',
    title: 'Research 25 stocks',
    target: 25,
    current: 25,
    unit: 'tickers',
    deadline: 'Jan 2026',
  },
];

export const Default: Story = {
  args: {
    goals: baseGoals,
  },
};

export const AllCompleted: Story = {
  args: {
    goals: baseGoals.map(goal => ({
      ...goal,
      current: goal.target,
    })),
  },
};
