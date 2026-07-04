import type {
  MessageDispatchConfig,
  MessageDispatchMode,
  StoreMessageContentSetting,
} from "./settings.types"

export const DISPATCH_MODES = ["memory", "provider"] as const
export const STORE_CONTENT_MODES = ["true", "false"] as const


export function normalizeDispatchMode(raw?: string): MessageDispatchMode {
  switch (raw?.trim().toLowerCase()) {
    case "provider":
      return "provider"
    case "memory":
    default:
      return "memory"
  }
}

export function normalizeStoreMessageContent(raw?: string): boolean {
  switch (raw?.trim().toLowerCase()) {
    case "true":
      return true
    case "false":
    default:
      return false
  }
}

export function parseMessageDispatchConfig(
  modeRaw?: string,
  storeRaw?: string,
): MessageDispatchConfig {
  return {
    mode: normalizeDispatchMode(modeRaw),
    storeMessageContent: normalizeStoreMessageContent(storeRaw),
  }
}

export function toStoreMessageContentSetting(storeMessageContent: boolean): StoreMessageContentSetting {
  return storeMessageContent ? "true" : "false"
}

export function dispatchConfigsEqual(
  left: MessageDispatchConfig,
  right: MessageDispatchConfig,
): boolean {
  return left.mode === right.mode && left.storeMessageContent === right.storeMessageContent
}
