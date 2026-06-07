import { Icon } from "@/components/ui/icon"
import { useTheme } from "@/shared/context/theme-context"

export function ThemeToggle() {
  const { theme, toggleTheme } = useTheme()

  return (
    <button
      type="button"
      onClick={toggleTheme}
      className="rounded-lg p-2 transition-colors hover:bg-surface-hover"
      aria-label="Toggle theme"
    >
      <Icon
        name={theme === "light" ? "dark_mode" : "light_mode"}
        size="sm"
        className="text-text-secondary"
      />
    </button>
  )
}
