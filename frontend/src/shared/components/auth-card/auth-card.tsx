import type { ReactNode } from "react"

interface AuthCardProps {
  children: ReactNode
}

export function AuthCard({ children }: AuthCardProps) {
  return (
    <div
      className="w-full max-w-[560px] rounded-2xl border border-border bg-card p-10"
      style={{ boxShadow: "var(--shadow-card)" }}
    >
      {children}
    </div>
  )
}
