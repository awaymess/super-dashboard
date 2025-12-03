'use client';

import { GlassCard, Modal } from '@/components/ui';
import { Keyboard } from 'lucide-react';

interface KeyboardShortcutsHelpProps {
  isOpen: boolean;
  onClose: () => void;
}

const shortcuts = [
  { keys: ['Ctrl', 'K'], description: 'Open Command Palette' },
  { keys: ['D'], description: 'Go to Dashboard' },
  { keys: ['B'], description: 'Go to Betting' },
  { keys: ['S'], description: 'Go to Stocks' },
  { keys: ['P'], description: 'Go to Paper Trading' },
  { keys: ['A'], description: 'Go to Analytics' },
  { keys: ['Esc'], description: 'Close Modal / Cancel' },
  { keys: ['?'], description: 'Show Keyboard Shortcuts' },
];

export function KeyboardShortcutsHelp({ isOpen, onClose }: KeyboardShortcutsHelpProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose} title="Keyboard Shortcuts">
      <div className="space-y-3">
        {shortcuts.map((shortcut, i) => (
          <div key={i} className="flex items-center justify-between py-2 border-b border-white/10">
            <span className="text-gray-300">{shortcut.description}</span>
            <div className="flex items-center gap-1">
              {shortcut.keys.map((key, j) => (
                <span key={j}>
                  <kbd className="px-2 py-1 bg-white/10 rounded text-sm text-gray-300">
                    {key}
                  </kbd>
                  {j < shortcut.keys.length - 1 && <span className="mx-1 text-gray-500">+</span>}
                </span>
              ))}
            </div>
          </div>
        ))}
      </div>
    </Modal>
  );
}
