export const WORKSPACE_ICON_NAMES = {
  mail: "mail",
  package: "package_2",
  headphones: "headphones",
  shield: "shield",
} as const

export type WorkspaceIconKey = keyof typeof WORKSPACE_ICON_NAMES

export function resolveWorkspaceIconName(iconKey?: string): string {
  const key = (iconKey ?? "package").toLowerCase() as WorkspaceIconKey
  return WORKSPACE_ICON_NAMES[key] ?? WORKSPACE_ICON_NAMES.package
}

export const WORKSPACE_ICON_OPTIONS = [
  { key: "mail", label: "Email", iconName: WORKSPACE_ICON_NAMES.mail },
  { key: "package", label: "Product", iconName: WORKSPACE_ICON_NAMES.package },
  { key: "headphones", label: "Support", iconName: WORKSPACE_ICON_NAMES.headphones },
  { key: "shield", label: "Security", iconName: WORKSPACE_ICON_NAMES.shield },
] as const
