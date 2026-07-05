import { cn } from "@/lib/utils"

interface RadioOptionProps {
  id: string
  name: string
  label: string
  description: string
  checked: boolean
  onChange: () => void
  disabled?: boolean
}

export function RadioOption({
  id,
  name,
  label,
  description,
  checked,
  onChange,
  disabled = false,
}: RadioOptionProps) {
  return (
    <label
      htmlFor={id}
      className={cn(
        "flex gap-3 rounded-lg border p-4 transition-colors",
        disabled ? "cursor-not-allowed opacity-60" : "cursor-pointer",
        checked ? "border-primary-brand bg-input" : "border-border",
        !disabled && !checked && "hover:bg-muted/40",
      )}
    >
      <input
        id={id}
        type="radio"
        name={name}
        checked={checked}
        onChange={onChange}
        disabled={disabled}
        className="mt-1 accent-primary-brand"
      />
      <div>
        <p className="text-sm font-medium text-foreground">{label}</p>
        <p className="mt-1 text-sm text-text-secondary">{description}</p>
      </div>
    </label>
  )
}
