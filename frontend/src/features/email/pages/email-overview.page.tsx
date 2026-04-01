import { Button } from "@/components/ui/button"

export function EmailOverviewPage() {
  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-semibold tracking-tight">Email</h1>
        <p className="text-muted-foreground">
          Configure SMTP, SendGrid, SES, or Mailgun integrations and monitor request logs.
        </p>
      </div>
      <Button type="button">Send test request</Button>
      <p className="text-sm text-muted-foreground">
        Email channel is the supported MVP; additional channels will appear as the API expands.
      </p>
    </div>
  )
}
