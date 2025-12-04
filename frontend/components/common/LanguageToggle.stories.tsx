import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { action } from '@storybook/addon-actions';
import { LanguageToggle } from './LanguageToggle';

const meta = {
  title: 'Common/LanguageToggle',
  component: LanguageToggle,
  parameters: {
    layout: 'centered',
  },
  args: {
    currentLocale: 'en',
    onLocaleChange: action('language-change'),
  },
} satisfies Meta<typeof LanguageToggle>;

export default meta;
type Story = StoryObj<typeof meta>;

function LanguageTogglePreview(args: React.ComponentProps<typeof LanguageToggle>) {
  const [locale, setLocale] = useState(args.currentLocale);

  return (
    <div className="p-10 bg-gradient-to-br from-[#05060b] via-[#0b1120] to-[#05060b] rounded-3xl">
      <LanguageToggle
        {...args}
        currentLocale={locale}
        onLocaleChange={nextLocale => {
          setLocale(nextLocale);
          args.onLocaleChange?.(nextLocale);
        }}
      />
      <p className="text-white/70 text-sm mt-4">Current locale: {locale}</p>
    </div>
  );
}

export const Default: Story = {
  render: args => <LanguageTogglePreview {...args} />,
};
