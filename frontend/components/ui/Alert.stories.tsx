import type { Meta, StoryObj } from '@storybook/react';
import { Alert } from './Alert';

const meta: Meta<typeof Alert> = {
  title: 'UI/Alert',
  component: Alert,
};

export default meta;
type Story = StoryObj<typeof Alert>;

export const Info: Story = {
  args: {
    title: 'Information',
    message: 'This is an info alert.',
    variant: 'info',
  },
};

export const Success: Story = {
  args: {
    title: 'Success',
    message: 'Operation completed successfully.',
    variant: 'success',
  },
};

export const Error: Story = {
  args: {
    title: 'Error',
    message: 'Something went wrong.',
    variant: 'error',
  },
};
