import type { StorybookConfig } from '@storybook/nextjs';

const config: StorybookConfig = {
  stories: ['../components/**/*.stories.@(js|jsx|mjs|ts|tsx)'],
  addons: [
    '@storybook/addon-links',
    '@storybook/addon-essentials',
    '@storybook/addon-interactions',
  ],
  framework: {
    name: '@storybook/nextjs',
    options: {},
  },
  docs: {},
  core: {
    disableTelemetry: true,
  },
  // Note: Storybook 8 has known compatibility issues with Next.js 15
  // See: https://github.com/storybookjs/storybook/issues/26243
  // Storybook dev/build will work once updated to Storybook 8.1+ or using experimental-nextjs-vite
};

export default config;
