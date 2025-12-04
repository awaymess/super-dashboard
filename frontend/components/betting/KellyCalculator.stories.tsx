import type { Meta, StoryObj } from '@storybook/react';
import { KellyCalculator } from './KellyCalculator';

const meta = {
  title: 'Betting/KellyCalculator',
  component: KellyCalculator,
  parameters: {
    layout: 'centered',
  },
} satisfies Meta<typeof KellyCalculator>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
