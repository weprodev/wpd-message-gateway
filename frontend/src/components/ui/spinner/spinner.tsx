import type { HTMLAttributes } from "react"
import { cva, type VariantProps } from "class-variance-authority"

import { cn } from "@/lib/utils"

const spinnerVariants = cva(
  "animate-spin rounded-full inline-block border-t-transparent",
  {
    variants: {
      size: {
        sm: "size-4 border-2",
        md: "size-6 border-2",
        lg: "size-8 border-[3px]",
      },
      variant: {
        default: "border-primary-brand",
        onSolid: "border-current",
      },
    },
    defaultVariants: {
      size: "md",
      variant: "default",
    },
  },
)

export interface SpinnerProps
  extends HTMLAttributes<HTMLSpanElement>,
    VariantProps<typeof spinnerVariants> {}

export function Spinner({ size, variant, className, ...props }: SpinnerProps) {
  return (
    <span
      className={cn(spinnerVariants({ size, variant }), className)}
      role="status"
      aria-label="loading"
      {...props}
    />
  )
}
