import type { Meta, StoryObj } from '@storybook/react';
import { AnalystRatings } from './AnalystRatings';

const meta = {
  title: 'Stocks/AnalystRatings',
  component: AnalystRatings,
  parameters: {
    layout: 'padded',
  },
  args: {
    buy: 24,
    hold: 15,
    sell: 3,
    priceTarget: {
      low: 165,
      average: 195,
      high: 225,
      current: 192.45,
    },
  },
} satisfies Meta<typeof AnalystRatings>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const MixedConsensus: Story = {
  args: {
    buy: 12,
    hold: 18,
    sell: 8,
    priceTarget: {
      low: 170,
      average: 185,
      high: 210,
      current: 198.3,
    },
  },
};
