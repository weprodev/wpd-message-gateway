import type { LogRow } from "./inbox.types"

export function hasInboxContent(log: LogRow) {
  return Boolean(log.inbox_message_id)
}

export function isEmailInboxLink(log: LogRow) {
  return log.channel_type.toLowerCase() === "email" && hasInboxContent(log)
}
