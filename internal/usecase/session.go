package usecase

import (
	"chemistry-coach/internal/catalog"
	"chemistry-coach/internal/domain"
	"context"
	"encoding/json"
	"errors"
	"time"
)

var (
	ErrSessionNotFound = errors.New("session not found")
	ErrSessionFinished = errors.New("session already finished")
	ErrInvalidGoal     = errors.New("invalid goal")
	ErrInvalidPersona  = errors.New("invalid persona")
)

type SessionUseCase struct {
	sessions domain.SessionRepository
	messages domain.MessageRepository
	debriefs domain.DebriefRepository
	llm      domain.LLMService
}

func NewSessionUseCase(
	sessions domain.SessionRepository,
	messages domain.MessageRepository,
	debriefs domain.DebriefRepository,
	llm domain.LLMService,
) *SessionUseCase {
	return &SessionUseCase{sessions: sessions, messages: messages, debriefs: debriefs, llm: llm}
}

type CreateSessionInput struct {
	GoalID    string `json:"goalId"`
	PersonaID string `json:"personaId"`
}

type CreateSessionOutput struct {
	SessionID       string `json:"sessionId"`
	OpeningMessage  string `json:"openingMessage"`
	PersonaTitle    string `json:"personaTitle"`
	PersonaInitials string `json:"personaInitials"`
	GoalTitle       string `json:"goalTitle"`
}

func (uc *SessionUseCase) Create(ctx context.Context, userID string, in CreateSessionInput) (*CreateSessionOutput, error) {
	if _, ok := catalog.GetGoal(in.GoalID); !ok {
		return nil, ErrInvalidGoal
	}
	persona, ok := catalog.GetPersona(in.PersonaID)
	if !ok {
		return nil, ErrInvalidPersona
	}
	systemPrompt := catalog.BuildSystemPrompt(in.GoalID, in.PersonaID)
	opening, err := uc.llm.GenerateOpening(ctx, systemPrompt)
	if err != nil {
		return nil, err
	}
	sid := newSessionID()
	session := &domain.Session{
		ID: sid, UserID: userID, GoalID: in.GoalID, PersonaID: in.PersonaID,
		Status: domain.SessionStatusActive, SystemPrompt: systemPrompt,
	}
	if err := uc.sessions.Create(ctx, session); err != nil {
		return nil, err
	}
	personaMsg := &domain.Message{
		ID: newMessageID(), SessionID: sid, Sender: domain.SenderPersona, Text: opening,
	}
	if err := uc.messages.Create(ctx, personaMsg); err != nil {
		return nil, err
	}
	return &CreateSessionOutput{
		SessionID: sid, OpeningMessage: opening,
		PersonaTitle: persona.Title, PersonaInitials: persona.Initials,
		GoalTitle: goalTitle(in.GoalID),
	}, nil
}

type SendMessageInput struct {
	Text string `json:"text"`
}

type MessageDTO struct {
	ID          string  `json:"id"`
	Sender      string  `json:"sender"`
	Text        string  `json:"text"`
	Status      *string `json:"status,omitempty"`
	StatusLabel *string `json:"statusLabel,omitempty"`
}

type ConsentRiskDTO struct {
	Detected    bool    `json:"detected"`
	Severity    *string `json:"severity"`
	Explanation *string `json:"explanation"`
	Suggestion  *string `json:"suggestion"`
}

type SendMessageOutput struct {
	UserMessage      MessageDTO     `json:"userMessage"`
	PersonaMessage   MessageDTO     `json:"personaMessage"`
	ConsentRisk      ConsentRiskDTO `json:"consentRisk"`
	MessageCount     int            `json:"messageCount"`
	ShouldSuggestEnd bool           `json:"shouldSuggestEnd"`
}

func (uc *SessionUseCase) SendMessage(ctx context.Context, userID, sessionID string, in SendMessageInput) (*SendMessageOutput, error) {
	session, err := uc.loadActiveSession(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	history, err := uc.messages.ListBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	turn, err := uc.llm.ProcessMessage(ctx, session.SystemPrompt, toChatHistory(history), in.Text)
	if err != nil {
		return nil, err
	}
	status := turn.Score.Status
	label := turn.Score.StatusLabel
	userMsg := &domain.Message{
		ID: newMessageID(), SessionID: sessionID, Sender: domain.SenderUser, Text: in.Text,
		Status: &status, StatusLabel: &label,
		Clarity: &turn.Score.Clarity, Confidence: &turn.Score.Confidence,
		Respect: &turn.Score.Respect, Balance: &turn.Score.Balance,
		ConsentRisk: turn.Consent.Detected,
	}
	if err := uc.messages.Create(ctx, userMsg); err != nil {
		return nil, err
	}
	personaMsg := &domain.Message{
		ID: newMessageID(), SessionID: sessionID, Sender: domain.SenderPersona, Text: turn.PersonaText,
	}
	if err := uc.messages.Create(ctx, personaMsg); err != nil {
		return nil, err
	}
	count, _ := uc.messages.CountBySession(ctx, sessionID)
	return &SendMessageOutput{
		UserMessage: MessageDTO{
			ID: userMsg.ID, Sender: domain.SenderUser, Text: in.Text,
			Status: &status, StatusLabel: &label,
		},
		PersonaMessage: MessageDTO{
			ID: personaMsg.ID, Sender: domain.SenderPersona, Text: turn.PersonaText,
		},
		ConsentRisk: ConsentRiskDTO{
			Detected: turn.Consent.Detected, Severity: turn.Consent.Severity,
			Explanation: turn.Consent.Explanation, Suggestion: turn.Consent.Suggestion,
		},
		MessageCount: count, ShouldSuggestEnd: count >= 8,
	}, nil
}

type SuggestInput struct {
	Draft string `json:"draft"`
}

type SuggestOutput struct {
	Suggestion string `json:"suggestion"`
}

func (uc *SessionUseCase) Suggest(ctx context.Context, userID, sessionID string, in SuggestInput) (*SuggestOutput, error) {
	session, err := uc.loadActiveSession(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	history, err := uc.messages.ListBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	suggestion, err := uc.llm.SuggestReply(ctx, session.SystemPrompt, toChatHistory(history), in.Draft)
	if err != nil {
		return nil, err
	}
	return &SuggestOutput{Suggestion: suggestion}, nil
}

type FinishOutput struct {
	SessionID       string                 `json:"sessionId"`
	FinishedAt      string                 `json:"finishedAt"`
	Scenario        string                 `json:"scenario"`
	Persona         string                 `json:"persona"`
	Scores          domain.DebriefScores   `json:"scores"`
	Strengths       []string               `json:"strengths"`
	Weaknesses      []string               `json:"weaknesses"`
	RiskFlags       []domain.RiskFlag      `json:"riskFlags"`
	ImprovedReplies []domain.ImprovedReply `json:"improvedReplies"`
	TipForNext      string                 `json:"tipForNext"`
	HasRisk         bool                   `json:"hasRisk"`
}

func (uc *SessionUseCase) Finish(ctx context.Context, userID, sessionID string) (*FinishOutput, error) {
	session, err := uc.loadActiveSession(ctx, userID, sessionID)
	if err != nil {
		return nil, err
	}
	history, err := uc.messages.ListBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	summary, err := uc.llm.SummarizeSession(ctx, session.SystemPrompt, toChatHistory(history))
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	session.Status = domain.SessionStatusFinished
	session.FinishedAt = &now
	if err := uc.sessions.Update(ctx, session); err != nil {
		return nil, err
	}
	scoresJSON, _ := json.Marshal(summary.Scores)
	strengthsJSON, _ := json.Marshal(summary.Strengths)
	weaknessesJSON, _ := json.Marshal(summary.Weaknesses)
	riskJSON, _ := json.Marshal(summary.RiskFlags)
	improvedJSON, _ := json.Marshal(summary.ImprovedReplies)
	debrief := &domain.Debrief{
		ID: newDebriefID(), SessionID: sessionID,
		ScoresJSON: scoresJSON, StrengthsJSON: strengthsJSON, WeaknessesJSON: weaknessesJSON,
		RiskFlagsJSON: riskJSON, ImprovedRepliesJSON: improvedJSON,
		TipForNext: summary.TipForNext, HasRisk: summary.HasRisk,
	}
	if err := uc.debriefs.Create(ctx, debrief); err != nil {
		return nil, err
	}
	return buildFinishOutput(sessionID, now, session, summary), nil
}

func (uc *SessionUseCase) GetDebrief(ctx context.Context, userID, sessionID string) (*FinishOutput, error) {
	session, err := uc.sessions.GetByIDForUser(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	debrief, err := uc.debriefs.GetBySessionID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if debrief == nil {
		return nil, errors.New("debrief not found")
	}
	var scores domain.DebriefScores
	_ = json.Unmarshal(debrief.ScoresJSON, &scores)
	var strengths, weaknesses []string
	var riskFlags []domain.RiskFlag
	var improved []domain.ImprovedReply
	_ = json.Unmarshal(debrief.StrengthsJSON, &strengths)
	_ = json.Unmarshal(debrief.WeaknessesJSON, &weaknesses)
	_ = json.Unmarshal(debrief.RiskFlagsJSON, &riskFlags)
	_ = json.Unmarshal(debrief.ImprovedRepliesJSON, &improved)
	if strengths == nil {
		strengths = []string{}
	}
	if weaknesses == nil {
		weaknesses = []string{}
	}
	if riskFlags == nil {
		riskFlags = []domain.RiskFlag{}
	}
	if improved == nil {
		improved = []domain.ImprovedReply{}
	}
	finishedAt := time.Now().UTC()
	if session.FinishedAt != nil {
		finishedAt = *session.FinishedAt
	}
	summary := &domain.SessionSummary{
		Scores: scores, Strengths: strengths, Weaknesses: weaknesses,
		RiskFlags: riskFlags, ImprovedReplies: improved,
		TipForNext: debrief.TipForNext, HasRisk: debrief.HasRisk,
	}
	return buildFinishOutput(sessionID, finishedAt, session, summary), nil
}

type SessionListItem struct {
	ID       string             `json:"id"`
	Date     string             `json:"date"`
	Time     string             `json:"time"`
	Scenario string             `json:"scenario"`
	Persona  string             `json:"persona"`
	Scores   domain.SkillScores `json:"scores"`
	HasRisk  bool               `json:"hasRisk"`
}

type SessionsListOutput struct {
	FocusSkill        string                `json:"focusSkill"`
	FocusSkillLabel   string                `json:"focusSkillLabel"`
	FocusSessionCount int                   `json:"focusSessionCount"`
	ClarityHistory    []int                 `json:"clarityHistory"`
	DailyExercise     *domain.DailyExercise `json:"dailyExercise"`
	Sessions          []SessionListItem     `json:"sessions"`
}

func (uc *SessionUseCase) List(ctx context.Context, userID string) (*SessionsListOutput, error) {
	avg, err := uc.messages.AverageSkillsForUser(ctx, userID, 10)
	if err != nil {
		return nil, err
	}
	focus := weakestSkill(avg)
	clarityHist, err := uc.messages.ClarityHistory(ctx, userID, 8)
	if err != nil {
		return nil, err
	}
	exercise, err := uc.llm.GenerateDailyExercise(ctx, catalog.SkillLabel(focus))
	if err != nil {
		exercise = &domain.DailyExercise{
			Title: "Дай конкретику", Description: "Ответь одной яркой деталью.",
			Criterion: "ясность 8+",
		}
	}
	sessions, err := uc.sessions.ListByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	items := make([]SessionListItem, 0)
	focusCount := 0
	for _, s := range sessions {
		if s.Status != domain.SessionStatusFinished {
			continue
		}
		debrief, err := uc.debriefs.GetBySessionID(ctx, s.ID)
		if err != nil || debrief == nil {
			continue
		}
		var scores domain.DebriefScores
		_ = json.Unmarshal(debrief.ScoresJSON, &scores)
		sk := sessionScoresFromDebrief(scores)
		date, timeStr := formatSessionDate(s.StartedAt)
		if s.FinishedAt != nil {
			date, timeStr = formatSessionDate(*s.FinishedAt)
		}
		items = append(items, SessionListItem{
			ID: s.ID, Date: date, Time: timeStr,
			Scenario: goalTitle(s.GoalID), Persona: personaTitle(s.PersonaID),
			Scores: sk, HasRisk: debrief.HasRisk,
		})
	}
	for _, it := range items {
		if it.Scores.Clarity <= 6 && focus == domain.SkillClarity {
			focusCount++
		}
	}
	return &SessionsListOutput{
		FocusSkill: focus, FocusSkillLabel: catalog.SkillLabel(focus),
		FocusSessionCount: focusCount, ClarityHistory: clarityHist,
		DailyExercise: exercise, Sessions: items,
	}, nil
}

func (uc *SessionUseCase) loadActiveSession(ctx context.Context, userID, sessionID string) (*domain.Session, error) {
	session, err := uc.sessions.GetByIDForUser(ctx, sessionID, userID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	if session.Status == domain.SessionStatusFinished {
		return nil, ErrSessionFinished
	}
	return session, nil
}

func buildFinishOutput(sessionID string, finishedAt time.Time, session *domain.Session, summary *domain.SessionSummary) *FinishOutput {
	strengths := summary.Strengths
	if strengths == nil {
		strengths = []string{}
	}
	weaknesses := summary.Weaknesses
	if weaknesses == nil {
		weaknesses = []string{}
	}
	riskFlags := summary.RiskFlags
	if riskFlags == nil {
		riskFlags = []domain.RiskFlag{}
	}
	improved := summary.ImprovedReplies
	if improved == nil {
		improved = []domain.ImprovedReply{}
	}
	return &FinishOutput{
		SessionID: sessionID, FinishedAt: finishedAt.Format(time.RFC3339),
		Scenario: goalTitle(session.GoalID), Persona: personaTitle(session.PersonaID),
		Scores: summary.Scores, Strengths: strengths, Weaknesses: weaknesses,
		RiskFlags: riskFlags, ImprovedReplies: improved,
		TipForNext: summary.TipForNext, HasRisk: summary.HasRisk,
	}
}
