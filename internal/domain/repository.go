package domain

import "context"

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id string) (*User, error)
}

type SessionRepository interface {
	Create(ctx context.Context, session *Session) error
	GetByID(ctx context.Context, id string) (*Session, error)
	GetByIDForUser(ctx context.Context, id, userID string) (*Session, error)
	Update(ctx context.Context, session *Session) error
	ListByUser(ctx context.Context, userID string) ([]Session, error)
	CountByUser(ctx context.Context, userID string) (int, error)
}

type MessageRepository interface {
	Create(ctx context.Context, msg *Message) error
	ListBySession(ctx context.Context, sessionID string) ([]Message, error)
	CountBySession(ctx context.Context, sessionID string) (int, error)
	AverageSkillsForUser(ctx context.Context, userID string, limit int) (map[string]float64, error)
	ClarityHistory(ctx context.Context, userID string, limit int) ([]int, error)
}

type DebriefRepository interface {
	Create(ctx context.Context, debrief *Debrief) error
	GetBySessionID(ctx context.Context, sessionID string) (*Debrief, error)
}
