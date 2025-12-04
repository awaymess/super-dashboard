import type { Meta, StoryObj } from '@storybook/react';
import { Heatmap } from './Heatmap';

const xLabels = ['Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat', 'Sun'];
const yLabels = ['Value', 'Stocks', 'Paper', 'Analytics', 'News'];

const generateData = () =>
  yLabels.flatMap((_, y) =>
    xLabels.map((_, x) => ({
      x,
      y,
      value: Math.round(Math.random() * 100),
    }))
  );

const meta = {
  title: 'Charts/Heatmap',
  component: Heatmap,
  parameters: {
    layout: 'centered',
    backgrounds: { default: 'dark' },
  },
  args: {
    xLabels,
    yLabels,
    data: generateData(),
  },
} satisfies Meta<typeof Heatmap>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const CoolPalette: Story = {
  args: {
    data: generateData(),
    colorScale: { min: '#38bdf8', mid: '#818cf8', max: '#f472b6' },
  },
};
