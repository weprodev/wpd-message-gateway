interface PageHeaderProps {
  title: string
  description?: string
}

export function PageHeader({ title, description }: PageHeaderProps) {
  return (
    <div>
      <h1 className="text-xl font-semibold text-text-primary">{title}</h1>
      {description ? (
        <p className="mt-1 text-sm text-text-secondary">{description}</p>
      ) : null}
    </div>
  )
}
