'use client'
export const dynamic = 'force-dynamic'

import { useState, useMemo, useEffect } from 'react'
import { api, setUserId, getUserId } from '@/lib/api'
import { Button } from '@/components/ui/button'
import { Card, CardContent } from '@/components/ui/card'
import { Checkbox } from '@/components/ui/checkbox'
import { Label } from '@/components/ui/label'
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from '@/components/ui/select'

export default function OnboardingPage() {
  const [selectedYear, setSelectedYear] = useState<string>('')
  const [isAdult, setIsAdult] = useState(false)
  const [understandsAI, setUnderstandsAI] = useState(false)
  const [agreesToRules, setAgreesToRules] = useState(false)
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  // Redirect returning users to goal selection
  useEffect(() => {
    if (getUserId()) {
      window.location.href = '/goal'
    }
  }, [])

  const currentYear = new Date().getFullYear()
  const years = useMemo(() => {
    const result: number[] = []
    for (let year = currentYear; year >= currentYear - 100; year--) {
      result.push(year)
    }
    return result
  }, [currentYear])

  const isUnderage = useMemo(() => {
    if (!selectedYear) return false
    const age = currentYear - parseInt(selectedYear)
    return age < 18
  }, [selectedYear, currentYear])

  const isFormValid = selectedYear && !isUnderage && isAdult && understandsAI && agreesToRules

  const handleStart = async () => {
    if (!isFormValid) return
    setIsLoading(true)
    setError(null)
    try {
      const out = await api.authStart({
        birthYear: parseInt(selectedYear),
        consents: { isAdult, understandsAI, agreesToRules },
      })
      setUserId(out.userId)
      window.location.href = '/goal'
    } catch (err: unknown) {
      const msg = err instanceof Error ? err.message : 'Ошибка сервера'
      if (msg === 'underage') {
        setError('Сервис доступен только пользователям 18+')
      } else if (msg === 'consents required') {
        setError('Необходимо подтвердить все условия')
      } else {
        setError('Не удалось подключиться к серверу. Проверьте соединение.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-background flex flex-col animate-in fade-in duration-500">
      {/* Top bar */}
      <header className="flex items-center justify-between px-4 py-3 md:px-6">
        <div className="text-h3 font-bold text-foreground">ACC</div>
      </header>

      {/* Main content */}
      <main className="flex-1 flex items-center justify-center px-4 py-8">
        <Card className="w-full max-w-[480px]">
          <CardContent className="p-6 md:p-8">
            {/* Decorative amber gradient blob */}
            <div className="flex justify-center mb-6">
              <div 
                className="w-24 h-24 rounded-full opacity-80"
                style={{
                  background: 'radial-gradient(circle at 30% 30%, #FFB733 0%, #F5A524 40%, rgba(245, 165, 36, 0.3) 70%, transparent 100%)',
                  filter: 'blur(8px)',
                }}
              />
            </div>

            {/* Title and subtitle */}
            <div className="text-center mb-8">
              <h1 className="text-h1 text-foreground mb-3 text-balance">
                Тренажёр романтической коммуникации
              </h1>
              <p className="text-body text-foreground-secondary text-balance">
                Учим общаться ясно, уверенно и уважительно. Без манипуляций и розовых сердечек.
              </p>
            </div>

            {/* Year of birth selector */}
            <div className="mb-6">
              <Label htmlFor="year-select" className="text-small text-foreground-secondary mb-2 block">
                Год рождения
              </Label>
              <Select value={selectedYear} onValueChange={setSelectedYear}>
                <SelectTrigger id="year-select" className="w-full">
                  <SelectValue placeholder="Выберите год" />
                </SelectTrigger>
                <SelectContent>
                  {years.map((year) => (
                    <SelectItem key={year} value={year.toString()}>
                      {year}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
              {isUnderage && (
                <p className="text-small text-danger mt-2 animate-in fade-in slide-in-from-top-1 duration-200">
                  Сервис доступен только пользователям 18+
                </p>
              )}
            </div>

            {/* Checkboxes */}
            <div className="space-y-4 mb-8">
              <div className="flex items-start gap-3">
                <Checkbox
                  id="adult"
                  checked={isAdult}
                  onCheckedChange={(checked) => setIsAdult(checked === true)}
                  className="mt-0.5"
                />
                <Label 
                  htmlFor="adult" 
                  className="text-body text-foreground cursor-pointer leading-normal font-normal"
                >
                  Мне 18 лет или больше
                </Label>
              </div>

              <div className="flex items-start gap-3">
                <Checkbox
                  id="understands-ai"
                  checked={understandsAI}
                  onCheckedChange={(checked) => setUnderstandsAI(checked === true)}
                  className="mt-0.5"
                />
                <Label 
                  htmlFor="understands-ai" 
                  className="text-body text-foreground cursor-pointer leading-normal font-normal"
                >
                  Я понимаю, что AI — симулятор, а не реальный человек
                </Label>
              </div>

              <div className="flex items-start gap-3">
                <Checkbox
                  id="agrees-rules"
                  checked={agreesToRules}
                  onCheckedChange={(checked) => setAgreesToRules(checked === true)}
                  className="mt-0.5"
                />
                <Label 
                  htmlFor="agrees-rules" 
                  className="text-body text-foreground cursor-pointer leading-normal font-normal"
                >
                  Я согласен с правилами: уважение, согласие, никаких манипуляций
                </Label>
              </div>
            </div>

            {/* Error */}
            {error && (
              <p className="text-small text-danger mb-4 text-center animate-in fade-in duration-200">
                {error}
              </p>
            )}

            {/* Submit button */}
            <Button
              className="w-full"
              size="lg"
              disabled={!isFormValid || isLoading}
              onClick={handleStart}
            >
              {isLoading ? 'Подключаемся...' : 'Начать'}
            </Button>

            {/* Footer links */}
            <div className="flex items-center justify-center gap-1.5 mt-6 text-small">
              <span className="text-foreground-secondary">Политика конфиденциальности</span>
              <span className="text-foreground-muted">·</span>
              <span className="text-foreground-secondary">Условия использования</span>
            </div>
          </CardContent>
        </Card>
      </main>
    </div>
  )
}
