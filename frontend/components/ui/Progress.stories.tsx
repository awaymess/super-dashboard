import type { Meta, StoryObj } from '@storybook/react';
import { Progress } from './Progress';

const meta: Meta<typeof Progress> = {
  title: 'UI/Progress',
  component: Progress,
};

export default meta;
type Story = StoryObj<typeof Progress>;

export const Default: Story = {
  args: {
    value: 50,
  },
};

export const Indeterminate: Story = {
  args: {
    indeterminate: true,
  },
};
