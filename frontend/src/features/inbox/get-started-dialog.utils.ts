const GATEWAY_BASE_URL = "http://localhost:10101"
const SECRET_PLACEHOLDER = "YOUR_CLIENT_SECRET"

export interface GatewayCurlExamples {
  email: string
  sms: string
}

export function buildGatewayCurlExamples(
  workspaceId: string,
  clientId: string,
): GatewayCurlExamples {
  const auth = `${clientId}:${SECRET_PLACEHOLDER}`

  return {
    email: `curl -X POST ${GATEWAY_BASE_URL}/v1/email \\
  -u "${auth}" \\
  -H "X-Workspace-Key: ${workspaceId}" \\
  -H "Content-Type: application/json" \\
  -d '{"to":["user@example.com"],"subject":"Test","html":"<p>Hi</p>"}'`,
    sms: `curl -X POST ${GATEWAY_BASE_URL}/v1/sms \\
  -u "${auth}" \\
  -H "X-Workspace-Key: ${workspaceId}" \\
  -H "Content-Type: application/json" \\
  -d '{"to":["+1234567890"],"message":"Hello"}'`,
  }
}
