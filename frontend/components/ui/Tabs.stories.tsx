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
      { id: 'overview', label: 'Overview' },
      { id: 'stats', label: 'Stats' },
      { id: 'settings', label: 'Settings' },
    ],
    activeTab: 'overview',
    onChange: () => {},
  },
};
