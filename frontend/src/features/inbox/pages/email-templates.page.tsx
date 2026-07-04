import { useEffect, useState, useCallback } from "react"
import { useParams, useNavigate } from "react-router-dom"

import { Button } from "@/components/ui/button"
import { Icon } from "@/components/ui/icon"
import { Modal } from "@/components/ui/modal"
import { Input } from "@/components/ui/input"
import { SearchInput } from "@/components/ui/search-input"
import { Textarea } from "@/components/ui/textarea"
import { Spinner } from "@/components/ui/spinner"
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableHead,
  TableCell,
} from "@/components/ui/table"
import { ROUTES } from "@/core/router/routes"
import { Permission, useWorkspaceAuthorization } from "@/core/auth"
import {
  createEmailTemplate,
  deleteEmailTemplate,
  fetchEmailTemplates,
} from "../api/inbox.api"
import type { EmailTemplate } from "../inbox.types"

type TabType = "templates" | "headers" | "footers"

function getCategoryColor(category: string) {
  switch (category.toLowerCase()) {
    case "transactional":
      return "bg-indigo-500/10 text-indigo-600 dark:text-indigo-400 font-semibold"
    case "marketing":
      return "bg-emerald-500/10 text-emerald-600 dark:text-emerald-400 font-semibold"
    case "notification":
      return "bg-amber-500/10 text-amber-600 dark:text-amber-400 font-semibold"
    default:
      return "bg-muted text-text-secondary font-semibold"
  }
}

export function EmailTemplatesPage() {
  const { wid } = useParams<{ wid: string }>()
  const navigate = useNavigate()
  const { can } = useWorkspaceAuthorization()
  const canManageTemplates = can(Permission.TemplatesWrite)

  const [activeTab, setActiveTab] = useState<TabType>("templates")
  const [templates, setTemplates] = useState<EmailTemplate[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  
  const [searchQuery, setSearchQuery] = useState("")
  const [isCreateOpen, setIsCreateOpen] = useState(false)
  
  // Create Form State
  const [formName, setFormName] = useState("")
  const [formKey, setFormKey] = useState("")
  const [formSubject, setFormSubject] = useState("")
  const [formContent, setFormContent] = useState("")
  const [formCategory, setFormCategory] = useState("transactional")
  const [formErrors, setFormErrors] = useState<Record<string, string>>({})
  const [isSaving, setIsSaving] = useState(false)

  const fetchTemplates = useCallback(async () => {
    if (!wid) return
    await Promise.resolve() // yield to avoid synchronous state update in effect
    setIsLoading(true)
    setError(null)
    try {
      const result = await fetchEmailTemplates(wid)
      if (!result.ok) {
        setError(result.message ?? "Failed to load templates")
      } else {
        setTemplates(result.items)
      }
    } finally {
      setIsLoading(false)
    }
  }, [wid])

  useEffect(() => {
    Promise.resolve().then(() => {
      void fetchTemplates()
    })
  }, [fetchTemplates])

  const handleDelete = async (id: string) => {
    if (!wid || !confirm("Are you sure you want to delete this template?")) return
    const result = await deleteEmailTemplate(wid, id)
    if (!result.ok) {
      alert(result.message ?? "Failed to delete template")
      return
    }
    setTemplates((prev) => prev.filter((t) => t.id !== id))
  }

  const handleNameChange = (val: string) => {
    setFormName(val)
    const slug = val
      .toLowerCase()
      .replace(/[^a-z0-9]+/g, "_")
      .replace(/(^-|-$)/g, "")
    setFormKey(slug)
  }

  const handleCreateSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    
    // Validate
    const errors: Record<string, string> = {}
    if (!formName.trim()) errors.name = "Name is required"
    if (!formKey.trim()) errors.key = "Unique key is required"
    if (!formSubject.trim()) errors.subject = "Subject is required"
    if (!formContent.trim()) errors.content = "HTML content is required"
    setFormErrors(errors)
    if (Object.keys(errors).length > 0) return

    if (!wid) return

    setIsSaving(true)
    try {
      const result = await createEmailTemplate(wid, {
        name: formName.trim(),
        unique_key: formKey.trim(),
        category: formCategory,
        subject: formSubject.trim(),
        content_html: formContent.trim(),
      })

      if (!result.ok) {
        alert(result.message ?? "Failed to create template")
        return
      }

      setTemplates((prev) => [...prev, result.item])
      handleCreateClose()
    } finally {
      setIsSaving(false)
    }
  }

  const handleCreateClose = () => {
    setFormName("")
    setFormKey("")
    setFormSubject("")
    setFormContent("")
    setFormCategory("transactional")
    setFormErrors({})
    setIsCreateOpen(false)
  }

  const handleBack = () => {
    if (wid) {
      navigate(ROUTES.workspace.email(wid))
    }
  }

  const filtered = templates.filter((t) => {
    const q = searchQuery.toLowerCase()
    return (
      t.name.toLowerCase().includes(q) ||
      t.unique_key.toLowerCase().includes(q) ||
      (t.subject && t.subject.toLowerCase().includes(q))
    )
  })

  return (
    <div className="flex flex-col gap-6 w-full">
      {/* Back Button */}
      <button
        type="button"
        onClick={handleBack}
        className="flex items-center gap-2 text-sm font-semibold text-text-secondary hover:text-foreground transition-colors self-start outline-none"
      >
        <Icon name="arrow_back" size="sm" />
        Back to Email Inbox
      </button>

      {/* Header Panel */}
      <div className="flex flex-col gap-2">
        <h1 className="text-2xl font-bold tracking-tight text-foreground">
          Email Templates
        </h1>
        <p className="text-sm text-text-secondary max-w-2xl">
          Manage reusable layout assets. Headers and footers can be selected dynamically during sending payload dispatching.
        </p>
      </div>

      {/* Tab Selectors */}
      <div className="flex gap-2 items-center">
        {(["templates", "headers", "footers"] as const).map((tab) => (
          <button
            key={tab}
            type="button"
            onClick={() => setActiveTab(tab)}
            className={`px-4 py-2 rounded-lg text-sm font-semibold border transition-all ${
              activeTab === tab
                ? "bg-primary-brand border-primary-brand text-white shadow-xs"
                : "bg-card border-border text-text-secondary hover:bg-secondary"
            }`}
          >
            <span className="capitalize">{tab}</span>
          </button>
        ))}
      </div>

      {/* Action Bar */}
      <div className="flex flex-col sm:flex-row gap-3 sm:items-center sm:justify-between w-full">
        <SearchInput
          placeholder={`Search ${activeTab}...`}
          value={searchQuery}
          onChange={(e) => setSearchQuery(e.target.value)}
          className="w-64 shadow-xs shrink-0"
        />

        {activeTab === "templates" && canManageTemplates ? (
          <Button type="button" size="sm" onClick={() => setIsCreateOpen(true)} className="h-10 shrink-0">
            <Icon name="add" size="sm" data-icon="inline-start" />
            Create Template
          </Button>
        ) : null}
      </div>

      {/* Table Content */}
      {isLoading ? (
        <div className="flex flex-col items-center justify-center gap-3 py-20 text-sm text-text-secondary bg-card border border-border rounded-[16px]">
          <Spinner />
          <span>Loading assets...</span>
        </div>
      ) : error ? (
        <div className="bg-card border border-border rounded-[16px] py-16 text-center">
          <p className="text-sm text-text-secondary font-semibold mb-2">Error loading templates</p>
          <p className="text-xs text-destructive mb-4">{error}</p>
          <Button variant="outline" size="sm" onClick={() => void fetchTemplates()}>
            Try again
          </Button>
        </div>
      ) : activeTab !== "templates" ? (
        <div className="bg-card border border-border rounded-[16px] p-12 text-center shadow-xs">
          <div className="mb-3 flex size-12 items-center justify-center rounded-full border border-border bg-muted/40 text-text-tertiary mx-auto">
            <Icon name="dashboard_customize" size="md" />
          </div>
          <p className="text-sm font-semibold text-foreground capitalize">No {activeTab} defined yet</p>
          <p className="text-xs text-text-secondary mt-1">
            Layout extensions will be available in future releases.
          </p>
        </div>
      ) : filtered.length === 0 ? (
        <div className="bg-card border border-border rounded-[16px] p-12 text-center shadow-xs">
          <p className="text-sm font-semibold text-foreground">No templates found</p>
          <p className="text-xs text-text-secondary mt-1 mb-4">
            {searchQuery ? "No templates match your search criteria." : "Create your first reusable email template template."}
          </p>
          {!searchQuery && canManageTemplates ? (
            <Button size="sm" onClick={() => setIsCreateOpen(true)}>
              <Icon name="add" size="sm" data-icon="inline-start" />
              Create Template
            </Button>
          ) : null}
        </div>
      ) : (
        <div className="bg-card rounded-[16px] border border-border shadow-xs overflow-hidden">
          <Table>
            <TableHeader>
              <TableRow className="bg-muted/30 hover:bg-muted/30 border-b border-border">
                <TableHead className="px-6 py-3 h-auto w-[25%] font-semibold uppercase tracking-wider text-text-secondary">Asset Name / Key</TableHead>
                <TableHead className="px-6 py-3 h-auto w-[33%] font-semibold uppercase tracking-wider text-text-secondary">Subject</TableHead>
                <TableHead className="px-6 py-3 h-auto w-[17%] font-semibold uppercase tracking-wider text-text-secondary">Category</TableHead>
                <TableHead className="px-6 py-3 h-auto w-[17%] font-semibold uppercase tracking-wider text-text-secondary">Status</TableHead>
                <TableHead className="px-6 py-3 h-auto w-[8%] text-right font-semibold uppercase tracking-wider text-text-secondary">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((t) => (
                <TableRow
                  key={t.id}
                  className="hover:bg-muted/10 border-b border-border last:border-0"
                >
                  <TableCell className="px-6 py-4 font-semibold text-foreground truncate pr-4">
                    <div>
                      <p className="font-semibold text-foreground truncate">{t.name}</p>
                      <span className="text-xs text-text-secondary font-mono truncate block mt-0.5">
                        {t.unique_key}
                      </span>
                    </div>
                  </TableCell>
                  <TableCell className="px-6 py-4 text-text-secondary truncate pr-4" title={t.subject}>
                    {t.subject || "(No Subject)"}
                  </TableCell>
                  <TableCell className="px-6 py-4">
                    <span className={`rounded-full px-2.5 py-0.5 text-xs ${getCategoryColor(t.category)}`}>
                      {t.category}
                    </span>
                  </TableCell>
                  <TableCell className="px-6 py-4">
                    <span className={`text-xs font-semibold ${t.is_active ? "text-emerald-600" : "text-text-tertiary"}`}>
                      {t.is_active ? "Active" : "Inactive"}
                    </span>
                  </TableCell>
                  <TableCell className="px-6 py-4 text-right">
                    {canManageTemplates ? (
                      <button
                        type="button"
                        onClick={() => handleDelete(t.id)}
                        className="p-1.5 rounded bg-rose-500/10 text-rose-600 hover:bg-rose-500/20 transition-colors"
                        title="Delete Template"
                      >
                        <Icon name="delete" size="sm" />
                      </button>
                    ) : (
                      <span className="text-xs text-text-tertiary">—</span>
                    )}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </div>
      )}

      {/* Create Modal */}
      <Modal isOpen={isCreateOpen} onClose={handleCreateClose} title="Create Template">
        <form onSubmit={handleCreateSubmit} className="flex flex-col gap-4 max-h-[80vh] overflow-y-auto pr-1">
          <div className="flex flex-col gap-1">
            <label htmlFor="t-name" className="text-xs font-semibold text-text-secondary uppercase">
              Template Name
            </label>
            <Input
              id="t-name"
              type="text"
              value={formName}
              onChange={(e) => handleNameChange(e.target.value)}
              placeholder="e.g. Welcome Message"
              className="bg-secondary border-border"
            />
            {formErrors.name && <p className="text-xs text-destructive mt-0.5">{formErrors.name}</p>}
          </div>

          <div className="flex flex-col gap-1">
            <label htmlFor="t-key" className="text-xs font-semibold text-text-secondary uppercase">
              Unique Key
            </label>
            <Input
              id="t-key"
              type="text"
              value={formKey}
              onChange={(e) => setFormKey(e.target.value)}
              placeholder="e.g. welcome_message"
              className="bg-secondary border-border font-mono"
            />
            {formErrors.key && <p className="text-xs text-destructive mt-0.5">{formErrors.key}</p>}
          </div>

          <div className="flex flex-col gap-1">
            <label htmlFor="t-category" className="text-xs font-semibold text-text-secondary uppercase">
              Category
            </label>
            <select
              id="t-category"
              value={formCategory}
              onChange={(e) => setFormCategory(e.target.value)}
              className="bg-secondary border border-border rounded-lg h-10 px-3 text-sm text-foreground outline-none focus:border-primary-brand transition-colors"
            >
              <option value="transactional">Transactional</option>
              <option value="marketing">Marketing</option>
              <option value="notification">Notification</option>
            </select>
          </div>

          <div className="flex flex-col gap-1">
            <label htmlFor="t-subject" className="text-xs font-semibold text-text-secondary uppercase">
              Subject
            </label>
            <Input
              id="t-subject"
              type="text"
              value={formSubject}
              onChange={(e) => setFormSubject(e.target.value)}
              placeholder="Enter email subject"
              className="bg-secondary border-border"
            />
            {formErrors.subject && <p className="text-xs text-destructive mt-0.5">{formErrors.subject}</p>}
          </div>

          <div className="flex flex-col gap-1">
            <label htmlFor="t-content" className="text-xs font-semibold text-text-secondary uppercase">
              HTML Content
            </label>
            <Textarea
              id="t-content"
              rows={6}
              value={formContent}
              onChange={(e) => setFormContent(e.target.value)}
              placeholder="<h1>Welcome!</h1><p>Thanks for signing up.</p>"
              className="bg-secondary border-border font-mono resize-y"
            />
            {formErrors.content && <p className="text-xs text-destructive mt-0.5">{formErrors.content}</p>}
          </div>

          <div className="flex flex-col gap-3 pt-2">
            <Button type="submit" disabled={isSaving}>
              {isSaving ? "Saving..." : "Save Template"}
            </Button>
            <button
              type="button"
              onClick={handleCreateClose}
              className="w-full bg-secondary hover:bg-muted text-text-secondary border border-border h-10 rounded-lg font-semibold text-sm transition-all"
            >
              Cancel
            </button>
          </div>
        </form>
      </Modal>
    </div>
  )
}
