import type { Preview } from '@storybook/react';
import '../app/globals.css';
import '../styles/animations.css';

const preview: Preview = {
  parameters: {
    controls: {
      matchers: {
        color: /(background|color)$/i,
        date: /Date$/i,
      },
    },
    backgrounds: {
      default: 'dark',
      values: [
        { name: 'dark', value: '#0a0a0f' },
        { name: 'light', value: '#ffffff' },
      ],
    },
  },
};

export default preview;

// Provide a minimal `process.env` shim for stories that may access it
// (e.g., Next.js components expecting `process` in browser context).
// Storybook runs in the browser where `process` is undefined.
// This avoids runtime errors like 'process is not defined'.
if (typeof globalThis !== 'undefined' && typeof (globalThis as any).process === 'undefined') {
  (globalThis as any).process = { env: {} } as any;
}
