import type { Meta, StoryObj } from '@storybook/react';
import { LineChart } from './LineChart';

const meta = {
  title: 'Charts/LineChart',
  component: LineChart,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: {
    labels: Array.from({ length: 12 }, (_, i) => `Q${i + 1}`),
    data: [24, 35, 42, 31, 45, 58, 62, 55, 68, 75, 80, 92],
    label: 'Signals',
    color: '#10b981',
    fill: true,
  },
} satisfies Meta<typeof LineChart>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Filled: Story = {};

export const NoFill: Story = {
  args: {
    fill: false,
    color: '#f97316',
  },
};
