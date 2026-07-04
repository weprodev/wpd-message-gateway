interface AppFooterProps {
  variant?: "default" | "dashboard"
}

export function AppFooter({ variant = "default" }: AppFooterProps) {
  const linkClass =
    variant === "dashboard"
      ? "font-semibold text-indigo-400 transition-colors hover:text-indigo-300"
      : "font-semibold text-primary-brand underline transition-colors hover:text-primary-brand-hover"

  return (
    <footer className="w-full shrink-0 border-t border-divider bg-card">
      <div className="flex h-16 items-center justify-between px-12 py-6">
        <p className="text-sm text-text-secondary">
          Developed by{" "}
          <a
            href="https://weprodev.com"
            target="_blank"
            rel="noopener noreferrer"
            className={linkClass}
          >
            WeProDev
          </a>
        </p>
        <p className="text-sm text-text-tertiary">We develop for growth, with growth in mind!</p>
      </div>
    </footer>
  )
}
