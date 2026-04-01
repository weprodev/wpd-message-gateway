import * as React from "react"

import { cn } from "@/lib/utils"

const sizeClass = {
  sm: "text-[18px] h-[18px] w-[18px]",
  md: "text-2xl h-6 w-6",
  lg: "text-[28px] h-7 w-7",
} as const

export type IconSize = keyof typeof sizeClass

export interface IconProps extends Omit<React.HTMLAttributes<HTMLSpanElement>, "children"> {
  name: string
  size?: IconSize
}

export function Icon({
  name,
  size = "md",
  className,
  "aria-label": ariaLabel,
  "aria-labelledby": ariaLabelledBy,
  ...props
}: IconProps) {
  const labelled = Boolean(ariaLabel || ariaLabelledBy)

  return (
    <span
      className={cn(
        "material-symbols-outlined inline-flex shrink-0 items-center justify-center align-middle",
        sizeClass[size],
        className
      )}
      aria-label={ariaLabel}
      aria-labelledby={ariaLabelledBy}
      aria-hidden={!labelled}
      {...props}
    >
      {name}
    </span>
  )
}
