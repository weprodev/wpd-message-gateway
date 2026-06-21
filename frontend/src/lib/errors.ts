export function toUserMessage(error: unknown, fallback: string): string {
  if (error instanceof Error && error.message.trim()) {
    return error.message
  }
  return fallback
}

function isClientError(status: number): boolean {
  return status >= 400 && status < 500
}

export async function httpError(response: Response, fallback: string): Promise<Error> {
  if (!isClientError(response.status)) {
    return new Error(fallback)
  }

  const body = (await response.json().catch(() => ({}))) as { message?: string }
  const message = body.message?.trim()
  return new Error(message || fallback)
}

export const MISSING_CLIENT_SECRET_MESSAGE =
  "The server did not return a client secret. Please try again."

export function requireClientSecret(payload: { client_secret?: string }): string {
  if (!payload.client_secret?.trim()) {
    throw new Error(MISSING_CLIENT_SECRET_MESSAGE)
  }
  return payload.client_secret
}
