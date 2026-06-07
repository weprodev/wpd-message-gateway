import { forwardRef } from "react"
import type { InputHTMLAttributes } from "react"

import { Icon } from "@/components/ui/icon"

interface SearchInputProps extends InputHTMLAttributes<HTMLInputElement> {
  placeholder?: string
}

export const SearchInput = forwardRef<HTMLInputElement, SearchInputProps>(
  ({ placeholder = "Search...", className = "", ...props }, ref) => {
    return (
      <div className={`bg-card border border-border rounded-lg h-10 w-full ${className}`}>
        <div className="flex items-center gap-2.5 px-3.5 size-full">
          <Icon name="search" size="sm" className="text-muted-foreground shrink-0" />
          <input
            ref={ref}
            type="text"
            className="flex-1 bg-transparent text-sm text-foreground outline-none placeholder:text-muted-foreground/60"
            placeholder={placeholder}
            {...props}
          />
        </div>
      </div>
    )
  }
)

SearchInput.displayName = "SearchInput"
