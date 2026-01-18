import { useState } from 'react'
import { Mail, MessageSquare, Bell, MessagesSquare, Trash2, Moon, Sun, Archive } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useStats, useEmails, useSMS, usePush, useChat, useDeleteMessage, useClearAll, useSSE } from '@/hooks/useMessages'
import { useDarkMode } from '@/hooks/useDarkMode'
import { EmailList, EmailDetail } from '@/features/email'
import { SMSList } from '@/features/sms'
import { PushList } from '@/features/push'
import { ChatList, ChatDetail } from '@/features/chat'
import type { StoredEmail, StoredChat, MessageType } from '@/types/messages'
import { cn } from '@/lib/utils'
import { BADGE_MAX_COUNT } from '@/lib/constants'

const NAV_ITEMS = [
  { type: 'email' as const, icon: Mail, label: 'Email' },
  { type: 'sms' as const, icon: MessageSquare, label: 'SMS' },
  { type: 'push' as const, icon: Bell, label: 'Push' },
  { type: 'chat' as const, icon: MessagesSquare, label: 'Chat' },
]

export default function App() {
  const { darkMode, toggle: toggleDarkMode } = useDarkMode()
  const [activeType, setActiveType] = useState<MessageType>('email')
  const [selectedEmail, setSelectedEmail] = useState<StoredEmail | null>(null)
  const [selectedChat, setSelectedChat] = useState<StoredChat | null>(null)

  useSSE()

  const { data: stats } = useStats()
  const { data: emails = [] } = useEmails()
  const { data: sms = [] } = useSMS()
  const { data: push = [] } = usePush()
  const { data: chat = [] } = useChat()
  const deleteMessage = useDeleteMessage()
  const clearAll = useClearAll()

  const getCount = (type: MessageType): number => {
    const counts: Record<MessageType, number> = {
      email: stats?.emails ?? 0,
      sms: stats?.sms ?? 0,
      push: stats?.push ?? 0,
      chat: stats?.chat ?? 0,
    }
    return counts[type]
  }

  const handleNavClick = (type: MessageType) => {
    setActiveType(type)
    setSelectedEmail(null)
    setSelectedChat(null)
  }

  const handleDelete = (type: string, id: string) => {
    deleteMessage.mutate({ type, id })
  }

  const hasDetailPanel = activeType === 'email' || activeType === 'chat'
  const totalCount = stats?.total ?? 0

  return (
    <div className="h-screen flex flex-col bg-background">
      <Header
        totalCount={totalCount}
        darkMode={darkMode}
        onToggleDarkMode={toggleDarkMode}
        onClear={() => clearAll.mutate()}
        isClearDisabled={clearAll.isPending || totalCount === 0}
      />

      <div className="flex flex-1 overflow-hidden">
        <Sidebar
          activeType={activeType}
          getCount={getCount}
          onNavClick={handleNavClick}
        />

        <MessageListPanel hasDetailPanel={hasDetailPanel}>
          <div className="p-3 border-b shrink-0">
            <h2 className="font-semibold capitalize">{activeType}</h2>
          </div>
          <ScrollArea className="flex-1">
            {activeType === 'email' && (
              <EmailList
                emails={emails}
                selected={selectedEmail}
                onSelect={setSelectedEmail}
                onDelete={(id) => handleDelete('emails', id)}
              />
            )}
            {activeType === 'sms' && (
              <SMSList
                messages={sms}
                onDelete={(id) => handleDelete('sms', id)}
              />
            )}
            {activeType === 'push' && (
              <PushList
                notifications={push}
                onDelete={(id) => handleDelete('push', id)}
              />
            )}
            {activeType === 'chat' && (
              <ChatList
                messages={chat}
                selected={selectedChat}
                onSelect={setSelectedChat}
                onDelete={(id) => handleDelete('chat', id)}
              />
            )}
          </ScrollArea>
        </MessageListPanel>

        {activeType === 'email' && (
          <DetailPanel>
            {selectedEmail ? (
              <EmailDetail email={selectedEmail} />
            ) : (
              <EmptyDetail message="Select an email to view" />
            )}
          </DetailPanel>
        )}
        {activeType === 'chat' && (
          <DetailPanel>
            {selectedChat ? (
              <ChatDetail chat={selectedChat} />
            ) : (
              <EmptyDetail message="Select a message to view" />
            )}
          </DetailPanel>
        )}
      </div>
    </div>
  )
}

// --- Sub-components ---

interface HeaderProps {
  totalCount: number
  darkMode: boolean
  onToggleDarkMode: () => void
  onClear: () => void
  isClearDisabled: boolean
}

function Header({ totalCount, darkMode, onToggleDarkMode, onClear, isClearDisabled }: HeaderProps) {
  return (
    <header className="flex items-center justify-between px-4 py-2 border-b shrink-0">
      <div className="flex items-center gap-2">
        <Archive className="h-5 w-5 text-primary" />
        <span className="font-semibold">Message Gateway - DevBox</span>
        <Badge variant="secondary" className="ml-2">{totalCount}</Badge>
      </div>
      <div className="flex items-center gap-2">
        <Button
          variant="ghost"
          size="sm"
          onClick={onClear}
          disabled={isClearDisabled}
        >
          <Trash2 className="h-4 w-4 mr-1" />
          Clear
        </Button>
        <Button variant="ghost" size="icon" onClick={onToggleDarkMode}>
          {darkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
        </Button>
      </div>
    </header>
  )
}

interface SidebarProps {
  activeType: MessageType
  getCount: (type: MessageType) => number
  onNavClick: (type: MessageType) => void
}

function Sidebar({ activeType, getCount, onNavClick }: SidebarProps) {
  return (
    <aside className="w-14 border-r flex flex-col items-center py-2 gap-1 shrink-0">
      {NAV_ITEMS.map(({ type, icon: Icon, label }) => {
        const count = getCount(type)
        return (
          <Button
            key={type}
            variant={activeType === type ? 'secondary' : 'ghost'}
            size="icon"
            className="relative"
            onClick={() => onNavClick(type)}
            title={label}
          >
            <Icon className="h-4 w-4" />
            {count > 0 && <CountBadge count={count} />}
          </Button>
        )
      })}
    </aside>
  )
}

function CountBadge({ count }: { count: number }) {
  const displayCount = count > BADGE_MAX_COUNT ? `${BADGE_MAX_COUNT}+` : count
  return (
    <span className="absolute -top-1 -right-1 h-4 w-4 text-[10px] bg-primary text-primary-foreground rounded-full flex items-center justify-center">
      {displayCount}
    </span>
  )
}

interface MessageListPanelProps {
  hasDetailPanel: boolean
  children: React.ReactNode
}

function MessageListPanel({ hasDetailPanel, children }: MessageListPanelProps) {
  return (
    <div className={cn('border-r flex flex-col', hasDetailPanel ? 'w-80' : 'flex-1')}>
      {children}
    </div>
  )
}

function DetailPanel({ children }: { children: React.ReactNode }) {
  return <div className="flex-1 flex flex-col overflow-hidden">{children}</div>
}

function EmptyDetail({ message }: { message: string }) {
  return (
    <div className="flex-1 flex items-center justify-center text-muted-foreground">
      <p>{message}</p>
    </div>
  )
}
