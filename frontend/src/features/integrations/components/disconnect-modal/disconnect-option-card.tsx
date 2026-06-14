import { Icon } from "@/components/ui/icon"
import { cn } from "@/lib/utils"

type DisconnectOptionCardProps = {
  icon: string
  title: string
  description: string
} & (
  | { variant: "neutral" }
  | { variant: "danger" }
)

export function DisconnectOptionCard({
  variant,
  icon,
  title,
  description,
}: DisconnectOptionCardProps) {
  const isDanger = variant === "danger"

  return (
    <div
      className={cn(
        "flex items-start gap-3 rounded-lg border p-4",
        isDanger
          ? "border-destructive/20 bg-destructive/5"
          : "border-border bg-muted/40",
      )}
    >
      <div
        className={cn(
          "flex size-10 shrink-0 items-center justify-center rounded-full border bg-card",
          isDanger ? "border-destructive/20" : "border-border",
        )}
      >
        <Icon
          name={icon}
          size="sm"
          className={isDanger ? "text-destructive" : "text-text-secondary"}
        />
      </div>
      <div className="min-w-0 flex-1">
        <h3 className="text-sm font-semibold text-foreground">{title}</h3>
        <p className="mt-1 text-[13px] leading-normal text-text-secondary">{description}</p>
      </div>
    </div>
  )
}
