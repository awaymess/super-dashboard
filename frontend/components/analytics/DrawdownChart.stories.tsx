import type { Meta, StoryObj } from '@storybook/react';
import { DrawdownChart } from './DrawdownChart';

const meta = {
  title: 'Analytics/DrawdownChart',
  component: DrawdownChart,
  parameters: {
    layout: 'padded',
  },
} satisfies Meta<typeof DrawdownChart>;

export default meta;
type Story = StoryObj<typeof meta>;

const sample = {
  labels: ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug'],
  data: [0, -2.5, -1.2, -4.7, -3.1, -6.2, -2.8, -1.5],
};

export const Default: Story = {
  args: sample,
};

export const MildDrawdown: Story = {
  args: {
    labels: sample.labels,
    data: [0, -0.8, -1.2, -0.5, -1.1, -0.3, -0.9, -0.4],
  },
};
