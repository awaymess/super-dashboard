import type { Meta, StoryObj } from '@storybook/react';
import { GlassSelect } from './GlassSelect';

const meta: Meta<typeof GlassSelect> = {
  title: 'UI/GlassSelect',
  component: GlassSelect,
};

export default meta;
type Story = StoryObj<typeof GlassSelect>;

export const Default: Story = {
  args: {
    options: [
      { value: 'aapl', label: 'AAPL' },
      { value: 'msft', label: 'MSFT' },
      { value: 'googl', label: 'GOOGL' },
    ],
    value: 'aapl',
  },
};
