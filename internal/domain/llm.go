package domain

import "context"

type ChatMessage struct {
	Role string
	Text string
}

type LLMService interface {
	GenerateOpening(ctx context.Context, systemPrompt string) (string, error)
	ProcessMessage(ctx context.Context, systemPrompt string, history []ChatMessage, userText string) (*ChatTurnResult, error)
	SuggestReply(ctx context.Context, systemPrompt string, history []ChatMessage, draft string) (string, error)
	SummarizeSession(ctx context.Context, systemPrompt string, history []ChatMessage) (*SessionSummary, error)
	GenerateDailyExercise(ctx context.Context, focusSkill string) (*DailyExercise, error)
}
