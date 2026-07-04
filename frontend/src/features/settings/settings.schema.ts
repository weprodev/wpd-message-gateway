import { z } from "zod"

import { DISPATCH_MODES, STORE_CONTENT_MODES } from "./message-dispatch-mode"
import type { MessageDispatchMode, StoreMessageContentSetting, WorkspaceSettings } from "./settings.types"

const messageDispatchModeSchema = z.enum(DISPATCH_MODES)

const storeMessageContentSchema = z.enum(STORE_CONTENT_MODES)

export const workspaceSettingsSchema = z.object({
  owner_email: z.string().optional(),
  pin_code: z.string().optional(),
  message_dispatch_mode: messageDispatchModeSchema.optional(),
  store_message_content: storeMessageContentSchema.optional(),
})

export function parseWorkspaceSettings(raw: unknown): WorkspaceSettings {
  const parsed = workspaceSettingsSchema.safeParse(raw)
  if (!parsed.success) {
    throw new Error("Invalid settings response from server")
  }
  return parsed.data satisfies WorkspaceSettings
}

export function parseMessageDispatchMode(raw: unknown): MessageDispatchMode {
  const parsed = messageDispatchModeSchema.safeParse(raw)
  return parsed.success ? parsed.data : "memory"
}

export function parseStoreMessageContentSetting(raw: unknown): StoreMessageContentSetting {
  const parsed = storeMessageContentSchema.safeParse(raw)
  return parsed.success ? parsed.data : "false"
}
