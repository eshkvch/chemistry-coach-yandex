package domain

import "time"

const (
	SessionStatusActive   = "active"
	SessionStatusFinished = "finished"

	SenderUser    = "user"
	SenderPersona = "persona"

	SkillClarity     = "clarity"
	SkillConfidence  = "confidence"
	SkillRespect     = "respect"
	SkillBalance     = "balance"

	StatusSuccess = "success"
	StatusWarning = "warning"
	StatusDanger  = "danger"
)

type User struct {
	ID        string
	BirthYear int
	Initials  string
	CreatedAt time.Time
}

type Session struct {
	ID           string
	UserID       string
	GoalID       string
	PersonaID    string
	Status       string
	SystemPrompt string
	StartedAt    time.Time
	FinishedAt   *time.Time
}

type Message struct {
	ID          string
	SessionID   string
	Sender      string
	Text        string
	Status      *string
	StatusLabel *string
	Clarity     *int
	Confidence  *int
	Respect     *int
	Balance     *int
	ConsentRisk bool
	CreatedAt   time.Time
}

type Debrief struct {
	ID                string
	SessionID         string
	ScoresJSON        []byte
	StrengthsJSON     []byte
	WeaknessesJSON    []byte
	RiskFlagsJSON     []byte
	ImprovedRepliesJSON []byte
	TipForNext        string
	HasRisk           bool
	CreatedAt         time.Time
}

type SkillScores struct {
	Clarity     int `json:"clarity"`
	Confidence  int `json:"confidence"`
	Respect     int `json:"respect"`
	Balance     int `json:"balance"`
}

type ScoreDetail struct {
	Value   int    `json:"value"`
	Comment string `json:"comment"`
}

type DebriefScores struct {
	Clarity    ScoreDetail `json:"clarity"`
	Confidence ScoreDetail `json:"confidence"`
	Respect    ScoreDetail `json:"respect"`
	Balance    ScoreDetail `json:"balance"`
}

type RiskFlag struct {
	Severity    string `json:"severity"`
	Quote       string `json:"quote"`
	Explanation string `json:"explanation"`
	Suggestion  string `json:"suggestion"`
}

type ImprovedReply struct {
	Original string `json:"original"`
	Improved string `json:"improved"`
	Reason   string `json:"reason"`
}

type MessageScore struct {
	Clarity     int    `json:"clarity"`
	Confidence  int    `json:"confidence"`
	Respect     int    `json:"respect"`
	Balance     int    `json:"balance"`
	StatusLabel string `json:"statusLabel"`
	Status      string `json:"status"`
}

type ConsentRisk struct {
	Detected    bool    `json:"detected"`
	Severity    *string `json:"severity"`
	Explanation *string `json:"explanation"`
	Suggestion  *string `json:"suggestion"`
}

type ChatTurnResult struct {
	PersonaText string
	Score       MessageScore
	Consent     ConsentRisk
}

type SessionSummary struct {
	Scenario        string
	Persona         string
	Scores          DebriefScores
	Strengths       []string
	Weaknesses      []string
	RiskFlags       []RiskFlag
	ImprovedReplies []ImprovedReply
	TipForNext      string
	HasRisk         bool
}

type DailyExercise struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Criterion   string `json:"criterion"`
}
