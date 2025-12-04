import type { Meta, StoryObj } from '@storybook/react';
import { GlassCard } from './GlassCard';

const meta: Meta<typeof GlassCard> = {
  title: 'UI/GlassCard',
  component: GlassCard,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'elevated', 'outlined'],
    },
    padding: {
      control: 'select',
      options: ['none', 'sm', 'md', 'lg'],
    },
    hover: {
      control: 'boolean',
    },
    glow: {
      control: 'boolean',
    },
    glowColor: {
      control: 'select',
      options: ['primary', 'success', 'danger', 'warning'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof GlassCard>;

export const Default: Story = {
  args: {
    children: (
      <div>
        <h3 className="text-white font-semibold mb-2">Default Card</h3>
        <p className="text-white/60">This is a default glass card with a subtle glassmorphism effect.</p>
      </div>
    ),
  },
};

export const Elevated: Story = {
  args: {
    variant: 'elevated',
    children: (
      <div>
        <h3 className="text-white font-semibold mb-2">Elevated Card</h3>
        <p className="text-white/60">This card has a stronger shadow for more depth.</p>
      </div>
    ),
  },
};

export const Outlined: Story = {
  args: {
    variant: 'outlined',
    children: (
      <div>
        <h3 className="text-white font-semibold mb-2">Outlined Card</h3>
        <p className="text-white/60">This card has a more visible border.</p>
      </div>
    ),
  },
};

export const WithHover: Story = {
  args: {
    hover: true,
    children: (
      <div>
        <h3 className="text-white font-semibold mb-2">Hoverable Card</h3>
        <p className="text-white/60">Hover over this card to see the effect.</p>
      </div>
    ),
  },
};

export const WithGlow: Story = {
  args: {
    glow: true,
    glowColor: 'primary',
    children: (
      <div>
        <h3 className="text-white font-semibold mb-2">Glowing Card</h3>
        <p className="text-white/60">This card has a subtle glow effect.</p>
      </div>
    ),
  },
};

export const SmallPadding: Story = {
  args: {
    padding: 'sm',
    children: (
      <div>
        <h3 className="text-white font-semibold mb-1">Small Padding</h3>
        <p className="text-white/60 text-sm">Compact card layout.</p>
      </div>
    ),
  },
};

export const LargePadding: Story = {
  args: {
    padding: 'lg',
    children: (
      <div>
        <h3 className="text-white font-semibold mb-2 text-lg">Large Padding</h3>
        <p className="text-white/60">This card has extra spacing inside.</p>
      </div>
    ),
  },
};
