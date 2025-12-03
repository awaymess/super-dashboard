'use client';

import { Globe } from 'lucide-react';
import { GlassButton, Dropdown } from '@/components/ui';

const languages = [
  { code: 'en', label: 'English', flag: '🇺🇸' },
  { code: 'th', label: 'ไทย', flag: '🇹🇭' },
];

interface LanguageToggleProps {
  currentLocale: string;
  onLocaleChange: (locale: string) => void;
}

export function LanguageToggle({ currentLocale, onLocaleChange }: LanguageToggleProps) {
  const currentLanguage = languages.find(l => l.code === currentLocale) || languages[0];

  return (
    <Dropdown
      trigger={
        <GlassButton variant="ghost" size="sm">
          <span className="mr-2">{currentLanguage.flag}</span>
          <span className="hidden sm:inline">{currentLanguage.label}</span>
        </GlassButton>
      }
      items={languages.map(lang => ({
        label: `${lang.flag} ${lang.label}`,
        onClick: () => onLocaleChange(lang.code),
      }))}
    />
  );
}
