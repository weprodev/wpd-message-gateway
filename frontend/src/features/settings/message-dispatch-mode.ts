import type {
  MessageDispatchApiValue,
  MessageDispatchConfig,
  MessageDispatchMode,
} from "./settings.types"

export const DISPATCH_MODES = ["memory", "provider"] as const

const API_VALUE_TO_CONFIG: Record<MessageDispatchApiValue, MessageDispatchConfig> = {
  memory_only: { mode: "memory", storeInDb: false },
  memory_and_provider: { mode: "memory", storeInDb: true },
  provider_only: { mode: "provider", storeInDb: false },
  provider_and_database: { mode: "provider", storeInDb: true },
}

export function normalizeDispatchMode(raw?: string): MessageDispatchMode {
  switch (raw?.toLowerCase()) {
    case "provider":
    case "providers":
    case "provider_only":
    case "provider_and_database":
    case "provider_database":
      return "provider"
    case "memory":
    case "memory_only":
    case "memory_and_provider":
    case "both":
    case "memory_database":
    default:
      return "memory"
  }
}

export function normalizeStoreInDb(raw?: boolean | string): boolean {
  if (typeof raw === "boolean") {
    return raw
  }

  switch (raw?.toLowerCase()) {
    case "true":
    case "memory_and_provider":
    case "provider_and_database":
    case "both":
    case "memory_database":
    case "provider_database":
      return true
    default:
      return false
  }
}

export function parseMessageDispatchConfig(raw?: string): MessageDispatchConfig {
  const value = raw?.toLowerCase()

  if (value && value in API_VALUE_TO_CONFIG) {
    return API_VALUE_TO_CONFIG[value as MessageDispatchApiValue]
  }

  return {
    mode: normalizeDispatchMode(value),
    storeInDb: normalizeStoreInDb(value),
  }
}

export function toMessageDispatchApiValue(config: MessageDispatchConfig): MessageDispatchApiValue {
  if (config.mode === "memory") {
    return config.storeInDb ? "memory_and_provider" : "memory_only"
  }

  return config.storeInDb ? "provider_and_database" : "provider_only"
}

export function dispatchConfigsEqual(
  left: MessageDispatchConfig,
  right: MessageDispatchConfig,
): boolean {
  return left.mode === right.mode && left.storeInDb === right.storeInDb
}
