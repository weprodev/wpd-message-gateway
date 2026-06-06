import { useState } from "react"
import { Icon } from "@/components/ui/icon"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog"
import { sendTestRequest } from "../inbox.api"
import type { MessageChannel } from "../inbox.types"

interface SendTestModalProps {
  workspaceId: string
  open: boolean
  onOpenChange: (open: boolean) => void
  onSent: () => void
  initialChannel?: MessageChannel
}

export function SendTestModal({
  workspaceId,
  open,
  onOpenChange,
  onSent,
  initialChannel = "email",
}: SendTestModalProps) {
  const [channel, setChannel] = useState<MessageChannel>(initialChannel)
  const [to, setTo] = useState("")
  const [subject, setSubject] = useState("")
  const [body, setBody] = useState("")
  const [html, setHtml] = useState("")
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setIsLoading(true)
    setError(null)

    let payload: Record<string, unknown> = {}
    if (channel === "email") {
      payload = {
        to: to.split(",").map((s) => s.trim()),
        subject,
        html,
      }
    } else if (channel === "sms") {
      payload = { to, body }
    } else if (channel === "push") {
      payload = { to, title: subject, body }
    } else if (channel === "chat") {
      payload = { to, body }
    }

    const res = await sendTestRequest(workspaceId, channel, payload)
    setIsLoading(false)

    if (!res.ok) {
      setError(res.message ?? "Failed to trigger test request")
    } else {
      onSent()
      onOpenChange(false)
    }
  }

  const handleChannelChange = (newChannel: "email" | "sms" | "push" | "chat") => {
    setChannel(newChannel)
    setError(null)
    if (newChannel === "email") {
      setTo("user@example.com")
      setSubject("Welcome to Message Gateway!")
    } else if (newChannel === "sms") {
      setTo("+15551234567")
      setBody("This is a test SMS message.")
    } else if (newChannel === "push") {
      setTo("device-token-12345")
      setSubject("New Notification")
      setBody("This is a test push notification.")
    } else if (newChannel === "chat") {
      setTo("slack-channel-id")
      setBody("Hello from the Chat API!")
    }
  }

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md gap-0 p-0 overflow-hidden">
        <DialogHeader className="border-b px-5 py-4 bg-muted/20 text-left">
          <DialogTitle>Send Test Request</DialogTitle>
          <DialogDescription>Trigger a gateway send for your workspace.</DialogDescription>
        </DialogHeader>

        <form onSubmit={handleSubmit} className="flex flex-col">
          <div className="flex flex-col gap-4 p-5 max-h-[70vh] overflow-y-auto">
            {error && (
              <div
                role="alert"
                className="flex items-center gap-2 rounded-lg border border-destructive/20 bg-destructive/10 p-3 text-xs text-destructive"
              >
                <Icon name="error" size="sm" className="shrink-0 text-destructive" />
                <span>{error}</span>
              </div>
            )}

            <div className="flex flex-col gap-1.5">
              <span className="text-xs font-semibold uppercase text-muted-foreground">Channel</span>
              <div className="grid grid-cols-4 gap-1 rounded-lg bg-muted p-1">
                {(["email", "sms", "push", "chat"] as const).map((ch) => (
                  <button
                    key={ch}
                    type="button"
                    onClick={() => handleChannelChange(ch)}
                    className={`rounded-md py-1.5 text-xs font-medium capitalize transition-all ${
                      channel === ch
                        ? "bg-card text-foreground shadow-xs"
                        : "text-muted-foreground hover:text-foreground"
                    }`}
                  >
                    {ch}
                  </button>
                ))}
              </div>
            </div>

            <div className="flex flex-col gap-1.5">
              <label htmlFor="send-test-to" className="text-xs font-semibold uppercase text-muted-foreground">
                {channel === "email" ? "To (comma-separated)" : "Recipient (To)"}
              </label>
              <Input
                id="send-test-to"
                required
                aria-invalid={error ? true : undefined}
                placeholder={
                  channel === "email"
                    ? "user@example.com, another@example.com"
                    : channel === "sms"
                      ? "+15551234567"
                      : "recipient-id"
                }
                value={to}
                onChange={(e) => setTo(e.target.value)}
                className="bg-background text-sm"
              />
            </div>

            {(channel === "email" || channel === "push") && (
              <div className="flex flex-col gap-1.5">
                <label htmlFor="send-test-subject" className="text-xs font-semibold uppercase text-muted-foreground">
                  {channel === "email" ? "Subject" : "Notification Title"}
                </label>
                <Input
                  id="send-test-subject"
                  required
                  placeholder={channel === "email" ? "Subject line" : "Title text"}
                  value={subject}
                  onChange={(e) => setSubject(e.target.value)}
                  className="bg-background text-sm"
                />
              </div>
            )}

            {channel === "email" && (
              <div className="flex flex-col gap-1.5">
                <label htmlFor="send-test-html" className="text-xs font-semibold uppercase text-muted-foreground">
                  HTML Content
                </label>
                <textarea
                  id="send-test-html"
                  required
                  rows={4}
                  placeholder="<h1>HTML...</h1>"
                  value={html}
                  onChange={(e) => setHtml(e.target.value)}
                  className="w-full rounded-lg border border-input bg-background px-3 py-2 font-mono text-sm focus-visible:outline-hidden focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>
            )}

            {channel !== "email" && (
              <div className="flex flex-col gap-1.5">
                <label htmlFor="send-test-body" className="text-xs font-semibold uppercase text-muted-foreground">
                  Message Body
                </label>
                <textarea
                  id="send-test-body"
                  required
                  rows={3}
                  placeholder="Enter message..."
                  value={body}
                  onChange={(e) => setBody(e.target.value)}
                  className="w-full rounded-lg border border-input bg-background px-3 py-2 text-sm focus-visible:outline-hidden focus-visible:ring-2 focus-visible:ring-ring"
                />
              </div>
            )}
          </div>

          <DialogFooter className="border-t bg-muted/20 px-5 py-4 sm:justify-end">
            <Button type="button" variant="outline" onClick={() => onOpenChange(false)} disabled={isLoading}>
              Cancel
            </Button>
            <Button type="submit" disabled={isLoading}>
              {isLoading ? (
                <>
                  <span className="size-3.5 animate-spin rounded-full border-2 border-primary-foreground border-t-transparent" />
                  Sending...
                </>
              ) : (
                <>
                  <Icon name="send" size="sm" data-icon="inline-start" />
                  Send Request
                </>
              )}
            </Button>
          </DialogFooter>
        </form>
      </DialogContent>
    </Dialog>
  )
}
