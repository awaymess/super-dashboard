'use client';

import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
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
  const router = useRouter();

  useKeyboardShortcuts([
    { key: 'k', ctrl: true, action: () => setCommandPaletteOpen(true) },
    { key: 'd', action: () => router.push('/dashboard') },
    { key: 'b', action: () => router.push('/betting') },
    { key: 's', action: () => router.push('/stocks') },
    { key: 'p', action: () => router.push('/paper-trading') },
    { key: 'a', action: () => router.push('/analytics') },
    { key: '?', action: () => setShortcutsHelpOpen(true) },
  ]);

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
