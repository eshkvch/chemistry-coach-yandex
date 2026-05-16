'use client'
export const dynamic = 'force-dynamic'
import { useState, useEffect } from 'react'
import Link from 'next/link'
import { MessageCircle, Heart, Phone, Shield, X } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ProgressLine } from '@/components/ui/progress-line'
import { cn } from '@/lib/utils'
import { api, getUserId, type Goal } from '@/lib/api'

const GOAL_ICONS: Record<string, React.ElementType> = {
  'first-message': MessageCircle,
  'keep-interest': Heart,
  'suggest-call': Phone,
  'boundaries': Shield,
}

function DifficultyDots({ level }: { level: number }) {
  return (
    <div className="flex items-center gap-1.5">
      {[1, 2, 3].map((dot) => (
        <div
          key={dot}
          className={cn(
            'size-1.5 rounded-full transition-colors',
            dot <= level ? 'bg-accent' : 'bg-surface-elevated',
          )}
        />
      ))}
    </div>
  )
}

const stepLabels = ['Цель', 'Персона', 'Диалог', 'Разбор']

export default function GoalSelectionPage() {
  const [goals, setGoals] = useState<Goal[]>([])
  const [selectedGoal, setSelectedGoal] = useState<string | null>(null)
  const [showHint, setShowHint] = useState(false)
  const [hintText, setHintText] = useState('')
  const [isLoading, setIsLoading] = useState(true)
  const [userInitials, setUserInitials] = useState('ИВ')

  useEffect(() => {
    if (!getUserId()) {
      window.location.href = '/'
      return
    }
    const load = async () => {
      try {
        const profile = await api.getProfile().catch(() => null)
        const goalsData = await api.getGoals(profile?.recommendedGoalId)
        setGoals(goalsData)
        if (profile) {
          setUserInitials(profile.initials ?? 'ИВ')
          if (profile.totalSessions > 0 && profile.weakestSkill && profile.recommendedGoalId) {
            const skillLabels: Record<string, string> = {
              clarity: 'ясность', confidence: 'уверенность',
              respect: 'уважительность', balance: 'баланс инициативы',
            }
            const rec = goalsData.find((g) => g.id === profile.recommendedGoalId)
            if (rec) {
              setHintText(`Ваша слабая зона: ${skillLabels[profile.weakestSkill] ?? profile.weakestSkill}. Рекомендуем сценарий «${rec.title}».`)
              setShowHint(true)
              setSelectedGoal(profile.recommendedGoalId)
            }
          }
        }
      } catch { /* silent */ } finally {
        setIsLoading(false)
      }
    }
    load()
  }, [])

  const handleContinue = () => {
    if (!selectedGoal) return
    sessionStorage.setItem('acc_goal_id', selectedGoal)
    window.location.href = '/persona'
  }

  return (
    <div className="min-h-screen bg-background flex flex-col animate-in fade-in duration-500">
      <header className="flex items-center justify-between px-4 py-3 md:px-6 border-b border-border">
        <div className="text-h3 font-bold text-foreground">ACC</div>
        <div className="flex items-center justify-center size-9 rounded-full bg-surface-elevated text-small font-medium text-foreground">
          {userInitials}
        </div>
      </header>
      <div className="px-4 py-4 md:px-6 md:py-5 max-w-3xl mx-auto w-full">
        <ProgressLine steps={4} currentStep={1} labels={stepLabels} />
      </div>
      <main className="flex-1 px-4 pb-24 md:px-6 max-w-3xl mx-auto w-full">
        <div className="mb-6">
          <h1 className="text-h1 text-foreground mb-2 text-balance">С чем сегодня тренируемся?</h1>
          <p className="text-body text-foreground-secondary text-balance">
            Выберите ситуацию — мы подстроим разговор и разбор под неё.
          </p>
        </div>
        {showHint && hintText && (
          <div className="mb-6 animate-in fade-in slide-in-from-top-2 duration-300">
            <div className="flex items-center justify-between gap-3 rounded-xl border border-accent/20 bg-accent/5 px-4 py-3">
              <p className="text-small text-foreground-secondary">{hintText}</p>
              <button onClick={() => setShowHint(false)} className="shrink-0 text-foreground-muted hover:text-foreground transition-colors" aria-label="Закрыть">
                <X className="size-4" />
              </button>
            </div>
          </div>
        )}
        {isLoading && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {[1,2,3,4].map((i) => <div key={i} className="h-40 rounded-xl bg-surface-elevated animate-pulse" />)}
          </div>
        )}
        {!isLoading && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {goals.map((goal) => {
              const Icon = GOAL_ICONS[goal.id] ?? MessageCircle
              const isSelected = selectedGoal === goal.id
              return (
                <Card
                  key={goal.id}
                  role="button"
                  tabIndex={0}
                  onClick={() => setSelectedGoal(goal.id)}
                  onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); setSelectedGoal(goal.id) } }}
                  className={cn(
                    'relative cursor-pointer p-5 transition-all duration-200 outline-none hover:border-accent/50 focus-visible:ring-2 focus-visible:ring-accent/50',
                    isSelected && 'border-accent shadow-[0_0_20px_rgba(245,165,36,0.15)]',
                  )}
                >
                  {goal.recommended && (
                    <Badge variant="accent" className="absolute top-4 right-4 text-[11px] px-2 py-0.5">Рекомендуем</Badge>
                  )}
                  <div className={cn('flex items-center justify-center size-10 rounded-xl mb-4 transition-colors', isSelected ? 'bg-accent/15' : 'bg-surface-elevated')}>
                    <Icon className={cn('size-5 transition-colors', isSelected ? 'text-accent' : 'text-foreground-secondary')} />
                  </div>
                  <h3 className="text-h3 text-foreground mb-1.5">{goal.title}</h3>
                  <p className="text-body text-foreground-secondary mb-4">{goal.description}</p>
                  <div className="flex items-center gap-2">
                    <DifficultyDots level={goal.difficulty} />
                    <span className="text-small text-foreground-muted">Сложность</span>
                  </div>
                </Card>
              )
            })}
          </div>
        )}
      </main>
      <div className="fixed bottom-0 left-0 right-0 border-t border-border bg-background/80 backdrop-blur-sm">
        <div className="flex items-center justify-between gap-4 px-4 py-4 md:px-6 max-w-3xl mx-auto">
          <Button variant="secondary" size="lg" asChild><Link href="/">Назад</Link></Button>
          <Button size="lg" disabled={!selectedGoal} onClick={handleContinue}>Продолжить</Button>
        </div>
      </div>
    </div>
  )
}
