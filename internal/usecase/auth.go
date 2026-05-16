package usecase

import (
	"chemistry-coach/internal/domain"
	"context"
	"errors"
	"time"
)

var ErrUnderage = errors.New("underage")

type AuthUseCase struct {
	users domain.UserRepository
}

func NewAuthUseCase(users domain.UserRepository) *AuthUseCase {
	return &AuthUseCase{users: users}
}

type AuthStartInput struct {
	BirthYear int `json:"birthYear"`
	Consents  struct {
		IsAdult       bool `json:"isAdult"`
		UnderstandsAI bool `json:"understandsAI"`
		AgreesToRules bool `json:"agreesToRules"`
	}
	UserID string // optional: return existing
}

type AuthStartOutput struct {
	UserID    string `json:"userId"`
	IsNewUser bool   `json:"isNewUser"`
}

func (uc *AuthUseCase) Start(ctx context.Context, in AuthStartInput) (*AuthStartOutput, error) {
	age := time.Now().Year() - in.BirthYear
	if age < 18 {
		return nil, ErrUnderage
	}
	if !in.Consents.IsAdult || !in.Consents.UnderstandsAI || !in.Consents.AgreesToRules {
		return nil, errors.New("consents required")
	}
	if in.UserID != "" {
		u, err := uc.users.GetByID(ctx, in.UserID)
		if err != nil {
			return nil, err
		}
		if u != nil {
			return &AuthStartOutput{UserID: u.ID, IsNewUser: false}, nil
		}
	}
	id := newUserID()
	user := &domain.User{ID: id, BirthYear: in.BirthYear, Initials: "ИВ"}
	if err := uc.users.Create(ctx, user); err != nil {
		return nil, err
	}
	return &AuthStartOutput{UserID: id, IsNewUser: true}, nil
}
