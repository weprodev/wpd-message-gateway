import { z } from "zod"

import type { LogRow } from "./inbox.types"

const logRowSchema = z.object({
  id: z.string(),
  workspace_id: z.string(),
  api_key_id: z.string().optional(),
  channel_type: z.string(),
  http_method: z.string(),
  status_code: z.number(),
  endpoint: z.string(),
  provider_name: z.string().optional(),
  request_id: z.string().optional(),
  duration_ms: z.number().optional(),
  error_message: z.string().optional(),
  created_at: z.string(),
  source_name: z.string().optional(),
  client_id: z.string().optional(),
})

const logsResponseSchema = z.object({
  items: z.array(logRowSchema),
  total: z.number(),
})

export function parseLogsResponse(raw: unknown): { items: LogRow[]; total: number } {
  const parsed = logsResponseSchema.safeParse(raw)
  if (!parsed.success) {
    throw new Error("Invalid logs response from server")
  }
  return parsed.data
}
