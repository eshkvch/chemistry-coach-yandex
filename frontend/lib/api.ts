/**
 * API client for Alice Chemistry Coach backend.
 * Base URL is resolved from NEXT_PUBLIC_API_URL env var (default: http://localhost:8080).
 */

const BASE_URL = process.env.NEXT_PUBLIC_API_URL ?? 'http://localhost:8080'
const API = `${BASE_URL}/api/v1`

// ---------------------------------------------------------------------------
// Storage helpers (localStorage)
// ---------------------------------------------------------------------------

export function getUserId(): string | null {
  if (typeof window === 'undefined') return null
  return localStorage.getItem('acc_user_id')
}

export function setUserId(id: string): void {
  if (typeof window === 'undefined') return
  localStorage.setItem('acc_user_id', id)
}

export function getSessionId(): string | null {
  if (typeof window === 'undefined') return null
  return sessionStorage.getItem('acc_session_id')
}

export function setSessionId(id: string): void {
  if (typeof window === 'undefined') return
  sessionStorage.setItem('acc_session_id', id)
}

export function clearSessionId(): void {
  if (typeof window === 'undefined') return
  sessionStorage.removeItem('acc_session_id')
}

// ---------------------------------------------------------------------------
// Low-level fetch helper
// ---------------------------------------------------------------------------

async function request<T>(
  path: string,
  options: RequestInit = {},
  userId?: string | null,
): Promise<T> {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
    ...(options.headers as Record<string, string>),
  }
  const uid = userId ?? getUserId()
  if (uid) headers['X-User-Id'] = uid

  const res = await fetch(`${API}${path}`, { ...options, headers })
  if (!res.ok) {
    let message = `HTTP ${res.status}`
    try {
      const body = await res.json()
      message = body.error ?? message
    } catch {}
    throw new Error(message)
  }
  return res.json() as Promise<T>
}

// ---------------------------------------------------------------------------
// Types (mirror backend DTOs)
// ---------------------------------------------------------------------------

export interface AuthStartInput {
  birthYear: number
  consents: {
    isAdult: boolean
    understandsAI: boolean
    agreesToRules: boolean
  }
}

export interface AuthStartOutput {
  userId: string
  isNewUser: boolean
}

export interface Goal {
  id: string
  title: string
  description: string
  difficulty: number
  recommended: boolean
}

export interface Persona {
  id: string
  initials: string
  title: string
  description: string
  tags: string[]
  difficulty: number
  showDemoBadge?: boolean
}

export interface CreateSessionInput {
  goalId: string
  personaId: string
}

export interface CreateSessionOutput {
  sessionId: string
  openingMessage: string
  personaTitle: string
  personaInitials: string
  goalTitle: string
}

export interface MessageDTO {
  id: string
  sender: 'user' | 'persona'
  text: string
  status?: 'success' | 'warning' | 'danger'
  statusLabel?: string
}

export interface ConsentRiskDTO {
  detected: boolean
  severity?: string
  explanation?: string
  suggestion?: string
}

export interface SendMessageOutput {
  userMessage: MessageDTO
  personaMessage: MessageDTO
  consentRisk: ConsentRiskDTO
  messageCount: number
  shouldSuggestEnd: boolean
}

export interface ScoreDetail {
  value: number
  comment: string
}

export interface DebriefScores {
  clarity: ScoreDetail
  confidence: ScoreDetail
  respect: ScoreDetail
  balance: ScoreDetail
}

export interface RiskFlag {
  severity: string
  quote: string
  explanation: string
  suggestion: string
}

export interface ImprovedReply {
  original: string
  improved: string
  reason: string
}

export interface FinishOutput {
  sessionId: string
  finishedAt: string
  scenario: string
  persona: string
  scores: DebriefScores
  strengths: string[]
  weaknesses: string[]
  riskFlags: RiskFlag[]
  improvedReplies: ImprovedReply[]
  tipForNext: string
  hasRisk: boolean
}

/** Ensures list fields are arrays (API may return null for empty slices). */
export function normalizeDebrief(data: FinishOutput): FinishOutput {
  return {
    ...data,
    strengths: data.strengths ?? [],
    weaknesses: data.weaknesses ?? [],
    riskFlags: data.riskFlags ?? [],
    improvedReplies: data.improvedReplies ?? [],
    tipForNext: data.tipForNext ?? '',
  }
}

const DEBRIEF_CACHE_KEY = 'acc_debrief'

export function cacheDebrief(data: FinishOutput): void {
  if (typeof window === 'undefined') return
  sessionStorage.setItem(DEBRIEF_CACHE_KEY, JSON.stringify(data))
}

export function takeCachedDebrief(): FinishOutput | null {
  if (typeof window === 'undefined') return null
  const raw = sessionStorage.getItem(DEBRIEF_CACHE_KEY)
  if (!raw) return null
  sessionStorage.removeItem(DEBRIEF_CACHE_KEY)
  try {
    return normalizeDebrief(JSON.parse(raw) as FinishOutput)
  } catch {
    return null
  }
}

export interface SkillScores {
  clarity: number
  confidence: number
  respect: number
  balance: number
}

export interface SessionListItem {
  id: string
  date: string
  time: string
  scenario: string
  persona: string
  scores: SkillScores
  hasRisk: boolean
}

export interface DailyExercise {
  title: string
  description: string
  criterion: string
}

export interface SessionsListOutput {
  focusSkill: string
  focusSkillLabel: string
  focusSessionCount: number
  clarityHistory: number[]
  dailyExercise: DailyExercise | null
  sessions: SessionListItem[]
}

export interface ProfileOutput {
  userId: string
  weakestSkill: string
  recommendedGoalId: string
  totalSessions: number
  initials: string
}

export interface SuggestOutput {
  suggestion: string
}

// ---------------------------------------------------------------------------
// API methods
// ---------------------------------------------------------------------------

export const api = {
  /** Onboarding — create or restore anonymous user */
  authStart(input: AuthStartInput): Promise<AuthStartOutput> {
    return request('/auth/start', {
      method: 'POST',
      body: JSON.stringify(input),
    })
  },

  /** Get current user profile */
  getProfile(): Promise<ProfileOutput> {
    return request('/profile')
  },

  /** List available goals */
  async getGoals(recommendedGoalId?: string): Promise<Goal[]> {
    const qs = recommendedGoalId ? `?recommendedGoalId=${encodeURIComponent(recommendedGoalId)}` : ''
    const data = await request<{ goals: Goal[] }>(`/goals${qs}`)
    return data.goals
  },

  /** List available personas */
  async getPersonas(): Promise<Persona[]> {
    const data = await request<{ personas: Persona[] }>('/personas')
    return data.personas
  },

  /** Create a new session */
  createSession(input: CreateSessionInput): Promise<CreateSessionOutput> {
    return request('/sessions', {
      method: 'POST',
      body: JSON.stringify(input),
    })
  },

  /** Send a message in an active session */
  sendMessage(sessionId: string, text: string): Promise<SendMessageOutput> {
    return request(`/sessions/${sessionId}/messages`, {
      method: 'POST',
      body: JSON.stringify({ text }),
    })
  },

  /** Get a suggestion for the current draft */
  suggest(sessionId: string, draft: string): Promise<SuggestOutput> {
    return request(`/sessions/${sessionId}/suggest`, {
      method: 'POST',
      body: JSON.stringify({ draft }),
    })
  },

  /** Finish session and get debrief */
  finishSession(sessionId: string): Promise<FinishOutput> {
    return request(`/sessions/${sessionId}/finish`, {
      method: 'POST',
    })
  },

  /** Get debrief for a finished session */
  getDebrief(sessionId: string): Promise<FinishOutput> {
    return request(`/sessions/${sessionId}`)
  },

  /** List all sessions for history page */
  listSessions(): Promise<SessionsListOutput> {
    return request('/sessions')
  },
}
