import type { Meta, StoryObj } from '@storybook/react';
import { DoughnutChart } from './DoughnutChart';

const meta = {
  title: 'Charts/DoughnutChart',
  component: DoughnutChart,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: {
    labels: ['Equities', 'Crypto', 'Cash', 'Derivatives'],
    data: [42, 28, 18, 12],
    centerText: '$120k',
  },
} satisfies Meta<typeof DoughnutChart>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const ExtendedAssetClasses: Story = {
  args: {
    labels: ['AI', 'Green Energy', 'Fintech', 'Defense', 'Space'],
    data: [20, 25, 15, 18, 22],
    centerText: '5 Themes',
  },
};
