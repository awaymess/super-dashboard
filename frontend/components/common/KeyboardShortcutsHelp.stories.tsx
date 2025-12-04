import type { Meta, StoryObj } from '@storybook/react';
import { useState, useEffect } from 'react';
import { action } from '@storybook/addon-actions';
import { KeyboardShortcutsHelp } from './KeyboardShortcutsHelp';

const meta = {
  title: 'Common/KeyboardShortcutsHelp',
  component: KeyboardShortcutsHelp,
  parameters: {
    layout: 'fullscreen',
  },
  args: {
    isOpen: true,
    onClose: action('close-keyboard-shortcuts'),
  },
} satisfies Meta<typeof KeyboardShortcutsHelp>;

export default meta;
type Story = StoryObj<typeof meta>;

function KeyboardShortcutsPreview(args: React.ComponentProps<typeof KeyboardShortcutsHelp>) {
  const [isOpen, setIsOpen] = useState(args.isOpen);

  useEffect(() => {
    setIsOpen(args.isOpen);
  }, [args.isOpen]);

  return (
    <div className="min-h-screen bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] text-white flex items-center justify-center">
      <KeyboardShortcutsHelp
        {...args}
        isOpen={isOpen}
        onClose={() => {
          setIsOpen(false);
          args.onClose?.();
        }}
      />
      {!isOpen && (
        <button
          type="button"
          onClick={() => setIsOpen(true)}
          className="absolute bottom-10 px-4 py-2 rounded-lg bg-white/10 hover:bg-white/20 transition"
        >
          Show Shortcuts
        </button>
      )}
    </div>
  );
}

export const Default: Story = {
  render: args => <KeyboardShortcutsPreview {...args} />,
};
