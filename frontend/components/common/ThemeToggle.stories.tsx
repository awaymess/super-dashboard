import type { Meta, StoryObj } from '@storybook/react';
import { ThemeToggle } from './ThemeToggle';
import { useTheme } from '@/hooks/useTheme';

const meta = {
  title: 'Common/ThemeToggle',
  component: ThemeToggle,
  parameters: {
    layout: 'fullscreen',
  },
} satisfies Meta<typeof ThemeToggle>;

export default meta;
type Story = StoryObj<typeof meta>;

function ThemeTogglePreview() {
  const { theme } = useTheme();

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#04060d] via-[#0d1324] to-[#04060d] flex items-center justify-center text-white">
      <div className="bg-white/5 border border-white/10 rounded-3xl p-8 max-w-md text-center space-y-6 backdrop-blur-xl">
        <p className="text-sm uppercase tracking-[0.3em] text-white/60">Theme toggler</p>
        <h2 className="text-3xl font-semibold">Current theme: {theme}</h2>
        <p className="text-white/70">
          Toggling the theme updates the `document.documentElement` classes and persists your preference in
          local storage. Components reading the theme hook respond instantly.
        </p>
        <ThemeToggle />
      </div>
    </div>
  );
}

export const Default: Story = {
  render: () => <ThemeTogglePreview />,
};
