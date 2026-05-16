package usecase

import (
	"chemistry-coach/internal/catalog"
	"chemistry-coach/internal/domain"
	"context"
	"errors"
)

type ProfileUseCase struct {
	users    domain.UserRepository
	sessions domain.SessionRepository
	messages domain.MessageRepository
}

func NewProfileUseCase(users domain.UserRepository, sessions domain.SessionRepository, messages domain.MessageRepository) *ProfileUseCase {
	return &ProfileUseCase{users: users, sessions: sessions, messages: messages}
}

type ProfileOutput struct {
	UserID            string `json:"userId"`
	WeakestSkill      string `json:"weakestSkill"`
	RecommendedGoalID string `json:"recommendedGoalId"`
	TotalSessions     int    `json:"totalSessions"`
	Initials          string `json:"initials"`
}

func (uc *ProfileUseCase) Get(ctx context.Context, userID string) (*ProfileOutput, error) {
	user, err := uc.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	avg, err := uc.messages.AverageSkillsForUser(ctx, userID, 5)
	if err != nil {
		return nil, err
	}
	ws := weakestSkill(avg)
	total, err := uc.sessions.CountByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &ProfileOutput{
		UserID: userID, WeakestSkill: ws,
		RecommendedGoalID: catalog.RecommendedGoalForSkill(ws),
		TotalSessions:     total, Initials: user.Initials,
	}, nil
}
