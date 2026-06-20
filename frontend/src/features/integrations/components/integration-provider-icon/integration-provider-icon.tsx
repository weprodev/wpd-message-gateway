import { cn } from "@/lib/utils"

interface IntegrationProviderIconProps {
  icon: string
  name: string
  className?: string
}

export function IntegrationProviderIcon({ icon, name, className }: IntegrationProviderIconProps) {
  const isImage = icon.startsWith("/") || icon.startsWith("http")

  return (
    <div
      className={cn(
        "flex shrink-0 items-center justify-center rounded-lg bg-input overflow-hidden",
        className,
      )}
    >
      {isImage ? (
        <img src={icon} alt={`${name} logo`} className="size-full object-contain" />
      ) : (
        <span>{icon}</span>
      )}
    </div>
  )
}
