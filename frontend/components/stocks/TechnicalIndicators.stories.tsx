import type { Meta, StoryObj } from '@storybook/react';
import { TechnicalIndicators } from './TechnicalIndicators';

const meta = {
  title: 'Stocks/TechnicalIndicators',
  component: TechnicalIndicators,
  parameters: {
    layout: 'padded',
  },
  args: {
    rsi: 62,
    macd: { value: 1.12, signal: 0.95, histogram: 0.17 },
    sma20: 498.2,
    sma50: 485.2,
    sma200: 412.8,
    currentPrice: 505.8,
    bollingerBands: { upper: 520.5, middle: 505.8, lower: 491.2 },
  },
} satisfies Meta<typeof TechnicalIndicators>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const BearishSignals: Story = {
  args: {
    rsi: 34,
    macd: { value: -0.85, signal: -0.6, histogram: -0.25 },
    sma20: 200.1,
    sma50: 198.4,
    sma200: 205.2,
    currentPrice: 198.0,
    bollingerBands: { upper: 205.0, middle: 198.0, lower: 191.0 },
  },
};
