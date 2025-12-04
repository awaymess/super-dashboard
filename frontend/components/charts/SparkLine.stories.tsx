import type { Meta, StoryObj } from '@storybook/react';
import { SparkLine } from './SparkLine';

const meta = {
  title: 'Charts/SparkLine',
  component: SparkLine,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: {
    data: [42, 44, 41, 47, 52, 55, 53, 58],
    showChange: true,
  },
} satisfies Meta<typeof SparkLine>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Positive: Story = {};

export const Negative: Story = {
  args: {
    data: [58, 53, 55, 52, 47, 41, 44, 42],
  },
};
