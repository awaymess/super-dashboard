import type { Meta, StoryObj } from '@storybook/react';
import { Avatar } from './Avatar';

const meta: Meta<typeof Avatar> = {
  title: 'UI/Avatar',
  component: Avatar,
};

export default meta;
type Story = StoryObj<typeof Avatar>;

export const Default: Story = {
  args: {
    name: 'Jane Doe',
    src: 'https://i.pravatar.cc/150?img=3',
  },
};

export const Initials: Story = {
  args: {
    name: 'John Smith',
  },
};
