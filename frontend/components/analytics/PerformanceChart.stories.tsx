import type { Meta, StoryObj } from '@storybook/react';
import { PerformanceChart } from './PerformanceChart';

const meta = {
  title: 'Analytics/PerformanceChart',
  component: PerformanceChart,
  parameters: {
    layout: 'padded',
  },
  args: {
    title: 'Equity Curve',
  },
} satisfies Meta<typeof PerformanceChart>;

export default meta;
type Story = StoryObj<typeof meta>;

const baseArgs = {
  labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4', 'Week 5', 'Week 6', 'Week 7'],
  data: [25000, 26500, 27200, 26000, 27800, 29200, 30500],
};

export const Default: Story = {
  args: {
    ...baseArgs,
    color: '#3b82f6',
  },
};

export const CryptoPortfolio: Story = {
  args: {
    title: 'Crypto Portfolio',
    labels: baseArgs.labels,
    data: [18000, 19500, 21000, 18500, 23000, 24000, 25500],
    color: '#a855f7',
  },
};
