import type { Meta, StoryObj } from "@storybook/react-vite"

import { Button } from "../button"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "./dialog"

const meta = {
  title: "UI/Dialog",
  component: Dialog,
  tags: ["autodocs"],
  parameters: { layout: "centered" },
} satisfies Meta<typeof Dialog>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Open Dialog</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Are you absolutely sure?</DialogTitle>
          <DialogDescription>
            This action cannot be undone. This will permanently delete your account
            and remove your data from our servers.
          </DialogDescription>
        </DialogHeader>
      </DialogContent>
    </Dialog>
  ),
}

export const WithFooter: Story = {
  render: () => (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Dialog with Actions</Button>
      </DialogTrigger>
      <DialogContent>
        <DialogHeader>
          <DialogTitle>Edit Profile</DialogTitle>
          <DialogDescription>
            Make changes to your profile here. Click save when you're done.
          </DialogDescription>
        </DialogHeader>
        <div className="grid gap-4 py-4">
          <div className="grid grid-cols-4 items-center gap-4">
            <label htmlFor="name" className="text-right text-sm">
              Name
            </label>
            <input
              id="name"
              defaultValue="Pedro Duarte"
              className="col-span-3 border p-2 rounded text-sm bg-background text-foreground"
            />
          </div>
          <div className="grid grid-cols-4 items-center gap-4">
            <label htmlFor="username" className="text-right text-sm">
              Username
            </label>
            <input
              id="username"
              defaultValue="@peduarte"
              className="col-span-3 border p-2 rounded text-sm bg-background text-foreground"
            />
          </div>
        </div>
        <DialogFooter>
          <Button type="submit">Save changes</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  ),
}

export const Scrollable: Story = {
  render: () => (
    <Dialog>
      <DialogTrigger asChild>
        <Button variant="outline">Scrollable Dialog</Button>
      </DialogTrigger>
      <DialogContent className="max-h-[80vh] overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Terms of Service</DialogTitle>
          <DialogDescription>
            Please read our terms of service carefully.
          </DialogDescription>
        </DialogHeader>
        <div className="space-y-4 py-4 text-sm text-muted-foreground">
          <p>
            Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aliquam id
            purus eget purus sollicitudin ultrices. Mauris sollicitudin purus
            id magna gravida, a finibus velit ultrices.
          </p>
          <p>
            Integer convallis leo vel finibus egestas. Pellentesque sodales
            fermentum urna vel pellentesque. In sodales arcu sed lacus dictum
            luctus.
          </p>
          <p>
            Curabitur volutpat ligula quis lorem molestie, a efficitur ex
            efficitur. Nam sollicitudin lacinia lacus, eget varius leo varius
            non.
          </p>
          <p>
            Proin bibendum erat a ex bibendum, sed pretium nibh volutpat.
            Donec sed neque sit amet dui interdum luctus.
          </p>
          <p>
            Aenean condimentum erat non sem finibus elementum. Proin ac dolor
            in mi efficitur efficitur.
          </p>
          <p>
            Duis viverra arcu vitae erat tristique pretium. Sed porta risus at
            arcu sodales lacinia.
          </p>
        </div>
        <DialogFooter>
          <Button>Close</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  ),
}
