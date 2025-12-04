import type { Meta, StoryObj } from '@storybook/react';
import { PoissonCalculator } from './PoissonCalculator';

const meta = {
  title: 'Betting/PoissonCalculator',
  component: PoissonCalculator,
  parameters: {
    layout: 'centered',
  },
} satisfies Meta<typeof PoissonCalculator>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
