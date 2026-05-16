'use client'
export const dynamic = 'force-dynamic'
import { useState, useEffect } from 'react'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProgressLine } from '@/components/ui/progress-line'
import { cn } from '@/lib/utils'
import { api, getUserId, setSessionId, type Persona } from '@/lib/api'

const GRADIENT_MAP: Record<string, { from: string; to: string }> = {
  'calm-careful': { from: 'from-success/30', to: 'to-success/50' },
  'ironic-fast':  { from: 'from-accent/30',  to: 'to-warning/40' },
  'busy':         { from: 'from-foreground-muted/30', to: 'to-foreground-secondary/30' },
}

function PersonaAvatar({ initials, gradientFrom, gradientTo, isSelected }: {
  initials: string; gradientFrom: string; gradientTo: string; isSelected: boolean
}) {
  return (
    <div className={cn(
      'size-14 rounded-2xl flex items-center justify-center mb-4 bg-gradient-to-br transition-all duration-200',
      gradientFrom, gradientTo,
      isSelected && 'ring-2 ring-accent ring-offset-2 ring-offset-background',
    )}>
      <span className="text-[18px] font-semibold text-foreground">{initials}</span>
    </div>
  )
}

const stepLabels = ['Цель', 'Персона', 'Диалог', 'Разбор']

export default function PersonaPage() {
  const [personas, setPersonas] = useState<Persona[]>([])
  const [selectedPersona, setSelectedPersona] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)
  const [isStarting, setIsStarting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!getUserId()) { window.location.href = '/'; return }
    if (!sessionStorage.getItem('acc_goal_id')) { window.location.href = '/goal'; return }
    api.getPersonas()
      .then(setPersonas)
      .catch(() => {})
      .finally(() => setIsLoading(false))
  }, [])

  const handleStart = async () => {
    const goalId = sessionStorage.getItem('acc_goal_id')
    if (!selectedPersona || !goalId) return
    setIsStarting(true)
    setError(null)
    try {
      const out = await api.createSession({ goalId, personaId: selectedPersona })
      setSessionId(out.sessionId)
      sessionStorage.setItem('acc_persona_title', out.personaTitle)
      sessionStorage.setItem('acc_persona_initials', out.personaInitials)
      sessionStorage.setItem('acc_goal_title', out.goalTitle)
      sessionStorage.setItem('acc_opening_message', out.openingMessage)
      window.location.href = '/chat'
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : 'Ошибка при создании сессии')
    } finally {
      setIsStarting(false)
    }
  }

  return (
    <div className="min-h-screen bg-background flex flex-col animate-in fade-in duration-500">
      <header className="flex items-center justify-between px-4 py-3 md:px-6 border-b border-border">
        <div className="text-h3 font-bold text-foreground">ACC</div>
      </header>
      <div className="px-4 py-4 md:px-6 md:py-5 max-w-3xl mx-auto w-full">
        <ProgressLine steps={4} currentStep={2} labels={stepLabels} />
      </div>
      <main className="flex-1 px-4 pb-24 md:px-6 max-w-3xl mx-auto w-full">
        <div className="mb-6">
          <h1 className="text-h1 text-foreground mb-2">С кем тренируемся?</h1>
          <p className="text-body text-foreground-secondary">Выберите персонажа — каждый по-своему реагирует на ваши реплики.</p>
        </div>
        {isLoading && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {[1,2,3].map((i) => <div key={i} className="h-52 rounded-xl bg-surface-elevated animate-pulse" />)}
          </div>
        )}
        {!isLoading && (
          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {personas.map((p) => {
              const grad = GRADIENT_MAP[p.id] ?? { from: 'from-accent/20', to: 'to-accent/30' }
              const isSelected = selectedPersona === p.id
              return (
                <Card
                  key={p.id}
                  role="button"
                  tabIndex={0}
                  onClick={() => setSelectedPersona(p.id)}
                  onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); setSelectedPersona(p.id) } }}
                  className={cn(
                    'relative cursor-pointer p-5 transition-all duration-200 outline-none hover:border-accent/50 focus-visible:ring-2 focus-visible:ring-accent/50',
                    isSelected && 'border-accent shadow-[0_0_20px_rgba(245,165,36,0.15)]',
                  )}
                >
                  {p.showDemoBadge && (
                    <Badge variant="accent" className="absolute top-4 right-4 text-[11px] px-2 py-0.5">Демо</Badge>
                  )}
                  <PersonaAvatar initials={p.initials} gradientFrom={grad.from} gradientTo={grad.to} isSelected={isSelected} />
                  <h3 className="text-h3 text-foreground mb-1.5">{p.title}</h3>
                  <p className="text-small text-foreground-secondary mb-3">{p.description}</p>
                  <div className="flex flex-wrap gap-1.5">
                    {p.tags.map((tag) => (
                      <span key={tag} className="text-[11px] px-2 py-0.5 rounded-full bg-surface-elevated text-foreground-secondary">{tag}</span>
                    ))}
                  </div>
                </Card>
              )
            })}
          </div>
        )}
        {error && <p className="text-small text-danger mt-4 text-center">{error}</p>}
      </main>
      <div className="fixed bottom-0 left-0 right-0 border-t border-border bg-background/80 backdrop-blur-sm">
        <div className="flex items-center justify-between gap-4 px-4 py-4 md:px-6 max-w-3xl mx-auto">
          <Button variant="secondary" size="lg" asChild><Link href="/goal">Назад</Link></Button>
          <Button size="lg" disabled={!selectedPersona || isStarting} onClick={handleStart}>
            {isStarting ? 'Запускаем...' : 'Начать диалог'}
          </Button>
        </div>
      </div>
    </div>
  )
}
