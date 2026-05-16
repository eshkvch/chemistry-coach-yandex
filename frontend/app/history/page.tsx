'use client'
export const dynamic = 'force-dynamic'
import { useState, useEffect } from 'react'
import Link from 'next/link'
import { ShieldAlert, ChevronRight } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardFooter, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ScoreBar } from '@/components/ui/score-bar'
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip'
import { cn } from '@/lib/utils'
import { api, getUserId, setSessionId, type SessionsListOutput, type SessionListItem } from '@/lib/api'

function Sparkline({ data, className }: { data: number[]; className?: string }) {
  if (data.length < 2) return null
  const max = Math.max(...data)
  const min = Math.min(...data)
  const range = max - min || 1
  const width = 120; const height = 32; const padding = 4
  const points = data.map((value, index) => {
    const x = padding + (index / (data.length - 1)) * (width - padding * 2)
    const y = height - padding - ((value - min) / range) * (height - padding * 2)
    return `${x},${y}`
  }).join(' ')
  const last = data[data.length - 1]
  const lastX = padding + ((data.length - 1) / (data.length - 1)) * (width - padding * 2)
  const lastY = height - padding - ((last - min) / range) * (height - padding * 2)
  return (
    <svg viewBox={`0 0 ${width} ${height}`} className={cn('text-accent', className)} aria-hidden="true">
      <polyline points={points} fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" />
      <circle cx={lastX} cy={lastY} r="3" fill="currentColor" />
    </svg>
  )
}

const SCORE_LABELS: Record<string, string> = {
  clarity: 'Ясность', confidence: 'Уверенность', respect: 'Уважительность', balance: 'Баланс',
}

export default function HistoryPage() {
  const [data, setData] = useState<SessionsListOutput | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    if (!getUserId()) { window.location.href = '/'; return }
    api.listSessions()
      .then(setData)
      .catch(() => {})
      .finally(() => setIsLoading(false))
  }, [])

  if (isLoading) return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="size-8 border-2 border-accent border-t-transparent rounded-full animate-spin" />
    </div>
  )

  const sessions = data?.sessions ?? []

  return (
    <div className="min-h-screen bg-background flex flex-col animate-in fade-in duration-500">
      <header className="flex items-center justify-between px-4 py-3 md:px-6 border-b border-border">
        <div className="text-h3 font-bold text-foreground">ACC</div>
        <Button size="sm" asChild><Link href="/goal">Новая сессия</Link></Button>
      </header>
      <main className="flex-1 px-4 py-6 md:px-6 max-w-3xl mx-auto w-full space-y-6">
        <h1 className="text-h1 text-foreground">История</h1>

        {/* Focus skill */}
        {data && data.focusSkillLabel && (
          <Card className="border-accent/30 bg-accent/5">
            <CardContent className="pt-4 flex items-center justify-between gap-4">
              <div>
                <p className="text-small font-medium text-foreground mb-0.5">Фокус: {data.focusSkillLabel}</p>
                <p className="text-small text-foreground-secondary">
                  {data.focusSessionCount > 0
                    ? `${data.focusSessionCount} сессий с низкой оценкой`
                    : 'Всё идёт хорошо'}
                </p>
              </div>
              {data.clarityHistory.length > 1 && (
                <Sparkline data={data.clarityHistory} className="w-[120px] h-8 shrink-0" />
              )}
            </CardContent>
          </Card>
        )}

        {/* Daily exercise */}
        {data?.dailyExercise && (
          <Card>
            <CardHeader><CardTitle>Упражнение дня</CardTitle></CardHeader>
            <CardContent>
              <p className="text-body font-medium text-foreground mb-1">{data.dailyExercise.title}</p>
              <p className="text-small text-foreground-secondary mb-2">{data.dailyExercise.description}</p>
              <Badge variant="default" className="text-[11px]">Критерий: {data.dailyExercise.criterion}</Badge>
            </CardContent>
          </Card>
        )}

        {/* Sessions list */}
        {sessions.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16 gap-4">
            <p className="text-body text-foreground-secondary text-center">Пока нет завершённых сессий.</p>
            <Button asChild><Link href="/goal">Начать первую</Link></Button>
          </div>
        ) : (
          <div className="space-y-3">
            {sessions.map((session) => (
              <Card key={session.id} className="hover:border-accent/40 transition-colors">
                <CardHeader className="pb-2">
                  <div className="flex items-start justify-between gap-2">
                    <div>
                      <CardTitle className="text-body">{session.scenario}</CardTitle>
                      <p className="text-small text-foreground-muted mt-0.5">{session.persona} · {session.date} {session.time}</p>
                    </div>
                    <div className="flex items-center gap-2 shrink-0">
                      {session.hasRisk && (
                        <Tooltip>
                          <TooltipTrigger asChild>
                            <button aria-label="Риск давления"><ShieldAlert className="size-4 text-danger" /></button>
                          </TooltipTrigger>
                          <TooltipContent>Зафиксирован риск давления</TooltipContent>
                        </Tooltip>
                      )}
                    </div>
                  </div>
                </CardHeader>
                <CardContent className="pb-2">
                  <div className="grid grid-cols-2 gap-x-4 gap-y-2">
                    {Object.entries(session.scores).map(([key, value]) => (
                      <div key={key}>
                        <div className="flex items-center justify-between mb-1">
                          <span className="text-[11px] text-foreground-muted">{SCORE_LABELS[key] ?? key}</span>
                          <span className="text-[11px] font-medium text-foreground">{value}</span>
                        </div>
                        <ScoreBar value={value} label="" showValue={false} size="compact" />
                      </div>
                    ))}
                  </div>
                </CardContent>
                <CardFooter className="pt-0">
                  <Button variant="ghost" size="sm" className="ml-auto gap-1" asChild>
                    <Link href={`/debrief?session=${session.id}`}>
                      Смотреть разбор <ChevronRight className="size-3.5" />
                    </Link>
                  </Button>
                </CardFooter>
              </Card>
            ))}
          </div>
        )}
      </main>
    </div>
  )
}
