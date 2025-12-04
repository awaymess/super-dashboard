import type { Meta, StoryObj } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { Header } from './Header';

const meta = {
  title: 'Layout/Header',
  component: Header,
  parameters: {
    layout: 'fullscreen',
  },
  args: {
    onOpenCommandPalette: action('open-command-palette'),
  },
} satisfies Meta<typeof Header>;

export default meta;
type Story = StoryObj<typeof meta>;

function HeaderPreview(args: React.ComponentProps<typeof Header>) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] text-white">
      <Header {...args} />
      <div className="p-10 text-sm text-white/70 space-y-4">
        <p>Use the buttons on the right to toggle the theme, change language, or view notifications.</p>
        <p>The search input opens the command palette when you click it.</p>
      </div>
    </div>
  );
}

export const Default: Story = {
  render: args => <HeaderPreview {...args} />,
};

export const WithCommandPaletteHint: Story = {
  args: {
    onOpenCommandPalette: action('open-command-palette-from-hint'),
  },
  render: args => <HeaderPreview {...args} />,
};
