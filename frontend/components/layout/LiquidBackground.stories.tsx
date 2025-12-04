import type { Meta, StoryObj } from '@storybook/react';
import { LiquidBackground } from './LiquidBackground';

const meta = {
  title: 'Layout/LiquidBackground',
  component: LiquidBackground,
  parameters: {
    layout: 'fullscreen',
  },
} satisfies Meta<typeof LiquidBackground>;

export default meta;
type Story = StoryObj<typeof meta>;

function LiquidBackgroundPreview() {
  return (
    <div className="relative min-h-[60vh] overflow-hidden bg-gradient-to-br from-[#04060d] via-[#0a1428] to-[#04060d] text-white flex items-center justify-center">
      <LiquidBackground />
      <div className="relative z-10 text-center space-y-4 max-w-xl px-6">
        <p className="text-sm uppercase tracking-[0.3em] text-white/60">Ambient motion layer</p>
        <h2 className="text-4xl font-semibold">Liquid glass backdrop</h2>
        <p className="text-white/70">
          This animated canvas adds subtle depth to the dashboard using soft, shifting blobs with blur and
          gradients. Overlay your application surface on top to create a layered glassmorphism feel.
        </p>
      </div>
    </div>
  );
}

export const Default: Story = {
  render: () => <LiquidBackgroundPreview />,
};
