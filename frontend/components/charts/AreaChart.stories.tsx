import type { Meta, StoryObj } from '@storybook/react';
import { AreaChart } from './AreaChart';

const meta = {
  title: 'Charts/AreaChart',
  component: AreaChart,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: {
    labels: ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'],
    data: [12000, 13500, 12800, 15000, 16200, 17500, 19000],
    label: 'PnL',
    gradientFrom: '#3b82f6',
  },
} satisfies Meta<typeof AreaChart>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const SecondaryGradient: Story = {
  args: {
    gradientFrom: '#a855f7',
    data: [100, 120, 140, 180, 210, 230, 250],
    label: 'Engagement',
  },
};
