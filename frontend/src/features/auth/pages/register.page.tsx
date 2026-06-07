import { type FormEvent, useState } from "react"
import { Link, useNavigate } from "react-router-dom"

import { ROUTES } from "@/core/router/routes"
import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Input } from "@/components/ui/input"
import { AuthCard } from "@/shared/components/auth-card"
import { registerAccount } from "../auth.api"

export function RegisterPage() {
  const navigate = useNavigate()
  const [email, setEmail] = useState("")
  const [password, setPassword] = useState("")
  const [displayName, setDisplayName] = useState("")
  const [rememberMe, setRememberMe] = useState(true)
  const [showPassword, setShowPassword] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  async function onSubmit(e: FormEvent) {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      const result = await registerAccount(email, password, displayName)
      if (!result.ok) {
        setError(result.message)
        return
      }
      navigate(ROUTES.workspaces, { replace: true })
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="flex flex-1 items-center justify-center px-6 py-16">
      <AuthCard>
        <form onSubmit={onSubmit} className="flex w-full flex-col gap-6">
          <div className="flex flex-col gap-2">
            <h2 className="text-2xl font-semibold leading-8 text-foreground">Sign up</h2>
            <p className="text-sm leading-5 text-text-secondary">
              Create your account to start using Message Gateway.
            </p>
          </div>

          <div className="flex flex-col gap-5">
            <div className="flex flex-col gap-2">
              <label htmlFor="display" className="text-sm font-medium text-text-secondary">
                Full name
              </label>
              <div className="relative">
                <Icon
                  name="person"
                  size="sm"
                  className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-placeholder"
                />
                <Input
                  id="display"
                  type="text"
                  autoComplete="name"
                  placeholder="Jane Doe"
                  className="pl-10"
                  value={displayName}
                  onChange={(ev) => setDisplayName(ev.target.value)}
                />
              </div>
            </div>

            <div className="flex flex-col gap-2">
              <label htmlFor="email" className="text-sm font-medium text-text-secondary">
                Email address
              </label>
              <div className="relative">
                <Icon
                  name="mail"
                  size="sm"
                  className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-placeholder"
                />
                <Input
                  id="email"
                  type="email"
                  autoComplete="email"
                  required
                  placeholder="you@company.com"
                  className="pl-10"
                  value={email}
                  onChange={(ev) => setEmail(ev.target.value)}
                />
              </div>
            </div>

            <div className="flex flex-col gap-2">
              <label htmlFor="password" className="text-sm font-medium text-text-secondary">
                Password
              </label>
              <div className="relative">
                <Icon
                  name="lock"
                  size="sm"
                  className="pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 text-text-placeholder"
                />
                <Input
                  id="password"
                  type={showPassword ? "text" : "password"}
                  autoComplete="new-password"
                  required
                  placeholder="••••••••••••"
                  className="px-10"
                  value={password}
                  onChange={(ev) => setPassword(ev.target.value)}
                />
                <button
                  type="button"
                  onClick={() => setShowPassword((v) => !v)}
                  className="absolute right-3 top-1/2 -translate-y-1/2 text-text-placeholder"
                  aria-label={showPassword ? "Hide password" : "Show password"}
                >
                  <Icon name={showPassword ? "visibility_off" : "visibility"} size="sm" />
                </button>
              </div>
            </div>

            <label className="flex cursor-pointer items-center gap-2 text-sm text-text-secondary">
              <input
                type="checkbox"
                checked={rememberMe}
                onChange={(ev) => setRememberMe(ev.target.checked)}
                className="size-4 rounded border-border accent-primary-brand"
              />
              Remember me
            </label>
          </div>

          {error ? <p className="text-sm text-destructive">{error}</p> : null}

          <Button type="submit" className="h-11 w-full" disabled={loading}>
            {loading ? "Creating…" : "Sign up"}
          </Button>

          <div className="flex items-center gap-4">
            <div className="h-px flex-1 bg-divider" />
            <span className="text-xs font-medium text-text-tertiary">OR</span>
            <div className="h-px flex-1 bg-divider" />
          </div>

          <Button type="button" variant="secondary" className="h-11 w-full">
            Continue with Gmail
          </Button>

          <p className="text-center text-sm text-text-secondary">
            Already have an account?{" "}
            <Link
              to={ROUTES.login}
              className="font-semibold text-primary-brand underline transition-colors hover:text-primary-brand-hover"
            >
              Sign in
            </Link>
          </p>
        </form>
      </AuthCard>
    </div>
  )
}
