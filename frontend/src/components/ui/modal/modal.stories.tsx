import type { Meta, StoryObj } from "@storybook/react-vite"
import { useState } from "react"
import { Button } from "../button"
import { Modal } from "./modal"

const meta = {
  title: "Components/Modal",
  component: Modal,
  tags: ["autodocs"],
  parameters: {
    layout: "centered",
  },
  argTypes: {
    isOpen: { control: "boolean" },
    title: { control: "text" },
  },
} satisfies Meta<typeof Modal>

export default meta
type Story = StoryObj<typeof meta>

export const Interactive: Story = {
  render: (args) => {
    const [isOpen, setIsOpen] = useState(false)
    return (
      <>
        <Button onClick={() => setIsOpen(true)}>Open Modal</Button>
        <Modal {...args} isOpen={isOpen} onClose={() => setIsOpen(false)}>
          <p className="text-sm text-text-secondary">This is a custom modal content block.</p>
        </Modal>
      </>
    )
  },
  args: {
    isOpen: true,
    onClose: () => {},
    children: <p className="text-sm text-text-secondary">This is a custom modal content block.</p>,
    title: "Modal Title",
  },
}

export const WithDescription: Story = {
  render: (args) => {
    const [isOpen, setIsOpen] = useState(false)
    return (
      <>
        <Button onClick={() => setIsOpen(true)}>Open Modal</Button>
        <Modal {...args} isOpen={isOpen} onClose={() => setIsOpen(false)} />
      </>
    )
  },
  args: {
    isOpen: true,
    onClose: () => {},
    title: "Delete item",
    description: "Are you sure you want to continue? This action cannot be undone.",
    children: null,
  },
}
