import type { Meta, StoryObj } from '@storybook/react';
import { BarChart } from './BarChart';

const baseArgs = {
  labels: ['Value', 'Betting', 'Stocks', 'Paper', 'Analytics', 'News', 'Watchlist'],
  data: [68, 54, 72, 41, 83, 35, 64],
  label: 'Engagement',
};

const meta = {
  title: 'Charts/BarChart',
  component: BarChart,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: baseArgs,
} satisfies Meta<typeof BarChart>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Vertical: Story = {};

export const Horizontal: Story = {
  args: {
    ...baseArgs,
    horizontal: true,
    data: [18, 32, 55, 62, 40, 28, 50],
  },
};

export const CustomPalette: Story = {
  args: {
    ...baseArgs,
    colors: ['#0ea5e9', '#22d3ee', '#a3e635', '#f97316', '#fb7185', '#c084fc', '#38bdf8'],
  },
};
