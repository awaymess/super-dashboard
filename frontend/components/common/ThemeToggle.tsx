'use client';

import { Sun, Moon } from 'lucide-react';
import { useTheme } from '@/hooks/useTheme';
import { GlassButton } from '@/components/ui';

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme();

  return (
    <GlassButton variant="ghost" size="sm" onClick={toggleTheme}>
      {theme === 'dark' ? (
        <Sun className="w-5 h-5" />
      ) : (
        <Moon className="w-5 h-5" />
      )}
    </GlassButton>
  );
}
