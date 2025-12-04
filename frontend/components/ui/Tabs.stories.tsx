import type { Meta, StoryObj } from '@storybook/react';
import { Tabs } from './Tabs';

const meta: Meta<typeof Tabs> = {
  title: 'UI/Tabs',
  component: Tabs,
};

export default meta;
type Story = StoryObj<typeof Tabs>;

export const Default: Story = {
  args: {
    tabs: [
      { key: 'overview', label: 'Overview' },
      { key: 'stats', label: 'Stats' },
      { key: 'settings', label: 'Settings' },
    ],
    value: 'overview',
  },
};
