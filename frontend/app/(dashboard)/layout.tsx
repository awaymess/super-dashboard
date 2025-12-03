'use client';

import { useState } from 'react';
import { Sidebar, Header, CommandPalette, LiquidBackground } from '@/components/layout';
import { KeyboardShortcutsHelp } from '@/components/common';
import { useKeyboardShortcuts } from '@/hooks/useKeyboardShortcuts';

export default function DashboardLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(false);
  const [commandPaletteOpen, setCommandPaletteOpen] = useState(false);
  const [shortcutsHelpOpen, setShortcutsHelpOpen] = useState(false);

  useKeyboardShortcuts();

  return (
    <div className="min-h-screen bg-background">
      <LiquidBackground />
      
      <Sidebar collapsed={sidebarCollapsed} onToggle={() => setSidebarCollapsed(!sidebarCollapsed)} />
      
      <main className={`transition-all duration-300 ${sidebarCollapsed ? 'ml-20' : 'ml-[280px]'}`}>
        <Header onOpenCommandPalette={() => setCommandPaletteOpen(true)} />
        
        <div className="p-6 relative z-10">
          {children}
        </div>
      </main>

      <CommandPalette isOpen={commandPaletteOpen} onClose={() => setCommandPaletteOpen(false)} />
      <KeyboardShortcutsHelp isOpen={shortcutsHelpOpen} onClose={() => setShortcutsHelpOpen(false)} />
    </div>
  );
}
