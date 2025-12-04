import type { Meta, StoryObj } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import React, { createContext } from 'react';
import { Sidebar } from './Sidebar';

const PathnameContext = createContext<string>('/dashboard');

type SidebarStoryProps = React.ComponentProps<typeof Sidebar> & {
  pathname?: string;
};

const meta = {
  title: 'Layout/Sidebar',
  component: SidebarPreview,
  parameters: {
    layout: 'fullscreen',
  },
  argTypes: {
    pathname: {
      control: 'text',
      description: 'Mocked pathname for Storybook preview',
    },
  },
  args: {
    collapsed: false,
    onToggle: action('toggle-sidebar'),
    pathname: '/dashboard',
  },
} satisfies Meta<SidebarStoryProps>;

export default meta;
type Story = StoryObj<typeof meta>;

function SidebarPreview({ pathname = '/dashboard', ...props }: SidebarStoryProps) {
  return (
    <PathnameContext.Provider value={pathname}>
      <div className="min-h-screen bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] text-white">
        <Sidebar {...props} />
        <main
          className="transition-all duration-300 p-10 text-white/80"
          style={{ marginLeft: props.collapsed ? 80 : 280 }}
        >
          <h2 className="text-2xl font-semibold mb-4">Preview Content</h2>
          <p className="max-w-2xl">
            This area represents the main dashboard content. Use the controls to collapse the sidebar or
            change the active route to see the highlighted navigation state update in real time.
          </p>
        </main>
      </div>
    </PathnameContext.Provider>
  );
}

export const Expanded: Story = {
  args: {
    collapsed: false,
    pathname: '/dashboard',
  },
};

export const Collapsed: Story = {
  args: {
    collapsed: true,
    pathname: '/dashboard',
  },
};

export const StocksSectionActive: Story = {
  args: {
    collapsed: false,
    pathname: '/stocks',
  },
};
