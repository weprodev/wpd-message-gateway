import { render, screen } from "@testing-library/react"
import { describe, expect, it } from "vitest"

import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "./table"

describe("Table Component Suite", () => {
  it("renders table structures correctly", () => {
    render(
      <Table>
        <TableHeader>
          <TableRow>
            <TableHead>Col 1</TableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          <TableRow>
            <TableCell>Val 1</TableCell>
          </TableRow>
        </TableBody>
      </Table>
    )

    expect(screen.getByText("Col 1")).toBeInTheDocument()
    expect(screen.getByText("Val 1")).toBeInTheDocument()
  })
})
