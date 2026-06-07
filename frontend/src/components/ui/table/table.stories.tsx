import type { Meta, StoryObj } from "@storybook/react-vite"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./table"

const meta = {
  title: "Components/Table",
  component: Table,
  tags: ["autodocs"],
  parameters: {
    layout: "padded",
    docs: {
      description: {
        component: "Accessible and clean HTML-based table elements.",
      },
    },
  },
} satisfies Meta<typeof Table>

export default meta

type Story = StoryObj<typeof meta>

export const Default: Story = {
  render: () => (
    <div className="rounded-xl border border-border bg-card shadow-sm overflow-hidden">
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Header 1</TableHead>
            <TableHead>Header 2</TableHead>
            <TableHead className="text-right">Header 3</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableCell className="font-medium">Data 1A</TableCell>
            <TableCell>Data 1B</TableCell>
            <TableCell className="text-right">Data 1C</TableCell>
          </TableRow>
          <TableRow>
            <TableCell className="font-medium">Data 2A</TableCell>
            <TableCell>Data 2B</TableCell>
            <TableCell className="text-right">Data 2C</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    </div>
  ),
}
