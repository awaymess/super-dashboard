import type { Meta, StoryObj } from '@storybook/react';
import { Search, Plus, ArrowRight, Settings } from 'lucide-react';
import { GlassButton } from './GlassButton';

const meta: Meta<typeof GlassButton> = {
  title: 'UI/GlassButton',
  component: GlassButton,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'primary', 'secondary', 'success', 'danger', 'ghost', 'outline'],
    },
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg', 'icon'],
    },
    loading: {
      control: 'boolean',
    },
    disabled: {
      control: 'boolean',
    },
    iconPosition: {
      control: 'select',
      options: ['left', 'right'],
    },
  },
};

export default meta;
type Story = StoryObj<typeof GlassButton>;

export const Default: Story = {
  args: {
    children: 'Default Button',
    variant: 'default',
  },
};

export const Primary: Story = {
  args: {
    children: 'Primary Button',
    variant: 'primary',
  },
};

export const Secondary: Story = {
  args: {
    children: 'Secondary Button',
    variant: 'secondary',
  },
};

export const Success: Story = {
  args: {
    children: 'Success Button',
    variant: 'success',
  },
};

export const Danger: Story = {
  args: {
    children: 'Danger Button',
    variant: 'danger',
  },
};

export const Ghost: Story = {
  args: {
    children: 'Ghost Button',
    variant: 'ghost',
  },
};

export const Outline: Story = {
  args: {
    children: 'Outline Button',
    variant: 'outline',
  },
};

export const Small: Story = {
  args: {
    children: 'Small Button',
    size: 'sm',
    variant: 'primary',
  },
};

export const Large: Story = {
  args: {
    children: 'Large Button',
    size: 'lg',
    variant: 'primary',
  },
};

export const Loading: Story = {
  args: {
    children: 'Loading...',
    loading: true,
    variant: 'primary',
  },
};

export const Disabled: Story = {
  args: {
    children: 'Disabled Button',
    disabled: true,
    variant: 'primary',
  },
};

export const WithIconLeft: Story = {
  args: {
    children: 'Search',
    icon: <Search className="w-4 h-4" />,
    iconPosition: 'left',
    variant: 'primary',
  },
};

export const WithIconRight: Story = {
  args: {
    children: 'Continue',
    icon: <ArrowRight className="w-4 h-4" />,
    iconPosition: 'right',
    variant: 'primary',
  },
};

export const IconOnly: Story = {
  args: {
    icon: <Settings className="w-5 h-5" />,
    size: 'icon',
    variant: 'ghost',
    'aria-label': 'Settings',
  },
};

export const AllVariants: Story = {
  render: () => (
    <div className="flex flex-wrap gap-3">
      <GlassButton variant="default">Default</GlassButton>
      <GlassButton variant="primary">Primary</GlassButton>
      <GlassButton variant="secondary">Secondary</GlassButton>
      <GlassButton variant="success">Success</GlassButton>
      <GlassButton variant="danger">Danger</GlassButton>
      <GlassButton variant="ghost">Ghost</GlassButton>
      <GlassButton variant="outline">Outline</GlassButton>
    </div>
  ),
};

export const AllSizes: Story = {
  render: () => (
    <div className="flex items-center gap-3">
      <GlassButton variant="primary" size="sm">Small</GlassButton>
      <GlassButton variant="primary" size="md">Medium</GlassButton>
      <GlassButton variant="primary" size="lg">Large</GlassButton>
      <GlassButton variant="primary" size="icon"><Plus className="w-5 h-5" /></GlassButton>
    </div>
  ),
};
