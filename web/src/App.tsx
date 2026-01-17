import { useState } from 'react'
import { Mail, MessageSquare, Bell, MessagesSquare, Trash2, Moon, Sun, Archive } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { ScrollArea } from '@/components/ui/scroll-area'
import { useStats, useEmails, useSMS, usePush, useChat, useDeleteMessage, useClearAll, useSSE } from '@/hooks/useMessages'
import { EmailList, EmailDetail } from '@/features/email'
import { SMSList } from '@/features/sms'
import { PushList } from '@/features/push'
import { ChatList, ChatDetail } from '@/features/chat'
import type { StoredEmail, StoredChat } from '@/types/messages'
import { cn } from '@/lib/utils'

type MessageType = 'email' | 'sms' | 'push' | 'chat'

const NAV_ITEMS = [
  { type: 'email' as const, icon: Mail, label: 'Email' },
  { type: 'sms' as const, icon: MessageSquare, label: 'SMS' },
  { type: 'push' as const, icon: Bell, label: 'Push' },
  { type: 'chat' as const, icon: MessagesSquare, label: 'Chat' },
]

export default function App() {
  const [darkMode, setDarkMode] = useState(true)
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

  if (darkMode) {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }

  const getCount = (type: MessageType) => {
    switch (type) {
      case 'email': return stats?.emails ?? 0
      case 'sms': return stats?.sms ?? 0
      case 'push': return stats?.push ?? 0
      case 'chat': return stats?.chat ?? 0
    }
  }

  const handleNavClick = (type: MessageType) => {
    setActiveType(type)
    setSelectedEmail(null)
    setSelectedChat(null)
  }

  const hasDetailPanel = activeType === 'email' || activeType === 'chat'

  return (
    <div className="h-screen flex flex-col bg-background">
      {/* Header */}
      <header className="flex items-center justify-between px-4 py-2 border-b shrink-0">
        <div className="flex items-center gap-2">
          <Archive className="h-5 w-5 text-primary" />
          <span className="font-semibold">DevBox</span>
          <Badge variant="secondary" className="ml-2">{stats?.total ?? 0}</Badge>
        </div>
        <div className="flex items-center gap-2">
          <Button
            variant="ghost"
            size="sm"
            onClick={() => clearAll.mutate()}
            disabled={clearAll.isPending || (stats?.total ?? 0) === 0}
          >
            <Trash2 className="h-4 w-4 mr-1" />
            Clear
          </Button>
          <Button variant="ghost" size="icon" onClick={() => setDarkMode(!darkMode)}>
            {darkMode ? <Sun className="h-4 w-4" /> : <Moon className="h-4 w-4" />}
          </Button>
        </div>
      </header>

      <div className="flex flex-1 overflow-hidden">
        {/* Sidebar */}
        <aside className="w-14 border-r flex flex-col items-center py-2 gap-1 shrink-0">
          {NAV_ITEMS.map(({ type, icon: Icon, label }) => {
            const count = getCount(type)
            return (
              <Button
                key={type}
                variant={activeType === type ? 'secondary' : 'ghost'}
                size="icon"
                className="relative"
                onClick={() => handleNavClick(type)}
                title={label}
              >
                <Icon className="h-4 w-4" />
                {count > 0 && (
                  <span className="absolute -top-1 -right-1 h-4 w-4 text-[10px] bg-primary text-primary-foreground rounded-full flex items-center justify-center">
                    {count > 99 ? '99+' : count}
                  </span>
                )}
              </Button>
            )
          })}
        </aside>

        {/* Message List */}
        <div className={cn("border-r flex flex-col", hasDetailPanel ? "w-80" : "flex-1")}>
          <div className="p-3 border-b shrink-0">
            <h2 className="font-semibold capitalize">{activeType}</h2>
          </div>
          <ScrollArea className="flex-1">
            {activeType === 'email' && (
              <EmailList
                emails={emails}
                selected={selectedEmail}
                onSelect={setSelectedEmail}
                onDelete={(id) => deleteMessage.mutate({ type: 'emails', id })}
              />
            )}
            {activeType === 'sms' && (
              <SMSList
                messages={sms}
                onDelete={(id) => deleteMessage.mutate({ type: 'sms', id })}
              />
            )}
            {activeType === 'push' && (
              <PushList
                notifications={push}
                onDelete={(id) => deleteMessage.mutate({ type: 'push', id })}
              />
            )}
            {activeType === 'chat' && (
              <ChatList
                messages={chat}
                selected={selectedChat}
                onSelect={setSelectedChat}
                onDelete={(id) => deleteMessage.mutate({ type: 'chat', id })}
              />
            )}
          </ScrollArea>
        </div>

        {/* Detail Panel */}
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
