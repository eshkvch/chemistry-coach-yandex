package postgres

import (
	"time"

	"chemistry-coach/internal/domain"
)

type UserModel struct {
	ID        string    `gorm:"primaryKey;column:id"`
	BirthYear int       `gorm:"column:birth_year;not null"`
	Initials  string    `gorm:"column:initials;not null;default:ИВ"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (UserModel) TableName() string { return "users" }

func (m *UserModel) ToDomain() *domain.User {
	return &domain.User{
		ID: m.ID, BirthYear: m.BirthYear, Initials: m.Initials, CreatedAt: m.CreatedAt,
	}
}

type SessionModel struct {
	ID           string     `gorm:"primaryKey;column:id"`
	UserID       string     `gorm:"column:user_id;not null"`
	GoalID       string     `gorm:"column:goal_id;not null"`
	PersonaID    string     `gorm:"column:persona_id;not null"`
	Status       string     `gorm:"column:status;default:active"`
	SystemPrompt string     `gorm:"column:system_prompt"`
	StartedAt    time.Time  `gorm:"column:started_at;autoCreateTime"`
	FinishedAt   *time.Time `gorm:"column:finished_at"`
}

func (SessionModel) TableName() string { return "sessions" }

func (m *SessionModel) ToDomain() *domain.Session {
	return &domain.Session{
		ID: m.ID, UserID: m.UserID, GoalID: m.GoalID, PersonaID: m.PersonaID,
		Status: m.Status, SystemPrompt: m.SystemPrompt, StartedAt: m.StartedAt, FinishedAt: m.FinishedAt,
	}
}

type MessageModel struct {
	ID          string    `gorm:"primaryKey;column:id"`
	SessionID   string    `gorm:"column:session_id;not null"`
	Sender      string    `gorm:"column:sender;not null"`
	Text        string    `gorm:"column:text;not null"`
	Status      *string   `gorm:"column:status"`
	StatusLabel *string   `gorm:"column:status_label"`
	Clarity     *int      `gorm:"column:clarity"`
	Confidence  *int      `gorm:"column:confidence"`
	Respect     *int      `gorm:"column:respect"`
	Balance     *int      `gorm:"column:balance"`
	ConsentRisk bool      `gorm:"column:consent_risk;default:false"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (MessageModel) TableName() string { return "messages" }

func (m *MessageModel) ToDomain() *domain.Message {
	return &domain.Message{
		ID: m.ID, SessionID: m.SessionID, Sender: m.Sender, Text: m.Text,
		Status: m.Status, StatusLabel: m.StatusLabel,
		Clarity: m.Clarity, Confidence: m.Confidence, Respect: m.Respect, Balance: m.Balance,
		ConsentRisk: m.ConsentRisk, CreatedAt: m.CreatedAt,
	}
}

type DebriefModel struct {
	ID                  string    `gorm:"primaryKey;column:id"`
	SessionID           string    `gorm:"column:session_id;unique;not null"`
	ScoresJSON          []byte    `gorm:"column:scores_json;type:jsonb;not null"`
	StrengthsJSON       []byte    `gorm:"column:strengths_json;type:jsonb;not null"`
	WeaknessesJSON      []byte    `gorm:"column:weaknesses_json;type:jsonb;not null"`
	RiskFlagsJSON       []byte    `gorm:"column:risk_flags_json;type:jsonb;not null"`
	ImprovedRepliesJSON []byte    `gorm:"column:improved_replies_json;type:jsonb;not null"`
	TipForNext          string    `gorm:"column:tip_for_next"`
	HasRisk             bool      `gorm:"column:has_risk;default:false"`
	CreatedAt           time.Time `gorm:"column:created_at;autoCreateTime"`
}

func (DebriefModel) TableName() string { return "debriefs" }
