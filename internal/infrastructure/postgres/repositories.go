package postgres

import (
	"context"
	"errors"
	"fmt"

	"chemistry-coach/internal/domain"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db: db} }

func (r *UserRepo) Create(ctx context.Context, user *domain.User) error {
	m := &UserModel{ID: user.ID, BirthYear: user.BirthYear, Initials: user.Initials}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	var m UserModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

type SessionRepo struct{ db *gorm.DB }

func NewSessionRepo(db *gorm.DB) *SessionRepo { return &SessionRepo{db: db} }

func (r *SessionRepo) Create(ctx context.Context, s *domain.Session) error {
	m := &SessionModel{
		ID: s.ID, UserID: s.UserID, GoalID: s.GoalID, PersonaID: s.PersonaID,
		Status: s.Status, SystemPrompt: s.SystemPrompt,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *SessionRepo) GetByID(ctx context.Context, id string) (*domain.Session, error) {
	var m SessionModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *SessionRepo) GetByIDForUser(ctx context.Context, id, userID string) (*domain.Session, error) {
	var m SessionModel
	if err := r.db.WithContext(ctx).First(&m, "id = ? AND user_id = ?", id, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return m.ToDomain(), nil
}

func (r *SessionRepo) Update(ctx context.Context, s *domain.Session) error {
	return r.db.WithContext(ctx).Model(&SessionModel{}).Where("id = ?", s.ID).Updates(map[string]interface{}{
		"status":        s.Status,
		"system_prompt": s.SystemPrompt,
		"finished_at":   s.FinishedAt,
	}).Error
}

func (r *SessionRepo) ListByUser(ctx context.Context, userID string) ([]domain.Session, error) {
	var models []SessionModel
	if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("started_at DESC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Session, len(models))
	for i, m := range models {
		out[i] = *m.ToDomain()
	}
	return out, nil
}

func (r *SessionRepo) CountByUser(ctx context.Context, userID string) (int, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&SessionModel{}).Where("user_id = ?", userID).Count(&n).Error
	return int(n), err
}

type MessageRepo struct{ db *gorm.DB }

func NewMessageRepo(db *gorm.DB) *MessageRepo { return &MessageRepo{db: db} }

func (r *MessageRepo) Create(ctx context.Context, msg *domain.Message) error {
	m := &MessageModel{
		ID: msg.ID, SessionID: msg.SessionID, Sender: msg.Sender, Text: msg.Text,
		Status: msg.Status, StatusLabel: msg.StatusLabel,
		Clarity: msg.Clarity, Confidence: msg.Confidence, Respect: msg.Respect, Balance: msg.Balance,
		ConsentRisk: msg.ConsentRisk,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *MessageRepo) ListBySession(ctx context.Context, sessionID string) ([]domain.Message, error) {
	var models []MessageModel
	if err := r.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("created_at ASC").Find(&models).Error; err != nil {
		return nil, err
	}
	out := make([]domain.Message, len(models))
	for i, m := range models {
		out[i] = *m.ToDomain()
	}
	return out, nil
}

func (r *MessageRepo) CountBySession(ctx context.Context, sessionID string) (int, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&MessageModel{}).Where("session_id = ? AND sender IN ?", sessionID, []string{domain.SenderUser, domain.SenderPersona}).Count(&n).Error
	return int(n), err
}

func (r *MessageRepo) AverageSkillsForUser(ctx context.Context, userID string, limit int) (map[string]float64, error) {
	type row struct {
		Clarity    *float64
		Confidence *float64
		Respect    *float64
		Balance    *float64
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT m.clarity, m.confidence, m.respect, m.balance
		FROM messages m
		JOIN sessions s ON s.id = m.session_id
		WHERE s.user_id = ? AND s.status = ? AND m.sender = ?
		ORDER BY m.created_at DESC
		LIMIT ?
	`, userID, domain.SessionStatusFinished, domain.SenderUser, limit*4).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	sum := map[string]float64{"clarity": 0, "confidence": 0, "respect": 0, "balance": 0}
	cnt := map[string]int{"clarity": 0, "confidence": 0, "respect": 0, "balance": 0}
	for _, rw := range rows {
		if rw.Clarity != nil {
			sum["clarity"] += *rw.Clarity
			cnt["clarity"]++
		}
		if rw.Confidence != nil {
			sum["confidence"] += *rw.Confidence
			cnt["confidence"]++
		}
		if rw.Respect != nil {
			sum["respect"] += *rw.Respect
			cnt["respect"]++
		}
		if rw.Balance != nil {
			sum["balance"] += *rw.Balance
			cnt["balance"]++
		}
	}
	avg := make(map[string]float64)
	for k, c := range cnt {
		if c > 0 {
			avg[k] = sum[k] / float64(c)
		}
	}
	return avg, nil
}

func (r *MessageRepo) ClarityHistory(ctx context.Context, userID string, limit int) ([]int, error) {
	type row struct {
		Clarity int
	}
	var rows []row
	err := r.db.WithContext(ctx).Raw(`
		SELECT AVG(m.clarity)::int AS clarity
		FROM sessions s
		JOIN messages m ON m.session_id = s.id AND m.sender = ?
		WHERE s.user_id = ? AND s.status = ? AND m.clarity IS NOT NULL
		GROUP BY s.id
		ORDER BY MAX(s.finished_at) DESC NULLS LAST
		LIMIT ?
	`, domain.SenderUser, userID, domain.SessionStatusFinished, limit).Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	out := make([]int, len(rows))
	for i, rw := range rows {
		out[i] = rw.Clarity
	}
	// reverse to chronological order for sparkline
	for i, j := 0, len(out)-1; i < j; i, j = i+1, j-1 {
		out[i], out[j] = out[j], out[i]
	}
	return out, nil
}

type DebriefRepo struct{ db *gorm.DB }

func NewDebriefRepo(db *gorm.DB) *DebriefRepo { return &DebriefRepo{db: db} }

func (r *DebriefRepo) Create(ctx context.Context, d *domain.Debrief) error {
	m := &DebriefModel{
		ID: d.ID, SessionID: d.SessionID, ScoresJSON: d.ScoresJSON,
		StrengthsJSON: d.StrengthsJSON, WeaknessesJSON: d.WeaknessesJSON,
		RiskFlagsJSON: d.RiskFlagsJSON, ImprovedRepliesJSON: d.ImprovedRepliesJSON,
		TipForNext: d.TipForNext, HasRisk: d.HasRisk,
	}
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *DebriefRepo) GetBySessionID(ctx context.Context, sessionID string) (*domain.Debrief, error) {
	var m DebriefModel
	if err := r.db.WithContext(ctx).First(&m, "session_id = ?", sessionID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &domain.Debrief{
		ID: m.ID, SessionID: m.SessionID, ScoresJSON: m.ScoresJSON,
		StrengthsJSON: m.StrengthsJSON, WeaknessesJSON: m.WeaknessesJSON,
		RiskFlagsJSON: m.RiskFlagsJSON, ImprovedRepliesJSON: m.ImprovedRepliesJSON,
		TipForNext: m.TipForNext, HasRisk: m.HasRisk, CreatedAt: m.CreatedAt,
	}, nil
}

var _ domain.UserRepository = (*UserRepo)(nil)
var _ domain.SessionRepository = (*SessionRepo)(nil)
var _ domain.MessageRepository = (*MessageRepo)(nil)
var _ domain.DebriefRepository = (*DebriefRepo)(nil)

func ErrNotFound() error { return fmt.Errorf("not found") }
