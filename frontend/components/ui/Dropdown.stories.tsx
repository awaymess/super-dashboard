import type { Meta, StoryObj } from '@storybook/react';
import Dropdown, { Dropdown as DropdownComponent, DropdownButton } from './Dropdown';
import { Trash2, Edit, MoreVertical, Settings } from 'lucide-react';

const baseItems = [
  { id: 'edit', label: 'Edit', icon: <Edit className="w-4 h-4" /> },
  { id: 'settings', label: 'Settings', icon: <Settings className="w-4 h-4" /> },
  { id: 'divider-1', label: '', divider: true },
  { id: 'delete', label: 'Delete', icon: <Trash2 className="w-4 h-4" />, danger: true },
  { id: 'disabled', label: 'Disabled', disabled: true },
];

const meta = {
  title: 'UI/Dropdown',
  component: DropdownComponent,
  parameters: { layout: 'padded' },
  args: {
    trigger: (
      <button className="flex items-center gap-2 px-4 py-2 rounded-xl border bg-white/5 border-white/10 text-white hover:bg-white/10">
        Actions <MoreVertical className="w-4 h-4" />
      </button>
    ),
    items: baseItems,
    align: 'left',
  },
} satisfies Meta<typeof DropdownComponent>;

export default meta;
type Story = StoryObj<typeof meta>;

export const LeftAligned: Story = {};

export const RightAligned: Story = {
  args: { align: 'right' },
};

export const WithDividerAndDanger: Story = {
  args: {
    items: baseItems,
  },
};

export const DropdownButtonDefault: Story = {
  render: () => (
    <DropdownButton label="Options" items={baseItems} />
  ),
};

export const DropdownButtonPrimary: Story = {
  render: () => (
    <DropdownButton label="Primary" items={baseItems} variant="primary" />
  ),
};
