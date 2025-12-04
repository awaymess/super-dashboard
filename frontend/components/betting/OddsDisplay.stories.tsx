import type { Meta, StoryObj } from '@storybook/react';
import { OddsDisplay } from './OddsDisplay';

const meta = {
  title: 'Betting/OddsDisplay',
  component: OddsDisplay,
  parameters: {
    layout: 'centered',
  },
  args: {
    home: 1.95,
    draw: 3.6,
    away: 4.1,
    previousHome: 2.02,
    previousDraw: 3.45,
    previousAway: 3.9,
    size: 'md',
  },
} satisfies Meta<typeof OddsDisplay>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Compact: Story = {
  args: {
    size: 'sm',
  },
};

export const Large: Story = {
  args: {
    size: 'lg',
    home: 2.4,
    draw: 3.1,
    away: 2.9,
  },
};
