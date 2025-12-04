import type { Meta, StoryObj } from '@storybook/react';
import { Leaderboard } from './Leaderboard';

const entries = [
  { rank: 1, username: 'quantumtrader', totalReturn: 18450, totalReturnPercent: 46.2, winRate: 68, totalTrades: 142, sharpeRatio: 2.3 },
  { rank: 2, username: 'alphaoracle', totalReturn: 15480, totalReturnPercent: 38.1, winRate: 62, totalTrades: 118, sharpeRatio: 1.9 },
  { rank: 3, username: 'gammaflow', totalReturn: 12340, totalReturnPercent: 31.7, winRate: 58, totalTrades: 133, sharpeRatio: 1.6 },
  { rank: 4, username: 'you', totalReturn: 9800, totalReturnPercent: 25.4, winRate: 55, totalTrades: 102, sharpeRatio: 1.4 },
  { rank: 5, username: 'betabase', totalReturn: 7650, totalReturnPercent: 19.8, winRate: 52, totalTrades: 89, sharpeRatio: 1.2 },
];

const meta = {
  title: 'PaperTrading/Leaderboard',
  component: Leaderboard,
  parameters: {
    layout: 'padded',
  },
  args: {
    entries,
  },
} satisfies Meta<typeof Leaderboard>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const HighlightCurrentUser: Story = {
  args: {
    currentUser: 'you',
  },
};
