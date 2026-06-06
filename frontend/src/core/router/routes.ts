export const ROUTES = {
  root: "/",
  login: "/login",
  register: "/register",
  workspaces: "/workspaces",
  workspace: {
    pattern: "/workspaces/:wid",
    overview: (workspaceId: string) => `/workspaces/${workspaceId}/overview`,
    email: (workspaceId: string) => `/workspaces/${workspaceId}/email`,
    emailTemplates: (workspaceId: string) => `/workspaces/${workspaceId}/email/templates`,
    sms: (workspaceId: string) => `/workspaces/${workspaceId}/sms`,
    push: (workspaceId: string) => `/workspaces/${workspaceId}/push`,
    chat: (workspaceId: string) => `/workspaces/${workspaceId}/chat`,
  },
} as const
