import type { Meta, StoryObj } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { TradeForm } from './TradeForm';

const meta = {
  title: 'PaperTrading/TradeForm',
  component: TradeForm,
  parameters: {
    layout: 'centered',
  },
  args: {
    symbol: 'AAPL',
    currentPrice: 192.45,
    onSubmit: action('submit-trade'),
  },
} satisfies Meta<typeof TradeForm>;

export default meta;
type Story = StoryObj<typeof meta>;

export const BuyMarket: Story = {};

export const SellLimit: Story = {
  args: {
    symbol: 'TSLA',
    currentPrice: 205.12,
  },
  render: args => <TradeForm {...args} />,
};
