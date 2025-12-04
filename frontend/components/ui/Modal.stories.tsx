import type { Meta, StoryObj } from '@storybook/react';
import { useState } from 'react';
import { Modal } from './Modal';
import { GlassButton } from './GlassButton';
import { GlassInput } from './GlassInput';

const meta: Meta<typeof Modal> = {
  title: 'UI/Modal',
  component: Modal,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    size: {
      control: 'select',
      options: ['sm', 'md', 'lg', 'xl', 'full'],
    },
    showCloseButton: {
      control: 'boolean',
    },
  },
};

export default meta;
type Story = StoryObj<typeof Modal>;

const ModalTemplate = ({ ...args }) => {
  const [isOpen, setIsOpen] = useState(false);
  
  return (
    <>
      <GlassButton variant="primary" onClick={() => setIsOpen(true)}>
        Open Modal
      </GlassButton>
      <Modal {...args} isOpen={isOpen} onClose={() => setIsOpen(false)}>
        {args.children}
      </Modal>
    </>
  );
};

export const Default: Story = {
  render: () => (
    <ModalTemplate
      title="Default Modal"
      description="This is a default modal with a title and description."
    >
      <p className="text-white/70">
        Modal content goes here. You can add any content including forms, lists, or other components.
      </p>
    </ModalTemplate>
  ),
};

export const Small: Story = {
  render: () => (
    <ModalTemplate title="Small Modal" size="sm">
      <p className="text-white/70">
        This is a small modal for quick confirmations or alerts.
      </p>
      <div className="flex gap-3 mt-6">
        <GlassButton variant="ghost">Cancel</GlassButton>
        <GlassButton variant="primary">Confirm</GlassButton>
      </div>
    </ModalTemplate>
  ),
};

export const Large: Story = {
  render: () => (
    <ModalTemplate
      title="Large Modal"
      description="This modal has more space for complex content."
      size="lg"
    >
      <div className="space-y-4">
        <p className="text-white/70">
          Large modals are great for forms, detailed information, or complex interactions
          that require more space.
        </p>
        <div className="grid grid-cols-2 gap-4">
          <GlassInput label="First Name" placeholder="John" />
          <GlassInput label="Last Name" placeholder="Doe" />
        </div>
        <GlassInput label="Email" placeholder="john@example.com" type="email" />
        <div className="flex justify-end gap-3 mt-6">
          <GlassButton variant="ghost">Cancel</GlassButton>
          <GlassButton variant="primary">Save Changes</GlassButton>
        </div>
      </div>
    </ModalTemplate>
  ),
};

export const ExtraLarge: Story = {
  render: () => (
    <ModalTemplate
      title="Extra Large Modal"
      description="For complex content that needs lots of space."
      size="xl"
    >
      <div className="grid grid-cols-3 gap-6">
        {[1, 2, 3, 4, 5, 6].map((i) => (
          <div key={i} className="p-4 bg-white/5 rounded-lg">
            <h4 className="font-semibold text-white mb-2">Item {i}</h4>
            <p className="text-sm text-white/60">
              Description for item {i}. This is some placeholder content.
            </p>
          </div>
        ))}
      </div>
    </ModalTemplate>
  ),
};

export const WithoutCloseButton: Story = {
  render: () => (
    <ModalTemplate
      title="Modal Without Close Button"
      showCloseButton={false}
    >
      <p className="text-white/70">
        This modal has no close button in the header. Users must use the action buttons below.
      </p>
      <div className="flex gap-3 mt-6">
        <GlassButton variant="ghost">Cancel</GlassButton>
        <GlassButton variant="primary">Confirm</GlassButton>
      </div>
    </ModalTemplate>
  ),
};

export const Confirmation: Story = {
  render: () => (
    <ModalTemplate title="Delete Item?" size="sm">
      <p className="text-white/70">
        Are you sure you want to delete this item? This action cannot be undone.
      </p>
      <div className="flex gap-3 mt-6 justify-end">
        <GlassButton variant="ghost">Cancel</GlassButton>
        <GlassButton variant="danger">Delete</GlassButton>
      </div>
    </ModalTemplate>
  ),
};

export const FormModal: Story = {
  render: () => (
    <ModalTemplate
      title="Create New Item"
      description="Fill out the form below to create a new item."
      size="md"
    >
      <div className="space-y-4">
        <GlassInput label="Title" placeholder="Enter title" />
        <GlassInput label="Description" placeholder="Enter description" />
        <GlassInput label="Category" placeholder="Select category" />
        <div className="flex justify-end gap-3 pt-4 border-t border-white/10">
          <GlassButton variant="ghost">Cancel</GlassButton>
          <GlassButton variant="success">Create Item</GlassButton>
        </div>
      </div>
    </ModalTemplate>
  ),
};
