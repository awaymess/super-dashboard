import type { Meta, StoryObj } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { useEffect, useState } from 'react';
import type { AppRouterInstance } from 'next/dist/shared/lib/app-router-context.shared-runtime';
import { AppRouterContext } from 'next/dist/shared/lib/app-router-context.shared-runtime';
import { CommandPalette } from './CommandPalette';

const mockRouter: AppRouterInstance = {
  back: () => action('router.back')(),
  forward: () => action('router.forward')(),
  refresh: () => action('router.refresh')(),
  push: (href, options) => action('router.push')({ href, options }),
  replace: (href, options) => action('router.replace')({ href, options }),
  prefetch: (href, options) => action('router.prefetch')({ href, options }),
};

const meta = {
  title: 'Layout/CommandPalette',
  component: CommandPalette,
  parameters: {
    layout: 'fullscreen',
  },
} satisfies Meta<typeof CommandPalette>;

export default meta;
type Story = StoryObj<typeof meta>;

type CommandPaletteStoryProps = React.ComponentProps<typeof CommandPalette>;

function CommandPalettePreview(args: CommandPaletteStoryProps) {
  const [isOpen, setIsOpen] = useState(args.isOpen ?? true);

  useEffect(() => {
    setIsOpen(args.isOpen ?? false);
  }, [args.isOpen]);

  return (
    <AppRouterContext.Provider value={mockRouter}>
      <div className="min-h-screen bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] text-white flex items-center justify-center relative p-6">
        <CommandPalette {...args} isOpen={isOpen} onClose={() => setIsOpen(false)} />
        <button
          type="button"
          onClick={() => setIsOpen(true)}
          className="absolute bottom-10 right-10 px-4 py-2 rounded-lg bg-white/10 hover:bg-white/20 transition text-sm"
        >
          Open Command Palette
        </button>
      </div>
    </AppRouterContext.Provider>
  );
}

export const Default: Story = {
  args: {
    isOpen: true,
    onClose: () => {},
  },
  render: args => <CommandPalettePreview {...args} />,
};

export const InitiallyClosed: Story = {
  args: {
    isOpen: false,
    onClose: () => {},
  },
  render: args => <CommandPalettePreview {...args} />,
};
