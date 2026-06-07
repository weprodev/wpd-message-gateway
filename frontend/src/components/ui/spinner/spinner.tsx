import { cn } from "@/lib/utils"

export interface SpinnerProps extends React.HTMLAttributes<HTMLSpanElement> {
  size?: "sm" | "md" | "lg"
}

export function Spinner({ size = "md", className, ...props }: SpinnerProps) {
  const sizeClasses = {
    sm: "size-4 border-2",
    md: "size-6 border-2",
    lg: "size-8 border-[3px]",
  }

  return (
    <span
      className={cn(
        "animate-spin rounded-full border-primary-brand border-t-transparent inline-block",
        sizeClasses[size],
        className
      )}
      role="status"
      aria-label="loading"
      {...props}
    />
  )
}
