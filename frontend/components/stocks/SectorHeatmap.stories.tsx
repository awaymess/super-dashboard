import type { Meta, StoryObj } from "@storybook/react";
import React from "react";
import { SectorHeatmap } from "./SectorHeatmap";

const meta: Meta<typeof SectorHeatmap> = {
  title: "Stocks/SectorHeatmap",
  component: SectorHeatmap,
  parameters: {
    backgrounds: {
      default: "dark",
      values: [
        { name: "dark", value: "#0b0f19" },
        { name: "light", value: "#f8fafc" },
      ],
    },
  },
};

export default meta;
type Story = StoryObj<typeof SectorHeatmap>;

const sampleSectors = [
  { name: "Technology", change: 1.8 },
  { name: "Healthcare", change: -0.6 },
  { name: "Financials", change: 0.4 },
  { name: "Consumer Discretionary", change: 2.1 },
  { name: "Consumer Staples", change: -0.3 },
  { name: "Energy", change: -2.4 },
  { name: "Industrials", change: 0.9 },
  { name: "Materials", change: -1.1 },
  { name: "Utilities", change: 0.2 },
  { name: "Real Estate", change: 0.7 },
  { name: "Communication Services", change: 1.2 },
];

export const Default: Story = {
  args: {
    sectors: sampleSectors,
  },
};

export const BearishDay: Story = {
  args: {
    sectors: sampleSectors.map((s) => ({ ...s, change: s.change - 2 })),
  },
};

export const LightMode: Story = {
  parameters: { backgrounds: { default: "light" } },
  args: {
    sectors: sampleSectors,
  },
};
