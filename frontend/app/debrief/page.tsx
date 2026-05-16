'use client'
export const dynamic = 'force-dynamic'
import * as React from 'react'
import Link from 'next/link'
import { Copy, Check, CheckCircle, XCircle, AlertTriangle } from 'lucide-react'
import { Button } from '@/components/ui/button'
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { ScoreBar } from '@/components/ui/score-bar'
import { ProgressLine } from '@/components/ui/progress-line'
import { cn } from '@/lib/utils'
import { api, getSessionId, clearSessionId, takeCachedDebrief, normalizeDebrief, type FinishOutput, type ScoreDetail } from '@/lib/api'
import { useSearchParams } from 'next/navigation'
import { Suspense } from 'react'

const PROGRESS_LABELS = ['Цель', 'Персона', 'Диалог', 'Разбор']

const SCORE_LABELS: Record<string, string> = {
  clarity: 'Ясность',
  confidence: 'Уверенность',
  respect: 'Уважительность',
  balance: 'Баланс инициативы',
}

function CopyButton({ text }: { text: string }) {
  const [copied, setCopied] = React.useState(false)
  const handleCopy = async () => {
    await navigator.clipboard.writeText(text)
    setCopied(true)
    setTimeout(() => setCopied(false), 2000)
  }
  return (
    <button onClick={handleCopy} className="text-foreground-muted hover:text-foreground transition-colors" aria-label="Копировать">
      {copied ? <Check className="size-4 text-success" /> : <Copy className="size-4" />}
    </button>
  )
}

function RiskIcon({ severity }: { severity: string }) {
  if (severity === 'high') return <XCircle className="size-4 text-danger shrink-0 mt-0.5" />
  if (severity === 'medium') return <AlertTriangle className="size-4 text-warning shrink-0 mt-0.5" />
  return <AlertTriangle className="size-4 text-foreground-muted shrink-0 mt-0.5" />
}

function DebriefContent() {
  const [data, setData] = React.useState<FinishOutput | null>(null)
  const [isLoading, setIsLoading] = React.useState(true)
  const [error, setError] = React.useState<string | null>(null)
  const searchParams = useSearchParams()

  React.useEffect(() => {
    const querySessionId = searchParams.get('session')
    const cached = !querySessionId ? takeCachedDebrief() : null
    if (cached) {
      setData(cached)
      clearSessionId()
      setIsLoading(false)
      return
    }
    const sessionId = querySessionId ?? getSessionId()
    if (!sessionId) { window.location.href = '/'; return }
    api.getDebrief(sessionId)
      .then((d) => {
        setData(normalizeDebrief(d))
        if (!querySessionId) clearSessionId()
      })
      .catch((err) => setError(err instanceof Error ? err.message : 'Ошибка загрузки разбора'))
      .finally(() => setIsLoading(false))
  }, [searchParams])

  if (isLoading) return (
    <div className="min-h-screen bg-background flex items-center justify-center">
      <div className="size-8 border-2 border-accent border-t-transparent rounded-full animate-spin" />
    </div>
  )

  if (error || !data) return (
    <div className="min-h-screen bg-background flex flex-col items-center justify-center gap-4 px-4">
      <p className="text-body text-foreground-secondary text-center">{error ?? 'Разбор не найден'}</p>
      <Button asChild><Link href="/">На главную</Link></Button>
    </div>
  )

  const strengths = data.strengths ?? []
  const weaknesses = data.weaknesses ?? []
  const riskFlags = data.riskFlags ?? []
  const improvedReplies = data.improvedReplies ?? []

  const emptyScore: ScoreDetail = { value: 0, comment: '' }
  const scores = [
    { key: 'clarity', ...(data.scores?.clarity ?? emptyScore) },
    { key: 'confidence', ...(data.scores?.confidence ?? emptyScore) },
    { key: 'respect', ...(data.scores?.respect ?? emptyScore) },
    { key: 'balance', ...(data.scores?.balance ?? emptyScore) },
  ]

  const avgScore = Math.round(scores.reduce((s, sc) => s + sc.value, 0) / scores.length)

  return (
    <div className="min-h-screen bg-background flex flex-col animate-in fade-in duration-500">
      <header className="flex items-center justify-between px-4 py-3 md:px-6 border-b border-border">
        <div className="text-h3 font-bold text-foreground">ACC</div>
        <Button variant="ghost" size="sm" asChild><Link href="/history">История</Link></Button>
      </header>
      <div className="px-4 py-4 md:px-6 md:py-5 max-w-3xl mx-auto w-full">
        <ProgressLine steps={4} currentStep={4} labels={PROGRESS_LABELS} />
      </div>
      <main className="flex-1 px-4 pb-16 md:px-6 max-w-3xl mx-auto w-full space-y-6">
        {/* Header */}
        <div className="flex items-start justify-between gap-4">
          <div>
            <h1 className="text-h1 text-foreground mb-1">Разбор сессии</h1>
            <p className="text-body text-foreground-secondary">{data.scenario} · {data.persona}</p>
          </div>
          <div className="flex flex-col items-center shrink-0">
            <span className="text-[40px] font-bold text-foreground leading-none">{avgScore}</span>
            <span className="text-small text-foreground-muted">из 10</span>
          </div>
        </div>

        {/* Risk banner */}
        {data.hasRisk && (
          <div className="flex items-start gap-3 rounded-xl border border-danger/30 bg-danger/5 px-4 py-3">
            <XCircle className="size-5 text-danger shrink-0 mt-0.5" />
            <p className="text-small text-foreground-secondary">В этой сессии были зафиксированы риски давления. Смотри раздел «Флаги» ниже.</p>
          </div>
        )}

        {/* Scores */}
        <Card>
          <CardHeader><CardTitle>Навыки</CardTitle></CardHeader>
          <CardContent className="space-y-4">
            {scores.map((sc) => (
              <div key={sc.key}>
                <div className="flex items-center justify-between mb-1.5">
                  <span className="text-body text-foreground">{SCORE_LABELS[sc.key]}</span>
                  <span className="text-body font-medium text-foreground">{sc.value}/10</span>
                </div>
                <ScoreBar value={sc.value} label="" showValue={false} />
                <p className="text-small text-foreground-secondary mt-1">{sc.comment}</p>
              </div>
            ))}
          </CardContent>
        </Card>

        {/* Strengths & Weaknesses */}
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {strengths.length > 0 && (
            <Card>
              <CardHeader><CardTitle className="text-success">Сильные стороны</CardTitle></CardHeader>
              <CardContent className="space-y-2">
                {strengths.map((s, i) => (
                  <div key={i} className="flex items-start gap-2">
                    <CheckCircle className="size-4 text-success shrink-0 mt-0.5" />
                    <p className="text-small text-foreground-secondary">{s}</p>
                  </div>
                ))}
              </CardContent>
            </Card>
          )}
          {weaknesses.length > 0 && (
            <Card>
              <CardHeader><CardTitle className="text-warning">Зоны роста</CardTitle></CardHeader>
              <CardContent className="space-y-2">
                {weaknesses.map((w, i) => (
                  <div key={i} className="flex items-start gap-2">
                    <XCircle className="size-4 text-warning shrink-0 mt-0.5" />
                    <p className="text-small text-foreground-secondary">{w}</p>
                  </div>
                ))}
              </CardContent>
            </Card>
          )}
        </div>

        {/* Risk flags */}
        {riskFlags.length > 0 && (
          <Card>
            <CardHeader><CardTitle>Флаги</CardTitle></CardHeader>
            <CardContent className="space-y-4">
              {riskFlags.map((rf, i) => (
                <div key={i} className="space-y-1">
                  <div className="flex items-start gap-2">
                    <RiskIcon severity={rf.severity} />
                    <p className="text-small text-foreground italic">"{rf.quote}"</p>
                  </div>
                  <p className="text-small text-foreground-secondary pl-6">{rf.explanation}</p>
                  {rf.suggestion && <p className="text-small text-foreground-muted pl-6">→ {rf.suggestion}</p>}
                </div>
              ))}
            </CardContent>
          </Card>
        )}

        {/* Improved replies */}
        {improvedReplies.length > 0 && (
          <Card>
            <CardHeader><CardTitle>Как можно было сказать лучше</CardTitle></CardHeader>
            <CardContent className="space-y-4">
              {improvedReplies.map((ir, i) => (
                <div key={i} className="space-y-2">
                  <div className="rounded-lg bg-danger/5 border border-danger/20 px-3 py-2">
                    <p className="text-small text-foreground-secondary line-through">{ir.original}</p>
                  </div>
                  <div className="rounded-lg bg-success/5 border border-success/20 px-3 py-2 flex items-start justify-between gap-2">
                    <p className="text-small text-foreground-secondary">{ir.improved}</p>
                    <CopyButton text={ir.improved} />
                  </div>
                  <p className="text-small text-foreground-muted">{ir.reason}</p>
                </div>
              ))}
            </CardContent>
          </Card>
        )}

        {/* Tip for next */}
        {data.tipForNext && (
          <Card className="border-accent/30 bg-accent/5">
            <CardContent className="pt-4">
              <p className="text-small font-medium text-foreground mb-1">Совет на следующую сессию</p>
              <p className="text-small text-foreground-secondary">{data.tipForNext}</p>
            </CardContent>
          </Card>
        )}

        {/* Actions */}
        <div className="flex flex-col sm:flex-row gap-3 pb-8">
          <Button size="lg" className="flex-1" asChild><Link href="/goal">Новая сессия</Link></Button>
          <Button variant="secondary" size="lg" className="flex-1" asChild><Link href="/history">История</Link></Button>
        </div>
      </main>
    </div>
  )
}

export default function DebriefPage() {
  return (
    <Suspense fallback={<div className="min-h-screen bg-background flex items-center justify-center"><div className="size-8 border-2 border-accent border-t-transparent rounded-full animate-spin" /></div>}>
      <DebriefContent />
    </Suspense>
  )
}
