'use client'
export const dynamic = 'force-dynamic'
import * as React from 'react'
import { ShieldAlert, Wand2, Send, Info, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Textarea } from '@/components/ui/textarea'
import { ChatBubble } from '@/components/ui/chat-bubble'
import { StatusDot } from '@/components/ui/status-dot'
import { Popover, PopoverContent, PopoverTrigger } from '@/components/ui/popover'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'
import { api, cacheDebrief, getUserId, getSessionId, type MessageDTO, type ConsentRiskDTO } from '@/lib/api'

type MessageStatus = 'success' | 'warning' | 'danger'

interface Message {
  id: string
  sender: 'user' | 'persona'
  text: string
  status?: MessageStatus
  statusLabel?: string
}

function TypingIndicator() {
  return (
    <div className="flex items-center gap-1.5 px-1 py-1">
      <span className="size-1.5 rounded-full bg-foreground-muted animate-pulse" />
      <span className="size-1.5 rounded-full bg-foreground-muted animate-pulse [animation-delay:150ms]" />
      <span className="size-1.5 rounded-full bg-foreground-muted animate-pulse [animation-delay:300ms]" />
    </div>
  )
}

function ConsentGuardBanner({ risk, onReplace, onKeep }: {
  risk: ConsentRiskDTO; onReplace: () => void; onKeep: () => void
}) {
  return (
    <div className="bg-danger/10 border border-danger/30 rounded-xl p-4 mb-4">
      <div className="flex items-start gap-3">
        <ShieldAlert className="size-5 text-danger shrink-0 mt-0.5" />
        <div className="flex-1 min-w-0">
          <h4 className="text-body font-medium text-foreground mb-1">Это могло прозвучать как давление</h4>
          <p className="text-small text-foreground-secondary mb-3">{risk.explanation ?? 'Фраза может восприниматься как попытка ускорить события.'}</p>
          {risk.suggestion && (
            <p className="text-small text-foreground-muted italic mb-4">"{risk.suggestion}"</p>
          )}
          <div className="flex items-center gap-3">
            <Button size="sm" onClick={onReplace}>Заменить на это</Button>
            <Button variant="ghost" size="sm" onClick={onKeep}>Оставить как есть</Button>
          </div>
        </div>
      </div>
    </div>
  )
}

function LegendPopover({ children }: { children: React.ReactNode }) {
  return (
    <Popover>
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent side="left" align="start" className="w-64 p-4">
        <h4 className="text-small font-medium text-foreground mb-3">Цвета оценки</h4>
        <div className="space-y-2">
          {[
            { color: 'bg-success', label: 'Отлично — сильная реплика' },
            { color: 'bg-warning', label: 'Есть что улучшить' },
            { color: 'bg-danger', label: 'Риск — стоит переформулировать' },
          ].map(({ color, label }) => (
            <div key={label} className="flex items-center gap-2">
              <div className={cn('size-2.5 rounded-full shrink-0', color)} />
              <span className="text-small text-foreground-secondary">{label}</span>
            </div>
          ))}
        </div>
      </PopoverContent>
    </Popover>
  )
}

function SuggestionPopover({ draft, onInsert, children }: {
  draft: string; onInsert: (s: string) => void; children: React.ReactNode
}) {
  const [suggestion, setSuggestion] = React.useState<string | null>(null)
  const [loading, setLoading] = React.useState(false)
  const [open, setOpen] = React.useState(false)

  const handleOpen = async (isOpen: boolean) => {
    setOpen(isOpen)
    if (!isOpen) return
    const sessionId = getSessionId()
    if (!sessionId) return
    setLoading(true)
    try {
      const out = await api.suggest(sessionId, draft)
      setSuggestion(out.suggestion)
    } catch { setSuggestion(null) } finally { setLoading(false) }
  }

  return (
    <Popover open={open} onOpenChange={handleOpen}>
      <PopoverTrigger asChild>{children}</PopoverTrigger>
      <PopoverContent side="top" align="start" className="w-72 p-4">
        <h4 className="text-small font-medium text-foreground mb-2">Подсказка</h4>
        {loading && <p className="text-small text-foreground-muted">Генерируем...</p>}
        {!loading && suggestion && (
          <>
            <p className="text-small text-foreground-secondary mb-3 italic">"{suggestion}"</p>
            <Button size="sm" className="w-full" onClick={() => { onInsert(suggestion); setOpen(false) }}>
              Вставить
            </Button>
          </>
        )}
        {!loading && !suggestion && <p className="text-small text-foreground-muted">Не удалось получить подсказку</p>}
      </PopoverContent>
    </Popover>
  )
}

export default function ChatPage() {
  const [messages, setMessages] = React.useState<Message[]>([])
  const [inputValue, setInputValue] = React.useState('')
  const [isSending, setIsSending] = React.useState(false)
  const [isTyping, setIsTyping] = React.useState(false)
  const [consentRisk, setConsentRisk] = React.useState<ConsentRiskDTO | null>(null)
  const [pendingText, setPendingText] = React.useState<string | null>(null)
  const [shouldSuggestEnd, setShouldSuggestEnd] = React.useState(false)
  const [isFinishing, setIsFinishing] = React.useState(false)
  const [finishError, setFinishError] = React.useState<string | null>(null)
  const [sendError, setSendError] = React.useState<string | null>(null)
  const messagesEndRef = React.useRef<HTMLDivElement>(null)

  const personaTitle = typeof window !== 'undefined' ? sessionStorage.getItem('acc_persona_title') ?? 'Персонаж' : 'Персонаж'
  const personaInitials = typeof window !== 'undefined' ? sessionStorage.getItem('acc_persona_initials') ?? '?' : '?'
  const goalTitle = typeof window !== 'undefined' ? sessionStorage.getItem('acc_goal_title') ?? '' : ''

  React.useEffect(() => {
    if (!getUserId() || !getSessionId()) { window.location.href = '/'; return }
    const opening = sessionStorage.getItem('acc_opening_message')
    if (opening) {
      setMessages([{ id: 'opening', sender: 'persona', text: opening }])
      sessionStorage.removeItem('acc_opening_message')
    }
  }, [])

  React.useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages, isTyping])

  const userMessageCount = messages.filter((m) => m.sender === 'user').length
  const totalSteps = 8

  const doSend = async (text: string) => {
    const sessionId = getSessionId()
    if (!sessionId || !text.trim()) return
    const optimisticId = `pending-${Date.now()}`
    const optimisticMsg: Message = { id: optimisticId, sender: 'user', text }
    setSendError(null)
    setMessages((prev) => [...prev, optimisticMsg])
    setIsSending(true)
    setIsTyping(true)
    try {
      const out = await api.sendMessage(sessionId, text)
      const userMsg: Message = {
        id: out.userMessage.id,
        sender: 'user',
        text: out.userMessage.text,
        status: out.userMessage.status as MessageStatus | undefined,
        statusLabel: out.userMessage.statusLabel,
      }
      const personaMsg: Message = { id: out.personaMessage.id, sender: 'persona', text: out.personaMessage.text }
      setMessages((prev) => [
        ...prev.filter((m) => m.id !== optimisticId),
        userMsg,
        personaMsg,
      ])
      setIsTyping(false)
      if (out.consentRisk?.detected) {
        setConsentRisk(out.consentRisk)
      }
      setShouldSuggestEnd(out.shouldSuggestEnd)
    } catch {
      setIsTyping(false)
      setMessages((prev) => prev.filter((m) => m.id !== optimisticId))
      setSendError('Не удалось отправить сообщение. Попробуй ещё раз.')
    } finally {
      setIsSending(false)
    }
  }

  const handleSend = async () => {
    const text = inputValue.trim()
    if (!text || isSending) return
    if (consentRisk?.detected) {
      setPendingText(text)
      return
    }
    setInputValue('')
    await doSend(text)
  }

  const handleReplaceMessage = () => {
    if (consentRisk?.suggestion) setInputValue(consentRisk.suggestion)
    setConsentRisk(null)
    setPendingText(null)
  }

  const handleKeepMessage = async () => {
    const text = pendingText
    setConsentRisk(null)
    setPendingText(null)
    if (text) { setInputValue(''); await doSend(text) }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) { e.preventDefault(); handleSend() }
  }

  const handleFinish = async () => {
    const sessionId = getSessionId()
    if (!sessionId) return
    setFinishError(null)
    setIsFinishing(true)
    try {
      const debrief = await api.finishSession(sessionId)
      cacheDebrief(debrief)
      window.location.href = '/debrief'
    } catch (err) {
      setFinishError(err instanceof Error ? err.message : 'Не удалось завершить сессию')
      setIsFinishing(false)
    }
  }

  return (
    <div className="min-h-screen bg-background flex flex-col">
      <header className="sticky top-0 z-10 bg-background border-b border-border">
        <div className="max-w-[720px] mx-auto px-4 py-3 flex items-center justify-between gap-4">
          <div className="flex items-center gap-3 min-w-0">
            <div className="size-9 rounded-xl bg-gradient-to-br from-accent/30 to-warning/40 flex items-center justify-center shrink-0">
              <span className="text-small font-semibold text-foreground">{personaInitials}</span>
            </div>
            <div className="min-w-0">
              <h2 className="text-body font-medium text-foreground truncate">{personaTitle}</h2>
              <p className="text-small text-foreground-muted truncate">{goalTitle}</p>
            </div>
          </div>
          <div className="hidden sm:flex items-center gap-2 text-small text-foreground-secondary">
            <span>Реплика {userMessageCount} из ~{totalSteps}</span>
          </div>
          <Button variant="ghost" size="sm" className="shrink-0" onClick={handleFinish} disabled={isFinishing}>
            <span className="hidden sm:inline">{isFinishing ? 'Завершаем...' : 'Завершить и получить разбор'}</span>
            <span className="sm:hidden">{isFinishing ? '...' : 'Разбор'}</span>
          </Button>
        </div>
        {finishError && (
          <div className="max-w-[720px] mx-auto px-4 pb-2">
            <div className="rounded-xl border border-danger/30 bg-danger/5 px-4 py-2 text-small text-foreground-secondary">
              {finishError}
            </div>
          </div>
        )}
        {shouldSuggestEnd && (
          <div className="max-w-[720px] mx-auto px-4 pb-2">
            <div className="flex items-center justify-between gap-3 rounded-xl border border-accent/20 bg-accent/5 px-4 py-2">
              <p className="text-small text-foreground-secondary">Хороший момент для завершения — диалог достиг своей цели.</p>
              <Button size="sm" onClick={handleFinish} disabled={isFinishing}>Завершить</Button>
            </div>
          </div>
        )}
      </header>
      <main className="flex-1 overflow-y-auto">
        <div className="max-w-[720px] mx-auto px-4 py-6 relative">
          <div className="absolute top-2 right-4">
            <LegendPopover>
              <button className="size-8 rounded-full bg-surface-elevated border border-border flex items-center justify-center text-foreground-secondary hover:text-foreground transition-colors">
                <Info className="size-4" />
              </button>
            </LegendPopover>
          </div>
          <div className="space-y-4 pt-8">
            {messages.map((message) => (
              <div key={message.id} className={cn('flex items-end gap-2', message.sender === 'user' && 'justify-end')}>
                {message.sender === 'user' && message.status && (
                  <Tooltip>
                    <TooltipTrigger asChild>
                      <button className="mb-2 cursor-help"><StatusDot status={message.status} /></button>
                    </TooltipTrigger>
                    <TooltipContent side="left" className="max-w-[200px]">{message.statusLabel}</TooltipContent>
                  </Tooltip>
                )}
                <ChatBubble variant={message.sender}>
                  <p className="whitespace-pre-wrap break-words">{message.text}</p>
                </ChatBubble>
              </div>
            ))}
            {isTyping && (
              <div className="flex items-end gap-2">
                <ChatBubble variant="persona"><TypingIndicator /></ChatBubble>
              </div>
            )}
            <div ref={messagesEndRef} />
          </div>
        </div>
      </main>
      <footer className="sticky bottom-0 bg-background border-t border-border">
        <div className="max-w-[720px] mx-auto px-4 py-4">
          {consentRisk?.detected && (
            <ConsentGuardBanner risk={consentRisk} onReplace={handleReplaceMessage} onKeep={handleKeepMessage} />
          )}
          {sendError && (
            <p className="text-small text-danger mb-3">{sendError}</p>
          )}
          <div className="flex items-end gap-2">
            <SuggestionPopover draft={inputValue} onInsert={setInputValue}>
              <Button variant="ghost" size="icon" className="shrink-0 mb-0.5" type="button">
                <Wand2 className="size-5" />
                <span className="sr-only">Подсказать формулировку</span>
              </Button>
            </SuggestionPopover>
            <div className="flex-1 relative">
              <Textarea
                value={inputValue}
                onChange={(e) => setInputValue(e.target.value)}
                onKeyDown={handleKeyDown}
                placeholder="Напишите ей..."
                className="min-h-[48px] max-h-[160px] py-3 pr-12 resize-none"
                rows={1}
              />
            </div>
            <Button size="icon" className="shrink-0 mb-0.5" onClick={handleSend} disabled={!inputValue.trim() || isSending}>
              {!isSending && <Send className="size-5" />}
              {isSending && <span className="size-4 border-2 border-current border-t-transparent rounded-full animate-spin" />}
              <span className="sr-only">Отправить</span>
            </Button>
          </div>
          <p className="text-[12px] text-foreground-muted mt-2 text-center">Enter — отправить, Shift+Enter — новая строка</p>
        </div>
      </footer>
    </div>
  )
}
